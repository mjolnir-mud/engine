package templates

import "github.com/mjolnir-mud/engine/plugins/templates"

type promptUsernameTemplate struct{}

func (t *promptUsernameTemplate) Name() string {
	return "prompt_username"
}

func (t *promptUsernameTemplate) Style() string {
	return "default"
}

func (t *promptUsernameTemplate) Render(_ interface{}) (string, error) {
	th, err := templates.GetTheme("default")

	if err != nil {
		return "", err
	}

	style := th.GetStyleFor("command")

	return "Enter your username, or type '" + style.Render("create") + "' to create a new account:", nil
}

var PromptUsernameTemplate = &promptUsernameTemplate{}
