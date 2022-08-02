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
	Start()
}

func teardown() {
	engine.Redis.FlushDB(context.Background())
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
