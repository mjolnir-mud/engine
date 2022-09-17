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

package entity_registry

import (
	"fmt"
	"github.com/mjolnir-mud/engine/plugins/ecs/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/ecs/pkg/errors"
	"github.com/rs/zerolog"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/ecs/internal/constants"
	"github.com/mjolnir-mud/engine/plugins/ecs/pkg/entity_type"
)

type entityRegistry struct {
	types               map[string]entity_type.EntityType
	loadableDirectories []string
}

var registry = entityRegistry{
	types:               make(map[string]entity_type.EntityType),
	loadableDirectories: make([]string, 0),
}

var log zerolog.Logger

// Start starts the entity registry.
func Start() {
	log = logger.Instance.With().Str("component", "entity_registry").Logger()
	log.Info().Msg("starting")
}

// Stop stops the entity registry.
func Stop() {
	log.Info().Msg("stopping entity registry")
}

// AllComponents returns all the components for the entity.
func AllComponents(id string) map[string]interface{} {
	keys, err := engine.RedisKeys(fmt.Sprintf("%s:*", id)).Result()

	if err != nil {
		panic(err)
	}

	components := make(map[string]interface{})

	for _, key := range keys {
		name := strings.Replace(key, fmt.Sprintf("%s:", id), "", 1)
		t := engine.RedisGet(fmt.Sprintf("%s:%s", constants.ComponentTypePrefix, key)).Val()

		switch t {
		case "map":
			m, err := engine.RedisHGetAll(key).Result()
			if err != nil {
				panic(err)
			}

			var final = make(map[string]interface{})

			for k, v := range m {
				final[k] = v
			}

			components[name] = final
		case "set":
			s := engine.RedisSMembers(key).Val()
			set := make([]interface{}, 0)

			for _, v := range s {
				set = append(set, v)
			}

			components[name] = set

		default:
			components[name] = engine.RedisGet(key).Val()
		}

	}

	return components
}

// Add adds an entity to the entity registry. It takes the entity id, and a map of arguments to be passed to the entity
// type's constructor. If the entity type is not registered, an error will be returned. If the entity already exists,
// an error will be returned.
func Add(entityType string, args map[string]interface{}) (string, error) {
	id := generateID()
	err := AddWithId(entityType, id, args)

	if err != nil {
		return "", err
	}

	return id, nil
}

// AddWithId adds an entity with the provided id to the entity registry. It takes the entity id, and a map of arguments
// to be passed to the entity type's constructor. If the entity type is not registered, an error will be returned. If
// the entity already exists, an error will be returned.
func AddWithId(entityType string, id string, args map[string]interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if exists {
		return errors.EntityExistsError{ID: id}
	}

	if !IsEntityTypeRegistered(entityType) {
		return errors.EntityTypeNotRegisteredError{Type: entityType}
	}

	log.Debug().Str("id", id).Msg("adding entity")
	setEntityType(id, entityType)

	args, err = Create(entityType, args)

	if err != nil {
		return err
	}

	// add the entities components to the world
	err = setComponentsFromMap(id, args)

	if err != nil {
		return err
	}

	return nil
}

// AddBoolComponent adds a boolean component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func AddBoolComponent(id string, name string, value bool) error {
	if value {
		return addComponent(id, name, true)
	} else {
		return addComponent(id, name, false)
	}
}

// AddBoolToMapComponent adds a boolean component to a map component. It takes the entity ID, component name, the key
// to which to add the value, and the value to add to the map. If an entity with the same id does not exist an error
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key already exists
// an error will be thrown. Once a value is added to the map, the type of that key is enforced. Attempting to change
// the type of key will result in an error in later updated.
func AddBoolToMapComponent(id string, name string, key string, value bool) error {
	if value {
		return addToMapComponent(id, name, key, 1)
	} else {
		return addToMapComponent(id, name, key, 0)
	}
}

// AddIntComponent adds an integer component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func AddIntComponent(id string, name string, value int) error {
	return addComponent(id, name, value)
}

// AddIntToMapComponent adds an integer component to a map component. It takes the entity ID, component name, the key
// to which to add the value, and the value to add to the map. If an entity with the same id does not exist an error
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key already exists
// an error will be thrown. Once a value is added to the map, the type of that key is enforced. Attempting to change
// the type of key will result in an error in later updated.
func AddIntToMapComponent(id string, name string, key string, value int) error {
	return addToMapComponent(id, name, key, value)
}

// AddInt64Component adds an integer64 component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func AddInt64Component(id string, name string, value int64) error {
	return addComponent(id, name, value)
}

// AddInt64ToMapComponent adds an integer64 component to a map component. It takes the entity ID, component name, the
// key to which to add the value, and the value to add to the map. If an entity with the same id does not exist an error
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key already exists
// an error will be thrown. Once a value is added to the map, the type of that key is enforced. Attempting to change
//// the type of key will result in an error in later updated.
func AddInt64ToMapComponent(id string, name string, key string, value int64) error {
	return addToMapComponent(id, name, key, value)
}

