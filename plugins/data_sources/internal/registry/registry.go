package registry

import (
	"fmt"

	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/constants"
	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/data_source"
	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/errors"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/rs/zerolog/log"
)

type InvalidDataSourceError struct {
	Source string
}

func (e InvalidDataSourceError) Error() string {
	return fmt.Sprintf("data source %s does not exist", e.Source)
}

type MetadataTypeRequiredError struct {
	ID string
}

func (e MetadataTypeRequiredError) Error() string {
	return fmt.Sprintf("data source %s does not return an entity with the type set in the metadata", e.ID)
}

type registry struct {
	dataSources map[string]data_source.DataSource
}

func Start() error {
	log.Info().Msg("starting")
	for _, d := range r.dataSources {
		log.Info().Msgf("starting data source %s", d.Name())
		err := d.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

func Stop() error {
	log.Info().Msg("stopping")
	for _, d := range r.dataSources {
		log.Info().Msgf("stopping data source %s", d.Name())
		err := d.Stop()
		if err != nil {
			return err
		}
	}

	return nil
}

func Register(dataSource data_source.DataSource) {
	log.Info().Msgf("registering data source %s", dataSource.Name())
	r.dataSources[dataSource.Name()] = dataSource
}

// Load loads data from a data source for a given entity. It will call `ecs.Create` passing the map returned by
// the data source. If the data source does not return an entity with the type  set in the metadata, an error will be
// returned. If the data source does not exist, an error will be returned. If the data source does not reference the
// entity, an error will be returned.
func Load(source string, entityId string) (map[string]interface{}, error) {
	if d, ok := r.dataSources[source]; ok {
		data, err := d.Load(entityId)

		if err != nil {
			return nil, err
		}

		data["id"] = entityId
		return loadEntity(data)
	} else {
		return nil, InvalidDataSourceError{Source: source}
	}
}

// LoadAll loads all entities from a data source. It will call `ecs.Create` passing the map returned by the data source
// for each entity, and return a map of entities keyed by their ids.
func LoadAll(source string) (map[string]map[string]interface{}, error) {
	if d, ok := r.dataSources[source]; ok {
		entities, err := d.LoadAll()

		if err != nil {
			return nil, err
		}

		return loadEntities(entities)
	} else {
		return nil, InvalidDataSourceError{Source: source}
	}
}

// Find returns a list of entities from executing a search against a provided map. It returns a list of entities as a
// map keyed by their ids. It passes each entity found through `ecs.create`, If the entity does not have valid metadata,
// or that metadata does not define the entity type, an error is thrown. If the data source does not exist, an error
// will be thrown.
func Find(source string, search map[string]interface{}) (map[string]map[string]interface{}, error) {
	if d, ok := r.dataSources[source]; ok {
		entities, err := d.Find(search)

		if err != nil {
			return nil, err
		}

		return loadEntities(entities)
	} else {
		return nil, InvalidDataSourceError{Source: source}
	}
}

// Save saves data to a data source for a given entity. If the entity does not have a valid metadata field an error will
// be thrown. If the data source does not exist, an error will be thrown. If the metadata field does not have a type
// set, an error will be thrown. If the entity exists in the data source, it will be overwritten.
func Save(source string, entityId string, entity map[string]interface{}) error {
	if d, ok := r.dataSources[source]; ok {
		metadata, ok := entity[constants.MetadataKey].(map[string]interface{})

		if !ok {
			return errors.MetadataRequiredError{ID: entityId}
		}

		_, ok = metadata[constants.MetadataTypeKey].(string)

		if !ok {
			return MetadataTypeRequiredError{ID: entityId}
		}

		return d.Save(entityId, entity)
	} else {
		return InvalidDataSourceError{Source: source}
	}
}

func Count(source string, search map[string]interface{}) (int64, error) {
	if d, ok := r.dataSources[source]; ok {
		return d.Count(search)
	} else {
		return 0, InvalidDataSourceError{Source: source}
	}
}

func loadEntities(entities map[string]map[string]interface{}) (map[string]map[string]interface{}, error) {
	result := make(map[string]map[string]interface{})

	for id, entity := range entities {
		entity["id"] = id
		loadedEntity, err := loadEntity(entity)

		if err != nil {
			return nil, err
		}

		result[id] = loadedEntity
	}

	return result, nil
}

func loadEntity(entity map[string]interface{}) (map[string]interface{}, error) {
	id, ok := entity["id"]

	if !ok {
		return nil, errors.IDRequiredError{}
	}

	metadata, ok := entity[constants.MetadataKey].(map[string]interface{})

	if !ok {
		return nil, errors.MetadataRequiredError{ID: entity["id"].(string)}
	}
	delete(entity, constants.MetadataKey)
	delete(entity, "id")

	entityType, ok := metadata[constants.MetadataTypeKey].(string)

	if !ok {
		return nil, MetadataTypeRequiredError{ID: id.(string)}
	}

	return ecs.CreateEntity(entityType, entity)
}

var r = &registry{
	dataSources: make(map[string]data_source.DataSource),
}

var logger = log.With().Str("plugin", "data_source").Str("service", "registry").Logger()
