package ecs

import "github.com/mjolnir-mud/engine/plugins/ecs/internal/entity_registry"

type plugin struct{}

func (p plugin) Name() string {
	return "ecs"
}

func (p plugin) Start() error {
	return nil
}

func (p plugin) Stop() error {
	return nil
}

// AddEntity adds an entity to the entity registry. it takes a map of components to be added. it
// will automatically generate a unique id for the entity. If the entity args does not contain the type component
// an error will be thrown.
func AddEntity(args map[string]interface{}) (string, error) {
	return entity_registry.Add(args)
}

// AddEntityWithID adds an entity with the provided id to the entity registry. It takes the entity id,
// and a map of components to be added. If an entity with the same id already exists, an error will be thrown.
func AddEntityWithID(id string, args map[string]interface{}) error {
	return entity_registry.AddWithID(id, args)
}

// AddBoolComponentToEntity adds a boolean component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func AddBoolComponentToEntity(id string, component string, value bool) error {
	return entity_registry.AddBoolComponentToEntity(id, component, value)
}

// AddBoolToMapComponent adds a boolean component to a map component. It takes the entity ID, component name, the key
// to which to add the value, and the value to add to the map. If an entity with the same id does not exist an error
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key already exists
// an error will be thrown. Once a value is added to the map, the type of that key is enforced. Attempting to change
//// the type of a key will result in an error in later updated.
func AddBoolToMapComponent(id string, component string, key string, value bool) error {
	return entity_registry.AddBoolToMapComponent(id, component, key, value)
}

// AddIntComponentToEntity adds an integer component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func AddIntComponentToEntity(id string, component string, value int) error {
	return entity_registry.AddIntComponentToEntity(id, component, value)
}

// AddIntToMapComponent adds an integer component to a map component. It takes the entity ID, component name, the key
// to which to add the value, and the value to add to the map. If an entity with the same id does not exist an error
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key already exists
// an error will be thrown. Once a value is added to the map, the type of that key is enforced. Attempting to change
//// the type of a key will result in an error in later updated.
func AddIntToMapComponent(id string, component string, key string, value int) error {
	return entity_registry.AddIntToMapComponent(id, component, key, value)
}

// AddInt64ComponentToEntity adds an integer64 component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func AddInt64ComponentToEntity(id string, component string, value int64) error {
	return entity_registry.AddInt64ComponentToEntity(id, component, value)
}

// AddInt64ToMapComponent adds an integer64 component to a map component. It takes the entity ID, component name, the
// key to which to add the value, and the value to add to the map. If an entity with the same id does not exist an error
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key already exists
// an error will be thrown. Once a value is added to the map, the type of that key is enforced. Attempting to change
// the type of a key will result in an error in later updated.
func AddInt64ToMapComponent(id string, component string, key string, value int64) error {
	return entity_registry.AddInt64ToMapComponent(id, component, key, value)
}

// AddMapComponentToEntity adds a map component to an entity. It takes the entity ID, component name, and the value of
// the component. If an entity with the same id does not exist an error will be thrown. If a component with the same name
// already exists, an error will be thrown. Once a value is added to the map, the type of that key is enforced.
// Attempting to change the type of a key will result in an error in later updated.
func AddMapComponentToEntity(id string, component string, value map[string]interface{}) error {
	return entity_registry.AddMapComponentToEntity(id, component, value)
}

// AddStringComponentToEntity adds a string component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func AddStringComponentToEntity(id string, component string, value string) error {
	return entity_registry.AddStringComponentToEntity(id, component, value)
}

// AddStringToMapComponent adds a string component to a map component. It takes the entity ID, component name, the
// key to which to add the value, and the value to add to the map. If an entity with the same id does not exist an error
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key already exists
// an error will be thrown. Once a value is added to the map, the type of that key is enforced. Attempting to change
//// the type of a key will result in an error in later updated.
func AddStringToMapComponent(id string, component string, key string, value string) error {
	return entity_registry.AddStringToMapComponent(id, component, key, value)
}

// CreateEntity will create an entity of the given entity type, without adding it to the entity registry. it takes the
// entity type and a map of components. It will merge the provided components with the default components for the
// entity type returning the merged components as a map.
func CreateEntity(entityType string, args map[string]interface{}) (map[string]interface{}, error) {
	return entity_registry.Create(entityType, args)
}

// CreateAndAddEntity creates an entity of the given entity type, adds it to the entity registry, and returns the
// id of the entity. It takes the entity type and a map of components. It will merge the provided components with the
// default components for the entity type returning the merged components as a map.
func CreateAndAddEntity(entityType string, args map[string]interface{}) (string, error) {
	return entity_registry.CreateAndAdd(entityType, args)
}

