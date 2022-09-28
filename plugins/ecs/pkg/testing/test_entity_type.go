package testing

type TestEntityType struct{}

func (t TestEntityType) Name() string {
	return "testing"
}

func (t TestEntityType) New(args map[string]interface{}) map[string]interface{} {
	if _, ok := args["testComponent"]; !ok {
		args["testComponent"] = "testing"
	}

	return args
}

func (t TestEntityType) Validate(args map[string]interface{}) error {
	return nil
}
