package entity_registry

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/ecs/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/ecs/pkg/entity_type"
)

const EntityTypePrefix = "__entityType"
const ComponentTypePrefix = "__type"
const MapTypePrefix = "__map"
const SetTypePrefix = "__set"

// AddComponentErrors is a collection of errors that occurred while adding components to an entity.
type AddComponentErrors struct {
	Errors []error
}

func (e AddComponentErrors) Error() string {
	errorStrings := make([]string, len(e.Errors))

	for i, err := range e.Errors {
		errorStrings[i] = err.Error()
	}

	return fmt.Sprintf(
		"%d errors occurred while adding components to an entity: %s", len(e.Errors),
		strings.Join(errorStrings, ", "),
	)
}

// ComponentNotFoundError is returned when a component is not found.
type ComponentNotFoundError struct {
	ID   string
	Name string
}

func (e ComponentNotFoundError) Error() string {
	return fmt.Sprintf("component %s not found for entity %s", e.Name, e.ID)
}

// ComponentTypeMismatchError is an error type that is returned when the component type does not match the entity type.
type ComponentTypeMismatchError struct {
	ID       string
	Name     string
	Expected string
	Actual   string
}

func (e ComponentTypeMismatchError) Error() string {
	return fmt.Sprintf(
		"component type mismatch for entity %s, component %s, expected %s, actual %s",
		e.ID,
		e.Name,
		e.Expected,
		e.Actual,
	)
}

// EmptySetError is an error type that is returned when an empty set is attempted to be added to a entity .
type EmptySetError struct {
	ID string
}

func (e EmptySetError) Error() string {
	return fmt.Sprintf("empty set cannot be added %s", e.ID)
}

// EntityExistsError is an error type that is returned when an entity with the given id already exists.
type EntityExistsError struct {
	ID string
}

func (e EntityExistsError) Error() string {
	return fmt.Sprintf("entity with id %s already exists", e.ID)
}

// EntityNotFoundError is an error type that is returned when an entity with the given id does not exist.
type EntityNotFoundError struct {
	ID string
}

func (e EntityNotFoundError) Error() string {
	return fmt.Sprintf("entity with id %s does not exist", e.ID)
}

// EntityTypeNotRegisteredError is called when an entity type is not registered.
type EntityTypeNotRegisteredError struct {
	Type string
}

func (e EntityTypeNotRegisteredError) Error() string {
	return fmt.Sprintf("entity type %s is not registered", e.Type)
}

// HasComponentError is an error type that is returned when the entity already has the given component.
type HasComponentError struct {
	ID   string
	Name string
}

func (e HasComponentError) Error() string {
	return fmt.Sprintf("entity %s already has component %s", e.ID, e.Name)
}

// MapHasKeyError is an error type that is returned when the map already has the given key.
type MapHasKeyError struct {
	ID   string
	Name string
	Key  string
}

func (e MapHasKeyError) Error() string {
	return fmt.Sprintf("map %s for entity %s already has key %s", e.Name, e.ID, e.Key)
}

// MapValueTypeMismatchError is an error type that is returned when the map key does not match the entity type.
type MapValueTypeMismatchError struct {
	ID       string
	Name     string
	Key      string
	Expected string
	Actual   string
}

func (e MapValueTypeMismatchError) Error() string {
	return fmt.Sprintf(
		"map value type mismatch for entity %s, component %s, key %s, expected %s, actual %s",
		e.ID,
		e.Name,
		e.Key,
		e.Expected,
		e.Actual,
	)
}

// MissingComponentError is an error type that is returned when the entity does not have the given component.
type MissingComponentError struct {
	ID   string
	Name string
}

func (e MissingComponentError) Error() string {
	return fmt.Sprintf("entity %s does not have component %s", e.ID, e.Name)
}

type SetValueTypeError struct {
	ID       string
	Name     string
	Expected string
	Actual   string
}

func (e SetValueTypeError) Error() string {
	return fmt.Sprintf(
		"set value type mismatch for entity %s, component %s, expected %s, actual %s",
		e.ID,
		e.Name,
		e.Expected,
		e.Actual,
	)
}