// AddMapComponent adds a map component to an entity. It takes the entity ID, component name, and the value of
// the component. If an entity with the same id already exists an error will be thrown. If a component with the same name
// already exists, an error will be thrown.
func AddMapComponent(id string, name string, value map[string]interface{}) error {
	return addComponent(id, name, value)
}

// AddOrUpdateStringInMapComponent adds or updates a string component to a map component. It takes the entity ID,
// component name, the key to which to add the value, and the value to add to the map. If an entity with the same id
// does not exist an error will be thrown. If a component with the same name does not exist, an error will be thrown.
// If the key already exists, the value will be updated. Once a value is added to the map, the type of that key is
// enforced. Attempting to change the type of key will result in an error in later updated.
func AddOrUpdateStringInMapComponent(id string, name string, key string, value string) error {
	return addOrUpdateToMapComponent(id, name, key, value)
}

// AddOrUpdateIntInMapComponent adds or updates an integer component to a map component. It takes the entity ID,
// component name, the key to which to add the value, and the value to add to the map. If an entity with the same id
// does not exist an error will be thrown. If a component with the same name does not exist, an error will be thrown.
// If the key already exists, the value will be updated. Once a value is added to the map, the type of that key is
// enforced. Attempting to change the type of key will result in an error in later updated.
func AddOrUpdateIntInMapComponent(id string, name string, key string, value int) error {
	return addOrUpdateToMapComponent(id, name, key, value)
}

// AddSetComponent adds a set component to an entity. It takes the entity ID, component name, and the value of
// the component. If an entity with the same id already exists error will be thrown. If a component with the same name
// already exists, an error will be thrown. Once a value is added to the set, the type of that value is enforced for
// all members of the set. Attempting to change the type of value will result in an error in later updates. One cannot
// add empty sets, so if the value is empty, an error will be thrown.
func AddSetComponent(id string, name string, value []interface{}) error {
	return addComponent(id, name, value)
}

// AddStringComponent adds a string component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func AddStringComponent(id string, name string, value string) error {
	return addComponent(id, name, value)
}

// AddOrUpdateStringComponent adds or updates a string component to an entity. It takes the entity ID, component name,
// and the value of the component. If an entity with the same id does not exist an error will be thrown. If a component
// with the same name already exists, it will be updated.
func AddOrUpdateStringComponent(id string, name string, value string) error {
	_, err := getCurrentValue(id, name)

	if err != nil {
		return addComponent(id, name, value)
	} else {
		return updateComponent(id, name, value)
	}
}

// AddStringToMapComponent adds a string component to a map component. It takes the entity ID, component name, the
// key to which to add the value, and the value to add to the map. If an entity with the same id does not exist an error
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key already exists
// an error will be thrown. Once a value is added to the map, the type of that key is enforced. Attempting to change
// the type of key will result in an error in later updated.
func AddStringToMapComponent(id string, name string, key string, value string) error {
	return addToMapComponent(id, name, key, value)
}

// AddToStringSetComponent adds a string value to a set component. It takes the entity ID, component name, and the
// value to add to the set. If an entity with the same id does not exist an error will be thrown. If a component with
// the same name does not exist, an error will be thrown. If the value type is not a string, an error will be thrown.
func AddToStringSetComponent(id string, name string, value string) error {
	return addToSetComponent(id, name, value)
}

// AddToIntSetComponent adds an integer value to a set component. It takes the entity ID, component name, and the
// value to add to the set. If an entity with the same id does not exist an error will be thrown. If a component with
// the same name does not exist, an error will be thrown. If the value type is not an integer, an error will be thrown.
func AddToIntSetComponent(id string, name string, value int) error {
	return addToSetComponent(id, name, value)
}

// AddToInt64SetComponent adds an integer64 value to a set component. It takes the entity ID, component name, and the
// value to add to the set. If an entity with the same id does not exist an error will be thrown. If a component with
// the same name does not exist, an error will be thrown. If the value type is not an integer64, an error will be thrown.
func AddToInt64SetComponent(id string, name string, value int64) error {
	return addToSetComponent(id, name, value)
}

// ComponentExists checks if a component exists. It takes the entity ID and the component name.
func ComponentExists(id string, name string) (bool, error) {
	i, err := engine.RedisExists(componentTypeID(id, name)).Result()

	if err != nil {
		return false, err
	}

	if i == 1 {
		return true, nil
	}

	return false, nil
}

// Create will create an entity of the given entity type, without adding it to the entity registry. it takes the
// entity type and a map of components. It will pass the component arguments to the entity type's `Create` method,
// and return the result.
func Create(entityType string, components map[string]interface{}) (map[string]interface{}, error) {
	t := getEntityTypeByName(entityType)

	if t == nil {
		return nil, fmt.Errorf("entity type %s not found", entityType)
	}

	return t.Create(components), nil
}

// CreateAndAdd creates an entity of the given entity type, adds it to the entity registry, and returns the
// id of the entity. It takes the entity type and a map of components. It will merge the provided components with the
// default components for the entity type returning the merged components as a map.
func CreateAndAdd(entityType string, components map[string]interface{}) (string, error) {
	m, err := Create(entityType, components)

	if err != nil {
		return "", err
	}

	id, err := Add(entityType, m)

	if err != nil {
		return "", err
	}

	return id, nil
}

