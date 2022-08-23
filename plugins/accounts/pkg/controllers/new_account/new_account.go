package new_account

import (
	"fmt"
	"github.com/mjolnir-mud/engine/plugins/accounts/pkg/controllers/login_controller"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/systems/session"
	"net/mail"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type controller struct{}

var Controller = controller{}

var UsernameValidator = func(username string) error {
	r, err := regexp.Compile("^[a-zA-Z0-9_\\-]+$")

	if err != nil {
		panic(err)
	}

	if len(username) < 4 {
		return fmt.Errorf("'%s' is not a valid username. It must be at least 4 characters long", username)
	}

	if len(username) > 10 {
		return fmt.Errorf("'%s' is not a valid username. It must be less than 10 characters long", username)
	}

	if !r.MatchString(username) {
		return fmt.Errorf(
			"'%s' is not a valid username. It must contain only alpha-numeric characters, dashes (-) or underscrores (_)",
			username,
		)
	}

	return nil
}
var PasswordValidator = func(password string) error {
	if len(password) < 8 {
		//goland:noinspection ALL
		return fmt.Errorf("Not a valid password. It must be at least 8 characters long.")
	}

	r, err := regexp.Compile("^(?=.{8,})(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*]).*$")

	if err != nil {
		panic(err)
	}

	if !r.MatchString(password) {
		//goland:noinspection GoErrorStringFormat
		return fmt.Errorf(
			"Not a valid password. It must contain at least one lowercase letter, one uppercase letter, one number, and one special character.",
		)
	}

	return nil
}

var AfterCreateCallback = func(sessId string, entityId string) error {
	return login_controller.Login(sessId, entityId)
}

func (n controller) Name() string {
	return "new_account"
}

func (n controller) Start(id string) error {
	return promptNewUsername(id)
}

func (n controller) Resume(_ string) error {
	return nil
}

func (n controller) Stop(_ string) error {
	return nil
}

func (n controller) HandleInput(id string, input string) error {
	step, err := session.GetIntFromFlashWithDefault(id, "step", 1)

	if err != nil {
		return err
	}

	switch step {
	case 1:
		return handleUsername(id, input)
	case 2:
		return handleEmail(id, input)
	case 3:
		return handlePassword(id, input)
	case 4:
		return handlePasswordConfirmation(id, input)
	}

	return nil
}

func handlePassword(id string, input string) error {
	err := PasswordValidator(input)

	if err != nil {
		return session.SendLine(id, err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input), 4)

	if err != nil {
		return err
	}

	err = session.SetStringInFlash(id, "password", string(hash))

	if err != nil {
		return err
	}

	return promptPasswordConfirmation(id)
}

func handleEmail(id string, input string) error {
	input = strings.ToLower(input)
	_, err := mail.ParseAddress(input)

	if err != nil {
		err := session.SendLineF(id, "'%s' is an invalid email address.", input)

		if err != nil {
			return err
		}

		return promptNewEmail(id)
	}

	count, err := data_sources.Count("accounts", map[string]interface{}{
		"email": input,
	})

	if err != nil {
		return err
	}

	if count > 0 {
		err := session.SendLineF(id, "'%s' is already registered to an account.")

		if err != nil {
			return err
		}

		return promptNewEmail(id)
	}

	err = session.SetStringInFlash(id, "email", input)

	if err != nil {
		return err
	}

	return promptNewPassword(id)
}

func handleUsername(id, input string) error {
	count, err := data_sources.Count("accounts", map[string]interface{}{
		"id": strings.ToLower(input),
	})

	if err != nil {
		return err
	}

	if count > int64(0) {
		err := session.RenderTemplate(id, "username_taken", input)

		if err != nil {
			return err
		}

		return promptNewUsername(id)
	}

	err = UsernameValidator(input)

	if err != nil {
		return session.SendLine(id, err.Error())
	}

	err = session.SetStringInFlash(id, "username", input)

	if err != nil {
		return err
	}

	return promptNewEmail(id)
}

func handlePasswordConfirmation(id string, input string) error {
	hashedPassword, err := session.GetStringFromFlash(id, "password")

	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input))

	if err != nil {
		err := session.RenderTemplate(id, "password_match_fail", nil)

		if err != nil {
			return err
		}

		return promptNewPassword(id)
	}

	username, err := session.GetStringFromFlash(id, "username")

	if err != nil {
		return err
	}

	email, err := session.GetStringFromFlash(id, "email")

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	// the downcased username is used as the account's id
	userId := strings.ToLower(username)

	err = data_sources.Save("accounts", userId, map[string]interface{}{
		"username": username,
		"email":    email,
		"password": hashedPassword,
	})

	return AfterCreateCallback(id, userId)
}

func promptNewUsername(id string) error {
	err := session.SetIntInFlash(id, "step", 1)

	if err != nil {
		return err
	}

	return session.RenderTemplate(id, "prompt_new_username", nil)
}

func promptNewEmail(id string) error {
	err := session.SetIntInFlash(id, "step", 2)

	if err != nil {
		return err
	}

	err = session.RenderTemplate(id, "prompt_new_email", nil)

	if err != nil {
		return err
	}

	return nil
}

func promptNewPassword(id string) error {
	err := session.SetIntInFlash(id, "step", 3)

	if err != nil {
		return err
	}

	err = session.RenderTemplate(id, "prompt_new_password", nil)

	if err != nil {
		return err
	}

	return nil
}

func promptPasswordConfirmation(id string) error {
	err := session.SetIntInFlash(id, "step", 4)

	if err != nil {
		return err
	}

	err = session.RenderTemplate(id, "prompt_password_confirmation", nil)

	if err != nil {
		return err
	}

	return nil
}
