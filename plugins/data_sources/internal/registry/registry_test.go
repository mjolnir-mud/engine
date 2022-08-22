package registry

import (
	"testing"

	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/data_sources/test"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/stretchr/testify/assert"
)

func setup() {
	engine.RegisterPlugin(ecs.Plugin)
	ecs.RegisterEntityType(test.FakeEntityType)
	Register(test.FakeDataSource())
	_ = Start()
	engine.Start("test")
}

func teardown() {
	_ = Stop()
	engine.Stop()
}

func TestRegister(t *testing.T) {
	setup()
	defer teardown()
	if _, ok := r.dataSources["fake"]; !ok {
		t.Errorf("Expected registry.dataSources to contain fake data source")
	}
}

func TestFindOne(t *testing.T) {
	setup()
	defer teardown()

	id, entity, err := FindOne("fake", map[string]interface{}{
		"id": "test1",
	})

	assert.Nil(t, err)
	assert.Equal(t, entity, map[string]interface{}{"testComponent": "test1"})
	assert.Equal(t, id, "entity_1")
}

func TestAll(t *testing.T) {
	setup()
	defer teardown()
	entities, err := All("fake")

	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"testComponent": "test1"}, entities["test1"])
}

func TestDelete(t *testing.T) {
	setup()
	defer teardown()
	err := Delete("fake", "test1")

	assert.Nil(t, err)

	_, _, err = FindOne("fake", map[string]interface{}{
		"id": "test1",
	})

	assert.NotNil(t, err)
}
