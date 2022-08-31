package system_registry

import (
	"github.com/mjolnir-mud/engine"
	testing2 "github.com/mjolnir-mud/engine/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/ecs/internal/entity_registry"
	testing3 "github.com/mjolnir-mud/engine/plugins/ecs/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/ecs/testing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() {
	testing2.Setup()

	entity_registry.Register(testing3.TestEntityType{})
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

	ts := testing3.NewTestSystem()

	Start()
	Register(ts)

	err := entity_registry.AddWithId("testing", "testing", map[string]interface{}{
		"testComponent": "testing",
	})

	assert.NoError(t, err)

	call := <-ts.ComponentAddedCalled

	assert.Equal(t, "testing", call.EntityId)
	assert.Equal(t, "testComponent", call.Key)
	assert.Equal(t, "testing", call.Value)
}

func Test_ComponentUpdatedEvents(t *testing.T) {
	setup()
	defer teardown()

	ts := testing3.NewTestSystem()

	Start()
	Register(ts)

	err := entity_registry.AddWithId("testing", "testing", map[string]interface{}{
		"testComponent": "testing",
	})

	assert.NoError(t, err)

	err = entity_registry.UpdateStringComponent("testing", "testComponent", "test2")

	assert.NoError(t, err)

	call := <-ts.ComponentUpdatedCalled

	assert.Equal(t, "testing", call.EntityId)
	assert.Equal(t, "testComponent", call.Key)
	assert.Equal(t, "testing", call.OldValue)
	assert.Equal(t, "test2", call.NewValue)
}

func Test_ComponentRemovedEvents(t *testing.T) {
	setup()
	defer teardown()

	ts := testing3.NewTestSystem()

	Start()
	Register(ts)

	err := entity_registry.AddWithId("testing", "testing", map[string]interface{}{
		"testComponent": "testing",
	})

	assert.NoError(t, err)

	err = entity_registry.RemoveComponent("testing", "testComponent")

	assert.NoError(t, err)

	call := <-ts.ComponentRemovedCalled

	assert.Equal(t, "testing", call.EntityId)
	assert.Equal(t, "testComponent", call.Key)
}

func TestMatchingComponentAddedEvent(t *testing.T) {
	setup()
	defer teardown()

	ts := testing3.NewTestSystem()

	Start()
	Register(ts)

	err := entity_registry.AddWithId("testing", "testing", map[string]interface{}{
		"testComponent": "testing",
	})

	assert.NoError(t, err)

	call := <-ts.ComponentAddedCalled

	assert.Equal(t, "testing", call.EntityId)
	assert.Equal(t, "testComponent", call.Key)
	assert.Equal(t, "testing", call.Value)
}

func TestMatchingComponentUpdatedEvent(t *testing.T) {
	setup()
	defer teardown()

	ts := testing3.NewTestSystem()

	Start()
	Register(ts)

	err := entity_registry.AddWithId("testing", "testing", map[string]interface{}{
		"testComponent": "testing",
	})

	assert.NoError(t, err)
	err = entity_registry.UpdateStringComponent("testing", "testComponent", "test2")

	assert.NoError(t, err)

	call := <-ts.ComponentUpdatedCalled

	assert.Equal(t, "testing", call.EntityId)
	assert.Equal(t, "testComponent", call.Key)
	assert.Equal(t, "testing", call.OldValue)
	assert.Equal(t, "test2", call.NewValue)
}

func TestMatchingComponentRemovedEvent(t *testing.T) {
	setup()
	defer teardown()

	ts := testing3.NewTestSystem()

	Start()
	Register(ts)

	err := entity_registry.AddWithId("testing", "testing", map[string]interface{}{
		"testComponent": "testing",
	})

	assert.NoError(t, err)
	err = entity_registry.RemoveComponent("testing", "testComponent")

	assert.NoError(t, err)

	call := <-ts.ComponentRemovedCalled

	assert.Equal(t, "testing", call.EntityId)
	assert.Equal(t, "testComponent", call.Key)
}
