package entity_registry

import (
	engineTesting "github.com/mjolnir-mud/engine/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/ecs/pkg/errors"
	"testing"

	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/ecs/test"
	"github.com/stretchr/testify/assert"
)

func setup() {
	Register(test.TestEntityType{})
	engineTesting.Setup()
	_ = engine.RedisFlushAll()
	Start()
}

func teardown() {
	_ = engine.RedisFlushAll()
	engineTesting.Teardown()
	Stop()
}

func TestAdd(t *testing.T) {
	setup()
	defer teardown()

	// test happy path
	id, err := Add("test", map[string]interface{}{})

	assert.Nil(t, err)
	assert.NotNil(t, id)

	ty, err := getEntityType(id)

	assert.Nil(t, err)
	assert.Equal(t, "test", ty)

	// test that an error is thrown if the entity type is not registered
	_, err = Add("notRegistered", map[string]interface{}{})

	assert.NotNil(t, err)
	teardown()
}

// Test that AddWithId adds an entity with the provided id to the entity registry. It takes the entity id,
// and a map of components to be added. If an entity with the same id already exists, an error will be thrown. If the
// type is not registered, an error will be thrown.
func TestAddWithID(t *testing.T) {
	setup()
	defer teardown()

	// test happy path
	err := AddWithId("test", "testId", map[string]interface{}{})

	assert.Nil(t, err)

	ty, err := getEntityType("testId")

	assert.Nil(t, err)
	assert.Equal(t, "test", ty)

	// test that an error is thrown if the entity type is not registered
	err = AddWithId("notRegistered", "testId", map[string]interface{}{})

	assert.NotNil(t, err)

	// test that an error is thrown if the id is already in use
	err = AddWithId("test", "testId", map[string]interface{}{})

	assert.NotNil(t, err)
	teardown()
}

// Test AddBoolComponent adds a boolean component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func TestAddBoolComponent(t *testing.T) {
	setup()
	defer teardown()

	err := AddWithId("test", "testEntity", map[string]interface{}{})

	assert.Nil(t, err)

	// test happy path
	err = AddBoolComponent("testEntity", "otherTestComponent", true)

	assert.Nil(t, err)

	componentValue, err := engine.RedisGet("testEntity:otherTestComponent").Bool()

	assert.Nil(t, err)
	assert.Equal(t, true, componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddBoolComponent("notRegistered", "otherTestComponent", true)

	assert.NotNil(t, err)

	// test that an error is thrown if the component already exists
	err = AddBoolComponent("test", "otherTestComponent", true)

	assert.NotNil(t, err)
	teardown()
}

func TestAddBoolToMapComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{},
	})

	assert.Nil(t, err)

	// test happy path
	err = AddBoolToMapComponent("testEntity", "testComponent", "testKey", true)

	assert.Nil(t, err)

	componentValue, err := engine.RedisHGet("testEntity:testComponent", "testKey").Bool()

	assert.Nil(t, err)
	assert.Equal(t, true, componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddBoolToMapComponent("notRegistered", "testComponent", "testKey", true)

	assert.NotNil(t, err)

	// test that an error is thrown if the component already exists
	err = AddBoolToMapComponent("test", "testComponent", "testKey", true)

	assert.NotNil(t, err)
	teardown()
}

// Test AddIntComponent adds an integer component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func TestAddIntComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{})

	assert.Nil(t, err)

	// test happy path
	err = AddIntComponent("testEntity", "otherTestComponent", 1)

	assert.Nil(t, err)

	componentValue, err := engine.RedisGet("testEntity:otherTestComponent").Int()

	assert.Nil(t, err)
	assert.Equal(t, 1, componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddIntComponent("notRegistered", "otherTestComponent", 1)

	assert.NotNil(t, err)

	// test that an error is thrown if the component already exists
	err = AddIntComponent("test", "otherTestComponent", 1)

	assert.NotNil(t, err)
	teardown()
}

func TestAddIntToMapComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{},
	})

	assert.Nil(t, err)

	// test happy path
	err = AddIntToMapComponent("testEntity", "testComponent", "testKey", 1)

	assert.Nil(t, err)

	componentValue, err := engine.RedisHGet("testEntity:testComponent", "testKey").Int()

	assert.Nil(t, err)
	assert.Equal(t, 1, componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddIntToMapComponent("notRegistered", "testComponent", "testKey", 1)

	assert.NotNil(t, err)

	// test that an error is thrown if the component does not exist
	err = AddIntToMapComponent("test", "notRegistered", "testKey", 1)

	assert.NotNil(t, err)

	// test that an error is thrown if the key already exists
	err = AddIntToMapComponent("test", "testComponent", "testKey", 1)

	assert.NotNil(t, err)
	teardown()
}

func TestAddInt64Component(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{})

	assert.Nil(t, err)

	// test happy path
	err = AddInt64Component("testEntity", "otherTestComponent", int64(1))

	assert.Nil(t, err)

	componentValue, err := engine.RedisGet("testEntity:otherTestComponent").Int64()

	assert.Nil(t, err)
	assert.Equal(t, int64(1), componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddInt64Component("notRegistered", "otherTestComponent", int64(1))

	assert.NotNil(t, err)

	// test that an error is thrown if the component already exists
	err = AddInt64Component("test", "otherTestComponent", int64(1))

	assert.NotNil(t, err)
	teardown()
}

func TestAddInt64ToMapComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{},
	})

	assert.Nil(t, err)

	// test happy path
	err = AddInt64ToMapComponent("testEntity", "testComponent", "testKey", int64(1))

	assert.Nil(t, err)

	componentValue, err := engine.RedisHGet("testEntity:testComponent", "testKey").Int64()

	assert.Nil(t, err)
	assert.Equal(t, int64(1), componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddInt64ToMapComponent("notRegistered", "testComponent", "testKey", int64(1))

	assert.NotNil(t, err)

	// test that an error is thrown if the component does not exist
	err = AddInt64ToMapComponent("test", "notRegistered", "testKey", int64(1))

	assert.NotNil(t, err)

	// test that an error is thrown if the key already exists
	err = AddInt64ToMapComponent("test", "testComponent", "testKey", int64(1))

	assert.NotNil(t, err)
	teardown()
}

func TestAddMapComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{})

	assert.Nil(t, err)

	// test happy path
	err = AddMapComponent("testEntity", "otherTestComponent", map[string]interface{}{
		"testKey": "testValue",
	})

	assert.Nil(t, err)

	componentValue, err := engine.RedisHGet("testEntity:otherTestComponent", "testKey").Result()

	assert.Nil(t, err)
	assert.Equal(t, "testValue", componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddMapComponent("notRegistered", "otherTestComponent", map[string]interface{}{
		"testKey": "testValue",
	})

	assert.NotNil(t, err)

	// test that an error is thrown if the component already exists
	err = AddMapComponent("test", "otherTestComponent", map[string]interface{}{
		"testKey": "testValue",
	})

	assert.NotNil(t, err)
	teardown()
}

func TestAddSetComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{})

	assert.Nil(t, err)

	// test happy path
	err = AddSetComponent("testEntity", "otherTestComponent", []interface{}{"testValue"})

	assert.Nil(t, err)

	componentValue, err := engine.RedisSMembers("testEntity:otherTestComponent").Result()

	assert.Nil(t, err)
	assert.Equal(t, []string{"testValue"}, componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddSetComponent("notRegistered", "otherTestComponent", []interface{}{"testValue"})

	assert.NotNil(t, err)

	// test that an error is thrown if the component already exists
	err = AddSetComponent("test", "otherTestComponent", []interface{}{"testValue"})

	assert.NotNil(t, err)
	teardown()
}

func TestAddStringComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{})

	assert.Nil(t, err)

	// test happy path
	err = AddStringComponent("testEntity", "otherTestComponent", "testValue")

	assert.Nil(t, err)

	componentValue, err := engine.RedisGet("testEntity:otherTestComponent").Result()

	assert.Nil(t, err)
	assert.Equal(t, "testValue", componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddStringComponent("notRegistered", "otherTestComponent", "testValue")

	assert.NotNil(t, err)

	// test that an error is thrown if the component already exists
	err = AddStringComponent("test", "otherTestComponent", "testValue")

	assert.NotNil(t, err)
	teardown()
}