// ElementInSetComponentExists checks if a value is in a set component. It takes the entity ID, component name, and
// the value to check for. If the component does not exist an error is thrown.
func ElementInSetComponentExists(id string, name string, element interface{}) (bool, error) {
	exists, err := Exists(id)

	if err != nil {
		return false, err
	}

	if !exists {
		return false, errors.EntityNotFoundError{ID: id}
	}

	isMember, err := engine.RedisSIsMember(componentId(id, name), element).Result()

	if err != nil {
		return false, err
	}

	return isMember, nil

}

// Exists checks if an entity with the given id exists. It takes the entity id and returns a boolean.
func Exists(id string) (bool, error) {
	_, err := getEntityType(id)

	if err != nil {
		if err.Error() == "redis: nil" {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

// GetBoolComponent returns a boolean component from an entity. It takes the entity ID and component name. If the
// entity does not exist or the component does not exist, an error will be thrown. If the component is not a boolean,
// an error will be thrown.
func GetBoolComponent(id string, name string) (bool, error) {
	exists, err := Exists(id)

	if err != nil {
		return false, err
	}

	if !exists {
		return false, errors.EntityNotFoundError{ID: id}
	}

	componentExists, err := ComponentExists(id, name)

	if err != nil {
		return false, err
	}

	if !componentExists {
		return false, errors.ComponentNotFoundError{ID: id, Name: name}
	}

	return engine.RedisGet(componentId(id, name)).Bool()
}

// GetInt64Component returns an integer64 component from an entity. It takes the entity ID and component name. If the
// entity does not exist or the component does not exist, an error will be thrown. If the component is not an integer64,
// an error will be thrown.
func GetInt64Component(id string, name string) (int64, error) {
	exists, err := Exists(id)

	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, errors.EntityNotFoundError{ID: id}
	}

	componentExists, err := ComponentExists(id, name)

	if err != nil {
		return 0, err
	}

	if !componentExists {
		return 0, errors.ComponentNotFoundError{ID: id, Name: name}
	}

	return engine.RedisGet(componentId(id, name)).Int64()
}

// GetIntComponent returns an integer component from an entity. It takes the entity ID and component name. If the
// entity does not exist or the component does not exist, an error will be thrown. If the component is not an integer,
// an error will be thrown.
func GetIntComponent(id string, name string) (int, error) {
	exists, err := Exists(id)

	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, errors.EntityNotFoundError{ID: id}
	}

	componentExists, err := ComponentExists(id, name)

	if err != nil {
		return 0, err
	}

	if !componentExists {
		return 0, errors.ComponentNotFoundError{ID: id, Name: name}
	}

	return engine.RedisGet(componentId(id, name)).Int()
}

// GetIntFromMapComponent returns the int value of an element in a map component. It takes the entity ID, component
// name, and the element name. If the entity does not exist or the component does not exist, an error will be thrown.
// If the component is not a map, an error will be thrown. If the element does not exist, an error will be thrown.
// If the element is not an integer, an error will be thrown.
func GetIntFromMapComponent(id string, name string, mapKey string) (int, error) {
	v, err := getValueFromMapComponent(id, name, mapKey)

	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(v.(string), 10, 64)

	if err != nil {
		return 0, errors.MapValueTypeMismatchError{
			ID:       id,
			Name:     name,
			Key:      mapKey,
			Expected: "int",
			Actual:   "unable to parse int",
		}
	}

	return int(i), nil
}

// GetInt64FromMapComponent returns the int64 value of an element in a map component. It takes the entity ID,
// component name, and the element name. If the entity does not exist or the component does not exist, an error will
// be thrown. If the component is not a map, an error will be thrown. If the element does not exist, an error will
// be thrown. If the element is not an integer64, an error will be thrown.
func GetInt64FromMapComponent(id string, name string, mapKey string) (int64, error) {
	v, err := getValueFromMapComponent(id, name, mapKey)

	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(v.(string), 10, 64)

	if err != nil {
		return 0, errors.MapValueTypeMismatchError{
			ID:       id,
			Name:     name,
			Key:      mapKey,
			Expected: "int64",
			Actual:   "unable to parse int64",
		}
	}

	return i, nil
}

// GetMapComponent returns a map component from an entity. It takes the entity ID and component name. If the
// entity does not exist or the component does not exist, an error will be thrown. If the component is not a map,
// an error will be thrown.
func GetMapComponent(id string, name string) (map[string]interface{}, error) {

	exists, err := Exists(id)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.EntityNotFoundError{ID: id}
	}

	componentExists, err := ComponentExists(id, name)

	if err != nil {
		return nil, err
	}

	if !componentExists {
		return nil, errors.ComponentNotFoundError{ID: id, Name: name}
	}

	t, err := getComponentType(id, name)

	if err != nil {
		return nil, err
	}

	if t != "map" {
		return nil, errors.ComponentTypeMismatchError{ID: id, Name: name, Expected: "map", Actual: t}
	}

	h, err := engine.RedisHGetAll(componentId(id, name)).Result()

	if err != nil {
		return nil, err
	}

	final := make(map[string]interface{})

	for k, v := range h {
		if k == "" {
			continue
		}

		final[k] = v
	}

	return final, nil
}

