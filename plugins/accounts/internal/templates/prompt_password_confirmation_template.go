package templates

type promptPasswordConfirmationTemplate struct{}

func (t promptPasswordConfirmationTemplate) Name() string {
	return "prompt_password_confirmation"
}

func (t promptPasswordConfirmationTemplate) Style() string {
	return "default"
}

func (t promptPasswordConfirmationTemplate) Render(_ interface{}) (string, error) {
	return "Confirm your password:", nil
}

var PromptPasswordConfirmationTemplate = &promptPasswordConfirmationTemplate{}
