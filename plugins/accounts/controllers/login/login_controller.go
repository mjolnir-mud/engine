/*
 * Copyright (c) 2022 eightfivefour llc. All rights reserved.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package login

import (
	accountSystem "github.com/mjolnir-mud/engine/plugins/accounts/systems/account"
	"github.com/mjolnir-mud/engine/plugins/controllers"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/templates"
	"github.com/mjolnir-mud/engine/systems"
)

// controller is the login controller, responsible handling user logins.
type controller struct{}

var Controller = controller{}

var AfterLoginCallback = func(id string, accountId string) error {
	err := ecs.AddMapComponentToEntity(id, "actingAs", map[string]interface{}{
		"data_source": "accounts",
		"id":          accountId,
	})

	if err != nil {
		return err
	}

	err = controllers.Set(id, "game")

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
	_ = systems.SetAccountId(id, accountId)

	err := AfterLoginCallback(id, accountId)

	if err != nil {
		return err
	}

	return nil
}

func handleInput(id string, input string) error {
	i, err := systems.GetIntFromFlashWithDefault(id, "step", 1)

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
		err := controllers.Set(id, "new_account")

		if err != nil {
			return err
		}

		return nil
	}

	err := systems.SetStringInFlash(id, "username", input)

	if err != nil {
		return err
	}

	return promptPassword(id)
}

func handlePassword(id string, input string) error {
	username, err := systems.GetStringFromFlash(id, "username")

	if err != nil {
		return err
	}

	accountId, err := accountSystem.CompareAccountCredentials(accountSystem.Credentials{
		Username: username,
		Password: input,
	})

	if err != nil {
		err := systems.Render(id, "login_invalid", nil)

		if err != nil {
			return err
		}

		return promptLoginUsername(id)
	}

	return Login(id, accountId)
}

func promptPassword(id string) error {
	err := systems.SetIntInFlash(id, "step", 2)

	if err != nil {
		return err
	}

	err = systems.Render(id, "prompt_password", nil)

	if err != nil {
		return err
	}
	return nil
}

func promptLoginUsername(id string) error {
	err := systems.SetIntInFlash(id, "step", 1)

	v, err := templates.RenderTemplate("prompt_username", nil)

	if err != nil {
		return err
	}

	err = systems.SendLine(id, v)

	if err != nil {
		return err
	}

	return nil
}
