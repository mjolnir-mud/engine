package templates

type promptEmailTemplate struct{}

func (t promptEmailTemplate) Name() string {
	return "prompt_email"
}

func (t promptEmailTemplate) Style() string {
	return "default"
}

func (t promptEmailTemplate) Render(_ interface{}) (string, error) {
	return "Enter your email:", nil
}

var PromptEmailTemplate = promptEmailTemplate{}
