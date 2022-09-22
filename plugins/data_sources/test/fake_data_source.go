package test

import (
	"github.com/google/uuid"
	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/data_source"
)

type fakeDataSource struct {
	entities map[string]map[string]interface{}
}

func (f fakeDataSource) Name() string {
	return "fake"
}

func (f fakeDataSource) All() (map[string]map[string]interface{}, error) {
	entities := make(map[string]map[string]interface{})

	for id, entity := range f.entities {
		e := make(map[string]interface{})

		for key, value := range entity {
			e[key] = value
		}

		entities[id] = e
	}

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
	delete(f.entities, entityId)

	return nil
}

func (f fakeDataSource) Find(search map[string]interface{}) (map[string]map[string]interface{}, error) {
	entities, _ := f.All()
	response := make(map[string]map[string]interface{})

	if search["id"] != nil {
		if entity, ok := entities[search["id"].(string)]; ok {
			entity["id"] = search["id"]
			response[search["id"].(string)] = entity
		}

		return response, nil
	}

	for id, entity := range entities {
		for key, value := range search {
			if entity[key] != value {
				continue
			}

			response[id] = entity
		}
	}

	return response, nil
}

func (f fakeDataSource) FindAndDelete(search map[string]interface{}) error {
	entities, err := f.Find(search)

	if err != nil {
		return err
	}

	for id, _ := range entities {
		err := f.Delete(id)

		if err != nil {
			return err
		}
	}

	return nil
}

func (t fakeDataSource) FindOne(search map[string]interface{}) (string, map[string]interface{}, error) {
	entities, err := t.Find(search)
	if err != nil {
		return "", nil, err
	}

	for id, entity := range entities {
		return id, entity, nil
	}

	return "", nil, nil
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
	f.entities[entityId] = entity
	return nil
}

func (f fakeDataSource) Start() error {
	return nil
}

func (f fakeDataSource) Stop() error {
	return nil
}

var FakeDataSource = func() data_source.DataSource {
	return fakeDataSource{
		entities: map[string]map[string]interface{}{
			"test1": map[string]interface{}{
				"__metadata": map[string]interface{}{
					"entityType": "fake",
				},
				"testComponent":  "test1",
				"otherComponent": "other",
			},
			"test2": map[string]interface{}{
				"__metadata": map[string]interface{}{
					"entityType": "fake",
				},
				"testComponent":  "test2",
				"otherComponent": "other",
			},
		},
	}
}
