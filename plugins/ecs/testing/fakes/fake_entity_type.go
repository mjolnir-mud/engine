package fakes

type FakeEntityType struct{}

func (t FakeEntityType) Name() string {
	return "testing"
}

func (t FakeEntityType) New(args map[string]interface{}) map[string]interface{} {
	if _, ok := args["testComponent"]; !ok {
		args["testComponent"] = "testing"
	}

	return args
}

func (t FakeEntityType) Validate(_ map[string]interface{}) error {
	return nil
}
