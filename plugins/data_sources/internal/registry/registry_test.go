package registry

import (
	"github.com/mjolnir-mud/engine/plugins/data_sources/constants"
	"github.com/mjolnir-mud/engine/plugins/data_sources/testing/fakes"
	testing2 "github.com/mjolnir-mud/engine/plugins/ecs/testing"
	engineTesting "github.com/mjolnir-mud/engine/testing"
	"testing"

	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/stretchr/testify/assert"
)

func setup() {
	testing2.Setup()

	engineTesting.RegisterSetupCallback("data_sources", func() {

		engine.RegisterBeforeServiceStartCallback("world", func() {
			Start()
			Register(fakes.FakeDataSource())
		})

		engine.RegisterAfterServiceStartCallback("world", func() {
			ecs.RegisterEntityType(fakes.FakeEntityType)
		})

		fakes.Reset()
	})

	engineTesting.Setup("world")
}

func teardown() {
	Stop()
	testing2.Teardown()
	engineTesting.Teardown()
}

func TestAll(t *testing.T) {
	setup()
	defer teardown()
	entities, err := All("fake")

	assert.Nil(t, err)
	assert.Equal(t, entities.Len(), 2)
	assert.NotNil(t, entities.Get("test1"))
}

func TestCreateEntity(t *testing.T) {
	setup()
	defer teardown()

	id, entity, err := CreateEntity("fake", "fake", map[string]interface{}{
		"testComponent": "test3",
	})

	assert.Nil(t, err)

	assert.Equal(t, "test3", entity["testComponent"])

	found, err := FindOne("fake", map[string]interface{}{
		"id": id,
	})

	assert.Nil(t, err)
	assert.Equal(t, "test3", found.Record["testComponent"])
}

func TestCreateEntityWithId(t *testing.T) {
	setup()
	defer teardown()

	entity, err := NewEntityWithId("fake", "fake", "test3", map[string]interface{}{
		"testComponent": "test3",
	})

	assert.Nil(t, err)

	assert.Equal(t, map[string]interface{}{
		"__metadata": map[string]interface{}{
			"entityType": "fake",
			"fake":       true,
		},
		"id":            "test3",
		"testComponent": "test3"}, entity)

	found, err := FindOne("fake", map[string]interface{}{
		"id": "test3",
	})

	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"testComponent": "test3"}, found.Record)
}

func TestCount(t *testing.T) {
	setup()
	defer teardown()

	count, err := Count("fake", map[string]interface{}{
		"testComponent": "test1",
	})

	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}

func TestDelete(t *testing.T) {
	setup()
	defer teardown()
	err := Delete("fake", "test1")

	assert.Nil(t, err)

	_, err = FindOne("fake", map[string]interface{}{
		"id": "test1",
	})

	assert.NotNil(t, err)
}

func TestFind(t *testing.T) {
	setup()
	defer teardown()

	entities, err := Find("fake", map[string]interface{}{
		"otherComponent": "other",
	})

	assert.Nil(t, err)
	assert.Equal(t, entities.Len(), 2)
}

func TestFindAndDelete(t *testing.T) {
	setup()
	defer teardown()

	err := FindAndDelete("fake", map[string]interface{}{
		"id": "test1",
	})

	entities, err := All("fake")

	assert.Nil(t, err)
	assert.Equal(t, entities.Len(), 1)
}

func TestFindOne(t *testing.T) {
	setup()
	defer teardown()

	entity, err := FindOne("fake", map[string]interface{}{
		"id": "test1",
	})

	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"testComponent":  "test1",
		"otherComponent": "other"}, entity.Record)
	assert.Equal(t, "test1", entity.Id)
}

func TestRegister(t *testing.T) {
	setup()
	defer teardown()

	_, ok := dataSources["fake"]

	assert.True(t, ok)
}

func TestSave(t *testing.T) {
	setup()
	defer teardown()

	id, err := Save("fake", map[string]interface{}{
		constants.MetadataKey: map[string]interface{}{
			constants.MetadataTypeKey: "fake",
		},
		"testComponent": "test3",
	})

	assert.Nil(t, err)

	entity, err := FindOne("fake", map[string]interface{}{
		"id": id,
	})

	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"testComponent": "test3"}, entity.Record)
}

func TestSaveWithId(t *testing.T) {
	setup()
	defer teardown()

	err := SaveWithId("fake", "test3", map[string]interface{}{
		constants.MetadataKey: map[string]interface{}{
			constants.MetadataTypeKey: "fake",
		},
		"testComponent": "test3",
	})

	assert.Nil(t, err)

	entity, err := FindOne("fake", map[string]interface{}{
		"id": "test3",
	})

	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"testComponent": "test3"}, entity.Record)
}