type entityRegistry struct {
	types               map[string]entity_type.EntityType
	loadableDirectories []string
}

var registry = entityRegistry{
	types:               make(map[string]entity_type.EntityType),
	loadableDirectories: make([]string, 0),
}

var log = logger.Logger.With().Str("service", "entity_registry").Logger()

// Start starts the entity registry.
func Start() {
	log.Info().Msg("starting entity registry")
	//
	//for _, dir := range registry.loadableDirectories {
	//	LoadDirectory(dir)
	//}
}

// Stop stops the entity registry.
func Stop() {
	log.Info().Msg("stopping entity registry")
}

// AddStringToSetComponent  adds a string to a set.
func AddStringToSetComponent(id string, name string, value string) {
	engine.Redis.SAdd(context.Background(), componentId(id, name), value)
}

// AllComponents returns all the components for the entity.
func AllComponents(id string) map[string]interface{} {
	keys, err := engine.Redis.Keys(context.Background(), fmt.Sprintf("%s:*", id)).Result()

	if err != nil {
		panic(err)
	}

	components := make(map[string]interface{})

	for _, key := range keys {
		name := strings.Replace(key, fmt.Sprintf("%s:", id), "", 1)
		t := engine.Redis.Get(context.Background(), fmt.Sprintf("__type:%s", key)).Val()

		switch t {
		case "map":
			m, err := engine.Redis.HGetAll(context.Background(), key).Result()
			if err != nil {
				panic(err)
			}

			var final = make(map[string]interface{})

			for k, v := range m {
				final[k] = v
			}

			components[name] = final
		case "set":
			s := engine.Redis.SMembers(context.Background(), key).Val()
			set := make([]interface{}, 0)

			for _, v := range s {
				set = append(set, v)
			}

			components[name] = set

		default:
			components[name] = engine.Redis.Get(context.Background(), key).Val()
		}

	}

	return components
}

//
//// AllEntitiesByType returns all entities of the given type.
//func AllEntitiesByType(entityType string) []string {
//	keys, err := engine.Redis.Keys(context.Background(), fmt.Sprintf("*:%s", TypeComponentName)).Result()
//
//	if err != nil {
//		panic(err)
//	}
//
//	// create a new slice which holds only entities whose type matches the given type
//	entities := make([]string, 0)
//	for _, key := range keys {
//		if strings.Contains(key, entityType) {
//			entities = append(entities, strings.Replace(key, fmt.Sprintf(":%s", TypeComponentName), "", 1))
//		}
//	}
//
//	return entities
//}
//
//func AllEntitiesByTypeWithComponent(entityType string, component string) []string {
//	ents := AllEntitiesByType(entityType)
//
//	entities := make([]string, 0)
//	for _, ent := range ents {
//		if ComponentExists(ent, component) {
//			entities = append(entities, ent)
//		}
//	}
//
//	return entities
//}
//
//func AllEntitiesByTypeWithComponentValue(entityType string, component string, value interface{}) []string {
//	ents := AllEntitiesByType(entityType)
//
//	entities := make([]string, 0)
//	for _, ent := range ents {
//		if HasComponentValue(ent, component, value) {
//			entities = append(entities, ent)
//		}
//	}
//
//	return entities
//}
//
//// HasComponentValue returns true if the component exists and has the given value.
//func HasComponentValue(id string, name string, value interface{}) bool {
//	exists := engine.Redis.Exists(context.Background(), componentId(id, name)).Val()
//
//	if exists == int64(0) {
//		return false
//	}
//
//	return compareValues(componentId(id, name), value)
//}

//// LoadDirectory loads entities from a directory.
//func LoadDirectory(dir string) {
//	// load all yml files in the directory
//	files, err := os.Open(dir)
//
//	if err != nil {
//		log.Error().Err(err).Msgf("failed to open directory %s", dir)
//		return
//	}
//
//	defer func() {
//		_ = files.Close()
//	}()
//
//	names, err := files.Readdirnames(0)
//
//	if err != nil {
//		log.Error().Err(err).Msgf("failed to read directory %s", dir)
//		return
//	}
//
//	for _, name := range names {
//		if strings.HasSuffix(name, ".yml") {
//			loadFromFile(fmt.Sprintf("%s/%s", dir, name))
//		}
//	}
//}

