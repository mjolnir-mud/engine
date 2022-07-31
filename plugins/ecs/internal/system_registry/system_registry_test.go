package system_registry

import (
	"context"
	"testing"

	"github.com/mjolnir-mud/engine"
	"github.com/stretchr/testify/assert"
)

type testSystemA struct {
	ComponentAddedCalled           map[string]map[string]interface{}
	MatchingComponentAddedCalled   map[string]map[string]interface{}
	ComponentUpdatedCalled         map[string]map[string]interface{}
	MatchingComponentUpdatedCalled map[string]map[string]interface{}
	ComponentRemovedCalled         map[string]map[string]interface{}
	MatchingComponentRemovedCalled map[string]map[string]interface{}
	Called                         chan bool
}

func (e testSystemA) Name() string {
	return "testSystemA"
}

func (e testSystemA) Component() string {
	return "testComponentA"
}

func (e testSystemA) Match(_ string, _ interface{}) bool {
	return true
}

func (e testSystemA) WorldStarted() {}

func (e testSystemA) ComponentAdded(entityId string, key string, value interface{}) error {
	e.ComponentAddedCalled[key] = map[string]interface{}{
		"entityId": entityId,
		"value":    value,
	}

	return nil
}

func (e testSystemA) ComponentUpdated(entityId string, key string, oldValue interface{}, newValue interface{}) error {
	e.ComponentUpdatedCalled[key] = map[string]interface{}{
		"entityId": entityId,
		"oldValue": oldValue,
		"newValue": newValue,
	}

	return nil
}

func (e testSystemA) ComponentRemoved(entityId string, key string, value interface{}) error {
	e.ComponentRemovedCalled[key] = map[string]interface{}{
		"entityId": entityId,
		"value":    value,
	}

	return nil
}

func (e testSystemA) MatchingComponentAdded(entityId string, key string, value interface{}) error {
	e.MatchingComponentAddedCalled[key] = map[string]interface{}{
		"entityId": entityId,
		"value":    value,
	}

	e.Called <- true

	return nil
}

func (e testSystemA) MatchingComponentUpdated(entityId string, key string, oldValue interface{}, newValue interface{}) error {
	e.MatchingComponentUpdatedCalled[key] = map[string]interface{}{
		"entityId": entityId,
		"oldValue": oldValue,
		"newValue": newValue,
	}

	e.Called <- true

	return nil
}

func (e testSystemA) MatchingComponentRemoved(entityId string, key string, value interface{}) error {
	e.MatchingComponentRemovedCalled[key] = map[string]interface{}{
		"entityId": entityId,
		"value":    value,
	}

	e.Called <- true

	return nil
}

func (e testSystemA) BeforeEntityRemoved(entityId string, _ map[string]interface{}) {}

type testSystemB struct {
	ComponentAddedCalled           map[string]map[string]interface{}
	MatchingComponentAddedCalled   map[string]map[string]interface{}
	ComponentUpdatedCalled         map[string]map[string]interface{}
	MatchingComponentUpdatedCalled map[string]map[string]interface{}
	ComponentRemovedCalled         map[string]map[string]interface{}
	MatchingComponentRemovedCalled map[string]map[string]interface{}
	Called                         chan bool
}

func (e testSystemB) Name() string {
	return "testSystemB"
}

func (e testSystemB) Component() string {
	return "testComponentB"
}

func (e testSystemB) Match(_ string, _ interface{}) bool {
	return true
}

func (e testSystemB) WorldStarted() {}

func (e testSystemB) ComponentAdded(entityId string, key string, value interface{}) error {
	e.ComponentAddedCalled[key] = map[string]interface{}{
		"entityId": entityId,
		"value":    value,
	}

	return nil
}

func (e testSystemB) ComponentUpdated(entityId string, key string, oldValue interface{}, newValue interface{}) error {
	e.ComponentUpdatedCalled[key] = map[string]interface{}{
		"entityId": entityId,
		"oldValue": oldValue,
		"newValue": newValue,
	}
	return nil
}

func (e testSystemB) ComponentRemoved(entityId string, key string, value interface{}) error {
	e.ComponentRemovedCalled[key] = map[string]interface{}{
		"entityId": entityId,
		"value":    value,
	}

	return nil
}

func (e testSystemB) MatchingComponentAdded(entityId string, key string, value interface{}) error {
	e.MatchingComponentAddedCalled[key] = map[string]interface{}{
		"entityId": entityId,
		"value":    value,
	}

	e.Called <- true

	return nil
}

