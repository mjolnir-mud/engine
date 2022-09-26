package templates

import "fmt"

type emailInvalidTemplate struct{}

func (t emailInvalidTemplate) Name() string {
	return "email_invalid"
}

func (t emailInvalidTemplate) Style() string {
	return "default"
}

func (t emailInvalidTemplate) Render(email interface{}) (string, error) {
	e := email.(string)

	return fmt.Sprintf("'%s' is not a valid email address.", e), nil
}

var EmailInvalidTemplate = &emailInvalidTemplate{}
