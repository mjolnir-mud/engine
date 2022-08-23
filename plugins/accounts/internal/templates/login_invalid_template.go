package templates

type loginInvalidTemplate struct{}

func (t loginInvalidTemplate) Name() string {
	return "login_invalid"
}

func (t loginInvalidTemplate) Style() string {
	return "default"
}

func (t loginInvalidTemplate) Render(_ interface{}) (string, error) {
	return "An account with that username and password combination was not found.", nil
}

var LoginInvalidTemplate = loginInvalidTemplate{}
