package ecs

import (
	"github.com/mjolnir-mud/engine/plugins/ecs/internal/entity_registry"
	"github.com/mjolnir-mud/engine/plugins/ecs/internal/system_registry"
	"github.com/mjolnir-mud/engine/plugins/ecs/pkg/entity_type"
	"github.com/mjolnir-mud/engine/plugins/ecs/pkg/system"
)

type plugin struct{}

func (p plugin) Name() string {
	return "ecs"
}

func (p plugin) Start() error {
	entity_registry.Start()
	system_registry.Start()
	return nil
}

func (p plugin) Stop() error {
	entity_registry.Stop()
	system_registry.Stop()
	return nil
}

// AddEntity adds an entity to the entity registry. It takes a type, a map of components to be added to the entity. If the
// entity already exists, an error will be returned. If the type is not registered, an error will be returned.
func AddEntity(entityType string, args map[string]interface{}) (string, error) {
	return entity_registry.Add(entityType, args)
}

// AddEntityWithID adds an entity with the provided id to the entity registry. It takes the entity id,
// and a map of components to be added. If an entity with the same id already exists, an error will be thrown. If the
// type is not registered, an error will be thrown.
func AddEntityWithID(entityType string, id string, args map[string]interface{}) error {
	return entity_registry.AddWithID(entityType, id, args)
}

// AddBoolComponent adds a boolean component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func AddBoolComponent(id string, component string, value bool) error {
	return entity_registry.AddBoolComponent(id, component, value)
}

// AddBoolToMapComponent adds a boolean component to a map component. It takes the entity ID, component name, the key
// to which to add the value, and the value to add to the map. If an entity with the same id does not exist an error
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key already exists
// an error will be thrown. Once a value is added to the map, the type of that key is enforced. Attempting to change
//// the type of a key will result in an error in later updated.
func AddBoolToMapComponent(id string, component string, key string, value bool) error {
	return entity_registry.AddBoolToMapComponent(id, component, key, value)
}

// AddIntComponent adds an integer component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func AddIntComponent(id string, component string, value int) error {
	return entity_registry.AddIntComponent(id, component, value)
}

// AddIntToMapComponent adds an integer component to a map component. It takes the entity ID, component name, the key
// to which to add the value, and the value to add to the map. If an entity with the same id does not exist an error
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key already exists
// an error will be thrown. Once a value is added to the map, the type of that key is enforced. Attempting to change
// the type of a key will result in an error in later updated.
func AddIntToMapComponent(id string, component string, key string, value int) error {
	return entity_registry.AddIntToMapComponent(id, component, key, value)
}

// AddInt64ComponentToEntity adds an integer64 component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func AddInt64ComponentToEntity(id string, component string, value int64) error {
	return entity_registry.AddInt64Component(id, component, value)
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
	return entity_registry.AddMapComponent(id, component, value)
}

// AddSetComponent adds a set component to an entity. It takes the entity ID, component name, and the value of
// the component. If an entity with the same id already exists error will be thrown. If a component with the same name
// already exists, an error will be thrown. Once a value is added to the set, the type of that value is enforced for
// all members of the set. Attempting to change the type of a value will result in an error in later updates.
func AddSetComponent(id string, component string, value []interface{}) error {
	return entity_registry.AddSetComponent(id, component, value)
}

// AddStringComponentToEntity adds a string component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func AddStringComponentToEntity(id string, component string, value string) error {
	return entity_registry.AddStringComponent(id, component, value)
}

// AddStringToMapComponent adds a string component to a map component. It takes the entity ID, component name, the
// key to which to add the value, and the value to add to the map. If an entity with the same id does not exist an error
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key already exists
// an error will be thrown. Once a value is added to the map, the type of that key is enforced. Attempting to change
//// the type of a key will result in an error in later updated.
func AddStringToMapComponent(id string, component string, key string, value string) error {
	return entity_registry.AddStringToMapComponent(id, component, key, value)
}

// AddToStringSetComponent adds a string value to a set component. It takes the entity ID, component name, and the
// value to add to the set. If an entity with the same id does not exist an error will be thrown. If a component with
// the same name does not exist, an error will be thrown. If the value type is not a string, an error will be thrown.
func AddToStringSetComponent(id string, component string, value string) error {
	return entity_registry.AddToStringSetComponent(id, component, value)
}

// AddToIntSetComponent adds an integer value to a set component. It takes the entity ID, component name, and the
// value to add to the set. If an entity with the same id does not exist an error will be thrown. If a component with
// the same name does not exist, an error will be thrown. If the value type is not an integer, an error will be thrown.
func AddToIntSetComponent(id string, component string, value int) error {
	return entity_registry.AddToIntSetComponent(id, component, value)
}

// AddToInt64SetComponent adds an integer64 value to a set component. It takes the entity ID, component name, and the
// value to add to the set. If an entity with the same id does not exist an error will be thrown. If a component with
// the same name does not exist, an error will be thrown. If the value type is not an integer64, an error will be thrown.
func AddToInt64SetComponent(id string, component string, value int64) error {
	return entity_registry.AddToInt64SetComponent(id, component, value)
}

