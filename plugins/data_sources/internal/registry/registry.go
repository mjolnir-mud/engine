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
	"github.com/mjolnir-mud/engine/plugins/data_sources/constants"
	"github.com/mjolnir-mud/engine/plugins/data_sources/data_source"
	dataSourceErrors "github.com/mjolnir-mud/engine/plugins/data_sources/errors"
	"github.com/mjolnir-mud/engine/plugins/data_sources/internal/logger"

	"github.com/rs/zerolog"

	"github.com/mjolnir-mud/engine/plugins/ecs"
)

var log zerolog.Logger
var dataSources map[string]data_source.Interface

func All(source string) (*data_source.FindResults, error) {
	d, err := getDataSource(source)

	if err != nil {
		return nil, err
	}

	entities, err := d.All()

	if err != nil {
		return nil, err
	}

	return loadEntities(entities)
}

func Count(source string, search map[string]interface{}) (int64, error) {
	d, err := getDataSource(source)

	if err != nil {
		return 0, err
	}

	return d.Count(search)
}

func CreateEntity(dataSource string, entityType string, data map[string]interface{}) (string, map[string]interface{}, error) {
	l := log.With().Str("data_source", dataSource).Str("entity_type", entityType).Logger()
	l.Debug().Msg("creating entity")

	d, err := getDataSource(dataSource)

	if err != nil {
		return "", nil, err
	}

	l.Trace().Msg("data source found, creating entity")
	entity, err := ecs.NewEntity(entityType, data)

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
}

func NewEntityWithId(dataSource string, entityType string, entityId string, data map[string]interface{}) (map[string]interface{}, error) {
	l := log.With().Str("data_source", dataSource).Str("entity_type", entityType).Logger()
	l.Debug().Msg("creating entity with id")
	d, err := getDataSource(dataSource)

	if err != nil {
		return nil, err
	}

	l.Trace().Msg("data source found, creating entity")
	entity, err := ecs.NewEntity(entityType, data)

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
}

func Find(source string, search map[string]interface{}) (*data_source.FindResults, error) {
	d, err := getDataSource(source)

	if err != nil {
		return nil, err
	}
	entities, err := d.Find(search)

	if err != nil {
		return nil, err
	}

	return loadEntities(entities)
}

func FindAndAddToECS(source string, search map[string]interface{}) ([]string, error) {
	//entites, err := Find(source, search)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//ids := make([]string, 0)
	//
	//for id, entity := range entites {
	//	ids = append(ids, id)
	//	err := ecs.AddEntityWithID(id, entity)
	//
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	return nil, nil
}

func FindOne(source string, search map[string]interface{}) (*data_source.FindResult, error) {
	d, err := getDataSource(source)

	if err != nil {
		return nil, err
	}

	entity, err := d.FindOne(search)

	if err != nil {
		return nil, err
	}

	e, err := newFindResult(entity)

	if err != nil {
		return nil, err
	}

	return e, nil
}

func Delete(source string, entityId string) error {
	d, err := getDataSource(source)

	if err != nil {
		return err
	}

	return d.Delete(entityId)
}

func FindAndDelete(source string, search map[string]interface{}) error {
	log.Debug().Str("source", source).Msg("finding and deleting")
	d, err := getDataSource(source)

	if err != nil {
		return err
	}

	return d.FindAndDelete(search)
}

func Register(dataSource data_source.Interface) {
	log.Info().Msgf("registering data source %s", dataSource.Name())
	dataSources[dataSource.Name()] = dataSource
}

func SaveWithId(source string, entityId string, entity map[string]interface{}) error {
	log.Trace().Str("source", source).Str("entityId", entityId).Msg("SaveWithId")
	d, err := getDataSource(source)

	if err != nil {
		return err
	}

	log.Trace().Str("source", source).Str("entityId", entityId).Msg("data source found")

	log.Trace().Str("source", source).Str("entityId", entityId).Msg("checking for metadata")
	metadata, ok := entity[constants.MetadataKey].(map[string]interface{})

	if !ok {
		return dataSourceErrors.MetadataRequiredError{ID: entityId}
	}

	log.Trace().Str("source", source).Str("entityId", entityId).Msg("checking for type")
	_, ok = metadata[constants.MetadataTypeKey].(string)

	if !ok {
		return dataSourceErrors.MetadataRequiredError{ID: entityId}
	}

	log.Trace().Str("source", source).Str("entityId", entityId).Msg("saving entity")
	return d.SaveWithId(entityId, entity)
}

func Save(source string, entity map[string]interface{}) (string, error) {
	log.Trace().Str("source", source).Msg("Save")
	d, err := getDataSource(source)

	if err != nil {
		return "", err
	}

	log.Trace().Str("source", source).Msg("data source found")

	log.Trace().Str("source", source).Msg("checking for metadata")
	metadata, ok := entity[constants.MetadataKey].(map[string]interface{})

	if !ok {
		return "", dataSourceErrors.MetadataRequiredError{}
	}

	log.Trace().Str("source", source).Msg("checking for type")
	_, ok = metadata[constants.MetadataTypeKey].(string)

	if !ok {
		return "", dataSourceErrors.MetadataRequiredError{}
	}

	log.Trace().Str("source", source).Msg("saving entity")
	return d.Save(entity)
}

func Start() {
	log = logger.Instance.With().Str("component", "registry").Logger()
	log.Info().Msg("starting")

	dataSources = make(map[string]data_source.Interface)
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

func loadEntities(entities []map[string]interface{}) (*data_source.FindResults, error) {
	results := make([]*data_source.FindResult, 0)
	for _, entity := range entities {
		e, err := newFindResult(entity)

		if err != nil {
			return nil, err
		}

		results = append(results, e)
	}

	return data_source.NewFindResults(results), nil

}

func newFindResult(entity map[string]interface{}) (*data_source.FindResult, error) {
	id, ok := entity["id"].(string)

	if !ok {
		return nil, dataSourceErrors.IDRequiredError{}
	}

	metadata, ok := entity[constants.MetadataKey].(map[string]interface{})

	if !ok {
		return nil, dataSourceErrors.MetadataRequiredError{ID: entity["id"].(string)}
	}

	delete(entity, constants.MetadataKey)
	delete(entity, "id")

	entityType, ok := metadata[constants.MetadataTypeKey].(string)

	if !ok {
		return nil, dataSourceErrors.MetadataRequiredError{ID: id}
	}

	newEntity, err := ecs.NewEntity(entityType, entity)

	if err != nil {
		return nil, err
	}

	return data_source.NewFindResult(entityType, id, newEntity, metadata), nil
}

func getDataSource(source string) (data_source.Interface, error) {
	if d, ok := dataSources[source]; ok {
		return d, nil
	} else {
		return nil, dataSourceErrors.InvalidDataSourceError{Source: source}
	}
}