func (e testSystemB) MatchingComponentUpdated(entityId string, key string, oldValue interface{}, newValue interface{}) error {
	e.MatchingComponentUpdatedCalled[key] = map[string]interface{}{
		"entityId": entityId,
		"oldValue": oldValue,
		"newValue": newValue,
	}

	e.Called <- true

	return nil
}

func (e testSystemB) MatchingComponentRemoved(entityId string, key string, value interface{}) error {
	e.MatchingComponentRemovedCalled[key] = map[string]interface{}{
		"entityId": entityId,
		"value":    value,
	}

	e.Called <- true

	return nil
}

func (e testSystemB) BeforeEntityRemoved(entityId string, _ map[string]interface{}) {}

func setup() (testSystemA, testSystemB) {
	redis.Start()
	engine.Redis.FlushDB(context.Background())

	tsa := testSystemA{
		ComponentAddedCalled:           make(map[string]map[string]interface{}),
		MatchingComponentAddedCalled:   make(map[string]map[string]interface{}),
		ComponentUpdatedCalled:         make(map[string]map[string]interface{}),
		MatchingComponentUpdatedCalled: make(map[string]map[string]interface{}),
		ComponentRemovedCalled:         make(map[string]map[string]interface{}),
		MatchingComponentRemovedCalled: make(map[string]map[string]interface{}),
		Called:                         make(chan bool),
	}

	tsb := testSystemB{
		ComponentAddedCalled:           make(map[string]map[string]interface{}),
		MatchingComponentAddedCalled:   make(map[string]map[string]interface{}),
		ComponentUpdatedCalled:         make(map[string]map[string]interface{}),
		MatchingComponentUpdatedCalled: make(map[string]map[string]interface{}),
		ComponentRemovedCalled:         make(map[string]map[string]interface{}),
		MatchingComponentRemovedCalled: make(map[string]map[string]interface{}),
		Called:                         make(chan bool),
	}

	Register(tsa)
	Register(tsb)

	Start()

	return tsa, tsb
}

func teardown() {
	Stop()
	engine.Redis.FlushDB(context.Background())
	redis.Stop()
}

// Test that the system registry handles regular value set callbacks correctly.
func TestRegister_Set(t *testing.T) {
	tsa, tsb := setup()

	engine.Redis.Set(context.Background(), "testEntity:testComponentA", "testValueA", 0)

	<-tsa.Called

	engine.Redis.Set(context.Background(), "testEntity:testComponentB", "testValueB", 0)

	<-tsb.Called

	// it calls ComponentAdded on all systems when a component is added
	assert.Equal(t, tsa.ComponentAddedCalled["testEntity:testComponentA"]["value"].(string), "testValueA")
	assert.Equal(t, tsa.ComponentAddedCalled["testEntity:testComponentB"]["value"].(string), "testValueB")

	assert.Equal(t, tsb.ComponentAddedCalled["testEntity:testComponentA"]["value"].(string), "testValueA")
	assert.Equal(t, tsb.ComponentAddedCalled["testEntity:testComponentB"]["value"].(string), "testValueB")

	assert.Equal(t, tsa.ComponentAddedCalled["testEntity:testComponentA"]["entityId"], "testEntity")
	assert.Equal(t, tsb.ComponentAddedCalled["testEntity:testComponentB"]["entityId"], "testEntity")

	// it calls MatchingComponentAdded on all systems when a matching component is added
	assert.Equal(t, tsa.MatchingComponentAddedCalled["testEntity:testComponentA"]["value"].(string), "testValueA")
	assert.Nil(t, tsa.MatchingComponentAddedCalled["testEntity:testComponentB"]["value"])

	assert.Nil(t, tsb.MatchingComponentAddedCalled["testEntity:testComponentA"]["value"])
	assert.Equal(t, tsb.MatchingComponentAddedCalled["testEntity:testComponentB"]["value"].(string), "testValueB")

	// test that component updates are handled correctly
	engine.Redis.Set(context.Background(), "__prev:testEntity:testComponentA", "testValueA", 0)
	engine.Redis.Set(context.Background(), "testEntity:testComponentA", "testValueA2", 0)

	<-tsa.Called

	assert.Equal(t, tsa.ComponentUpdatedCalled["testEntity:testComponentA"]["oldValue"].(string), "testValueA")
	assert.Equal(t, tsa.ComponentUpdatedCalled["testEntity:testComponentA"]["newValue"].(string), "testValueA2")

	assert.Nil(t, tsb.MatchingComponentUpdatedCalled["testEntity:testComponentA"])
	assert.Equal(t, tsa.ComponentUpdatedCalled["testEntity:testComponentA"]["oldValue"].(string), "testValueA")
	assert.Equal(t, tsa.MatchingComponentUpdatedCalled["testEntity:testComponentA"]["newValue"].(string), "testValueA2")

	teardown()
}

