package templates

import "fmt"

type usernameTakenTemplate struct{}

func (t usernameTakenTemplate) Name() string {
	return "username_taken"
}

func (t usernameTakenTemplate) Style() string {
	return "default"
}

func (t usernameTakenTemplate) Render(args interface{}) (string, error) {
	username := args.(string)

	return fmt.Sprintf("The username '%s' is already taken.", username), nil
}

var UsernameTakenTemplate = &usernameTakenTemplate{}
