package entity_registry

import (
	"context"
	"testing"

	"github.com/mjolnir-mud/engine"
	"github.com/stretchr/testify/assert"
)

type testEntityType struct{}

func (t testEntityType) Name() string {
	return "test"
}

func (t testEntityType) Create(args map[string]interface{}) map[string]interface{} {
	args["testComponent"] = "test"

	return args
}

func setup() {
	Register(testEntityType{})
	engine.Start("test")
	engine.Redis.FlushAll(context.Background())
	Start()
}

func teardown() {
	engine.Redis.FlushAll(context.Background())
	engine.Stop()
	Stop()
}

// Test that Add adds an entity to the entity registry. it takes a map of components to be added. It
// will automatically generate a unique id for the entity. If the passed entity type is not a valid registered type an
// error will be thrown.

func TestAdd(t *testing.T) {
	setup()

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

// Test that AddWithID adds an entity with the provided id to the entity registry. It takes the entity id,
// and a map of components to be added. If an entity with the same id already exists, an error will be thrown. If the
// type is not registered, an error will be thrown.
func TestAddWithID(t *testing.T) {
	setup()

	// test happy path
	err := AddWithID("test", "testId", map[string]interface{}{})

	assert.Nil(t, err)

	ty, err := getEntityType("testId")

	assert.Nil(t, err)
	assert.Equal(t, "test", ty)

	// test that an error is thrown if the entity type is not registered
	err = AddWithID("notRegistered", "testId", map[string]interface{}{})

	assert.NotNil(t, err)

	// test that an error is thrown if the id is already in use
	err = AddWithID("test", "testId", map[string]interface{}{})

	assert.NotNil(t, err)
	teardown()
}

// Test AddBoolComponent adds a boolean component to an entity. It takes the entity ID, component name, and the
// value of the component. If an entity with the same id does not exist an error will be thrown. If a component with the
// same name already exists, an error will be thrown.
func TestAddBoolComponent(t *testing.T) {
	setup()

	err := AddWithID("test", "testEntity", map[string]interface{}{})

	assert.Nil(t, err)

	// test happy path
	err = AddBoolComponent("testEntity", "testComponent", true)

	assert.Nil(t, err)

	componentValue, err := engine.Redis.Get(context.Background(), "testEntity:testComponent").Bool()

	assert.Nil(t, err)
	assert.Equal(t, true, componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddBoolComponent("notRegistered", "testComponent", true)

	assert.NotNil(t, err)

	// test that an error is thrown if the component already exists
	err = AddBoolComponent("test", "testComponent", true)

	assert.NotNil(t, err)
	teardown()
}

func TestAddBoolToMapComponent(t *testing.T) {
	setup()

	err := AddWithID("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{},
	})

	assert.Nil(t, err)

	// test happy path
	err = AddBoolToMapComponent("testEntity", "testComponent", "testKey", true)

	assert.Nil(t, err)

	componentValue, err := engine.Redis.HGet(context.Background(), "testEntity:testComponent", "testKey").Bool()

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

	err := AddWithID("test", "testEntity", map[string]interface{}{})

	assert.Nil(t, err)

	// test happy path
	err = AddIntComponent("testEntity", "testComponent", 1)

	assert.Nil(t, err)

	componentValue, err := engine.Redis.Get(context.Background(), "testEntity:testComponent").Int()

	assert.Nil(t, err)
	assert.Equal(t, 1, componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddIntComponent("notRegistered", "testComponent", 1)

	assert.NotNil(t, err)

	// test that an error is thrown if the component already exists
	err = AddIntComponent("test", "testComponent", 1)

	assert.NotNil(t, err)
	teardown()
}

func TestAddIntToMapComponent(t *testing.T) {
	setup()

	err := AddWithID("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{},
	})

	assert.Nil(t, err)

	// test happy path
	err = AddIntToMapComponent("testEntity", "testComponent", "testKey", 1)

	assert.Nil(t, err)

	componentValue, err := engine.Redis.HGet(context.Background(), "testEntity:testComponent", "testKey").Int()

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

	err := AddWithID("test", "testEntity", map[string]interface{}{})

	assert.Nil(t, err)

	// test happy path
	err = AddInt64Component("testEntity", "testComponent", int64(1))

	assert.Nil(t, err)

	componentValue, err := engine.Redis.Get(context.Background(), "testEntity:testComponent").Int64()

	assert.Nil(t, err)
	assert.Equal(t, int64(1), componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddInt64Component("notRegistered", "testComponent", int64(1))

	assert.NotNil(t, err)

	// test that an error is thrown if the component already exists
	err = AddInt64Component("test", "testComponent", int64(1))

	assert.NotNil(t, err)
	teardown()
}

func TestAddInt64ToMapComponent(t *testing.T) {
	setup()

	err := AddWithID("test", "testEntity", map[string]interface{}{
		"testComponent": map[string]interface{}{},
	})

	assert.Nil(t, err)

	// test happy path
	err = AddInt64ToMapComponent("testEntity", "testComponent", "testKey", int64(1))

	assert.Nil(t, err)

	componentValue, err := engine.Redis.HGet(context.Background(), "testEntity:testComponent", "testKey").Int64()

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

	err := AddWithID("test", "testEntity", map[string]interface{}{})

	assert.Nil(t, err)

	// test happy path
	err = AddMapComponent("testEntity", "testComponent", map[string]interface{}{
		"testKey": "testValue",
	})

	assert.Nil(t, err)

	componentValue, err := engine.Redis.HGet(context.Background(), "testEntity:testComponent", "testKey").Result()

	assert.Nil(t, err)
	assert.Equal(t, "testValue", componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddMapComponent("notRegistered", "testComponent", map[string]interface{}{
		"testKey": "testValue",
	})

	assert.NotNil(t, err)

	// test that an error is thrown if the component already exists
	err = AddMapComponent("test", "testComponent", map[string]interface{}{
		"testKey": "testValue",
	})

	assert.NotNil(t, err)
	teardown()
}

func TestAddSetComponent(t *testing.T) {
	setup()

	err := AddWithID("test", "testEntity", map[string]interface{}{})

	assert.Nil(t, err)

	// test happy path
	err = AddSetComponent("testEntity", "testComponent", []interface{}{"testValue"})

	assert.Nil(t, err)

	componentValue, err := engine.Redis.SMembers(context.Background(), "testEntity:testComponent").Result()

	assert.Nil(t, err)
	assert.Equal(t, []string{"testValue"}, componentValue)

	// test that an error is thrown if the entity does not exist
	err = AddSetComponent("notRegistered", "testComponent", []interface{}{"testValue"})

	assert.NotNil(t, err)

	// test that an error is thrown if the component already exists
	err = AddSetComponent("test", "testComponent", []interface{}{"testValue"})

	assert.NotNil(t, err)
	teardown()
}
