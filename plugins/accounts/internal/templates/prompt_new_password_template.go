package templates

type promptNewPasswordTemplate struct{}

func (t promptNewPasswordTemplate) Name() string {
	return "prompt_new_password"
}

func (t promptNewPasswordTemplate) Style() string {
	return "default"
}

func (t promptNewPasswordTemplate) Render(_ interface{}) (string, error) {
	return "Enter a password:", nil
}

var PromptNewPasswordTemplate = &promptNewPasswordTemplate{}