// GetStringComponent returns the value of the string component. It takes the entity ID and component name. If the
// entity does not exist or the component does not exist, an error will be thrown. If the component is not a string,
// an error will be thrown.
func GetStringComponent(id string, name string) (string, error) {

	exists, err := Exists(id)

	if err != nil {
		return "", err
	}

	if !exists {
		return "", errors.EntityNotFoundError{ID: id}
	}

	componentExists, err := ComponentExists(id, name)

	if err != nil {
		return "", err
	}

	if !componentExists {
		return "", errors.ComponentNotFoundError{ID: id, Name: name}
	}

	t, err := getComponentType(id, name)

	if err != nil {
		return "", err
	}

	if t != "string" {
		return "", errors.ComponentTypeMismatchError{ID: id, Name: name, Expected: "string", Actual: t}
	}

	return engine.RedisGet(componentId(id, name)).Result()
}

// GetStringFromMapComponent returns the string value of an element in a map component. It takes the entity ID,
// component name, and the element name. If the entity does not exist or the component does not exist, an error will
// be thrown. If the component is not a map, an error will be thrown. If the element does not exist, an error will
// be thrown. If the element is not a string, an error will be thrown.
func GetStringFromMapComponent(id string, name string, mapKey string) (string, error) {
	v, err := getValueFromMapComponent(id, name, mapKey)

	if err != nil {
		return "", err
	}

	s, ok := v.(string)

	if !ok {
		return "", errors.MapValueTypeMismatchError{
			ID:       id,
			Name:     name,
			Key:      mapKey,
			Expected: "string",
			Actual:   reflect.TypeOf(v).Kind().String(),
		}
	}

	return s, nil
}

// IsEntityTypeRegistered checks if an entity type is registered. It takes the entity type name.
func IsEntityTypeRegistered(name string) bool {
	_, ok := registry.types[name]

	return ok
}

// Register registers an entity type. Entity Types must implement the `EntityType` interface. It is expected that
// developers can override default EntityType implementations with their own implementations.
func Register(e entity_type.EntityType) {
	registry.types[e.Name()] = e
}

// Replace removes and then replaces an entity in the entity registry. It takes the entity id, and a map of
// components. It will remove the entity with the provided id and then add the provided components to the entity. If an
// entity with the same id does not exist.
func Replace(id string, components map[string]interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return errors.EntityNotFoundError{ID: id}
	}

	log.Debug().Str("id", id).Msg("replacing entity")
	t, err := getEntityType(id)

	if err != nil {
		return err
	}

	// remove the entity
	err = Remove(id)

	if err != nil {
		return err
	}

	// add the entity
	err = AddWithId(t, id, components)

	if err != nil {
		return err
	}

	return nil
}

// Remove removes an entity from the entity registry. It takes the entity type and id. If an entity with the same
// id does not exist an error will be thrown.
func Remove(id string) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("entity %s does not exist", id)
	}

	log.Debug().Str("id", id).Msg("removing entity")
	// grab all the entities components
	components := AllComponents(id)

	for name := range components {
		err := RemoveComponent(id, name)

		if err != nil {
			return err
		}
	}

	removeMetadata(id)

	return nil
}

// RemoveComponent removes the component from the entity. If an entity with the same id does not exist an error will be
// thrown. If a component with the same name does not exist, an error will be thrown.
func RemoveComponent(id string, name string) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return errors.EntityNotFoundError{ID: id}
	}

	componentExists, err := ComponentExists(id, name)

	if err != nil {
		return err
	}

	if !componentExists {
		return errors.ComponentNotFoundError{ID: id, Name: name}
	}

	log.Debug().Str("id", id).Str("name", name).Msg("removing component")

	engine.RedisDel(componentId(id, name))

	return nil
}

// RemoveFromStringSetComponent removes a string value from a set component. It takes the entity ID, component name, and
// the value to remove from the set. If an entity with the same id does not exist an error will be thrown. If a component
// with the same name does not exist, an error will be thrown. If the value type is not a string, an error will be thrown.
func RemoveFromStringSetComponent(id string, name string, value string) error {
	return removeFromSetComponent(id, name, value)
}

// RemoveFromIntSetComponent removes an integer value from a set component. It takes the entity ID, component name, and
// the value to remove from the set. If an entity with the same id does not exist an error will be thrown. If a component
// with the same name does not exist, an error will be thrown. If the value type is not an integer, an error will be thrown.
func RemoveFromIntSetComponent(id string, name string, value int) error {
	return removeFromSetComponent(id, name, value)
}

// RemoveFromInt64SetComponent removes an integer64 value from a set component. It takes the entity ID, component name,
// and the value to remove from the set. If an entity with the same id does not exist an error will be thrown. If a
// component with the same name does not exist, an error will be thrown. If the value type is not an integer64, an error
// will be thrown.
func RemoveFromInt64SetComponent(id string, name string, value int64) error {
	return removeFromSetComponent(id, name, value)
}

