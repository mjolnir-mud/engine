package test

type fakeEntityType struct{}

func (f fakeEntityType) Name() string {
	return "fake"
}

func (f fakeEntityType) Create(entity map[string]interface{}) map[string]interface{} {
	return entity
}

func (f fakeEntityType) Validate(entity map[string]interface{}) error {
	return nil
}

var FakeEntityType = fakeEntityType{}
