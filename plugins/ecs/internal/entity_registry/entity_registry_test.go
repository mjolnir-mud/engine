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