//// RegisterLoadableDirectory registers a directory to load entities from.
//func RegisterLoadableDirectory(dir string) {
//	log.Debug().Str("dir", dir).Msg("registering loadable directory")
//	registry.loadableDirectories = append(registry.loadableDirectories, dir)
//}

// Add adds an entity to the entity registry. It takes a type, a map of components to be added to the entity. If the
// entity already exists, an error will be returned. If the type is not registered, an error will be returned.
func Add(entityType string, args map[string]interface{}) (string, error) {
	id := generateID()
	err := AddWithID(entityType, id, args)

	if err != nil {
		return "", err
	}

	return id, nil
}

// AddWithID adds an entity with the provided id to the entity registry. It takes the entity id,
// and a map of components to be added. If an entity with the same id already exists, an error will be thrown. If the
// type is not registered, an error will be thrown.
func AddWithID(entityType string, id string, args map[string]interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if exists {
		return EntityExistsError{ID: id}
	}

	if !IsEntityTypeRegistered(entityType) {
		return EntityTypeNotRegisteredError{Type: entityType}
	}

	log.Debug().Str("id", id).Msg("adding entity")
	setEntityType(id, entityType)

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
		return addComponent(id, name, 1)
	} else {
		return addComponent(id, name, 0)
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
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return EntityNotFoundError{ID: id}
	}

	err = setComponentType(id, name, "map")

	if err != nil {
		return err
	}

	log.Debug().Str("id", id).Str("name", name).Msg("adding map component")

	engine.Redis.HSet(context.Background(), componentId(id, name), value)

	if err != nil {
		removeComponentMetadata(id, name)
		return err
	}

	return nil
}

// AddSetComponent adds a set component to an entity. It takes the entity ID, component name, and the value of
// the component. If an entity with the same id already exists error will be thrown. If a component with the same name
// already exists, an error will be thrown. Once a value is added to the set, the type of that value is enforced for
// all members of the set. Attempting to change the type of value will result in an error in later updates. One cannot
// add empty sets, so if the value is empty, an error will be thrown.
func AddSetComponent(id string, name string, value []interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return EntityExistsError{ID: id}
	}

	err = setComponentType(id, name, "set")

	if err != nil {
		return err
	}

	if len(value) == 0 {
		return EmptySetError{ID: id}
	}

	firstElement := value[0]

	setType := reflect.TypeOf(firstElement).String()

	for _, element := range value {
		if reflect.TypeOf(element).String() != setType {
			return SetValueTypeError{ID: id, Name: name, Expected: setType, Actual: reflect.TypeOf(element).String()}
		}
	}

	log.Debug().Str("id", id).Str("name", name).Msg("adding set component")

	err = engine.Redis.Set(
		context.Background(),
		fmt.Sprintf("%s:%s", SetTypePrefix, componentId(id, name)),
		setType,
		0,
	).Err()

	if err != nil {
		removeComponentMetadata(id, name)
		return err
	}

	err = engine.Redis.SAdd(context.Background(), componentId(id, name), value...).Err()

	if err != nil {
		removeComponentMetadata(id, name)
		return err
	}

	return nil
}

