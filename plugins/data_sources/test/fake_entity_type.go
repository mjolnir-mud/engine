package test

type fakeEntityType struct{}

func (f fakeEntityType) Name() string {
	return "fake"
}

func (f fakeEntityType) Create(entity map[string]interface{}) map[string]interface{} {
	return entity
}

var FakeEntityType = fakeEntityType{}