// Update updates an entity in the entity registry. It takes the id, and a map of components. It will
// apply the provided components to the entity. If an entity with the same id does not exist, or the entity type does
// not match an error will be thrown. Any components that are not provided will be removed from the entity.
func Update(id string, components map[string]interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return errors.EntityNotFoundError{ID: id}
	}

	log.Debug().Str("id", id).Msg("updating entity")

	// grab all the entities components
	currentComponents := AllComponents(id)

	// remove the components that are not provided
	for name := range currentComponents {
		if _, ok := components[name]; !ok {
			_ = RemoveComponent(id, name)
		}
	}

	// set the components that are provided
	for name, value := range components {
		if _, ok := components[name]; ok {
			err := updateComponent(id, name, value)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// UpdateBoolComponent updates a bool component in the entity registry. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name does not exist, an error will be thrown. If the value type does not meet the type of the component, an error
// will be thrown.
func UpdateBoolComponent(id string, name string, value bool) error {
	return updateComponent(id, name, value)
}

// UpdateBoolInMapComponent updates a bool component in a map component. It takes the entity ID, component name, the
// key to which to add the value, and the value to add to the map. If an entity with the same id does not exist an error
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key does not
// already exist an error will be thrown. Once a value is added to the map, the type of that key is enforced. If
// the value is not the correct type an error will be thrown.
func UpdateBoolInMapComponent(id string, name string, key string, value bool) error {
	return updateInMapComponent(id, name, key, value)
}

// UpdateIntComponent updates an integer component in the entity registry. It takes the entity ID, component name, and
// the value of the component. If an entity with the same id does not exist an error will be thrown. If a component with
// the same name does not exist, an error will be thrown.
func UpdateIntComponent(id string, name string, value int) error {
	return updateComponent(id, name, value)
}

// UpdateIntInMapComponent updates an integer component in a map component. It takes the entity ID, component name, the
// key to which to add the value, and the value to add to the map. If an entity with the same id does not exist an error
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key does not
// already exist an error will be thrown. Once a value is added to the map, the type of that key is enforced. If
// the value is not the correct type an error will be thrown.
func UpdateIntInMapComponent(id string, name string, key string, value int) error {
	return updateInMapComponent(id, name, key, value)
}

// UpdateInt64Component updates an integer64 component in the entity registry. It takes the entity ID, component name,
// and the value of the component. If an entity with the same id does not exist an error will be thrown. If a component
// with the same name does not exist, an error will be thrown.
func UpdateInt64Component(id string, name string, value int64) error {
	return updateComponent(id, name, value)
}

// UpdateInt64InMapComponent updates an integer64 component in a map component. It takes the entity ID, component name,
// the key to which to add the value, and the value to add to the map. If an entity with the same id does not exist an
// error will be thrown. If a component with the same name does not exist, an error will be thrown. If the key does not
// already exist an error will be thrown. Once a value is added to the map, the type of that key is enforced. If
// the value is not the correct type an error will be thrown.
func UpdateInt64InMapComponent(id string, name string, key string, value int64) error {
	return updateInMapComponent(id, name, key, value)
}

// UpdateOrAdd updates an entity in the entity registry. It takes the entity type, id, and a map of components.
// It will apply the provided components to the entity. If an entity with the same id does not exist, it will add the
// entity with the provided id and components. If the entity type does not match the existing an error will be thrown.
// Any components that are not provided will be removed from the entity.
func UpdateOrAdd(entityType string, id string, components map[string]interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if exists {
		return Update(id, components)
	}

	return AddWithId(entityType, id, components)
}

// UpdateStringComponent updates a string component in the entity registry. It takes the entity ID, component name, and
// the value of the component. If an entity with the same id does not exist an error will be thrown. If a component with
// the same name does not exist, an error will be thrown.
func UpdateStringComponent(id string, name string, value string) error {
	return updateComponent(id, name, value)
}

// UpdateStringInMapComponent updates a string component in a map component. It takes the entity ID, component name,
// the key to which to add the value, and the value to add to the map. If an entity with the same id does not exist an
// error will be thrown. If a component with the same name does not exist, an error will be thrown. If the key does not
// already exist an error will be thrown. Once a value is added to the map, the type of that key is enforced. If
// the value is not the correct type an error will be thrown.
func UpdateStringInMapComponent(id string, name string, key string, value string) error {
	return updateInMapComponent(id, name, key, value)
}

func addOrUpdateToMapComponent(id string, name string, key string, value interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return errors.EntityNotFoundError{
			ID: id,
		}
	}

	exists, err = ComponentExists(id, name)

	if err != nil {
		return err
	}

	// TODO add testing for this
	if !exists {
		return errors.ComponentNotFoundError{
			ID:   id,
			Name: name,
		}
	}

	hasKey, err := mapHasKey(id, name, key)

	if err != nil {
		return err
	}

	if !hasKey {
		return addToMapComponent(id, name, key, value)
	}

	return updateInMapComponent(id, name, key, value)
}

func componentId(id string, key string) string {
	return fmt.Sprintf("%s:%s", id, key)
}

func generateID() string {
	return uuid.NewString()
}

func getComponentType(id string, key string) (string, error) {
	return engine.RedisGet(fmt.Sprintf("%s:%s", constants.ComponentTypePrefix, componentId(id, key))).Result()
}

func getMapValueType(id string, name string, key string) string {
	return engine.RedisHGet(
		fmt.Sprintf("%s:%s", constants.MapTypePrefix, componentId(id, name)),
		key,
	).String()
}

func getEntityTypeByName(name string) entity_type.EntityType {
	return registry.types[name]
}

func addComponent(id string, name string, value interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("entity %s does not exist", id)
	}

	componentExists, err := ComponentExists(id, name)

	if err != nil {
		return err
	}

	if componentExists {
		return fmt.Errorf("component %s already exists", name)
	}

	log.Debug().Str("id", id).Str("name", name).Msg("adding component")

	valueType := reflect.TypeOf(value).Kind().String()

	if valueType == "slice" {
		valueType = "set"
	}

	err = setComponentType(id, name, valueType)

	if err != nil {
		log.Error().Err(err).Msg("error setting component type")
		return err
	}

	cleanup := func() {
		engine.RedisDel(fmt.Sprintf("%s:%s", constants.ComponentTypePrefix, componentId(id, name)))
		engine.RedisDel(componentId(id, name))
		removeMetadata(id)
	}

	switch valueType {
	case "map":

		if len(value.(map[string]interface{})) == 0 {
			return nil
		}

		err := engine.RedisHMSet(componentId(id, name), value.(map[string]interface{})).Err()

		if err != nil {
			log.Error().Err(err).Msg("error setting component value")
			cleanup()
			return err
		}

		err = setMapValueTypes(id, name, value.(map[string]interface{}))

		if err != nil {
			log.Error().Err(err).Msg("error setting component value")
			cleanup()
			return err
		}
	case "set":

		if len(value.([]interface{})) == 0 {
			return nil
		}

		for _, value := range value.([]interface{}) {
			err = addToSetComponent(id, name, value)
			if err != nil {
				log.Error().Err(err).Msg("error setting component value")
				cleanup()
				return err
			}
		}

		err := validateAndSetSetValueTypes(id, name, value.([]interface{}))

		if err != nil {
			log.Error().Err(err).Msg("error setting component value type")
			cleanup()
			return err
		}

	default:
		err = engine.RedisSet(componentId(id, name), value).Err()

		if err != nil {
			return err
		}
	}

	err = updatePreviousValueToCurrent(id, name)

	if err != nil {
		log.Error().Err(err).Msg("error setting previous value")
		engine.RedisDel(fmt.Sprintf("%s:%s", constants.ComponentTypePrefix, componentId(id, name)))
		engine.RedisDel(componentId(id, name))
		return err
	}

	return nil
}