// Test that the system registry handles map set callbacks correctly.
func TestRegister_HSet(t *testing.T) {
	tsa, tsb := setup()

	engine.Redis.HSet(context.Background(), "testEntity:testComponentA", "testKeyA", "testValueA")

	<-tsa.Called

	engine.Redis.HSet(context.Background(), "testEntity:testComponentB", "testKeyB", "testValueB")

	<-tsb.Called

	// it calls ComponentAdded on all systems when a component is added
	assert.Equal(t, tsa.ComponentAddedCalled["testEntity:testComponentA"]["value"].(map[string]interface{}), map[string]interface{}{
		"testKeyA": "testValueA",
	})
	assert.Equal(t, tsa.ComponentAddedCalled["testEntity:testComponentB"]["value"].(map[string]interface{}), map[string]interface{}{
		"testKeyB": "testValueB",
	})

	assert.Equal(t, tsb.ComponentAddedCalled["testEntity:testComponentA"]["value"].(map[string]interface{}), map[string]interface{}{
		"testKeyA": "testValueA",
	})

	assert.Equal(t, tsb.ComponentAddedCalled["testEntity:testComponentB"]["value"].(map[string]interface{}), map[string]interface{}{
		"testKeyB": "testValueB",
	})

	assert.Equal(t, tsa.ComponentAddedCalled["testEntity:testComponentA"]["entityId"], "testEntity")
	assert.Equal(t, tsb.ComponentAddedCalled["testEntity:testComponentB"]["entityId"], "testEntity")

	// it calls MatchingComponentAdded on all systems when a matching component is added
	assert.Equal(t, tsa.MatchingComponentAddedCalled["testEntity:testComponentA"]["value"].(map[string]interface{}), map[string]interface{}{
		"testKeyA": "testValueA",
	})
	assert.Nil(t, tsa.MatchingComponentAddedCalled["testEntity:testComponentB"]["value"])

	assert.Nil(t, tsb.MatchingComponentAddedCalled["testEntity:testComponentA"]["value"])
	assert.Equal(t, tsb.MatchingComponentAddedCalled["testEntity:testComponentB"]["value"].(map[string]interface{}), map[string]interface{}{
		"testKeyB": "testValueB",
	})

	teardown()
}

func TestRegister_SAdd(t *testing.T) {
	tsa, tsb := setup()

	engine.Redis.SAdd(context.Background(), "testEntity:testComponentA", "testValueA")

	<-tsa.Called

	engine.Redis.SAdd(context.Background(), "testEntity:testComponentB", "testValueB")

	<-tsb.Called

	// it calls ComponentAdded on all systems when a component is added
	assert.Equal(t, tsa.ComponentAddedCalled["testEntity:testComponentA"]["value"].([]interface{}), []interface{}{
		"testValueA",
	})
	assert.Equal(t, tsa.ComponentAddedCalled["testEntity:testComponentB"]["value"].([]interface{}), []interface{}{
		"testValueB",
	})

	assert.Equal(t, tsb.ComponentAddedCalled["testEntity:testComponentA"]["value"].([]interface{}), []interface{}{
		"testValueA",
	})
	assert.Equal(t, tsb.ComponentAddedCalled["testEntity:testComponentB"]["value"].([]interface{}), []interface{}{
		"testValueB",
	})

	assert.Equal(t, tsa.ComponentAddedCalled["testEntity:testComponentA"]["entityId"], "testEntity")
	assert.Equal(t, tsb.ComponentAddedCalled["testEntity:testComponentB"]["entityId"], "testEntity")

	// it calls MatchingComponentAdded on all systems when a matching component is added
	assert.Equal(t, tsa.MatchingComponentAddedCalled["testEntity:testComponentA"]["value"].([]interface{}), []interface{}{
		"testValueA",
	})
	assert.Nil(t, tsa.MatchingComponentAddedCalled["testEntity:testComponentB"]["value"])

	assert.Nil(t, tsb.MatchingComponentAddedCalled["testEntity:testComponentA"]["value"])
	assert.Equal(t, tsb.MatchingComponentAddedCalled["testEntity:testComponentB"]["value"].([]interface{}), []interface{}{
		"testValueB",
	})

	teardown()
}

