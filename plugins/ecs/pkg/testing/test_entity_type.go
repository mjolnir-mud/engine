package testing

type TestEntityType struct{}

func (t TestEntityType) Name() string {
	return "testing"
}

func (t TestEntityType) Create(args map[string]interface{}) map[string]interface{} {
	if _, ok := args["testComponent"]; !ok {
		args["testComponent"] = "testing"
	}

	return args
}
