package fakes

import (
	"github.com/google/uuid"
	"github.com/mjolnir-mud/engine/plugins/data_sources/data_source"
)

type fakeDataSource struct{}

var entities []map[string]interface{}

func Reset() {
	entities = []map[string]interface{}{
		{
			"__metadata": map[string]interface{}{
				"entityType": "fake",
			},
			"id":             "test1",
			"testComponent":  "test1",
			"otherComponent": "other",
		},
		{
			"__metadata": map[string]interface{}{
				"entityType": "fake",
			},
			"id":             "test2",
			"testComponent":  "test2",
			"otherComponent": "other",
		},
	}
}

func (f fakeDataSource) Name() string {
	return "fake"
}

func (f fakeDataSource) All() ([]map[string]interface{}, error) {
	return entities, nil
}

func (f fakeDataSource) AppendMetadata(metadata map[string]interface{}) map[string]interface{} {
	metadata["fake"] = true

	return metadata
}

func (f fakeDataSource) Count(filter map[string]interface{}) (int64, error) {
	entities, err := f.Find(filter)

	if err != nil {
		return 0, err
	}
	return int64(len(entities)), nil
}

func (f fakeDataSource) Delete(entityId string) error {
	for i, entity := range entities {
		if entity["id"] == entityId {
			entities = append(entities[:i], entities[i+1:]...)
			return nil
		}
	}

	return nil
}

func (f fakeDataSource) Find(search map[string]interface{}) ([]map[string]interface{}, error) {
	entities, _ := f.All()

	response := make([]map[string]interface{}, 0)

	for _, entity := range entities {
		for key, value := range search {
			if entity[key] != value {
				continue
			}

			response = append(response, entity)
		}
	}

	return response, nil
}

func (f fakeDataSource) FindAndDelete(search map[string]interface{}) error {
	entities, err := f.Find(search)

	if err != nil {
		return err
	}

	for _, e := range entities {
		err := f.Delete(e["id"].(string))

		if err != nil {
			return err
		}
	}

	return nil
}

func (t fakeDataSource) FindOne(search map[string]interface{}) (map[string]interface{}, error) {
	entities, err := t.Find(search)
	if err != nil {
		return nil, err
	}

	for _, entity := range entities {
		return entity, nil
	}

	return nil, nil
}

func (f fakeDataSource) Save(entity map[string]interface{}) (string, error) {
	uid, _ := uuid.NewUUID()
	uidStr := uid.String()

	err := f.SaveWithId(uidStr, entity)

	if err != nil {
		return "", err
	}

	return uidStr, nil
}

func (f fakeDataSource) SaveWithId(entityId string, entity map[string]interface{}) error {
	entity["id"] = entityId
	entities = append(entities, entity)

	return nil
}

func (f fakeDataSource) Start() error {
	return nil
}

func (f fakeDataSource) Stop() error {
	return nil
}

var FakeDataSource = func() data_source.Interface {
	return fakeDataSource{}
}
