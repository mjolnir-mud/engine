package templates

import "fmt"

type invalidEmailAddressTemplate struct{}

func (t invalidEmailAddressTemplate) Name() string {
	return "invalid_email_address"
}

func (t invalidEmailAddressTemplate) Style() string {
	return "default"
}

func (t invalidEmailAddressTemplate) Render(args interface{}) (string, error) {
	email := args.(string)

	return fmt.Sprintf("The provided email address '%s' is invalid.", email), nil
}

var InvalidEmailAddressTemplate = &invalidEmailAddressTemplate{}
