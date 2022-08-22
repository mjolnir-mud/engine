package login_controller

import (
	"github.com/mjolnir-mud/engine/plugins/templates"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/systems/session"
)

// controller is the login controller, responsible handling user logins.
type controller struct{}

func (l controller) Name() string {
	return "login"
}

func (l controller) Start(id string) error {
	return promptLoginUsername(id)
}

func (l controller) Resume(_ string) error {
	return nil
}

func (l controller) Stop(_ string) error {
	return nil
}

func (l controller) HandleInput(id string, input string) error {
	return handleInput(id, input)
}

func handleInput(id string, input string) error {
	i, err := session.GetIntFromFlashWithDefault(id, "step", 1)

	if err != nil {
		return err
	}

	switch i {
	case 1:
		return handleUsername(id, input)
	case 2:
		return handlePassword(id, input)
	}

	return nil
}

func handleUsername(id string, input string) error {
	if input == "create" {
		err := session.SetController(id, "new_account")

		if err != nil {
			return err
		}

		return nil
	}

	err := session.SetStringInFlash(id, "username", input)

	if err != nil {
		return err
	}

	return promptPassword(id)
}

func handlePassword(id string, input string) error {
	return nil
}

func promptPassword(id string) error {
	err := session.SetIntInFlash(id, "step", 2)

	if err != nil {
		return err
	}

	err = session.RenderTemplate(id, "prompt_password", nil)

	if err != nil {
		return err
	}
	return nil
}

func promptLoginUsername(id string) error {
	err := session.SetIntInFlash(id, "step", 1)

	v, err := templates.RenderTemplate("prompt_username", nil)

	if err != nil {
		return err
	}

	err = session.SendLine(id, v)

	if err != nil {
		return err
	}

	return nil
}
