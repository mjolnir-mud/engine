package registry

import (
	engineTesting "github.com/mjolnir-mud/engine/pkg/testing"
	"testing"

	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/data_sources/test"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/stretchr/testify/assert"
)

func setup() {
	engineTesting.Setup("world", func() {
		engine.RegisterPlugin(ecs.Plugin)

		engine.RegisterBeforeServiceStartCallback("world", func() {
			Start()
			Register(test.FakeDataSource())
		})

		engine.RegisterAfterServiceStartCallback("world", func() {
			ecs.RegisterEntityType(test.FakeEntityType)
		})
	})
}

func teardown() {
	Stop()
	engineTesting.Teardown()
}

func TestRegister(t *testing.T) {
	setup()
	defer teardown()

	_, ok := dataSources["fake"]

	assert.True(t, ok)
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
