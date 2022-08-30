package test

type TestEntityType struct {
}

func (t TestEntityType) Name() string {
	return "test"
}

func (t TestEntityType) Create(args map[string]interface{}) map[string]interface{} {
	if _, ok := args["testComponent"]; !ok {
		args["testComponent"] = "test"
	}

	return args
}
