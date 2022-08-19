package test

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

func (f fakeDataSource) Load(entityId string) (map[string]interface{}, error) {
	return f.entities[entityId], nil
}

func (f fakeDataSource) LoadAll() (map[string]map[string]interface{}, error) {
	return f.entities, nil
}

func (f fakeDataSource) Find(search map[string]interface{}) (map[string]map[string]interface{}, error) {
	entities := make(map[string]map[string]interface{})

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

func (f fakeDataSource) Count(filter map[string]interface{}) (int64, error) {
	return int64(len(f.entities)), nil
}

func (f fakeDataSource) Save(entityId string, entity map[string]interface{}) error {
	return nil
}

var FakeDataSource = fakeDataSource{
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
