package system_registry

import (
	"context"
	"testing"
	"time"

	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/ecs/internal/entity_registry"
	"github.com/mjolnir-mud/engine/plugins/ecs/test"
	"github.com/stretchr/testify/assert"
)

var testSystem = test.NewTestSystem()

func setup() {
	engine.Start("test")

	entity_registry.Register(test.TestEntityType{})
	Register(testSystem)

	entity_registry.Start()

	Start()

	engine.Redis.FlushDB(context.Background())
}

func teardown() {
	engine.Redis.FlushDB(context.Background())
	Stop()
	entity_registry.Stop()
	engine.Stop()
}

func waitUntilCalledWith(system test.TestSystem, function string, args map[string]interface{}) chan bool {
	called := make(chan bool)

	go func() {
		for !system.HasBeenCalledWith(function, args) {
			time.Sleep(time.Millisecond * 10)
		}

		called <- true
	}()

	return called
}

func Test_ComponentAddedEvents(t *testing.T) {
	setup()
	defer teardown()

	called := waitUntilCalledWith(testSystem, "MatchingComponentAdded", map[string]interface{}{
		"entityId": "test",
		"key":      "testComponent",
		"value":    "test",
	})

	err := entity_registry.AddWithID("test", "test", map[string]interface{}{
		"testComponent": "test",
	})

	assert.NoError(t, err)

	<-called

	assert.True(t, testSystem.HasBeenCalledWith("MatchingComponentAdded", map[string]interface{}{
		"entityId": "test",
		"key":      "testComponent",
		"value":    "test",
	}))

	assert.True(t, testSystem.HasBeenCalledWith("ComponentAdded", map[string]interface{}{
		"entityId": "test",
		"key":      "testComponent",
		"value":    "test",
	}))
}

func Test_ComponentUpdatedEvents(t *testing.T) {
	setup()
	defer teardown()

	called := waitUntilCalledWith(testSystem, "MatchingComponentUpdated", map[string]interface{}{
		"entityId": "test",
		"key":      "testComponent",
		"oldValue": "test",
		"newValue": "test2",
	})

	err := entity_registry.AddWithID("test", "test", map[string]interface{}{
		"testComponent": "test",
	})

	assert.NoError(t, err)

	err = entity_registry.Update("test", map[string]interface{}{
		"testComponent": "test2",
	})

	assert.NoError(t, err)

	<-called

	assert.True(t, testSystem.HasBeenCalledWith("MatchingComponentUpdated", map[string]interface{}{
		"entityId": "test",
		"key":      "testComponent",
		"oldValue": "test",
		"newValue": "test2",
	}))
}

func Test_ComponentRemovedEvents(t *testing.T) {
	setup()
	defer teardown()

	called := waitUntilCalledWith(testSystem, "MatchingComponentRemoved", map[string]interface{}{
		"entityId": "test",
		"key":      "testComponent",
	})

	err := entity_registry.AddWithID("test", "test", map[string]interface{}{
		"testComponent": "test",
	})

	assert.NoError(t, err)

	err = entity_registry.RemoveComponent("test", "testComponent")

	assert.NoError(t, err)

	<-called

	assert.True(t, testSystem.HasBeenCalledWith("MatchingComponentRemoved", map[string]interface{}{
		"entityId": "test",
		"key":      "testComponent",
	}))

	// ensure no meta keys exist
	keys := engine.Redis.Keys(context.Background(), "__*:test:*").Val()

	assert.Len(t, keys, 0)
}
