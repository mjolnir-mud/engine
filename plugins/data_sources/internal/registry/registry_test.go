package registry

import (
	engineTesting "github.com/mjolnir-mud/engine/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/constants"
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

func TestAll(t *testing.T) {
	setup()
	defer teardown()
	entities, err := All("fake")

	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"__metadata": map[string]interface{}{
			"entityType": "fake",
		},
		"testComponent":  "test1",
		"otherComponent": "other"}, entities["test1"])
}

func TestCreateEntity(t *testing.T) {
	setup()
	defer teardown()

	id, entity, err := NewEntity("fake", "fake", map[string]interface{}{
		"testComponent": "test3",
	})

	assert.Nil(t, err)

	assert.Equal(t, map[string]interface{}{
		"__metadata": map[string]interface{}{
			"entityType": "fake",
			"fake":       true,
		},
		"testComponent": "test3"}, entity)

	_, entity, err = FindOne("fake", map[string]interface{}{
		"id": id,
	})

	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"__metadata": map[string]interface{}{
			"entityType": "fake",
			"fake":       true,
		},
		"testComponent": "test3"}, entity)
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
		"testComponent": "test3"}, entity)

	_, entity, err = FindOne("fake", map[string]interface{}{
		"id": "test3",
	})

	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"__metadata": map[string]interface{}{
			"entityType": "fake",
			"fake":       true,
		},
		"testComponent": "test3"}, entity)
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

	_, _, err = FindOne("fake", map[string]interface{}{
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
	assert.Len(t, entities, 2)
}

func TestFindAndDelete(t *testing.T) {
	setup()
	defer teardown()

	err := FindAndDelete("fake", map[string]interface{}{
		"id": "test1",
	})

	entities, err := All("fake")

	assert.Nil(t, err)
	assert.Len(t, entities, 1)
}

func TestFindOne(t *testing.T) {
	setup()
	defer teardown()

	id, entity, err := FindOne("fake", map[string]interface{}{
		"id": "test1",
	})

	assert.Nil(t, err)
	assert.Equal(t, entity, map[string]interface{}{
		"__metadata": map[string]interface{}{
			"entityType": "fake",
		},
		"testComponent":  "test1",
		"otherComponent": "other"}, entity)
	assert.Equal(t, "test1", id)
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

	_, entity, err := FindOne("fake", map[string]interface{}{
		"id": id,
	})

	assert.Nil(t, err)
	assert.Equal(t, entity, map[string]interface{}{
		"__metadata": map[string]interface{}{
			"entityType": "fake",
		},
		"testComponent": "test3"}, entity)
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

	_, entity, err := FindOne("fake", map[string]interface{}{
		"id": "test3",
	})

	assert.Nil(t, err)
	assert.Equal(t, entity, map[string]interface{}{
		"__metadata": map[string]interface{}{
			"entityType": "fake",
		},
		"testComponent": "test3"}, entity)
}
