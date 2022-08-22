package test

import "github.com/mjolnir-mud/engine/plugins/data_sources/pkg/data_source"

type fakeDataSource struct {
	entities map[string]map[string]interface{}
}

func (f fakeDataSource) Name() string {
	return "fake"
}

func (f fakeDataSource) Start() error {
	return nil
}

func (f fakeDataSource) Stop() error {
	return nil
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

func (f fakeDataSource) Find(search map[string]interface{}) (map[string]map[string]interface{}, error) {
	entities, _ := f.All()

	if search["id"] != nil {
		if entity, ok := f.entities[search["id"].(string)]; ok {
			entity["id"] = search["id"]
			entities[search["id"].(string)] = entity
		}
	}

	for id, entity := range f.entities {
		for key, value := range search {
			if entity[key] != value {
				continue
			}

			entities[id] = entity
		}
	}

	return entities, nil
}

func (t fakeDataSource) FindOne(search map[string]interface{}) (string, map[string]interface{}, error) {
	entities, err := t.Find(search)
	if err != nil {
		return "", nil, err
	}

	for _, entity := range entities {
		return "", entity, nil
	}

	return "", nil, nil
}

func (f fakeDataSource) Count(filter map[string]interface{}) (int64, error) {
	return int64(len(f.entities)), nil
}

func (f fakeDataSource) Save(entityId string, entity map[string]interface{}) error {
	return nil
}

func (f fakeDataSource) Delete(entityId string) error {
	delete(f.entities, entityId)

	return nil
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

var FakeDataSource = func() data_source.DataSource {
	return fakeDataSource{
		entities: map[string]map[string]interface{}{
			"test1": map[string]interface{}{
				"__metadata": map[string]interface{}{
					"entityType": "fake",
				},
				"testComponent": "test1",
			},
			"test2": map[string]interface{}{
				"__metadata": map[string]interface{}{
					"entityType": "fake",
				},
				"testComponent": "test2",
			},
		},
	}
}