// ComponentExists checks if a component exists. It takes the entity ID and the component name.
func ComponentExists(id string, component string) bool {
	return entity_registry.ComponentExists(id, component)
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

// EntityExists checks if an entity with the given id exists. It takes the entity id and returns a boolean.
func EntityExists(id string) (bool, error) {
	return entity_registry.Exists(id)
}

// GetBoolComponent returns a boolean component from an entity. It takes the entity ID and component name. If the
// entity does not exist or the component does not exist, an error will be thrown. If the component is not a boolean,
// an error will be thrown.
func GetBoolComponent(id string, component string) (bool, error) {
	return entity_registry.GetBoolComponent(id, component)
}

// GetInt64Component returns an integer64 component from an entity. It takes the entity ID and component name. If the
// entity does not exist or the component does not exist, an error will be thrown. If the component is not an integer64,
// an error will be thrown.
func GetInt64Component(id string, component string) (int64, error) {
	return entity_registry.GetInt64Component(id, component)
}

// GetIntComponent returns an integer component from an entity. It takes the entity ID and component name. If the
// entity does not exist or the component does not exist, an error will be thrown. If the component is not an integer,
// an error will be thrown.
func GetIntComponent(id string, component string) (int, error) {
	return entity_registry.GetIntComponent(id, component)
}

// GetIntFromMapComponent returns the int value of an element in a map component. It takes the entity ID, component
// name, and the element name. If the entity does not exist or the component does not exist, an error will be thrown.
// If the component is not a map, an error will be thrown. If the element does not exist, an error will be thrown.
// If the element is not an integer, an error will be thrown.
func GetIntFromMapComponent(id string, component string, element string) (int, error) {
	return entity_registry.GetIntFromMapComponent(id, component, element)
}

// GetInt64FromMapComponent returns the int64 value of an element in a map component. It takes the entity ID,
// component name, and the element name. If the entity does not exist or the component does not exist, an error will
// be thrown. If the component is not a map, an error will be thrown. If the element does not exist, an error will
// be thrown. If the element is not an integer64, an error will be thrown.
func GetInt64FromMapComponent(id string, component string, element string) (int64, error) {
	return entity_registry.GetInt64FromMapComponent(id, component, element)
}

// GetHashComponent returns a hash component from an entity. It takes the entity ID and component name. If the
// entity does not exist or the component does not exist, an error will be thrown. If the component is not a hash,
// an error will be thrown.
func GetHashComponent(id string, component string) (map[string]interface{}, error) {
	return entity_registry.GetMapComponent(id, component)
}

// GetStringComponent returns the value of the string component. It takes the entity ID and component name. If the
// entity does not exist or the component does not exist, an error will be thrown. If the component is not a string,
// an error will be thrown.
func GetStringComponent(id string, component string) (string, error) {
	return entity_registry.GetStringComponent(id, component)
}

// GetStringFromMapComponent returns the string value of an element in a map component. It takes the entity ID,
// component name, and the element name. If the entity does not exist or the component does not exist, an error will
// be thrown. If the component is not a map, an error will be thrown. If the element does not exist, an error will
// be thrown. If the element is not a string, an error will be thrown.
func GetStringFromMapComponent(id string, component string, element string) (string, error) {
	return entity_registry.GetStringFromMapComponent(id, component, element)
}

// IsEntityTypeRegistered checks if an entity type is registered. It takes the entity type name.
func IsEntityTypeRegistered(entityType string) bool {
	return entity_registry.IsEntityTypeRegistered(entityType)
}

// RegisterSystem registers a system with the registry. If a system with the same name is already registered, it will be
// overwritten.
func RegisterSystem(system system.System) {
	system_registry.Register(system)
}

// RegisterEntityType registers an entity type. Entity Types must implmeent the `EntityType` interface. It is
// expected that developers can override default EntityType implementations with their own implementations.Q
func RegisterEntityType(entityType entity_type.EntityType) {
	entity_registry.Register(entityType)
}

// RemoveComponent removes the component from the entity. If an entity with the same id does not exist an error will be
// thrown. If a component with the same name does not exist, an error will be thrown.
func RemoveComponent(id string, name string) error {
	return entity_registry.RemoveComponent(id, name)
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

// RemoveFromStringSetComponent removes a string value from a set component. It takes the entity ID, component name, and
// the value to remove from the set. If an entity with the same id does not exist an error will be thrown. If a component
// with the same name does not exist, an error will be thrown. If the value type is not a string, an error will be thrown.
func RemoveFromStringSetComponent(id string, component string, value string) error {
	return entity_registry.RemoveFromStringSetComponent(id, component, value)
}

// RemoveFromIntSetComponent removes an integer value from a set component. It takes the entity ID, component name, and
// the value to remove from the set. If an entity with the same id does not exist an error will be thrown. If a component
// with the same name does not exist, an error will be thrown. If the value type is not an integer, an error will be thrown.
func RemoveFromIntSetComponent(id string, component string, value int) error {
	return entity_registry.RemoveFromIntSetComponent(id, component, value)
}

// RemoveFromInt64SetComponent removes an integer64 value from a set component. It takes the entity ID, component name,
// and the value to remove from the set. If an entity with the same id does not exist an error will be thrown. If a
// component with the same name does not exist, an error will be thrown. If the value type is not an integer64, an error
// will be thrown.
func RemoveFromInt64SetComponent(id string, component string, value int64) error {
	return entity_registry.RemoveFromInt64SetComponent(id, component, value)
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
func UpdateOrAddEntity(entityType string, id string, args map[string]interface{}) error {
	return entity_registry.UpdateOrAdd(entityType, id, args)
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

var Plugin = plugin{}
