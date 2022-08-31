package mongo_data_source

import (
	"context"
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/constants"
	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/errors"
	"github.com/mjolnir-mud/engine/plugins/mongo_data_source/internal/logger"
	constants2 "github.com/mjolnir-mud/engine/plugins/mongo_data_source/pkg/constants"
	errors2 "github.com/mjolnir-mud/engine/plugins/mongo_data_source/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Configuration struct {
	MongoURL string
	Database string
}

var Configs = map[string]func(c *Configuration) *Configuration{
	"default": func(c *Configuration) *Configuration {
		c.MongoURL = "mongodb://localhost:27017"
		c.Database = "mjolnir_dev"

		return c
	},
}

type plugin struct {
	database *mongo.Database
}

// ConfigureForEnv sets the config for the plugin for the specified environment.
func ConfigureForEnv(env string, cb func(c *Configuration) *Configuration) {
	Configs[env] = cb
}

func (p *plugin) Name() string {
	return "mongo_data_source"
}

func (p *plugin) Registered() error {
	engine.RegisterOnServiceStartCallback("world", func() {
		logger.Start()
		log = logger.Instance
		env := engine.GetEnv()
		defaultConfig := Configs["default"](&Configuration{})
		config := Configs[env](defaultConfig)

		if config == nil {
			log.Fatal().Msg("no config for environment")
			panic("no config for environment")
		}

		log = logger.
			Instance.
			With().
			Str("mongo_url", config.MongoURL).
			Str("database", config.Database).
			Logger()

		log.Info().Msg("starting mongo connection")
		c, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURL))

		if err != nil {
			log.Fatal().Err(err).Msg("error connecting to mongo")
			panic(err)
		}

		err = c.Connect(context.Background())

		if err != nil {
			log.Fatal().Err(err).Msg("error connecting to mongo")
			panic(err)
		}

		p.database = c.Database(config.Database)

		log.Info().Msg("mongo connection started")
	})

	engine.RegisterBeforeStopCallback(func() {
		log.Info().Msg("stopping mongo connection")
		_ = p.database.Client().Disconnect(context.Background())
	})

	return nil
}

func (p *plugin) collection(name string) *mongo.Collection {
	return p.database.Collection(name)
}

var log zerolog.Logger

var Plugin = &plugin{}

type MongoDataSource struct {
	collectionName string
	collection     *mongo.Collection
	logger         zerolog.Logger
}

func New(collection string) *MongoDataSource {
	return &MongoDataSource{
		collectionName: collection,
		logger:         log.With().Str("collection", collection).Logger(),
	}
}

func (m *MongoDataSource) Name() string {
	return m.collectionName
}

func (m *MongoDataSource) Start() error {
	m.collection = Plugin.collection(m.collectionName)

	return nil
}

func (m *MongoDataSource) Stop() error {
	return nil
}

func (m *MongoDataSource) All() (map[string]map[string]interface{}, error) {
	m.logger.Debug().Msg("loading all entities")

	cursor, err := m.collection.Find(context.Background(), map[string]interface{}{})

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

func (m *MongoDataSource) Find(search map[string]interface{}) (map[string]map[string]interface{}, error) {
	m.logger.Debug().Interface("search", search).Msg("searching entities")

	cursor, err := m.collection.Find(context.Background(), search)

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

func (m *MongoDataSource) FindOne(search map[string]interface{}) (string, map[string]interface{}, error) {
	m.logger.Debug().Interface("search", search).Msg("searching one entity")

	result := map[string]interface{}{}

	err := m.collection.FindOne(context.Background(), search).Decode(&result)

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

func (m *MongoDataSource) Delete(id string) error {
	m.logger.Debug().Str("id", id).Msg("deleting entity")

	_ = m.collection.FindOneAndDelete(context.Background(), map[string]interface{}{"_id": id})

	return nil
}

func (m *MongoDataSource) FindAndDelete(search map[string]interface{}) error {
	m.logger.Debug().Interface("search", search).Msg("searching and deleting entity")

	result := map[string]interface{}{}

	err := m.collection.FindOneAndDelete(context.Background(), search).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.EntityNotFoundError{}
		}
		return err
	}

	return nil
}

func (m *MongoDataSource) Count(search map[string]interface{}) (int64, error) {
	m.logger.Debug().Interface("search", search).Msg("counting entities")

	count, err := m.collection.CountDocuments(context.Background(), search)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *MongoDataSource) Save(entityId string, entity map[string]interface{}) error {
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

	_, err := m.collection.InsertOne(context.Background(), entity)
	return err
}

func cleanId(entity map[string]interface{}) {
	delete(entity, "_id")
}
