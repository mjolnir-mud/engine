/*
 * Copyright (c) 2022 eightfivefour llc. All rights reserved.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package engine

import (
	"context"
	"fmt"
	"github.com/fatih/structs"
	engineErrors "github.com/mjolnir-engine/engine/errors"
	"github.com/mjolnir-engine/engine/uid"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

type MongoConfiguration struct {
	Host     string
	Port     int
	Database string
}

// MongoDataSource is a data source that uses MongoDB as the backing store. Instances of this data source can be
// created and registered to store data within a specific collection.
type MongoDataSource struct {
	engine     *Engine
	logger     zerolog.Logger
	client     *mongo.Client
	collection string
}

// NewMongoDataSource creates a new MongoDataSource instance for the provided collection.
func NewMongoDataSource(collection string, engine *Engine) DataSource {
	logger := engine.logger.With().Str("collection", collection).Logger()

	c, err := mongo.NewClient(options.
		Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%d", engine.config.Mongo.Host, engine.config.Mongo.Port)),
	)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create new Mongo client")
		panic(err)
	}

	return MongoDataSource{
		collection: collection,
		engine:     engine,
		client:     c,
		logger:     engine.logger.With().Str("collection", collection).Logger(),
	}
}

func (m MongoDataSource) All(results interface{}) error {
	m.logger.Info().Msg("getting all documents")

	err := m.Find(map[string]interface{}{}, results)

	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to get all documents")
		return err
	}

	return nil
}

func (m MongoDataSource) Count(search interface{}) (int64, error) {
	m.logger.Debug().Interface("search", search).Msg("counting entities")
	search = parseMongoSearch(search)

	count, err := m.getCollection().CountDocuments(context.Background(), search)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m MongoDataSource) Delete(search interface{}) error {
	m.logger.Debug().Msg("deleting entity")
	search = parseMongoSearch(search)

	_, _ = m.getCollection().DeleteMany(context.Background(), search)

	return nil
}

func (m MongoDataSource) FindOne(search interface{}, entity interface{}) error {
	m.logger.Debug().Interface("search", search).Msg("searching one entity")
	search = parseMongoSearch(search)

	err := m.getCollection().FindOne(context.Background(), search).Decode(entity)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			if reflect.TypeOf(search).Kind() == reflect.Map {
				return engineErrors.EntityNotFoundError{
					Id: search.(map[string]interface{})["_id"].(uid.UID),
				}
			} else {
				m := structs.Map(search)

				return engineErrors.EntityNotFoundError{
					Id: m["Id"].(uid.UID),
				}
			}
		}
		return err
	}

	return nil
}

func (m MongoDataSource) Find(search interface{}, entities interface{}) error {
	m.logger.Debug().Msg("searching entities")
	search = parseMongoSearch(search)

	cursor, err := m.getCollection().Find(context.Background(), search)

	if err != nil {
		return err
	}

	err = cursor.All(context.Background(), entities)

	if err != nil {
		return err
	}

	return nil
}

func (m MongoDataSource) Name() string {
	return m.collection
}

func (m MongoDataSource) Start() error {
	m.logger.Info().Msg("Starting MongoDataSource")

	err := m.client.Connect(context.Background())

	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to connect to Mongo")
		return err
	}

	connected := make(chan bool)

	go func() {
		for {
			err := m.client.Ping(context.Background(), nil)
			if err == nil {
				connected <- true
				return
			}
		}
	}()

	<-connected

	return nil
}

func (m MongoDataSource) Save(entity interface{}) (uid.UID, error) {
	m.logger.Debug().Interface("entity", entity).Msg("saving entity")
	e := structs.New(entity)
	field, ok := e.FieldOk("Id")
	var id uid.UID

	if ok {
		id = field.Value().(uid.UID)
	} else {
		return (uid.UID)(""), engineErrors.EntityStructMustHaveIdFieldError{}
	}

	if id == (uid.UID)("") {
		one, err := m.getCollection().InsertOne(context.Background(), entity)

		if err != nil {
			return (uid.UID)(""), err
		}

		oid := one.InsertedID.(primitive.ObjectID)

		return uid.FromBSON(oid), nil
	} else {
		idValue := e.Field("Id").Value()
		id := idValue.(uid.UID)

		_, err := m.getCollection().ReplaceOne(context.Background(), bson.M{"_id": id}, entity)

		if err != nil {
			return (uid.UID)(""), err
		}

		return id, nil
	}
}

func (m MongoDataSource) Stop() error {
	m.logger.Info().Msg("Stopping MongoDataSource")
	return m.client.Disconnect(context.Background())
}

func (m MongoDataSource) getCollection() *mongo.Collection {
	return m.getClient().Database(m.engine.config.Mongo.Database).Collection(m.collection)
}

func (m MongoDataSource) setClient(client *mongo.Client) {
	m.client = client
}

func (m MongoDataSource) getClient() *mongo.Client {
	return m.client
}

func parseMongoSearch(search interface{}) interface{} {
	switch search.(type) {
	case map[string]interface{}:
		parseSearchMap(search.(map[string]interface{}))
	}

	return search
}

func parseSearchMap(search map[string]interface{}) map[string]interface{} {
	if search["id"] != nil {
		search["_id"] = search["id"]
		delete(search, "id")
	}

	switch search["_id"].(type) {
	case uid.UID:
		search["_id"] = search["_id"].(uid.UID).ToBSON()
	case string:
		id, err := primitive.ObjectIDFromHex(search["_id"].(string))

		if err == nil {
			panic(err)
		}

		search["_id"] = id
	}

	return search
}
