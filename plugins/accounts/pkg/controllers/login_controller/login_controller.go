package login_controller

import (
	"github.com/mjolnir-mud/engine/plugins/accounts/pkg/entities/account"
	"github.com/mjolnir-mud/engine/plugins/templates"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/systems/session"
)

// controller is the login controller, responsible handling user logins.
type controller struct{}

var AfterLoginCallback = func(id string) error {
	err := session.SetController(id, "game")

	if err != nil {
		return err
	}

	return nil
}

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

func Login(id string, accountId string) error {
	err := session.SetStringInStore(id, "accountId", accountId)

	if err != nil {
		return err
	}
	err = AfterLoginCallback(id)

	if err != nil {
		return err
	}

	return nil
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
	username, err := session.GetStringFromFlash(id, "username")

	if err != nil {
		return err
	}

	accountId, err := account.ValidateAccount(account.Credentials{
		Username: username,
		Password: input,
	})

	if err != nil {
		err := session.RenderTemplate(id, "login_invalid", nil)

		if err != nil {
			return err
		}

		return promptLoginUsername(id)
	}

	return Login(id, accountId)
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