func TestRegister_Del(t *testing.T) {
	tsa, tsb := setup()

	// test that string deletes are handled correctly
	// set the value
	engine.Redis.Set(context.Background(), "testEntity:testComponentA", "testValueA", 0)
	<-tsa.Called

	// set the previous value
	engine.Redis.Set(context.Background(), "__prev:testEntity:testComponentA", "testValueA", 0)

	// set the value type to string
	engine.Redis.Set(context.Background(), "__type:testEntity:testComponentA", "string", 0)

	// Now delete the value triggering the callbacks
	engine.Redis.Del(context.Background(), "testEntity:testComponentA")
	<-tsa.Called

	// it calls ComponentRemoved on all systems when a component is removed
	assert.Equal(t, tsa.ComponentRemovedCalled["testEntity:testComponentA"]["value"].(string), "testValueA")
	assert.Equal(t, tsb.ComponentRemovedCalled["testEntity:testComponentA"]["value"].(string), "testValueA")

	assert.
		Equal(t, tsa.MatchingComponentRemovedCalled["testEntity:testComponentA"]["value"].(string), "testValueA")

	metaType := engine.Redis.Exists(context.Background(), "__type:testEntity:testComponentA").Val()

	assert.Equal(t, metaType, int64(0))

	prevValue := engine.Redis.Exists(context.Background(), "__prev:testEntity:testComponentA").Val()

	assert.Equal(t, prevValue, int64(0))

	teardown()
}

func TestRegister_DelHash(t *testing.T) {
	tsa, tsb := setup()

	// test that hash deletes are handled correctly
	// set the value
	engine.Redis.HSet(context.Background(), "testEntity:testComponentA", "testKeyA", "testValueA")
	<-tsa.Called

	// set the previous value
	engine.Redis.HSet(context.Background(), "__prev:testEntity:testComponentA", "testKeyA", "testValueA")

	// set the value type to hash
	engine.Redis.Set(context.Background(), "__type:testEntity:testComponentA", "map", 0)

	// Now delete the value triggering the callbacks
	engine.Redis.Del(context.Background(), "testEntity:testComponentA", "testKeyA")
	<-tsa.Called

	// it calls ComponentRemoved on all systems when a component is removed
	assert.Equal(t, tsa.ComponentRemovedCalled["testEntity:testComponentA"]["value"].(map[string]interface{}), map[string]interface{}{
		"testKeyA": "testValueA",
	})
	assert.Equal(t, tsb.ComponentRemovedCalled["testEntity:testComponentA"]["value"].(map[string]interface{}), map[string]interface{}{
		"testKeyA": "testValueA",
	})

	assert.
		Equal(t, tsa.MatchingComponentRemovedCalled["testEntity:testComponentA"]["value"].(map[string]interface{}), map[string]interface{}{
			"testKeyA": "testValueA",
		})

	metaType := engine.Redis.Exists(context.Background(), "__type:testEntity:testComponentA", "testKeyA").Val()

	assert.Equal(t, metaType, int64(0))

	prevValue := engine.Redis.Exists(context.Background(), "__prev:testEntity:testComponentA", "testKeyA").Val()

	assert.Equal(t, prevValue, int64(0))
	teardown()
}

func TestRegister_DelSet(t *testing.T) {
	tsa, tsb := setup()

	// test that set deletes are handled correctly
	// set the value
	engine.Redis.SAdd(context.Background(), "testEntity:testComponentA", "testValueA")
	<-tsa.Called

	// set the previous value
	engine.Redis.SAdd(context.Background(), "__prev:testEntity:testComponentA", "testValueA")

	// set the value type to set
	engine.Redis.Set(context.Background(), "__type:testEntity:testComponentA", "set", 0)

	// Now delete the value triggering the callbacks
	engine.Redis.Del(context.Background(), "testEntity:testComponentA", "testValueA")
	<-tsa.Called

	// it calls ComponentRemoved on all systems when a component is removed
	assert.Equal(t, tsa.ComponentRemovedCalled["testEntity:testComponentA"]["value"].([]interface{}), []interface{}{
		"testValueA",
	})
	assert.Equal(t, tsb.ComponentRemovedCalled["testEntity:testComponentA"]["value"].([]interface{}), []interface{}{
		"testValueA",
	})

	assert.
		Equal(t, tsa.MatchingComponentRemovedCalled["testEntity:testComponentA"]["value"].([]interface{}), []interface{}{
			"testValueA",
		})

	metaType := engine.Redis.Exists(context.Background(), "__type:testEntity:testComponentA").Val()

	assert.Equal(t, metaType, int64(0))

	prevValue := engine.Redis.Exists(context.Background(), "__prev:testEntity:testComponentA").Val()

	assert.Equal(t, prevValue, int64(0))
	teardown()
}