// AddToStringSetComponent adds a string value to a set component. It takes the entity ID, component name, and the
// value to add to the set. If an entity with the same id does not exist an error will be thrown. If a component with
// the same name does not exist, an error will be thrown. If the value type is not a string, an error will be thrown.
func TestAddToStringSetComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": []interface{}{
			"testValue",
		},
	})

	assert.Nil(t, err)

	// test happy path
	err = AddToStringSetComponent("testEntity", "testComponent", "otherTestValue")

	assert.Nil(t, err)

	componentValue, err := engine.RedisSMembers("testEntity:testComponent").Result()

	assert.Nil(t, err)
	assert.Subset(t, []string{"testValue", "otherTestValue"}, componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddToStringSetComponent("notRegistered", "testComponent", "testValue")

	assert.NotNil(t, err)

	// test that an error is thrown if the component already exists
	err = AddToStringSetComponent("test", "testComponent", "testValue")

	assert.NotNil(t, err)
	teardown()
}

func TestAddOrUpdateStringComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{})

	assert.Nil(t, err)

	// test happy path
	err = AddOrUpdateStringComponent("testEntity", "testComponent", "testValue")

	assert.Nil(t, err)

	componentValue, err := engine.RedisGet("testEntity:testComponent").Result()

	assert.Nil(t, err)
	assert.Equal(t, "testValue", componentValue)
	teardown()
}

func TestAddStringToMapComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{},
	})

	assert.Nil(t, err)

	// test happy path
	err = AddStringToMapComponent("testEntity", "testComponent", "testKey", "testValue")

	assert.Nil(t, err)

	componentValue, err := engine.RedisHGet("testEntity:testComponent", "testKey").Result()

	assert.Nil(t, err)
	assert.Equal(t, "testValue", componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddStringToMapComponent("notRegistered", "testComponent", "testKey", "testValue")

	assert.NotNil(t, err)

	// test that an error is thrown if the component does not exist
	err = AddStringToMapComponent("test", "notRegistered", "testKey", "testValue")

	assert.NotNil(t, err)

	// test that an error is thrown if the key already exists
	err = AddStringToMapComponent("test", "testComponent", "testKey", "testValue")

	assert.NotNil(t, err)
	teardown()
}

func TestAddToIntSetComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": []interface{}{
			1,
		},
	})

	assert.Nil(t, err)

	// test happy path
	err = AddToIntSetComponent("testEntity", "testComponent", 2)

	assert.Nil(t, err)

	componentValue, err := engine.RedisSMembers("testEntity:testComponent").Result()

	assert.Nil(t, err)
	assert.Subset(t, []string{"1", "2"}, componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddToIntSetComponent("notRegistered", "testComponent", 2)

	assert.NotNil(t, err)

	// test that an error is thrown if the component already exists
	err = AddToIntSetComponent("test", "testComponent", 2)

	assert.NotNil(t, err)
	teardown()
}

func TestAddToInt64SetComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": []interface{}{
			int64(1),
		},
	})

	assert.Nil(t, err)

	// test happy path
	err = AddToInt64SetComponent("testEntity", "testComponent", int64(2))

	assert.Nil(t, err)

	componentValue, err := engine.RedisSMembers("testEntity:testComponent").Result()

	assert.Nil(t, err)
	assert.Subset(t, []string{"1", "2"}, componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddToInt64SetComponent("notRegistered", "testComponent", int64(2))

	assert.NotNil(t, err)

	// test that an error is thrown if the component already exists
	err = AddToInt64SetComponent("test", "testComponent", int64(2))

	assert.NotNil(t, err)
	teardown()
}

