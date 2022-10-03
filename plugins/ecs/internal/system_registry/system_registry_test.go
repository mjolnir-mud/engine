package system_registry

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/ecs/internal/entity_registry"
	"github.com/mjolnir-mud/engine/plugins/ecs/pkg/testing/fakes"
	engineTesting "github.com/mjolnir-mud/engine/testing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() {
	engineTesting.RegisterSetupCallback("ecs", func() {
		engine.RegisterBeforeStartCallback(func() {
			Start()
			entity_registry.Start()
		})

		engine.RegisterAfterStartCallback(func() {
			_ = engine.RedisFlushAll()
			entity_registry.Register(fakes.FakeEntityType{})
		})
	})
	engineTesting.Setup("world")
}

func teardown() {
	Stop()
	engineTesting.Teardown()
}

func Test_ComponentAddedEvents(t *testing.T) {
	//setup()
	//defer teardown()
	//
	//ts := ecsTesting.NewFakeSystem()
	//
	//Start()
	//Register(ts)
	//
	//err := plugin("testing", "testing", map[string]data_source{}{
	//	"testComponent": "testing",
	//})
	//
	//assert.NoError(t, err)
	//
	//call := <-ts.ComponentAddedCalled
	//
	//assert.Equal(t, "testing", call.EntityId)
	//assert.Equal(t, "testComponent", call.Key)
	//assert.Equal(t, "testing", call.Value)
}

func Test_ComponentUpdatedEvents(t *testing.T) {
	setup()
	defer teardown()

	ts := fakes.NewFakeSystem()

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

	ts := fakes.NewFakeSystem()

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

	ts := fakes.NewFakeSystem()

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

	ts := fakes.NewFakeSystem()

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

	ts := fakes.NewFakeSystem()

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
