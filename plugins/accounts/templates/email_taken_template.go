package templates

type emailTakenTemplate struct{}

func (t emailTakenTemplate) Name() string {
	return "email_taken"
}

func (t emailTakenTemplate) Style() string {
	return "default"
}

func (t emailTakenTemplate) Render(_ interface{}) (string, error) {
	return "That email already belongs to an account.", nil
}

var EmailTakenTemplate = &emailTakenTemplate{}