func validateAndSetSetValueTypes(id string, name string, values []interface{}) error {
	err := validateSetValueTypes(id, name, values)

	if err != nil {
		return err
	}

	return setSetValueType(id, name, reflect.TypeOf(values[0]).Kind().String())
}

func validateSetValueTypes(id string, name string, values []interface{}) error {
	baseType := reflect.TypeOf(values[0]).Kind().String()

	for _, value := range values {
		if reflect.TypeOf(value).Kind().String() != baseType {
			return errors.SetValueTypeError{
				ID:       id,
				Name:     name,
				Actual:   reflect.TypeOf(value).Kind().String(),
				Expected: baseType,
			}
		}
	}

	return nil
}

func setSetValueType(id string, name string, t string) error {
	return engine.RedisSet(
		fmt.Sprintf("%s:%s", constants.SetTypePrefix, componentId(id, name)),
		t,
	).Err()
}

func setMapValueTypes(id string, name string, value map[string]interface{}) error {
	for key, value := range value {
		err := setMapValueType(id, name, key, reflect.TypeOf(value).Kind().String())
		if err != nil {
			return err
		}
	}
	return nil
}

func setPreviousValue(id string, name string, value interface{}) error {
	t, err := getComponentType(id, name)

	if err != nil {
		return err
	}

	switch t {
	case "map":
		return engine.RedisHMSet(
			fmt.Sprintf("%s:%s", constants.PreviousValuePrefix, componentId(id, name)),
			value.(map[string]interface{}),
		).Err()
	case "set":
		return engine.RedisSAdd(
			fmt.Sprintf("%s:%s", constants.PreviousValuePrefix, componentId(id, name)),
			value.([]interface{}),
		).Err()
	default:
		return engine.RedisSet(
			fmt.Sprintf("%s:%s", constants.PreviousValuePrefix, componentId(id, name)),
			value,
		).Err()
	}
}

func getPreviousValue(id string, name string) (interface{}, error) {
	t, err := getComponentType(id, name)

	if err != nil {
		return nil, err
	}

	switch t {
	case "map":
		return engine.RedisHGetAll(fmt.Sprintf("%s:%s", constants.PreviousValuePrefix, componentId(id, name))).Result()
	case "set":
		return engine.RedisSMembers(fmt.Sprintf("%s:%s", constants.PreviousValuePrefix, componentId(id, name))).Result()
	default:
		return engine.RedisGet(fmt.Sprintf("%s:%s", constants.PreviousValuePrefix, componentId(id, name))).Result()
	}
}

