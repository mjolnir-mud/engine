package test

type fakeEntityType struct{}

func (e fakeEntityType) Name() string {
	return "fake"
}

func (e fakeEntityType) Create(args map[string]interface{}) map[string]interface{} {
	return args
}
