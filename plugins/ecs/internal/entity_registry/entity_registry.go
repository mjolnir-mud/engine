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

const TypeComponentName = "type"
const ComponentTypePrefix = "__type"
const MapTypePrefix = "__map"

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

// EntityTypeRequired is an error type that is returned when the entity type is not specified.
type EntityTypeRequired struct {
	ID string
}

func (e EntityTypeRequired) Error() string {
	return fmt.Sprintf("component type required for component %s", e.ID)
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

// Add adds an entity to the entity registry. it takes a map of components to be added. it
// will automatically generate a unique id for the entity. If the entity args does not contain the type component
// an error will be thrown.
func Add(args map[string]interface{}) (string, error) {
	id := generateID()
	err := AddWithID(id, args)

	if err != nil {
		return "", err
	}

	return id, nil
}

// AddWithID adds an entity with the provided id to the entity registry. It takes the entity id,
// and a map of components to be added. If an entity with the same id already exists, an error will be thrown. If the
// args map does not contain the type component, an error will be thrown.
func AddWithID(id string, args map[string]interface{}) error {
	if _, ok := args[TypeComponentName]; !ok {
		return fmt.Errorf("entity type not specified")
	}

	if EntityExists(id) {
		return fmt.Errorf("entity with id %s already exists", id)
	}

	log.Debug().Str("id", id).Msg("adding entity")

	//if t, ok := registry.types[entityType]; ok {
	//// check to see if an entity of this type already exists
	//	return
	//	// get the entity type, if it doesn't match throw an error
	//	t := GetType(id)
	//
	//	if t != entityType {
	//		return fmt.Errorf("entity already exists with type %s", t)
	//	}
	//
	//	// get all the components for this entity, if they exist and a value for the same component is not
	//	// already set in the args, set it to the value in the registry
	//	components := AllComponents(id)
	//
	//	for component, value := range components {
	//		if _, ok := args[component]; !ok {
	//			// get the actual component name
	//			componentName := strings.Replace(component, fmt.Sprintf("%s:", id), "", 1)
	//			args[componentName] = value
	//		}
	//	}
	//}

	// add the entities components to the world
	setComponentsFromMap(id, args)

	return nil
}

// Create will create an entity of the given entity type, without adding it to the entity registry. it takes the
// entity type and a map of components. It will pass the component arguments to the entity type's `Create` method,
// and return the result.
func Create(entityType string, components map[string]interface{}) (map[string]interface{}, error) {
	t := getEntityType(entityType)

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

	id, err := Add(m)

	if err != nil {
		return "", err
	}

	return id, nil
}

func SetStringInSetComponent(id string, name string, value string) {
	SetInSetComponent(id, name, value)
}

// SetInSetComponent adds a value to a set.
func SetInSetComponent(id string, name string, value interface{}) {
	engine.Redis.Set(context.Background(), fmt.Sprintf("__type:%s", componentId(id, name)), "set", 0)
	engine.Redis.SAdd(context.Background(), componentId(id, name), value)
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

// AllEntitiesByType returns all entities of the given type.
func AllEntitiesByType(entityType string) []string {
	keys, err := engine.Redis.Keys(context.Background(), fmt.Sprintf("*:%s", TypeComponentName)).Result()

	if err != nil {
		panic(err)
	}

	// create a new slice which holds only entities whose type matches the given type
	entities := make([]string, 0)
	for _, key := range keys {
		if strings.Contains(key, entityType) {
			entities = append(entities, strings.Replace(key, fmt.Sprintf(":%s", TypeComponentName), "", 1))
		}
	}

	return entities
}

func AllEntitiesByTypeWithComponent(entityType string, component string) []string {
	ents := AllEntitiesByType(entityType)

	entities := make([]string, 0)
	for _, ent := range ents {
		if HasComponent(ent, component) {
			entities = append(entities, ent)
		}
	}

	return entities
}

func AllEntitiesByTypeWithComponentValue(entityType string, component string, value interface{}) []string {
	ents := AllEntitiesByType(entityType)

	entities := make([]string, 0)
	for _, ent := range ents {
		if HasComponentValue(ent, component, value) {
			entities = append(entities, ent)
		}
	}

	return entities
}

// EntityExists returns true if the entity exists.
func EntityExists(id string) bool {
	l := len(engine.Redis.Keys(context.Background(), fmt.Sprintf("__type:%s:*", id)).Val())

	return l > 0
}

// GetBoolComponent returns the value of the bool component.
func GetBoolComponent(id string, name string) (bool, error) {
	return engine.Redis.Get(context.Background(), componentId(id, name)).Bool()
}

// GetInt64Component returns the value of the int64 component.
func GetInt64Component(id string, name string) (int64, error) {
	return engine.Redis.Get(context.Background(), componentId(id, name)).Int64()
}

// GetIntComponent returns the value of the int component.
func GetIntComponent(id string, name string) (int, error) {
	return engine.Redis.Get(context.Background(), componentId(id, name)).Int()
}

func GetHashComponent(id string, name string) (map[string]interface{}, error) {
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

// GetIntFromHashComponent returns the int value of an element in a hash component.
func GetIntFromHashComponent(id string, name string, hashKey string) (int, error) {
	return engine.Redis.HGet(context.Background(), componentId(id, name), hashKey).Int()
}

// GetStringComponent returns the value of the string component.
func GetStringComponent(id string, name string) (string, error) {
	return engine.Redis.Get(context.Background(), componentId(id, name)).Result()
}

// GetStringFromHashComponent returns the string value of an element in a hash component.
func GetStringFromHashComponent(id string, name string, hashKey string) (string, error) {
	str, err := getFromHashComponent(id, name, hashKey)

	if err != nil {
		return "", err
	}

	return str.(string), nil
}

// GetStringsFromSetComponent GetStringSetComponent returns the value of the string set component.
func GetStringsFromSetComponent(id string, name string) ([]string, error) {
	return engine.Redis.SMembers(context.Background(), componentId(id, name)).Result()
}

// HasComponent returns true if the component exists.
func HasComponent(id string, name string) bool {
	exists := engine.Redis.Exists(context.Background(), componentId(id, name)).Val()

	return exists > int64(0)
}

// HasComponentValue returns true if the component exists and has the given value.
func HasComponentValue(id string, name string, value interface{}) bool {
	exists := engine.Redis.Exists(context.Background(), componentId(id, name)).Val()

	if exists == int64(0) {
		return false
	}

	return compareValues(componentId(id, name), value)
}

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

// Replace removes and then replaces an entity in the entity registry. It takes the entity type, id, and a map of
// components. It will remove the entity with the provided id and then add the provided components to the entity. If an
// entity with the same id does not exist, or the entity type does not match an error will be thrown.
func Replace(id string, components map[string]interface{}) error {
	// if the components map does not have a type component, throw an error
	if _, ok := components[TypeComponentName]; !ok {
		return fmt.Errorf("'type' component is required")
	}

	t, err := GetStringComponent(id, TypeComponentName)

	if err != nil {
		return err
	}

	if t != components[TypeComponentName].(string) {
		return fmt.Errorf(
			"entity type for the replace ment (%s) does not match %s",
			components[TypeComponentName].(string),
			t,
		)
	}

	log.Debug().Str("id", id).Msg("replacing entity")

	// remove the entity
	err = Remove(id)

	if err != nil {
		return err
	}

	// add the entity
	err = AddWithID(id, components)

	if err != nil {
		return err
	}

	return nil
}

// Register registers an entity type.
func Register(e entity_type.EntityType) {
	registry.types[e.Name()] = e
}

// RegisterLoadableDirectory registers a directory to load entities from.
func RegisterLoadableDirectory(dir string) {
	log.Debug().Str("dir", dir).Msg("registering loadable directory")
	registry.loadableDirectories = append(registry.loadableDirectories, dir)
}

// Remove removes an entity from the entity registry. It takes the entity type and id. If an entity with the same
// id does not exist an error will be thrown.
func Remove(id string) error {
	exists := EntityExists(id)

	if !exists {
		return fmt.Errorf("entity %s does not exist", id)
	}

	log.Debug().Str("id", id).Msg("removing entity")
	// grab all the entities components
	components := AllComponents(id)

	for name := range components {
		RemoveComponent(id, name)
	}

	return nil
}

// RemoveComponent removes the component from the entity.
func RemoveComponent(id string, name string) {
	log.Debug().Str("id", id).Str("name", name).Msg("removing component")

	engine.Redis.Del(context.Background(), componentId(id, name))
}

//RemoveStringFromSetComponent removes a string from a set component.
func RemoveFromSetComponent(id string, name string, value interface{}) error {
	return engine.Redis.SRem(context.Background(), componentId(id, name), value).Err()
}

// AddBoolComponentToEntity adds a boolean component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func AddBoolComponentToEntity(id string, name string, value bool) error {
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
// the type of a key will result in an error in later updated.
func AddBoolToMapComponent(id string, name string, key string, value bool) error {
	if value {
		return addToMapComponent(id, name, key, 1)
	} else {
		return addToMapComponent(id, name, key, 0)
	}
}

// AddIntComponentToEntity adds an integer component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func AddIntComponentToEntity(id string, name string, value int) error {
	return addComponent(id, name, value)
}

// AddIntToMapComponent adds an integer component to a map component. It takes the entity ID, component name, the key
// to which to add the value, and the value to add to the map. If an entity with the same id does not exist an error
// will be thrown. If a component with the same name does not exist, an error will be thrown. If the key already exists
// an error will be thrown. Once a value is added to the map, the type of that key is enforced. Attempting to change
//// the type of key will result in an error in later updated.
func AddIntToMapComponent(id string, name string, key string, value int) error {
	return addToMapComponent(id, name, key, value)
}

// AddInt64ComponentToEntity adds an integer64 component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func AddInt64ComponentToEntity(id string, name string, value int64) error {
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

// AddMapComponentToEntity adds a map component to an entity. It takes the entity ID, component name, and the value of
// the component. If an entity with the same id does not exist an error will be thrown. If a component with the same name
// already exists, an error will be thrown.
func AddMapComponentToEntity(id string, name string, value map[string]interface{}) error {
	if EntityExists(id) {
		return EntityExistsError{ID: id}
	}

	err := setType(id, name, "map")

	if err != nil {
		return err
	}

	log.Debug().Str("id", id).Str("name", name).Msg("adding map component")

	err = engine.Redis.MSet(context.Background(), componentId(id, name), value).Err()

	if err != nil {
		delType(id, name)
		return err
	}

	return nil
}

// AddStringComponentToEntity adds a string component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func AddStringComponentToEntity(id string, name string, value string) error {
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

// Update updates an entity in the entity registry. It takes the entity type, id, and a map of components. It will
// apply the provided components to the entity. If an entity with the same id does not exist, or the entity type does
// not match an error will be thrown. Any components that are not provided will be removed from the entity.
func Update(id string, components map[string]interface{}) error {
	// if the components map does not have a type component, throw an error
	if _, ok := components[TypeComponentName]; !ok {
		return EntityTypeRequired{ID: id}
	}

	exists := EntityExists(id)

	if !exists {
		return EntityNotFoundError{ID: id}
	}

	t, err := GetStringComponent(id, TypeComponentName)

	if err != nil {
		return err
	}

	if t != components[TypeComponentName].(string) {
		return fmt.Errorf(
			"entity type for the update (%s) does not match %s",
			components[TypeComponentName].(string),
			t,
		)
	}

	log.Debug().Str("id", id).Msg("updating entity")

	// grab all the entities components
	currentComponents := AllComponents(id)

	// remove the components that are not provided
	for name := range currentComponents {
		if _, ok := components[name]; !ok {
			RemoveComponent(id, name)
		}
	}

	// set the components that are provided
	for name, value := range components {
		if _, ok := components[name]; ok {
			err = updateComponent(id, name, value)

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
func UpdateOrAdd(id string, components map[string]interface{}) error {
	if EntityExists(id) {
		return Update(id, components)
	}

	return AddWithID(id, components)
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

func getEntityType(name string) entity_type.EntityType {
	return registry.types[name]
}

func getFromHashComponent(id string, name string, hashKey string) (interface{}, error) {
	return engine.Redis.HGet(context.Background(), componentId(id, name), hashKey).Result()
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
	if !EntityExists(id) {
		return fmt.Errorf("entity %s does not exist", id)
	}

	if HasComponent(id, name) {
		return fmt.Errorf("component %s already exists", name)
	}

	log.Debug().Str("id", id).Str("name", name).Msg("adding component")

	valueType := reflect.TypeOf(value).Kind().String()

	err := engine.
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

func addToMapComponent(id string, name string, hashKey string, value interface{}) error {
	if !EntityExists(id) {
		return EntityExistsError{ID: id}
	}

	if !HasComponent(id, name) {
		return MissingComponentError{ID: id, Name: name}
	}

	hasKey, err := mapHasKey(id, name, hashKey)

	if err != nil {
		return err
	}

	if hasKey {
		return MapHasKeyError{ID: id, Name: name, Key: hashKey}
	}

	log.Debug().Str("id", id).Str("name", name).Str("hashKey", hashKey).Msg("adding to map component")

	err = setMapValueType(id, name, hashKey, value)

	if err != nil {
		return err
	}

	err = engine.Redis.HSet(context.Background(), componentId(id, name), hashKey, value).Err()

	if err != nil {
		return err
	}

	return nil
}

func componentTypeMatch(id string, name string, value interface{}) bool {
	return reflect.TypeOf(value).Kind().String() == getComponentType(id, name)
}

func mapHasKey(id string, name string, hashKey string) (bool, error) {
	return engine.Redis.HExists(context.Background(), componentId(id, name), hashKey).Result()
}

func mapValueMatch(id string, name string, hashKey string, value interface{}) bool {
	return reflect.TypeOf(value).Kind().String() == getMapValueType(id, name, hashKey)
}

func updateComponent(id string, name string, value interface{}) error {
	if !EntityExists(id) {
		return EntityExistsError{ID: id}
	}

	if !HasComponent(id, name) {
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

	err := engine.Redis.Set(context.Background(), componentId(id, name), value, 0).Err()

	if err != nil {
		return err
	}

	return nil
}

func updateInMapComponent(id string, name string, hashKey string, value interface{}) error {
	// check if the entity exists
	if !EntityExists(id) {
		return EntityExistsError{ID: id}
	}

	// check if the component exists
	if !HasComponent(id, name) {
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
	hasKey, err := mapHasKey(id, name, hashKey)

	if err != nil {
		return err
	}

	if !hasKey {
		return MapHasKeyError{ID: id, Name: name, Key: hashKey}
	}

	// check if the value type matches the expected type
	if mapValueMatch(id, name, hashKey, value) {
		return MapValueTypeMismatchError{
			ID:       id,
			Name:     name,
			Key:      hashKey,
			Expected: getMapValueType(id, name, hashKey),
			Actual:   reflect.TypeOf(value).Kind().String(),
		}
	}

	log.Debug().Str("id", id).Str("name", name).Str("hashKey", hashKey).Msg("updating in map component")

	return engine.Redis.HSet(context.Background(), componentId(id, name), hashKey, value).Err()
}

func setMapValueType(id string, name string, hashKey string, value interface{}) error {
	if !EntityExists(id) {
		return EntityExistsError{ID: id}
	}

	if !HasComponent(id, name) {
		return MissingComponentError{ID: id, Name: name}
	}

	hasKey, err := mapHasKey(id, name, hashKey)

	if err != nil {
		return err
	}

	if !hasKey {
		return MapHasKeyError{ID: id, Name: name, Key: hashKey}
	}

	log.Trace().Str("id", id).Str("name", name).Str("hashKey", hashKey).Msg("setting map component value type")

	err = engine.Redis.HSet(
		context.Background(),
		fmt.Sprintf("%s:%s", MapTypePrefix, componentId(id, name)),
		hashKey,
		reflect.TypeOf(value).Kind().String(),
	).Err()

	if err != nil {
		return err
	}

	return nil
}

func setType(id string, name string, value interface{}) error {
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

func delType(id string, name string) {
	engine.Redis.Del(context.Background(), fmt.Sprintf("%s:%s", ComponentTypePrefix, componentId(id, name)))
}

func setComponentsFromMap(id string, components map[string]interface{}) {
	for name, value := range components {
		kind := reflect.TypeOf(value).Kind().String()
		switch kind {
		case "string":
			AddStringComponentToEntity(id, name, value.(string))
		case "bool":
			SetBoolComponent(id, name, value.(bool))
		case "int":
			SetIntComponent(id, name, value.(int))
		case "int64":
			AddInt64ComponentToEntity(id, name, value.(int64))
		case "map":
			for k, v := range value.(map[string]interface{}) {
				SetInHashComponent(id, name, k, v)
			}
		case "slice":
			for _, v := range value.([]interface{}) {
				SetInSetComponent(id, name, v)
			}
		}
	}
}