func getCurrentValue(id string, name string) (interface{}, error) {
	exists, err := Exists(id)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.EntityNotFoundError{
			ID: id,
		}
	}

	exists, err = ComponentExists(id, name)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.ComponentNotFoundError{
			ID:   id,
			Name: name,
		}
	}

	t, err := getComponentType(id, name)

	if err != nil {
		return nil, err
	}

	switch t {
	case "map":
		return engine.RedisHGetAll(componentId(id, name)).Result()
	case "set":
		return engine.RedisSMembers(componentId(id, name)).Result()
	default:
		return engine.RedisGet(componentId(id, name)).Result()
	}
}

func updatePreviousValueToCurrent(id string, name string) error {
	currentValue, err := getCurrentValue(id, name)

	if err != nil {
		return err
	}

	t, err := getComponentType(id, name)

	if err != nil {
		return err
	}

	switch t {
	case "map":
		m := make(map[string]interface{})

		for key, value := range currentValue.(map[string]string) {
			m[key] = value
		}

		return setPreviousValue(id, name, m)
	case "set":
		newSet := make([]interface{}, 0)

		for _, value := range currentValue.([]string) {
			newSet = append(newSet, value)
		}

		return setPreviousValue(id, name, newSet)
	default:
		return setPreviousValue(id, name, currentValue)
	}
}

func addToMapComponent(id string, name string, mapKey string, value interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return errors.EntityExistsError{ID: id}
	}

	componentExists, err := ComponentExists(id, name)

	if err != nil {
		return err
	}

	if !componentExists {
		return errors.MissingComponentError{ID: id, Name: name}
	}

	hasKey, err := mapHasKey(id, name, mapKey)

	if err != nil {
		return err
	}

	if hasKey {
		return errors.MapHasKeyError{ID: id, Name: name, Key: mapKey}
	}

	log.Debug().Str("id", id).Str("name", name).Str("mapKey", mapKey).Msg("adding to map component")

	err = engine.RedisHSet(componentId(id, name), mapKey, value).Err()

	if err != nil {
		return err
	}

	err = setMapValueType(id, name, mapKey, value)

	if err != nil {
		return err
	}

	return nil
}

func addToSetComponent(id string, name string, value interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return errors.EntityExistsError{ID: id}
	}

	componentExists, err := ComponentExists(id, name)

	if err != nil {
		return err
	}

	if !componentExists {
		return errors.MissingComponentError{ID: id, Name: name}
	}

	// check the value type against the set value type
	valueType := reflect.TypeOf(value).Kind().String()
	setValueType, err := getSetValueType(id, name)

	if err != nil && err.Error() != "redis: nil" {
		return err
	} else if err == nil {
		if valueType != setValueType {
			return errors.SetValueTypeError{ID: id, Name: name, Expected: valueType, Actual: setValueType}
		}
	}

	log.Debug().Str("id", id).Str("name", name).Msg("adding to set component")

	err = engine.RedisSAdd(componentId(id, name), value).Err()

	if err != nil {
		return err
	}

	return nil
}

func componentTypeMatch(id string, name string, value interface{}) (bool, error) {
	t, err := getComponentType(id, name)

	if err != nil {
		return false, err
	}

	return reflect.TypeOf(value).Kind().String() == t, nil
}

func getEntityType(id string) (string, error) {
	return engine.RedisGet(entityTypeID(id)).Result()
}

func getSetValueType(id string, name string) (string, error) {
	return engine.RedisGet(fmt.Sprintf("%s:%s", constants.SetTypePrefix, componentId(id, name))).
		Result()
}

func getValueFromMapComponent(id string, name string, mapKey string) (interface{}, error) {
	exists, err := Exists(id)

	if err != nil {
		return false, err
	}

	if !exists {
		return nil, errors.EntityNotFoundError{ID: id}
	}

	componentExists, err := ComponentExists(id, name)

	if err != nil {
		return nil, err
	}

	if !componentExists {
		return nil, errors.MissingComponentError{ID: id, Name: name}
	}

	t, err := getComponentType(id, name)

	if err != nil {
		return nil, err
	}

	if t != "map" {
		return nil, errors.ComponentTypeMismatchError{ID: id, Name: name, Expected: "map", Actual: t}
	}

	mhk, err := mapHasKey(id, name, mapKey)

	if err != nil {
		return nil, err
	}

	if !mhk {
		return nil, errors.MapKeyNotFoundError{ID: id, Name: name, Key: mapKey}
	}

	return engine.RedisHGet(componentId(id, name), mapKey).Result()
}

func mapHasKey(id string, name string, mapKey string) (bool, error) {
	return engine.RedisHExists(componentId(id, name), mapKey).Result()
}

func mapValueMatch(id string, name string, mapKey string, value interface{}) bool {
	return reflect.TypeOf(value).Kind().String() == getMapValueType(id, name, mapKey)
}

func removeMetadata(id string) {
	keys := engine.RedisKeys(fmt.Sprintf("__*:%s*", id)).Val()

	for _, key := range keys {
		engine.RedisDel(key)
	}
}