func TestCreate(t *testing.T) {
	setup()

	// test happy path
	entity, err := Create("test", map[string]interface{}{})

	assert.Nil(t, err)
	assert.Equal(t, "test", entity["testComponent"])

	// test that an error is thrown if the entity type does not exist
	_, err = Create("notRegistered", map[string]interface{}{
		"testComponent": "testValue",
	})

	assert.NotNil(t, err)
	teardown()
}

// CreateAndAdd creates an entity of the given entity type, adds it to the entity registry, and returns the
// id of the entity. It takes the entity type and a map of components. It will merge the provided components with the
// default components for the entity type returning the merged components as a map.
func TestCreateAndAdd(t *testing.T) {
	setup()

	// test happy path
	id, err := CreateAndAdd("test", map[string]interface{}{})

	assert.Nil(t, err)

	v := engine.RedisGet(componentId(id, "testComponent")).Val()
	assert.Equal(t, "test", v)

	// test that an error is thrown if the entity type does not exist
	_, err = CreateAndAdd("notRegistered", map[string]interface{}{
		"testComponent": "testValue",
	})

	assert.NotNil(t, err)
	teardown()
}

func TestGetBoolComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": true,
	})

	assert.Nil(t, err)

	// test happy path
	value, err := GetBoolComponent("testEntity", "testComponent")

	assert.Nil(t, err)
	assert.Equal(t, true, value)

	// test that an error is thrown if the entity does not exist
	_, err = GetBoolComponent("notRegistered", "testComponent")

	assert.NotNil(t, err)

	// test that an error is thrown if the component does not exist
	_, err = GetBoolComponent("test", "notRegistered")

	assert.NotNil(t, err)
	teardown()
}

func TestGetIntComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": 1,
	})

	assert.Nil(t, err)

	// test happy path
	value, err := GetIntComponent("testEntity", "testComponent")

	assert.Nil(t, err)
	assert.Equal(t, 1, value)

	// test that an error is thrown if the entity does not exist
	_, err = GetIntComponent("notRegistered", "testComponent")

	assert.NotNil(t, err)

	// test that an error is thrown if the component does not exist
	_, err = GetIntComponent("test", "notRegistered")

	assert.NotNil(t, err)
	teardown()
}

func TestGetInt64FromMapComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{
			"testKey": int64(1),
		},
	})

	assert.Nil(t, err)

	// test happy path
	value, err := GetInt64FromMapComponent("testEntity", "testComponent", "testKey")

	assert.Nil(t, err)
	assert.Equal(t, int64(1), value)

	// test that an error is thrown if the entity does not exist
	_, err = GetInt64FromMapComponent("notRegistered", "testComponent", "testKey")

	assert.NotNil(t, err)

	// test that an error is thrown if the component does not exist
	_, err = GetInt64FromMapComponent("test", "notRegistered", "testKey")

	assert.NotNil(t, err)

	// test that an error is thrown if the key does not exist
	_, err = GetInt64FromMapComponent("test", "testComponent", "notRegistered")

	assert.NotNil(t, err)

	// test get when entity does not exist
	_, err = GetInt64FromMapComponent("notRegistered", "testComponent", "testKey")

	assert.NotNil(t, err)
	assert.IsTypef(t, errors.EntityNotFoundError{}, err, "error should be of type EntityNotFoundError")

	teardown()
}

func TestGetIntFromMapComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{
			"testKey": 1,
		},
	})

	assert.Nil(t, err)

	// test happy path
	value, err := GetIntFromMapComponent("testEntity", "testComponent", "testKey")

	assert.Nil(t, err)
	assert.Equal(t, 1, value)

	// test that an error is thrown if the entity does not exist
	_, err = GetIntFromMapComponent("notRegistered", "testComponent", "testKey")

	assert.NotNil(t, err)

	// test that an error is thrown if the component does not exist
	_, err = GetIntFromMapComponent("test", "notRegistered", "testKey")

	assert.NotNil(t, err)

	// test that an error is thrown if the key does not exist
	_, err = GetIntFromMapComponent("test", "testComponent", "notRegistered")

	assert.NotNil(t, err)
	teardown()
}

func TestGetMapComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{
			"testKey": "testValue",
		},
	})

	assert.Nil(t, err)

	// test happy path
	value, err := GetMapComponent("testEntity", "testComponent")

	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"testKey": "testValue",
	}, value)

	// test that an error is thrown if the entity does not exist
	_, err = GetMapComponent("notRegistered", "testComponent")

	assert.NotNil(t, err)

	// test that an error is thrown if the component does not exist
	_, err = GetMapComponent("test", "notRegistered")

	assert.NotNil(t, err)
	teardown()
}

func TestGetStringComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": "testValue",
	})

	assert.Nil(t, err)

	// test happy path
	value, err := GetStringComponent("testEntity", "testComponent")

	assert.Nil(t, err)
	assert.Equal(t, "testValue", value)

	// test that an error is thrown if the entity does not exist
	_, err = GetStringComponent("notRegistered", "testComponent")

	assert.NotNil(t, err)

	// test that an error is thrown if the component does not exist
	_, err = GetStringComponent("test", "notRegistered")

	assert.NotNil(t, err)
	teardown()
}

func TestGetStringFromMapComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{
			"testKey": "testValue",
		},
	})

	assert.Nil(t, err)

	// test happy path
	value, err := GetStringFromMapComponent("testEntity", "testComponent", "testKey")

	assert.Nil(t, err)
	assert.Equal(t, "testValue", value)

	// test that an error is thrown if the entity does not exist
	_, err = GetStringFromMapComponent("notRegistered", "testComponent", "testKey")

	assert.NotNil(t, err)

	// test that an error is thrown if the component does not exist
	_, err = GetStringFromMapComponent("test", "notRegistered", "testKey")

	assert.NotNil(t, err)

	// test that an error is thrown if the key does not exist
	_, err = GetStringFromMapComponent("test", "testComponent", "notRegistered")

	assert.NotNil(t, err)
	teardown()
}

func TestIsEntityTypeRegistered(t *testing.T) {
	setup()

	// test happy path
	value := IsEntityTypeRegistered("test")

	assert.True(t, value)

	// test that an error is thrown if the entity does not exist
	value = IsEntityTypeRegistered("notRegistered")

	assert.False(t, value)
	teardown()
}

// Replace removes and then replaces an entity in the entity registry. It takes the entity id, and a map of
// components. It will remove the entity with the provided id and then add the provided components to the entity. If an
// entity with the same id does not exist, it with throw an error.
func TestReplace(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{
			"testKey": "testValue",
		},
	})

	assert.Nil(t, err)

	// test happy path
	err = Replace("testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{
			"testKey": "testValue2",
		},
	})

	assert.Nil(t, err)

	// test that an error is thrown if the entity does not exist
	err = Replace("notRegistered", map[string]interface{}{
		"testComponent": map[string]interface{}{
			"testKey": "testValue2",
		},
	})

	assert.NotNil(t, err)
	teardown()
}

func TestRemove(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{
			"testKey": "testValue",
		},
	})

	assert.Nil(t, err)

	// test happy path
	err = Remove("testEntity")

	assert.Nil(t, err)

	exists, err := Exists("testEntity")

	assert.Nil(t, err)
	assert.False(t, exists)

	// test that an error is thrown if the entity does not exist
	err = Remove("notRegistered")

	assert.NotNil(t, err)
	teardown()
}

func TestRemoveComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{
			"testKey": "testValue",
		},
	})

	assert.Nil(t, err)

	// test happy path
	err = RemoveComponent("testEntity", "testComponent")

	assert.Nil(t, err)

	e, err := engine.RedisExists("testEntity:testComponent").Result()

	assert.Nil(t, err)
	assert.Equal(t, int64(0), e)

	// test that an error is thrown if the entity does not exist
	err = RemoveComponent("notRegistered", "testComponent")

	assert.NotNil(t, err)

	// test that an error is thrown if the component does not exist
	err = RemoveComponent("testEntity", "notRegistered")

	assert.NotNil(t, err)
	teardown()
}

func TestRemoveFromStringSetComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": []interface{}{"testValue"},
	})

	assert.Nil(t, err)

	// test happy path
	err = RemoveFromStringSetComponent("testEntity", "testComponent", "testValue")

	assert.Nil(t, err)

	exists, err := ElementInSetComponentExists("testEntity", "testComponent", "testValue")
	assert.Nil(t, err)
	assert.False(t, exists)

	// test that an error is thrown if the entity does not exist
	err = RemoveFromStringSetComponent("notRegistered", "testComponent", "testValue")

	assert.NotNil(t, err)

	// test that an error is thrown if the component does not exist
	err = RemoveFromStringSetComponent("testEntity", "notRegistered", "testValue")

	assert.NotNil(t, err)

	teardown()
}

func TestRemoveFromInt64SetComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": []interface{}{int64(1)},
	})

	assert.Nil(t, err)

	// test happy path
	err = RemoveFromInt64SetComponent("testEntity", "testComponent", int64(1))

	assert.Nil(t, err)

	exists, err := ElementInSetComponentExists("testEntity", "testComponent", int64(1))
	assert.Nil(t, err)
	assert.False(t, exists)

	// test that an error is thrown if the entity does not exist
	err = RemoveFromInt64SetComponent("notRegistered", "testComponent", int64(1))

	assert.NotNil(t, err)

	// test that an error is thrown if the component does not exist
	err = RemoveFromInt64SetComponent("testEntity", "notRegistered", int64(1))

	assert.NotNil(t, err)
	teardown()
}

func TestRemoveFromIntSetComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": []interface{}{int(1)},
	})

	assert.Nil(t, err)

	// test happy path
	err = RemoveFromIntSetComponent("testEntity", "testComponent", int(1))

	assert.Nil(t, err)

	exists, err := ElementInSetComponentExists("testEntity", "testComponent", int(1))
	assert.Nil(t, err)
	assert.False(t, exists)

	// test that an error is thrown if the entity does not exist
	err = RemoveFromIntSetComponent("notRegistered", "testComponent", int(1))

	assert.NotNil(t, err)

	// test that an error is thrown if the component does not exist
	err = RemoveFromIntSetComponent("testEntity", "notRegistered", int(1))

	assert.NotNil(t, err)
	teardown()
}

func TestUpdate(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": "testValue",
	})

	assert.Nil(t, err)

	// test happy path
	err = Update("testEntity", map[string]interface{}{
		"testComponent": "testValue2",
	})

	assert.Nil(t, err)

	// test that an error is thrown if the entity does not exist
	err = Update("notRegistered", map[string]interface{}{
		"testComponent": "testValue2",
	})

	assert.NotNil(t, err)
	teardown()
}

func TestUpdateBoolComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": true,
	})

	assert.Nil(t, err)

	// test happy path
	err = UpdateBoolComponent("testEntity", "testComponent", false)

	assert.Nil(t, err)

	// test that an error is thrown if the entity does not exist
	err = UpdateBoolComponent("notRegistered", "testComponent", false)

	assert.NotNil(t, err)
	teardown()
}

func TestUpdateInt64Component(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": int64(1),
	})

	assert.Nil(t, err)

	// test happy path
	err = UpdateInt64Component("testEntity", "testComponent", int64(2))

	assert.Nil(t, err)

	// test that an error is thrown if the entity does not exist
	err = UpdateInt64Component("notRegistered", "testComponent", int64(2))

	assert.NotNil(t, err)
	teardown()
}

func TestUpdateIntComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": int(1),
	})

	assert.Nil(t, err)

	// test happy path
	err = UpdateIntComponent("testEntity", "testComponent", int(2))

	assert.Nil(t, err)

	// test that an error is thrown if the entity does not exist
	err = UpdateIntComponent("notRegistered", "testComponent", int(2))

	assert.NotNil(t, err)
	teardown()
}

func TestUpdateBoolInMapComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{
			"testKey": true,
		},
	})

	assert.Nil(t, err)

	// test happy path
	err = UpdateBoolInMapComponent("testEntity", "testComponent", "testKey", false)

	assert.Nil(t, err)

	// test that an error is thrown if the entity does not exist
	err = UpdateBoolInMapComponent("notRegistered", "testComponent", "testKey", false)

	assert.NotNil(t, err)
	teardown()
}

