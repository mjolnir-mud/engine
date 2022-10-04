package mongo_data_source

import (
	"context"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	engineTesting "github.com/mjolnir-mud/engine/testing"
	"testing"

	"github.com/mjolnir-mud/engine"
	"github.com/stretchr/testify/assert"
)

func setup() {
	engineTesting.RegisterSetupCallback("mongo_data_source", func() {
		engine.RegisterPlugin(data_sources.Plugin)
		engine.RegisterPlugin(Plugin)

		engine.RegisterAfterServiceStartCallback("world", func() {

			// seed data
			_, _ = Plugin.Collection("entities").InsertOne(context.Background(), map[string]interface{}{
				"_id": sanitizeId("entity_1"),
				"__metadata": map[string]interface{}{
					"collection": "entities",
					"type":       "fake",
				},
				"testComponent": "testing",
			})

			_, _ = Plugin.Collection("entities").InsertOne(context.Background(), map[string]interface{}{
				"_id": sanitizeId("entity_2"),
				"__metadata": map[string]interface{}{
					"collection": "entities",
					"type":       "fake",
				},
				"testComponent": "test2",
			})
		})
	})

	engineTesting.Setup("world")
}

func teardown() {
	_ = Plugin.Collection("entities").Drop(context.Background())
	engineTesting.Teardown()
}

func TestMongoDataSource_All(t *testing.T) {
	setup()
	defer teardown()

	dataSource := New("entities")
	_ = dataSource.Start()

	entities, err := dataSource.All()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(entities))
}

func TestMongoDataSource_AppendMetadata(t *testing.T) {
	setup()
	defer teardown()

	dataSource := New("entities")
	_ = dataSource.Start()

	metadata := dataSource.AppendMetadata(map[string]interface{}{})

	assert.Equal(t, map[string]interface{}{
		"collection": "entities",
	}, metadata)
}

func TestMongoDataSource_Find(t *testing.T) {
	setup()
	defer teardown()

	dataSource := New("entities")
	_ = dataSource.Start()

	entities, err := dataSource.Find(map[string]interface{}{
		"testComponent": "testing",
	})

	assert.Nil(t, err)
	assert.Equal(t, 1, len(entities))
}

func TestMongoDataSource_SaveWithId(t *testing.T) {
	setup()
	defer teardown()

	dataSource := New("entities")
	_ = dataSource.Start()

	err := dataSource.SaveWithId("entity_3", map[string]interface{}{
		"__metadata": map[string]interface{}{
			"collection": "entities",
			"type":       "fake",
		},
		"testComponent": "testing",
	})

	assert.Nil(t, err)

	c, err := dataSource.Count(map[string]interface{}{
		"_id": "entity_3",
	})

	assert.Nil(t, err)
	assert.Equal(t, int64(1), c)
}

func TestMongoDataSource_Save(t *testing.T) {
	setup()
	defer teardown()

	dataSource := New("entities")
	_ = dataSource.Start()

	id, err := dataSource.Save(map[string]interface{}{
		"__metadata": map[string]interface{}{
			"collection": "entities",
			"type":       "fake",
		},
		"testComponent": "testing",
	})

	assert.Nil(t, err)

	c, err := dataSource.Count(map[string]interface{}{
		"_id": id,
	})

	assert.Nil(t, err)
	assert.Equal(t, int64(1), c)
}

func TestMongoDataSource_FindOne(t *testing.T) {
	setup()
	defer teardown()

	dataSource := New("entities")
	_ = dataSource.Start()

	entity, err := dataSource.FindOne(map[string]interface{}{
		"testComponent": "testing",
	})

	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"__metadata": map[string]interface{}{
			"collection": "entities",
			"type":       "fake",
		},
		"id":            sanitizeId("entity_1").Hex(),
		"testComponent": "testing",
	}, entity)
}

func TestMongoDataSource_Count(t *testing.T) {
	setup()
	defer teardown()

	dataSource := New("entities")
	_ = dataSource.Start()

	c, err := dataSource.Count(map[string]interface{}{
		"testComponent": "testing",
	})

	assert.Nil(t, err)
	assert.Equal(t, int64(1), c)
}

func TestMongoDataSource_Delete(t *testing.T) {
	setup()
	defer teardown()

	dataSource := New("entities")
	_ = dataSource.Start()

	err := dataSource.Delete("entity_1")

	assert.Nil(t, err)

	c, err := Plugin.Collection("entities").CountDocuments(context.Background(), map[string]interface{}{
		"_id": "entity_1",
	})

	assert.Nil(t, err)
	assert.Equal(t, int64(0), c)
}

func TestMongoDataSource_FindAndDelete(t *testing.T) {
	setup()
	defer teardown()

	dataSource := New("entities")
	_ = dataSource.Start()

	err := dataSource.FindAndDelete(map[string]interface{}{
		"testComponent": "testing",
	})

	assert.Nil(t, err)

	c, err := Plugin.Collection("entities").CountDocuments(context.Background(), map[string]interface{}{
		"_id": "entity_1",
	})

	assert.Nil(t, err)
	assert.Equal(t, int64(0), c)
}
