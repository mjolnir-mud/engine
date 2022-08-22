package data_sources

import (
	"github.com/mjolnir-mud/engine/plugins/data_sources/internal/registry"
	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/data_source"
)

type plugin struct{}

func (p plugin) Name() string {
	return "data_sources"
}

func (p plugin) Start() error {
	return registry.Start()
}

func (p plugin) Stop() error {
	return registry.Stop()
}

// Register registers a data source with the registry. If a data source with the same name is already registered,
//i it will be overwritten.
func Register(source data_source.DataSource) {
	registry.Register(source)
}

// Find returns all entities in a data source that match the provided filter. If the data source does not exist, an
// error will be thrown.
func Find(source string, filter map[string]interface{}) (map[string]map[string]interface{}, error) {
	return registry.Find(source, filter)
}

// FindOne returns the first entity in a data source that matches the provided filter. If the data source does not
// exist, an error will be thrown.
func FindOne(source string, filter map[string]interface{}) (map[string]interface{}, error) {
	return registry.FindOne(source, filter)
}

// LoadAll loads all entities from a data source. It will call `ecs.Create` passing the map returned by the data source
// for each entity, and return a map of entities keyed by their ids.
func All(source string) (map[string]map[string]interface{}, error) {
	return registry.All(source)
}

// Count returns the number of entities in a data source using the provided map as a filter. If the data source does not
// exist, an error will be returned.
func Count(source string, filter map[string]interface{}) (int64, error) {
	return registry.Count(source, filter)
}

// Save saves data to a data source for a given entity. If the entity does not have a valid metadata field an error will
// be thrown. If the data source does not exist, an error will be thrown. If the metadata field does not have a type
// set, an error will be thrown. If the entity exists in the data source, it will be overwritten.
func Save(source string, entityId string, entity map[string]interface{}) error {
	return registry.Save(source, entityId, entity)
}

var Plugin = plugin{}