func TestUpdateInt64InMapComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{
			"testKey": int64(1),
		},
	})

	assert.Nil(t, err)

	// test happy path
	err = UpdateInt64InMapComponent("testEntity", "testComponent", "testKey", int64(2))

	assert.Nil(t, err)

	// test that an error is thrown if the entity does not exist
	err = UpdateInt64InMapComponent("notRegistered", "testComponent", "testKey", int64(2))

	assert.NotNil(t, err)
	teardown()
}

func TestUpdateIntInMapComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{
			"testKey": int(1),
		},
	})

	assert.Nil(t, err)

	// test happy path
	err = UpdateIntInMapComponent("testEntity", "testComponent", "testKey", int(2))

	assert.Nil(t, err)

	// test that an error is thrown if the entity does not exist
	err = UpdateIntInMapComponent("notRegistered", "testComponent", "testKey", int(2))

	assert.NotNil(t, err)
	teardown()
}

func TestUpdateStringInMapComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{
			"testKey": "testValue",
		},
	})

	assert.Nil(t, err)

	// test happy path
	err = UpdateStringInMapComponent("testEntity", "testComponent", "testKey", "testValue2")

	assert.Nil(t, err)

	// test that an error is thrown if the entity does not exist
	err = UpdateStringInMapComponent("notRegistered", "testComponent", "testKey", "testValue2")

	assert.NotNil(t, err)
	teardown()
}

// AddOrUpdateStringInMapComponent adds or updates a string component to a map component. It takes the entity ID,
// component name, the key to which to add the value, and the value to add to the map. If an entity with the same id
// does not exist an error will be thrown. If a component with the same name does not exist, an error will be thrown.
// If the key already exists, the value will be updated. Once a value is added to the map, the type of that key is
// enforced. Attempting to change the type of key will result in an error in later updated.
func TestAddOrUpdateStringInMapComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{
			"testKey": "testValue",
		},
	})

	assert.Nil(t, err)

	// test happy path for adding a new key
	err = AddOrUpdateStringInMapComponent("testEntity", "testComponent", "testKey", "testValue2")

	assert.Nil(t, err)

	// test happy path for updating an existing key
	err = AddOrUpdateStringInMapComponent("testEntity", "testComponent", "testKey", "testValue3")

	assert.Nil(t, err)

	// test that an error is thrown if the entity does not exist
	err = AddOrUpdateStringInMapComponent("notRegistered", "testComponent", "testKey", "testValue3")

	assert.NotNil(t, err)
	assert.IsTypef(t, errors.EntityNotFoundError{}, err, "error should be of type EntityNotFoundError")

	// test that an error is thrown if the component does not exist
	err = AddOrUpdateStringInMapComponent("testEntity", "notRegistered", "testKey", "testValue3")

	assert.NotNil(t, err)
	assert.IsTypef(t, errors.ComponentNotFoundError{}, err, "error should be of type ComponentNotFoundError")

	teardown()
}

func TestAddOrUpdateIntInMapComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{
			"testKey": 1,
		},
	})

	assert.Nil(t, err)

	// test happy path for adding a new key
	err = AddOrUpdateIntInMapComponent("testEntity", "testComponent", "testKey", 2)

	assert.Nil(t, err)

	// test happy path for updating an existing key
	err = AddOrUpdateIntInMapComponent("testEntity", "testComponent", "testKey", 3)

	assert.Nil(t, err)
	teardown()
}

func TestUpdateStringComponent(t *testing.T) {
	setup()

	err := AddWithId("test", "testEntity", map[string]interface{}{
		"testComponent": "testValue",
	})

	assert.Nil(t, err)

	// test happy path
	err = UpdateStringComponent("testEntity", "testComponent", "testValue2")

	assert.Nil(t, err)

	// test that an error is thrown if the entity does not exist
	err = UpdateStringComponent("notRegistered", "testComponent", "testValue2")

	assert.NotNil(t, err)
	teardown()
}
