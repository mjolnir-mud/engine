package data_source

// DataSource is a data source that can be used to retrieve and persist data. So long as the data source implements this
// interface, it can be registered and used.
type DataSource interface {
	// All loads all entities from the data source returning a map of entity ID to entities.
	All() (map[string]map[string]interface{}, error)

	// AppendMetadata appends metadata to the entity, returning the entity with the metadata appended.
	AppendMetadata(entity map[string]interface{}) map[string]interface{}

	// Count returns the number of entities in the data source using the provided map.
	Count(map[string]interface{}) (int64, error)

	// Delete deletes an entity from the data source.
	Delete(entityId string) error

	// Find returns a list of entities from executing a search against a provided map. It returns a list of entities as a
	// map keyed by their ids.
	Find(search map[string]interface{}) (map[string]map[string]interface{}, error)

	// FindOne returns a single id, and entity from executing a search against a provided map.
	FindOne(search map[string]interface{}) (string, map[string]interface{}, error)

	// FindAndDelete deletes all entities from the data source that match the provided filter.
	FindAndDelete(search map[string]interface{}) error

	// Name returns the name of the data source. The name must be unique. Registering a data source with the same name
	// will replace the existing data source of the same name.
	Name() string

	// SaveWithId saves an entity to the data source. The entity ID is the key used to save the entity. The entity is a map of
	// key/value pairs representing the component.
	SaveWithId(entityId string, entity map[string]interface{}) error

	// Save saves an entity to the data source. The entity is a map of key/value pairs representing the component. The
	// returned entity ID is the key used to save the entity.
	Save(entity map[string]interface{}) (string, error)

	// Start is called when the registry is started, and should be used to do any work to "start" the data source.
	Start() error

	// Stop is called when the registry is stopped, and should be used to do any work to "stop" the data source.
	Stop() error
}
