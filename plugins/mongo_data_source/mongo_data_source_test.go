package mongo_data_source

import (
	"context"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	"testing"

	"github.com/mjolnir-mud/engine"
	"github.com/stretchr/testify/assert"
)

func setup() {
	engine.RegisterPlugin(data_sources.Plugin)
	engine.RegisterPlugin(Plugin)
	engine.SetEnv("test")
	engine.Start("test")
	// seed data
	_, _ = Plugin.collection("entities").InsertOne(context.Background(), map[string]interface{}{
		"_id": "entity_1",
		"__metadata": map[string]interface{}{
			"collection": "entities",
			"type":       "fake",
		},
		"testComponent": "test",
	})

	_, _ = Plugin.collection("entities").InsertOne(context.Background(), map[string]interface{}{
		"_id": "entity_2",
		"__metadata": map[string]interface{}{
			"collection": "entities",
			"type":       "fake",
		},
		"testComponent": "test2",
	})
}

func teardown() {
	_ = Plugin.collection("entities").Drop(context.Background())
	engine.Stop()
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

func TestMongoDataSource_Find(t *testing.T) {
	setup()
	defer teardown()

	dataSource := New("entities")
	_ = dataSource.Start()

	entities, err := dataSource.Find(map[string]interface{}{
		"testComponent": "test",
	})

	assert.Nil(t, err)
	assert.Equal(t, 1, len(entities))
}

func TestMongoDataSource_Save(t *testing.T) {
	setup()
	defer teardown()

	dataSource := New("entities")
	_ = dataSource.Start()

	err := dataSource.Save("entity_3", map[string]interface{}{
		"__metadata": map[string]interface{}{
			"collection": "entities",
			"type":       "fake",
		},
		"testComponent": "test",
	})

	assert.Nil(t, err)

	c, err := Plugin.collection("entities").CountDocuments(context.Background(), map[string]interface{}{
		"_id": "entity_3",
	})

	assert.Nil(t, err)
	assert.Equal(t, int64(1), c)
}

func TestMongoDataSource_FindOne(t *testing.T) {
	setup()
	defer teardown()

	dataSource := New("entities")
	_ = dataSource.Start()

	id, entity, err := dataSource.FindOne(map[string]interface{}{
		"testComponent": "test",
	})

	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"__metadata": map[string]interface{}{
			"collection": "entities",
			"type":       "fake",
		},
		"testComponent": "test",
	}, entity)

	assert.Equal(t, "entity_1", id)
}

func TestMongoDataSource_Count(t *testing.T) {
	setup()
	defer teardown()

	dataSource := New("entities")
	_ = dataSource.Start()

	c, err := dataSource.Count(map[string]interface{}{
		"testComponent": "test",
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

	c, err := Plugin.collection("entities").CountDocuments(context.Background(), map[string]interface{}{
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
		"testComponent": "test",
	})

	assert.Nil(t, err)

	c, err := Plugin.collection("entities").CountDocuments(context.Background(), map[string]interface{}{
		"_id": "entity_1",
	})

	assert.Nil(t, err)
	assert.Equal(t, int64(0), c)
}
