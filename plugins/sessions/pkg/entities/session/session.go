package session

type sessionType struct{}

func (s sessionType) Name() string {
	return "session"
}

func (s sessionType) Create(args map[string]interface{}) map[string]interface{} {
	args["expireIn"] = 900

	// if the args does not have a store yet then create one
	if _, ok := args["store"]; !ok {
		args["store"] = map[string]interface{}{}
	}

	// if the store does not have a controller set
	if _, ok := args["store"].(map[string]interface{})["controller"]; !ok {
		args["store"].(map[string]interface{})["controller"] = "login"
	}

	if _, ok := args["flash"]; !ok {
		args["flash"] = map[string]interface{}{}
	}

	return args
}