func removeFromSetComponent(id string, name string, value interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return errors.EntityExistsError{ID: id}
	}

	componentExists, err := ComponentExists(id, name)

	if err != nil {
		return err
	}

	if !componentExists {
		return errors.MissingComponentError{ID: id, Name: name}
	}

	// check the value type against the set value type
	valueType := reflect.TypeOf(value).Kind().String()
	setValueType, err := getSetValueType(id, name)

	if err != nil {
		return err
	}

	if valueType != setValueType {
		return errors.SetValueTypeError{ID: id, Name: name, Expected: valueType, Actual: setValueType}
	}

	log.Debug().Str("id", id).Str("name", name).Msg("removing from set component")

	err = engine.RedisSRem(componentId(id, name), value).Err()

	if err != nil {
		return err
	}

	return nil
}

func componentTypeID(id string, name string) string {
	return fmt.Sprintf("%s:%s", constants.ComponentTypePrefix, componentId(id, name))
}

func entityTypeID(id string) string {
	return fmt.Sprintf("%s:%s", constants.EntityTypePrefix, id)
}

func setEntityType(id string, entityType string) {
	entId := entityTypeID(id)
	log.Trace().Str("id", id).Str("entityType", entityType).Str("entId", entId).Msg("setting entity type")
	engine.RedisSet(entId, entityType)
}

func updateComponent(id string, name string, value interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return errors.EntityNotFoundError{ID: id}
	}

	componentExists, err := ComponentExists(id, name)

	if err != nil {
		return err
	}

	if !componentExists {
		return errors.MissingComponentError{ID: id, Name: name}
	}

	m, err := componentTypeMatch(id, name, value)

	if err != nil {
		return err
	}

	t, err := getComponentType(id, name)

	if err != nil {
		return err
	}

	if !m {
		return errors.ComponentTypeMismatchError{
			ID:       id,
			Name:     name,
			Expected: t,
			Actual:   reflect.TypeOf(value).Kind().String(),
		}
	}

	log.Debug().Str("id", id).Str("name", name).Msg("updating component")

	err = engine.RedisSet(componentId(id, name), value).Err()

	if err != nil {
		return err
	}

	return nil
}

func updateInMapComponent(id string, name string, mapKey string, value interface{}) error {
	// check if the entity exists
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return errors.EntityExistsError{ID: id}
	}

	// check if the component exists
	componentExists, err := ComponentExists(id, name)

	if err != nil {
		return err
	}

	if !componentExists {
		return errors.MissingComponentError{ID: id, Name: name}
	}

	t, err := getComponentType(id, name)

	if err != nil {
		return err
	}

	// check if the component is of type 'map
	if t != "map" {
		return errors.ComponentTypeMismatchError{
			ID:       id,
			Name:     name,
			Expected: "map",
			Actual:   t,
		}
	}

	// check if the key exists on the map
	hasKey, err := mapHasKey(id, name, mapKey)

	if err != nil {
		return err
	}

	if !hasKey {
		return errors.MapKeyNotFoundError{ID: id, Name: name, Key: mapKey}
	}

	// check if the value type matches the expected type
	if mapValueMatch(id, name, mapKey, value) {
		return errors.MapValueTypeMismatchError{
			ID:       id,
			Name:     name,
			Key:      mapKey,
			Expected: getMapValueType(id, name, mapKey),
			Actual:   reflect.TypeOf(value).Kind().String(),
		}
	}

	log.Debug().Str("id", id).Str("name", name).Str("mapKey", mapKey).Msg("updating in map component")

	return engine.RedisHSet(componentId(id, name), mapKey, value).Err()
}

func setMapValueType(id string, name string, mapKey string, value interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return errors.EntityExistsError{ID: id}
	}

	componentExists, err := ComponentExists(id, name)

	if err != nil {
		return err
	}

	if !componentExists {
		return errors.MissingComponentError{ID: id, Name: name}
	}

	hasKey, err := mapHasKey(id, name, mapKey)

	if err != nil {
		return err
	}

	if !hasKey {
		return errors.MapKeyMissingError{ID: id, Name: name, Key: mapKey}
	}

	log.Trace().Str("id", id).Str("name", name).Str("mapKey", mapKey).Msg("setting map component value type")

	err = engine.RedisHSet(
		fmt.Sprintf("%s:%s", constants.MapTypePrefix, componentId(id, name)),
		mapKey,
		reflect.TypeOf(value).Kind().String(),
	).Err()

	if err != nil {
		return err
	}

	return nil
}

func setComponentType(id string, name string, value interface{}) error {
	err := engine.RedisSet(
		fmt.Sprintf("%s:%s", constants.ComponentTypePrefix, componentId(id, name)),
		value,
	).
		Err()

	if err != nil {
		return err
	}

	return nil
}

func setComponentsFromMap(id string, components map[string]interface{}) error {
	errs := make([]error, 0)

	for name, value := range components {
		err := addComponent(id, name, value)

		if err != nil {
			errs = append(errs, err)
		}

	}

	if len(errs) > 0 {
		return errors.AddComponentErrors{Errors: errs}
	}

	return nil
}
