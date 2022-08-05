package data_source

// DataSource is a data source that can be used to retrieve and persist data. So long as the data source implements this
// interface, it can be registered and used.
type DataSource interface {
	// Name returns the name of the data source. The name must be unique. Registering a data source with the same name
	// will replace the existing data source of the same name.
	Name() string

	// Start is called when the registry is started, and should be used to do any work to "start" the data source.
	Start() error

	// Stop is called when the registry is stopped, and should be used to do any work to "stop" the data source.
	Stop() error

	// Load loads an entity from the data source. The entity ID is the key used to load the entity. If the entity does
	// not exist, an error should be returned. The method should return an the id,  a map of key/value pairs
	// representing the component, as well as an additional __metadata key, whose value is a map of
	// key/value pairs representing the metadata. The metadata MUST include the __entityType key, whose value is a
	// valid entity type registered with the ECS plugin. The data source registry will then call `ecs.CreateWithId`
	// when the `data_source.Load()` method is called.
	Load(entityId string) (map[string]interface{}, error)

	// LoadAll loads all entities from the data source returning a map of entity ID to entities.
	LoadAll() (map[string]map[string]interface{}, error)

	// Find returns a list of entities from executing a search against a provided map. It returns a list of entities as a
	// map keyed by their ids.
	Find(search map[string]interface{}) (map[string]map[string]interface{}, error)

	// Save saves an entity to the data source. The entity ID is the key used to save the entity. The entity is a map of
	// key/value pairs representing the component.
	Save(entityId string, entity map[string]interface{}) error
}
