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

package registry

import (
	"fmt"
	"github.com/mjolnir-mud/engine/plugins/data_sources/internal/logger"

	"github.com/rs/zerolog"

	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/constants"
	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/data_source"
	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/errors"
	"github.com/mjolnir-mud/engine/plugins/ecs"
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
	return fmt.Sprintf("entity %s does not return an entity with the type set in the metadata", e.ID)
}

var log zerolog.Logger
var dataSources map[string]data_source.DataSource

func All(source string) (map[string]map[string]interface{}, error) {
	if d, ok := dataSources[source]; ok {
		entities, err := d.All()

		if err != nil {
			return nil, err
		}

		return loadEntities(entities)
	} else {
		return nil, InvalidDataSourceError{Source: source}
	}
}

func Count(source string, search map[string]interface{}) (int64, error) {
	if d, ok := dataSources[source]; ok {
		return d.Count(search)
	} else {
		return 0, InvalidDataSourceError{Source: source}
	}
}

func CreateEntity(dataSource string, entityType string, data map[string]interface{}) (string, map[string]interface{}, error) {
	l := log.With().Str("data_source", dataSource).Str("entity_type", entityType).Logger()
	l.Debug().Msg("creating entity")
	if d, ok := dataSources[dataSource]; ok {
		l.Trace().Msg("data source found, creating entity")
		entity, err := ecs.CreateEntity(entityType, data)

		if err != nil {
			return "", nil, err
		}

		l.Trace().Msg("entity created, saving")

		metadata := map[string]interface{}{
			constants.MetadataTypeKey: entityType,
		}

		l.Trace().Msg("appending metadata from data source")
		metadata = d.AppendMetadata(metadata)

		entity[constants.MetadataKey] = metadata

		l.Trace().Msg("saving entity")
		id, err := d.Save(entity)

		if err != nil {
			return "", nil, err
		}

		return id, entity, nil
	} else {
		return "", nil, InvalidDataSourceError{Source: dataSource}
	}
}

func CreateEntityWithId(dataSource string, entityType string, entityId string, data map[string]interface{}) (map[string]interface{}, error) {
	l := log.With().Str("data_source", dataSource).Str("entity_type", entityType).Logger()
	l.Debug().Msg("creating entity with id")
	if d, ok := dataSources[dataSource]; ok {
		l.Trace().Msg("data source found, creating entity")
		entity, err := ecs.CreateEntity(entityType, data)

		if err != nil {
			return nil, err
		}

		l.Trace().Msg("entity created, saving")

		metadata := map[string]interface{}{
			constants.MetadataTypeKey: entityType,
		}

		l.Trace().Msg("appending metadata from data source")
		metadata = d.AppendMetadata(metadata)

		entity[constants.MetadataKey] = metadata

		l.Trace().Msg("saving entity")
		err = d.SaveWithId(entityId, entity)

		if err != nil {
			return nil, err
		}

		return entity, nil
	} else {
		return nil, InvalidDataSourceError{Source: dataSource}
	}
}

func Find(source string, search map[string]interface{}) (map[string]map[string]interface{}, error) {
	if d, ok := dataSources[source]; ok {
		entities, err := d.Find(search)

		if err != nil {
			return nil, err
		}

		return loadEntities(entities)
	} else {
		return nil, InvalidDataSourceError{Source: source}
	}
}

func FindOne(source string, search map[string]interface{}) (string, map[string]interface{}, error) {
	if d, ok := dataSources[source]; ok {
		id, entity, err := d.FindOne(search)

		if err != nil {
			return "", nil, err
		}

		e, err := loadEntity(entity)

		if err != nil {
			return "", nil, err
		}

		return id, e, nil
	} else {
		return "", nil, InvalidDataSourceError{Source: source}
	}
}

func Delete(source string, entityId string) error {
	if d, ok := dataSources[source]; ok {
		return d.Delete(entityId)
	} else {
		return InvalidDataSourceError{Source: source}
	}
}

func FindAndDelete(source string, search map[string]interface{}) error {
	log.Debug().Str("source", source).Msg("finding and deleting")
	if d, ok := dataSources[source]; ok {
		return d.FindAndDelete(search)
	} else {
		return InvalidDataSourceError{Source: source}
	}
}

func Register(dataSource data_source.DataSource) {
	log.Info().Msgf("registering data source %s", dataSource.Name())
	dataSources[dataSource.Name()] = dataSource
}

func SaveWithId(source string, entityId string, entity map[string]interface{}) error {
	log.Trace().Str("source", source).Str("entityId", entityId).Msg("SaveWithId")

	if d, ok := dataSources[source]; ok {
		log.Trace().Str("source", source).Str("entityId", entityId).Msg("data source found")

		log.Trace().Str("source", source).Str("entityId", entityId).Msg("checking for metadata")
		metadata, ok := entity[constants.MetadataKey].(map[string]interface{})

		if !ok {
			return errors.MetadataRequiredError{ID: entityId}
		}

		log.Trace().Str("source", source).Str("entityId", entityId).Msg("checking for type")
		_, ok = metadata[constants.MetadataTypeKey].(string)

		if !ok {
			return MetadataTypeRequiredError{ID: entityId}
		}

		log.Trace().Str("source", source).Str("entityId", entityId).Msg("saving entity")
		return d.SaveWithId(entityId, entity)
	} else {
		return InvalidDataSourceError{Source: source}
	}
}

func Save(source string, entity map[string]interface{}) (string, error) {
	log.Trace().Str("source", source).Msg("Save")

	if d, ok := dataSources[source]; ok {
		log.Trace().Str("source", source).Msg("data source found")

		log.Trace().Str("source", source).Msg("checking for metadata")
		metadata, ok := entity[constants.MetadataKey].(map[string]interface{})

		if !ok {
			return "", errors.MetadataRequiredError{}
		}

		log.Trace().Str("source", source).Msg("checking for type")
		_, ok = metadata[constants.MetadataTypeKey].(string)

		if !ok {
			return "", MetadataTypeRequiredError{}
		}

		log.Trace().Str("source", source).Msg("saving entity")
		return d.Save(entity)
	} else {
		return "", InvalidDataSourceError{Source: source}
	}
}

func Start() {
	log = logger.Instance.With().Str("component", "registry").Logger()
	log.Info().Msg("starting")

	dataSources = make(map[string]data_source.DataSource)
}

func StartDataSources() {
	log.Info().Msgf("%d", len(dataSources))
	for _, d := range dataSources {
		log.Info().Msgf("starting data source %s", d.Name())
		err := d.Start()
		if err != nil {
			log.Error().Err(err).Msg("error starting data source")
			panic(err)
		}
	}
}

func Stop() {
	log.Info().Msg("stopping")
	for _, d := range dataSources {
		log.Info().Msgf("stopping data source %s", d.Name())
		_ = d.Stop()
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
	delete(entity, "id")

	entityType, ok := metadata[constants.MetadataTypeKey].(string)

	if !ok {
		return nil, MetadataTypeRequiredError{ID: id.(string)}
	}

	return ecs.CreateEntity(entityType, entity)
}
