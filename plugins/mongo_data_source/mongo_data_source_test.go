package mongo_data_source

import (
	"context"
	"testing"

	"github.com/mjolnir-mud/engine"
	"github.com/stretchr/testify/assert"
)

func setup() {
	engine.Start("test")
	engine.SetEnv("test")
	_ = Plugin.Start()
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

	_ = Plugin.Stop()
	engine.Stop()
}

func TestMongoDataSource_Load(t *testing.T) {
	setup()
	defer teardown()

	dataSource := New("entities")

	entity, err := dataSource.Load("entity_1")

	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"__metadata": map[string]interface{}{
			"collection": "entities",
			"type":       "fake",
		},
		"testComponent": "test",
	}, entity)

}

func TestMongoDataSource_LoadAll(t *testing.T) {
	setup()
	defer teardown()

	dataSource := New("entities")

	entities, err := dataSource.LoadAll()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(entities))
}

func TestMongoDataSource_Find(t *testing.T) {
	setup()
	defer teardown()

	dataSource := New("entities")

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
