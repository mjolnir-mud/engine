package registry

import (
	"testing"

	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/data_sources/test"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/stretchr/testify/assert"
)

func setup() {
	engine.Start("test")
	ecs.RegisterEntityType(test.FakeEntityType)
	Register(test.FakeDataSource)
	_ = ecs.Plugin.Start()
	_ = Start()
}

func teardown() {
	_ = Stop()
	_ = ecs.Plugin.Stop()
	engine.Stop()
}

func TestRegister(t *testing.T) {
	setup()
	defer teardown()
	if _, ok := r.dataSources["fake"]; !ok {
		t.Errorf("Expected registry.dataSources to contain fake data source")
	}
}

func TestLoad(t *testing.T) {
	setup()
	defer teardown()
	entity, err := Load("fake", "test1")

	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"testComponent": "test1"}, entity)
}

func TestLoadAll(t *testing.T) {
	setup()
	defer teardown()
	entities, err := LoadAll("fake")

	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"testComponent": "test"}, entities["fakeId"])
}
