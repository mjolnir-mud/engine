package templates

type promptPasswordTemplate struct{}

func (t promptPasswordTemplate) Name() string {
	return "prompt_password"
}

func (t promptPasswordTemplate) Style() string {
	return "default"
}

func (t promptPasswordTemplate) Render(_ interface{}) (string, error) {
	return "Enter your password:", nil
}

var PromptPasswordTemplate = &promptPasswordTemplate{}