// AddStringComponent adds a string component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func AddStringComponent(id string, name string, value string) error {
	return addComponent(id, name, value)
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
func ComponentExists(id string, name string) bool {
	i := engine.Redis.Exists(context.Background(), componentTypeID(id, name)).Val()

	if i == 1 {
		return true
	}

	return false
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
		return false, EntityNotFoundError{ID: id}
	}

	if !ComponentExists(id, name) {
		return false, ComponentNotFoundError{ID: id, Name: name}
	}

	return engine.Redis.Get(context.Background(), componentId(id, name)).Bool()
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
		return 0, EntityNotFoundError{ID: id}
	}

	if !ComponentExists(id, name) {
		return 0, ComponentNotFoundError{ID: id, Name: name}
	}

	return engine.Redis.Get(context.Background(), componentId(id, name)).Int64()
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
		return 0, EntityNotFoundError{ID: id}
	}

	if !ComponentExists(id, name) {
		return 0, ComponentNotFoundError{ID: id, Name: name}
	}

	return engine.Redis.Get(context.Background(), componentId(id, name)).Int()
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

	i, ok := v.(int)

	if !ok {
		return 0, MapValueTypeMismatchError{
			ID:       id,
			Name:     name,
			Key:      mapKey,
			Expected: "int",
			Actual:   reflect.TypeOf(v).String(),
		}
	}

	return i, nil
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

	i, ok := v.(int64)

	if !ok {
		return 0, MapValueTypeMismatchError{
			ID:       id,
			Name:     name,
			Key:      mapKey,
			Expected: "int64",
			Actual:   reflect.TypeOf(v).String(),
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
		return nil, EntityNotFoundError{ID: id}
	}

	if !ComponentExists(id, name) {
		return nil, ComponentNotFoundError{ID: id, Name: name}
	}

	if getComponentType(id, name) != "map" {
		return nil, ComponentTypeMismatchError{ID: id, Name: name, Expected: "map", Actual: getComponentType(id, name)}
	}

	h, err := engine.Redis.HGetAll(context.Background(), componentId(id, name)).Result()

	if err != nil {
		return nil, err
	}

	final := make(map[string]interface{})

	for k, v := range h {
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
		return "", EntityNotFoundError{ID: id}
	}

	if !ComponentExists(id, name) {
		return "", ComponentNotFoundError{ID: id, Name: name}
	}

	if getComponentType(id, name) != "string" {
		return "", ComponentTypeMismatchError{ID: id, Name: name, Expected: "string", Actual: getComponentType(id, name)}
	}

	return engine.Redis.Get(context.Background(), componentId(id, name)).Result()
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
		return "", MapValueTypeMismatchError{
			ID:       id,
			Name:     name,
			Key:      mapKey,
			Expected: "string",
			Actual:   reflect.TypeOf(v).String(),
		}
	}

	return s, nil
}

// GetStringsFromSetComponent GetStringSetComponent returns the value of the string set component.
func GetStringsFromSetComponent(id string, name string) ([]string, error) {
	return engine.Redis.SMembers(context.Background(), componentId(id, name)).Result()
}

func IsEntityTypeRegistered(name string) bool {
	_, ok := registry.types[name]

	return ok
}

// Register registers an entity type. Entity Types must implmeent the `EntityType` interface. It is expected that
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
		return EntityNotFoundError{ID: id}
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
	err = AddWithID(t, id, components)

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
		return EntityNotFoundError{ID: id}
	}

	if !ComponentExists(id, name) {
		return ComponentNotFoundError{ID: id, Name: name}
	}

	log.Debug().Str("id", id).Str("name", name).Msg("removing component")

	engine.Redis.Del(context.Background(), componentId(id, name))
	removeComponentMetadata(id, name)

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
		return EntityNotFoundError{ID: id}
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
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key soes not
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
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key soes not
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
// error will be thrown. If a component with the same name does not exist, an error will be thrown. If the key soes not
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

	return AddWithID(entityType, id, components)
}

// UpdateStringComponent updates a string component in the entity registry. It takes the entity ID, component name, and
// the value of the component. If an entity with the same id does not exist an error will be thrown. If a component with
// the same name does not exist, an error will be thrown.
func UpdateStringComponent(id string, name string, value string) error {
	return updateComponent(id, name, value)
}

// UpdateStringInMapComponent updates a string component in a map component. It takes the entity ID, component name,
// the key to which to add the value, and the value to add to the map. If an entity with the same id does not exist an
// error will be thrown. If a component with the same name does not exist, an error will be thrown. If the key soes not
// already exist an error will be thrown. Once a value is added to the map, the type of that key is enforced. If
// the value is not the correct type an error will be thrown.
func UpdateStringInMapComponent(id string, name string, key string, value string) error {
	return updateInMapComponent(id, name, key, value)
}

