package entity_registry

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/world/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/entity_type"
	"gopkg.in/yaml.v3"
)

const TYPE_COMPONENT_NAME = "type"

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

	for _, dir := range registry.loadableDirectories {
		LoadDirectory(dir)
	}
}

// Stop stops the entity registry.
func Stop() {
	log.Info().Msg("stopping entity registry")
}

// Add adds an entity to the registry, as well as the game world. The `entityType` represents the type of entity to
// add, and the `args are the arguments to pass to the entity's constructor.
func Add(entityType string, id string, args map[string]interface{}) error {
	if t, ok := registry.types[entityType]; ok {
		// check to see if an entity of this type already exists
		if EntityExists(id) {
			// get all the components for this entity, if they exist and a value for the same component is not
			// already set in the args, set it to the value in the registry
			components := AllComponents(id)

			for component, value := range components {
				if _, ok := args[component]; !ok {
					// get the actual component name
					componentName := strings.Replace(component, fmt.Sprintf("%s:", id), "", 1)
					args[componentName] = value
				}
			}
		}

		log.Debug().Str("entityType", entityType).Str("id", id).Msg("creating entity")
		ent := t.Create(id, args)

		// ensure the type is set
		ent[TYPE_COMPONENT_NAME] = entityType

		// add the entities components to the world
		for name, value := range ent {
			kind := reflect.TypeOf(value).Kind().String()
			switch kind {
			case "string":
				SetStringComponent(id, name, value.(string))
			case "bool":
				SetBoolComponent(id, name, value.(bool))
			case "int":
				SetIntComponent(id, name, value.(int))
			case "int64":
				SetInt64Component(id, name, value.(int64))
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

		return nil
	}

	return fmt.Errorf("entity type %s not found", entityType)
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
	keys, err := engine.Redis.Keys(context.Background(), fmt.Sprintf("*:%s", TYPE_COMPONENT_NAME)).Result()

	if err != nil {
		panic(err)
	}

	// create a new slice which holds only entities whose type matches the given type
	entities := make([]string, 0)
	for _, key := range keys {
		if strings.Contains(key, entityType) {
			entities = append(entities, strings.Replace(key, fmt.Sprintf(":%s", TYPE_COMPONENT_NAME), "", 1))
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

// LoadDirectory loads entities from a directory.
func LoadDirectory(dir string) {
	// load all yml files in the directory
	files, err := os.Open(dir)

	if err != nil {
		log.Error().Err(err).Msgf("failed to open directory %s", dir)
		return
	}

	defer func() {
		_ = files.Close()
	}()

	names, err := files.Readdirnames(0)

	if err != nil {
		log.Error().Err(err).Msgf("failed to read directory %s", dir)
		return
	}

	for _, name := range names {
		if strings.HasSuffix(name, ".yml") {
			loadFromFile(fmt.Sprintf("%s/%s", dir, name))
		}
	}
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

// Remove removes the entity's components from the world and removes the entity from the registry.
func Remove(id string) {
	log.Debug().Str("id", id).Msg("removing entity")
	// grab all the entities components
	components := AllComponents(id)

	for name := range components {
		RemoveComponent(id, name)
	}

	engine.Redis.SRem(context.Background(), "__entities", id)
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

// SetBoolComponent sets the value of the bool component.
func SetBoolComponent(id string, name string, value bool) {
	if value {
		setComponent(id, name, 1)
	} else {
		setComponent(id, name, 0)
	}
}

// SetInt64Component sets the value of the int64 component.
func SetInt64Component(id string, name string, value int64) {
	setComponent(id, name, value)
}

// SetIntComponent sets the value of the int component.
func SetIntComponent(id string, name string, value int) {
	setComponent(id, name, value)
}

// SetStringComponent  adds a string component to the entity.
func SetStringComponent(id string, name string, value string) {
	setComponent(id, name, value)
}

// SetInHashComponent adds a value to a hash component.
func SetInHashComponent(id string, name string, hashKey string, value interface{}) {
	log.Debug().Str("id", id).Str("name", name).Str("hashKey", hashKey).Msg("setting in hash component")

	// add the component to the hash
	engine.Redis.HSet(context.Background(), componentId(id, name), hashKey, value)

	// make sure we set the type
	engine.Redis.Set(context.Background(), fmt.Sprintf("__type:%s", componentId(id, name)), "map", 0)

	log.Debug().Str("id", id).Str("name", name).Str("hashKey", hashKey).
		Msg("component already exists, calling ComponentAdded")
}

// SetStringInHashComponent adds a string component to the entity.
func SetStringInHashComponent(id string, name string, hashKey string, value string) {
	SetInHashComponent(id, name, hashKey, value)
}

// SetIntInHashComponent SetInHashComponent adds an int value to a hash component.
func SetIntInHashComponent(id string, name string, hashKey string, value int) {
	SetInHashComponent(id, name, hashKey, value)
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

func getFromHashComponent(id string, name string, hashKey string) (interface{}, error) {
	return engine.Redis.HGet(context.Background(), componentId(id, name), hashKey).Result()
}

func loadFromFile(file string) {
	// load from a yml file and process
	log.Debug().Str("file", file).Msg("loading entities from file")
	content, err := os.ReadFile(file)

	if err != nil {
		log.Error().Str("file", file).Err(err).Msgf("failed to read file %s", file)
		return
	}

	loadFromYaml(content)
}

func loadFromMap(id string, m map[string]interface{}) {
	// remove the entity, if we are loading from a map we want to load fresh every time
	Remove(id)
	// load from a map and process
	log.Debug().Str("id", id).Msg("loading entity from map")

	err := Add(m[TYPE_COMPONENT_NAME].(string), id, m)

	if err != nil {
		panic(err)
	}
}

func loadFromYaml(yml []byte) {
	// load yaml and process each entity using loadMap
	ymlEntries := make(map[string]interface{})

	err := yaml.Unmarshal(yml, ymlEntries)

	if err != nil {
		log.Error().Err(err).Msg("failed to parse yaml")
		return
	}

	for id, m := range ymlEntries {
		loadFromMap(id, m.(map[string]interface{}))
	}
}

func setComponent(id string, name string, value interface{}) {
	valueType := reflect.TypeOf(value).Kind().String()

	engine.Redis.Set(context.Background(), fmt.Sprintf("__type:%s", componentId(id, name)), valueType, 0)
	engine.Redis.Set(context.Background(), componentId(id, name), value, 0)
}
