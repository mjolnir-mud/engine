package templates

type promptNewUsernameTemplate struct{}

func (t promptNewUsernameTemplate) Name() string {
	return "prompt_new_username"
}

func (t promptNewUsernameTemplate) Style() string {
	return "default"
}

func (t promptNewUsernameTemplate) Render(_ interface{}) (string, error) {
	return "Enter a username:", nil
}

var PromptNewUsernameTemplate = &promptNewUsernameTemplate{}