// RemoveEntity removes an entity from the entity registry. It takes the entity type and id. If an entity with the same
// id does not exist an error will be thrown.
func RemoveEntity(id string) error {
	return entity_registry.Remove(id)
}

// ReplaceEntity removes and then replaces an entity in the entity registry. It takes the entity type, id, and a map of
// components. It will remove the entity with the provided id and then add the provided components to the entity. If an
// entity with the same id does not exist, or the entity type does not match an error will be thrown.
func ReplaceEntity(id string, args map[string]interface{}) error {
	return entity_registry.Replace(id, args)
}

// UpdateEntity updates an entity in the entity registry. It takes the entity type, id, and a map of components. It will
// apply the provided components to the entity. If an entity with the same id does not exist, or the entity type does
// not match an error will be thrown. If the entity type does not match the existing an error will be thrown.
// Any components that are not provided will be removed from the entity.
func UpdateEntity(id string, args map[string]interface{}) error {
	return entity_registry.Update(id, args)
}

// UpdateBoolComponent updates a bool component in the entity registry. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name does not exist, an error will be thrown.
func UpdateBoolComponent(id string, component string, value bool) error {
	return entity_registry.UpdateBoolComponent(id, component, value)
}

// UpdateBoolInMapComponent updates a bool component in a map component. It takes the entity ID, component name, the
// key to which to add the value, and the value to add to the map. If an entity with the same id does not exist an error
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key soes not
// already exist an error will be thrown. Once a value is added to the map, the type of that key is enforced. If
// the value is not the correct type an error will be thrown.
func UpdateBoolInMapComponent(id string, component string, key string, value bool) error {
	return entity_registry.UpdateBoolInMapComponent(id, component, key, value)
}

// UpdateIntComponent updates an integer component in the entity registry. It takes the entity ID, component name, and
// the value of the component. If an entity with the same id does not exist an error will be thrown. If a component with
// the same name does not exist, an error will be thrown.
func UpdateIntComponent(id string, component string, value int) error {
	return entity_registry.UpdateIntComponent(id, component, value)
}

// UpdateIntInMapComponent updates an integer component in a map component. It takes the entity ID, component name, the
// key to which to add the value, and the value to add to the map. If an entity with the same id does not exist an error
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key soes not
// already exist an error will be thrown. Once a value is added to the map, the type of that key is enforced. If
// the value is not the correct type an error will be thrown.
func UpdateIntInMapComponent(id string, component string, key string, value int) error {
	return entity_registry.UpdateIntInMapComponent(id, component, key, value)
}

// UpdateInt64Component updates an integer64 component in the entity registry. It takes the entity ID, component name,
// and the value of the component. If an entity with the same id does not exist an error will be thrown. If a component
// with the same name does not exist, an error will be thrown.
func UpdateInt64Component(id string, component string, value int64) error {
	return entity_registry.UpdateInt64Component(id, component, value)
}

// UpdateInt64InMapComponent updates an integer64 component in a map component. It takes the entity ID, component name,
// the key to which to add the value, and the value to add to the map. If an entity with the same id does not exist an
// error will be thrown. If a component with the same name does not exist, an error will be thrown. If the key soes not
// already exist an error will be thrown. Once a value is added to the map, the type of that key is enforced. If
// the value is not the correct type an error will be thrown.
func UpdateInt64InMapComponent(id string, component string, key string, value int64) error {
	return entity_registry.UpdateInt64InMapComponent(id, component, key, value)
}

// UpdateOrAddEntity updates an entity in the entity registry. It takes the entity type, id, and a map of components.
// It will apply the provided components to the entity. If an entity with the same id does not exist, it will add the
// entity with the provided id and components. If the entity type does not match the existing an error will be thrown.
// Any components that are not provided will be removed from the entity.
func UpdateOrAddEntity(id string, args map[string]interface{}) error {
	return entity_registry.UpdateOrAdd(id, args)
}

// UpdateStringComponent updates a string component in the entity registry. It takes the entity ID, component name, and
// the value of the component. If an entity with the same id does not exist an error will be thrown. If a component with
// the same name does not exist, an error will be thrown.
func UpdateStringComponent(id string, component string, value string) error {
	return entity_registry.UpdateStringComponent(id, component, value)
}

// UpdateStringInMapComponent updates a string component in a map component. It takes the entity ID, component name,
// the key to which to add the value, and the value to add to the map. If an entity with the same id does not exist an
// error will be thrown. If a component with the same name does not exist, an error will be thrown. If the key soes not
// already exist an error will be thrown. Once a value is added to the map, the type of that key is enforced. If
// the value is not the correct type an error will be thrown.
func UpdateStringInMapComponent(id string, component string, key string, value string) error {
	return entity_registry.UpdateStringInMapComponent(id, component, key, value)
}