func compareValues(componentId string, value interface{}) bool {
	comparisonType := reflect.TypeOf(value).Kind()

	switch comparisonType {
	case reflect.String:
		str, err := engine.Redis.Get(context.Background(), componentId).Result()

		if err != nil {
			return false
		}

		return str == value.(string)
	case reflect.Int:
		intValue, err := engine.Redis.Get(context.Background(), componentId).Int()

		if err != nil {
			return false
		}

		return intValue == value.(int)
	case reflect.Int64:
		intValue, err := engine.Redis.Get(context.Background(), componentId).Int64()

		if err != nil {
			return false
		}

		return intValue == value.(int64)
	case reflect.Bool:
		boolValue, err := engine.Redis.Get(context.Background(), componentId).Bool()

		if err != nil {
			return false
		}

		return boolValue == value.(bool)
	case reflect.Slice:
		slice, err := engine.Redis.SMembers(context.Background(), componentId).Result()

		if err != nil {
			return false
		}

		return reflect.DeepEqual(slice, value.([]string))
	case reflect.Map:
		mapValue, err := engine.Redis.HGetAll(context.Background(), componentId).Result()

		if err != nil {
			return false
		}

		imap := make(map[string]interface{})

		for k, v := range mapValue {
			imap[k] = v
		}

		return reflect.DeepEqual(imap, value.(map[string]interface{}))
	}

	return false
}

func componentId(id string, key string) string {
	return fmt.Sprintf("%s:%s", id, key)
}

func generateID() string {
	return uuid.NewString()
}

func getComponentType(id string, key string) string {
	return engine.Redis.Get(context.Background(), fmt.Sprintf("__type:%s", componentId(id, key))).String()
}

func getMapValueType(id string, name string, key string) string {
	return engine.Redis.HGet(
		context.Background(),
		fmt.Sprintf("%s:%s", MapTypePrefix, componentId(id, name)),
		key,
	).String()
}

func getEntityTypeByName(name string) entity_type.EntityType {
	return registry.types[name]
}

func getFromHashComponent(id string, name string, mapKey string) (interface{}, error) {
	return engine.Redis.HGet(context.Background(), componentId(id, name), mapKey).Result()
}

//
//func loadFromFile(file string) {
//	// load from a yml file and process
//	log.Debug().Str("file", file).Msg("loading entities from file")
//	content, err := os.ReadFile(file)
//
//	if err != nil {
//		log.Error().Str("file", file).Err(err).Msgf("failed to read file %s", file)
//		return
//	}
//
//	loadFromYaml(content)
//}

//func loadFromMap(id string, m map[string]interface{}) {
//	// remove the entity, if we are loading from a map we want to load fresh every time
//	Remove(id)
//	// load from a map and process
//	log.Debug().Str("id", id).Msg("loading entity from map")
//
//	err := AddWithID(m[TYPE_COMPONENT_NAME].(string), id, m)
//
//	if err != nil {
//		panic(err)
//	}
//}
//
//func loadFromYaml(yml []byte) {
//	// load yaml and process each entity using loadMap
//	ymlEntries := make(map[string]interface{})
//
//	err := yaml.Unmarshal(yml, ymlEntries)
//
//	if err != nil {
//		log.Error().Err(err).Msg("failed to parse yaml")
//		return
//	}
//
//	for id, m := range ymlEntries {
//		loadFromMap(id, m.(map[string]interface{}))
//	}
//}

func addComponent(id string, name string, value interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("entity %s does not exist", id)
	}

	if ComponentExists(id, name) {
		return fmt.Errorf("component %s already exists", name)
	}

	log.Debug().Str("id", id).Str("name", name).Msg("adding component")

	valueType := reflect.TypeOf(value).Kind().String()

	err = engine.
		Redis.
		Set(context.Background(), fmt.Sprintf("__type:%s", componentId(id, name)), valueType, 0).
		Err()

	if err != nil {
		return err
	}

	err = engine.Redis.Set(context.Background(), componentId(id, name), value, 0).Err()

	if err != nil {
		engine.Redis.Del(context.Background(), fmt.Sprintf("__type:%s", componentId(id, name)))
		return err
	}

	return nil
}

