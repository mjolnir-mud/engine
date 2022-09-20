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

package mongo_data_source

import (
	"context"

	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/constants"
	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/data_source"
	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/errors"
	"github.com/mjolnir-mud/engine/plugins/mongo_data_source/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/mongo_data_source/internal/plugin"
	"github.com/mjolnir-mud/engine/plugins/mongo_data_source/pkg/config"
	constants2 "github.com/mjolnir-mud/engine/plugins/mongo_data_source/pkg/constants"
	errors2 "github.com/mjolnir-mud/engine/plugins/mongo_data_source/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

var Plugin = plugin.Plugin

type MongoDataSource struct {
	collectionName string
	logger         zerolog.Logger
}

func New(collection string) data_source.DataSource {
	return MongoDataSource{
		collectionName: collection,
		logger:         logger.Instance.With().Str("collection", collection).Logger(),
	}
}

// ConfigureForEnv sets the config for the plugin for the specified environment.
func ConfigureForEnv(env string, cb func(c *config.Configuration) *config.Configuration) {
	plugin.ConfigureForEnv(env, cb)
}

func (m MongoDataSource) Name() string {
	return m.collectionName
}

func (m MongoDataSource) Start() error {
	return nil
}

func (m MongoDataSource) Stop() error {
	return nil
}

func (m MongoDataSource) All() (map[string]map[string]interface{}, error) {
	m.logger.Debug().Msg("loading all entities")

	cursor, err := m.getCollection().Find(context.Background(), map[string]interface{}{})

	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	err = cursor.All(context.Background(), &results)

	if err != nil {
		return nil, err
	}

	entities := make(map[string]map[string]interface{})

	for _, entity := range results {
		id := entity["_id"].(string)
		cleanId(entity)
		entities[id] = entity
	}

	return entities, nil
}

func (m MongoDataSource) Find(search map[string]interface{}) (map[string]map[string]interface{}, error) {
	m.logger.Debug().Interface("search", search).Msg("searching entities")

	cursor, err := m.getCollection().Find(context.Background(), search)

	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	err = cursor.All(context.Background(), &results)

	if err != nil {
		return nil, err
	}

	entities := make(map[string]map[string]interface{})

	for _, entity := range results {
		id := entity["_id"].(string)
		cleanId(entity)
		entities[id] = entity
	}

	return entities, nil
}

func (m MongoDataSource) FindOne(search map[string]interface{}) (string, map[string]interface{}, error) {
	m.logger.Debug().Interface("search", search).Msg("searching one entity")

	result := map[string]interface{}{}

	err := m.getCollection().FindOne(context.Background(), search).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil, errors.EntityNotFoundError{}
		}
		return "", nil, err
	}

	id := result["_id"].(string)

	cleanId(result)

	result["id"] = id

	return id, result, nil
}

func (m MongoDataSource) Delete(id string) error {
	m.logger.Debug().Str("id", id).Msg("deleting entity")

	_ = m.getCollection().FindOneAndDelete(context.Background(), map[string]interface{}{"_id": id})

	return nil
}

func (m MongoDataSource) FindAndDelete(search map[string]interface{}) error {
	m.logger.Debug().Interface("search", search).Msg("searching and deleting entity")

	result := map[string]interface{}{}

	err := m.getCollection().FindOneAndDelete(context.Background(), search).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.EntityNotFoundError{}
		}
		return err
	}

	return nil
}

func (m MongoDataSource) Count(search map[string]interface{}) (int64, error) {
	m.logger.Debug().Interface("search", search).Msg("counting entities")

	count, err := m.getCollection().CountDocuments(context.Background(), search)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m MongoDataSource) Save(entityId string, entity map[string]interface{}) error {
	metadata, ok := entity[constants.MetadataKey]

	if !ok {
		return errors.MetadataRequiredError{ID: entityId}
	}

	collection, ok := metadata.(map[string]interface{})[constants2.MetadataCollectionKey]

	if !ok {
		return errors2.CollectionMetadataRequiredError{ID: entityId}
	}

	m.logger.Debug().Interface("entity", entity).Msg("saving entity")

	if collection != m.collectionName {
		return errors2.CollectionMismatchError{
			SourceCollection: m.collectionName,
			TargetCollection: collection.(string),
		}
	}

	entity["_id"] = entityId

	_, err := m.getCollection().InsertOne(context.Background(), entity)
	return err
}

func (m MongoDataSource) getCollection() *mongo.Collection {
	return plugin.Plugin.Collection(m.collectionName)
}

func cleanId(entity map[string]interface{}) {
	delete(entity, "_id")
}
