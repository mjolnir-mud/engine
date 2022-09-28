package fakes

type fakeEntityType struct{}

func (f fakeEntityType) Name() string {
	return "fake"
}

func (f fakeEntityType) New(entity map[string]interface{}) map[string]interface{} {
	return entity
}

func (f fakeEntityType) Validate(_ map[string]interface{}) error {
	return nil
}

var FakeEntityType = fakeEntityType{}