func addToMapComponent(id string, name string, mapKey string, value interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return EntityExistsError{ID: id}
	}

	if !ComponentExists(id, name) {
		return MissingComponentError{ID: id, Name: name}
	}

	hasKey, err := mapHasKey(id, name, mapKey)

	if err != nil {
		return err
	}

	if hasKey {
		return MapHasKeyError{ID: id, Name: name, Key: mapKey}
	}

	log.Debug().Str("id", id).Str("name", name).Str("mapKey", mapKey).Msg("adding to map component")

	err = setMapValueType(id, name, mapKey, value)

	if err != nil {
		return err
	}

	err = engine.Redis.HSet(context.Background(), componentId(id, name), mapKey, value).Err()

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
		return EntityExistsError{ID: id}
	}

	if !ComponentExists(id, name) {
		return MissingComponentError{ID: id, Name: name}
	}

	// check the value type aginst the set value type
	valueType := reflect.TypeOf(value).Kind().String()
	setValueType := getSetValueType(id, name)

	if valueType != setValueType {
		return SetValueTypeError{ID: id, Name: name, Expected: valueType, Actual: setValueType}
	}

	log.Debug().Str("id", id).Str("name", name).Msg("adding to set component")

	err = engine.Redis.SAdd(context.Background(), componentId(id, name), value).Err()

	if err != nil {
		return err
	}

	return nil
}

func componentTypeMatch(id string, name string, value interface{}) bool {
	return reflect.TypeOf(value).Kind().String() == getComponentType(id, name)
}

func getEntityType(id string) (string, error) {
	return engine.Redis.Get(context.Background(), entityTypeID(id)).Result()
}

func getSetValueType(id string, name string) string {
	return engine.Redis.
		Type(context.Background(), fmt.Sprintf("%s:%s", SetTypePrefix, componentId(id, name))).
		String()
}

func getValueFromMapComponent(id string, name string, mapKey string) (interface{}, error) {
	exists, err := Exists(id)

	if err != nil {
		return false, err
	}

	if !exists {
		return nil, EntityExistsError{ID: id}
	}

	if !ComponentExists(id, name) {
		return nil, MissingComponentError{ID: id, Name: name}
	}

	if getComponentType(id, name) != "map" {
		return nil, ComponentTypeMismatchError{ID: id, Name: name, Expected: "map", Actual: getComponentType(id, name)}
	}

	mhk, err := mapHasKey(id, name, mapKey)

	if err != nil {
		return nil, err
	}

	if mhk {
		return nil, MapHasKeyError{ID: id, Name: name, Key: mapKey}
	}

	return engine.Redis.HGet(context.Background(), componentId(id, name), mapKey).Result()
}

func mapHasKey(id string, name string, mapKey string) (bool, error) {
	return engine.Redis.HExists(context.Background(), componentId(id, name), mapKey).Result()
}

func mapValueMatch(id string, name string, mapKey string, value interface{}) bool {
	return reflect.TypeOf(value).Kind().String() == getMapValueType(id, name, mapKey)
}

func removeComponentMetadata(id string, name string) {
	keys := engine.Redis.Keys(context.Background(), fmt.Sprintf("__*:%s", componentId(id, name))).Val()

	for _, key := range keys {
		engine.Redis.Del(context.Background(), key)
	}
}

func removeFromSetComponent(id string, name string, value interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return EntityExistsError{ID: id}
	}

	if !ComponentExists(id, name) {
		return MissingComponentError{ID: id, Name: name}
	}

	// check the value type aginst the set value type
	valueType := reflect.TypeOf(value).Kind().String()
	setValueType := getSetValueType(id, name)

	if valueType != setValueType {
		return SetValueTypeError{ID: id, Name: name, Expected: valueType, Actual: setValueType}
	}

	log.Debug().Str("id", id).Str("name", name).Msg("removing from set component")

	err = engine.Redis.SRem(context.Background(), componentId(id, name), value).Err()

	if err != nil {
		return err
	}

	return nil
}

func componentTypeID(id string, name string) string {
	return fmt.Sprintf("%s:%s", ComponentTypePrefix, componentId(id, name))
}

func entityTypeID(id string) string {
	return fmt.Sprintf("%s:%s", EntityTypePrefix, id)
}

func setEntityType(id string, entityType string) {
	etid := entityTypeID(id)
	log.Trace().Str("id", id).Str("entityType", entityType).Str("etid", etid).Msg("setting entity type")
	engine.Redis.Set(context.Background(), etid, entityType, 0)
}

