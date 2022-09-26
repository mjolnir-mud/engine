package templates

type passwordMatchFailTemplate struct{}

func (t passwordMatchFailTemplate) Name() string {
	return "password_match_fail"
}

func (t passwordMatchFailTemplate) Style() string {
	return "default"
}

func (t passwordMatchFailTemplate) Render(_ interface{}) (string, error) {
	return "The password and the confirmation do not match.", nil
}

var PasswordMatchFailTemplate = &passwordMatchFailTemplate{}
