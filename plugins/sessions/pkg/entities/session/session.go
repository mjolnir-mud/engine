package session

type sessionType struct{}

func (s sessionType) Name() string {
	return "session"
}

func (s sessionType) Create(args map[string]interface{}) map[string]interface{} {
	if _, ok := args["flash"]; !ok {
		args["flash"] = map[string]interface{}{}
	}

	return args
}

func (s sessionType) Validate(_ map[string]interface{}) error {
	return nil
}

var Type = sessionType{}
