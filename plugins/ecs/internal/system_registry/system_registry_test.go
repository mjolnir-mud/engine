package system_registry

import (
	"github.com/mjolnir-mud/engine"
	testing2 "github.com/mjolnir-mud/engine/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/ecs/internal/entity_registry"
	"github.com/mjolnir-mud/engine/plugins/ecs/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() {
	testing2.Setup()

	entity_registry.Register(test.TestEntityType{})
	entity_registry.Start()

	_ = engine.RedisFlushAll()
}

func teardown() {
	_ = engine.RedisFlushAll()
	Stop()
	entity_registry.Stop()
	testing2.Teardown()
}

func Test_ComponentAddedEvents(t *testing.T) {
	setup()
	defer teardown()

	ts := test.NewTestSystem()

	Start()
	Register(ts)

	err := entity_registry.AddWithId("test", "test", map[string]interface{}{
		"testComponent": "test",
	})

	assert.NoError(t, err)

	call := <-ts.ComponentAddedCalled

	assert.Equal(t, "test", call.EntityId)
	assert.Equal(t, "testComponent", call.Key)
	assert.Equal(t, "test", call.Value)
}

func Test_ComponentUpdatedEvents(t *testing.T) {
	setup()
	defer teardown()

	ts := test.NewTestSystem()

	Start()
	Register(ts)

	err := entity_registry.AddWithId("test", "test", map[string]interface{}{
		"testComponent": "test",
	})

	assert.NoError(t, err)

	err = entity_registry.UpdateStringComponent("test", "testComponent", "test2")

	assert.NoError(t, err)

	call := <-ts.ComponentUpdatedCalled

	assert.Equal(t, "test", call.EntityId)
	assert.Equal(t, "testComponent", call.Key)
	assert.Equal(t, "test", call.OldValue)
	assert.Equal(t, "test2", call.NewValue)
}

func Test_ComponentRemovedEvents(t *testing.T) {
	setup()
	defer teardown()

	ts := test.NewTestSystem()

	Start()
	Register(ts)

	err := entity_registry.AddWithId("test", "test", map[string]interface{}{
		"testComponent": "test",
	})

	assert.NoError(t, err)

	err = entity_registry.RemoveComponent("test", "testComponent")

	assert.NoError(t, err)

	call := <-ts.ComponentRemovedCalled

	assert.Equal(t, "test", call.EntityId)
	assert.Equal(t, "testComponent", call.Key)
}

func TestMatchingComponentAddedEvent(t *testing.T) {
	setup()
	defer teardown()

	ts := test.NewTestSystem()

	Start()
	Register(ts)

	err := entity_registry.AddWithId("test", "test", map[string]interface{}{
		"testComponent": "test",
	})

	assert.NoError(t, err)

	call := <-ts.ComponentAddedCalled

	assert.Equal(t, "test", call.EntityId)
	assert.Equal(t, "testComponent", call.Key)
	assert.Equal(t, "test", call.Value)
}

func TestMatchingComponentUpdatedEvent(t *testing.T) {
	setup()
	defer teardown()

	ts := test.NewTestSystem()

	Start()
	Register(ts)

	err := entity_registry.AddWithId("test", "test", map[string]interface{}{
		"testComponent": "test",
	})

	assert.NoError(t, err)
	err = entity_registry.UpdateStringComponent("test", "testComponent", "test2")

	assert.NoError(t, err)

	call := <-ts.ComponentUpdatedCalled

	assert.Equal(t, "test", call.EntityId)
	assert.Equal(t, "testComponent", call.Key)
	assert.Equal(t, "test", call.OldValue)
	assert.Equal(t, "test2", call.NewValue)
}

func TestMatchingComponentRemovedEvent(t *testing.T) {
	setup()
	defer teardown()

	ts := test.NewTestSystem()

	Start()
	Register(ts)

	err := entity_registry.AddWithId("test", "test", map[string]interface{}{
		"testComponent": "test",
	})

	assert.NoError(t, err)
	err = entity_registry.RemoveComponent("test", "testComponent")

	assert.NoError(t, err)

	call := <-ts.ComponentRemovedCalled

	assert.Equal(t, "test", call.EntityId)
	assert.Equal(t, "testComponent", call.Key)
}
