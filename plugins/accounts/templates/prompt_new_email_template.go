package templates

type promptNewEmailTemplate struct{}

func (t promptNewEmailTemplate) Name() string {
	return "prompt_new_email"
}

func (t promptNewEmailTemplate) Style() string {
	return "default"
}

func (t promptNewEmailTemplate) Render(_ interface{}) (string, error) {
	return "Enter an email:", nil
}

var PromptNewEmailTemplate = &promptNewEmailTemplate{}