func updateComponent(id string, name string, value interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return EntityExistsError{ID: id}
	}

	if !ComponentExists(id, name) {
		return MissingComponentError{ID: id, Name: name}
	}

	if !componentTypeMatch(id, name, value) {
		return ComponentTypeMismatchError{
			ID:       id,
			Name:     name,
			Expected: getComponentType(id, name),
			Actual:   reflect.TypeOf(value).Kind().String(),
		}
	}

	log.Debug().Str("id", id).Str("name", name).Msg("updating component")

	err = engine.Redis.Set(context.Background(), componentId(id, name), value, 0).Err()

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
		return EntityExistsError{ID: id}
	}

	// check if the component exists
	if !ComponentExists(id, name) {
		return MissingComponentError{ID: id, Name: name}
	}

	// check if the component is of type 'map
	if getComponentType(id, name) != "map" {
		return ComponentTypeMismatchError{
			ID:       id,
			Name:     name,
			Expected: "map",
			Actual:   getComponentType(id, name),
		}
	}

	// check if the key exists on the map
	hasKey, err := mapHasKey(id, name, mapKey)

	if err != nil {
		return err
	}

	if !hasKey {
		return MapHasKeyError{ID: id, Name: name, Key: mapKey}
	}

	// check if the value type matches the expected type
	if mapValueMatch(id, name, mapKey, value) {
		return MapValueTypeMismatchError{
			ID:       id,
			Name:     name,
			Key:      mapKey,
			Expected: getMapValueType(id, name, mapKey),
			Actual:   reflect.TypeOf(value).Kind().String(),
		}
	}

	log.Debug().Str("id", id).Str("name", name).Str("mapKey", mapKey).Msg("updating in map component")

	return engine.Redis.HSet(context.Background(), componentId(id, name), mapKey, value).Err()
}

func setMapValueType(id string, name string, mapKey string, value interface{}) error {
	exists, err := Exists(id)

	if err != nil {
		return err
	}

	if !exists {
		return EntityExistsError{ID: id}
	}

	if !ComponentExists(id, name) {
		return MissingComponentError{ID: id, Name: name}
	}

	hasKey, err := mapHasKey(id, name, mapKey)

	if err != nil {
		return err
	}

	if hasKey {
		return MapHasKeyError{ID: id, Name: name, Key: mapKey}
	}

	log.Trace().Str("id", id).Str("name", name).Str("mapKey", mapKey).Msg("setting map component value type")

	err = engine.Redis.HSet(
		context.Background(),
		fmt.Sprintf("%s:%s", MapTypePrefix, componentId(id, name)),
		mapKey,
		reflect.TypeOf(value).Kind().String(),
	).Err()

	if err != nil {
		return err
	}

	return nil
}

func setComponentType(id string, name string, value interface{}) error {
	err := engine.
		Redis.
		Set(
			context.Background(),
			fmt.Sprintf("%s:%s", ComponentTypePrefix, componentId(id, name)),
			value,
			0).
		Err()

	if err != nil {
		return err
	}

	return nil
}

func setComponentsFromMap(id string, components map[string]interface{}) error {
	errors := make([]error, 0)

	for name, value := range components {
		kind := reflect.TypeOf(value).Kind().String()
		switch kind {
		case "string":
			err := AddStringComponent(id, name, value.(string))
			if err != nil {
				errors = append(errors, err)
			}
		case "bool":
			err := AddBoolComponent(id, name, value.(bool))
			if err != nil {
				errors = append(errors, err)
			}
		case "int":
			err := AddIntComponent(id, name, value.(int))
			if err != nil {
				errors = append(errors, err)
			}
		case "int64":
			err := AddInt64Component(id, name, value.(int64))
			if err != nil {
				errors = append(errors, err)
			}
		case "map":
			err := AddMapComponent(id, name, value.(map[string]interface{}))
			if err != nil {
				errors = append(errors, err)
			}
		case "set":
			err := AddSetComponent(id, name, value.([]interface{}))
			if err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) > 0 {
		return AddComponentErrors{Errors: errors}
	}

	return nil
}
