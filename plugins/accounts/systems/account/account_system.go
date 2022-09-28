package account

import (
	"fmt"
	"github.com/mjolnir-mud/engine/plugins/accounts/errors"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"regexp"
	"strings"
)

const MinUsernameLength = 4
const MaxUsernameLength = 20
const UsernameRegex = `^[a-zA-Z0-9_]+$`

const PasswordEntropy = 40

type Credentials struct {
	Username string
	Password string
}

var UsernameValidator = func(username string) error {
	r, err := regexp.Compile(UsernameRegex)

	if err != nil {
		panic(err)
	}

	if len(username) < MinUsernameLength {
		//goland:noinspection GoErrorStringFormat
		return fmt.Errorf(
			"'%s' is not a valid username. It must be at least %d characters long.",
			username,
			MinUsernameLength,
		)
	}

	if len(username) >= MaxUsernameLength {
		//goland:noinspection GoErrorStringFormat
		return fmt.Errorf(
			"'%s' is not a valid username. It must be at most %d characters long.",
			username,
			MaxUsernameLength,
		)
	}

	if !r.MatchString(username) {
		return fmt.Errorf(
			"'%s' is not a valid username. It must contain only alpha-numeric characters, dashes (-) or underscores (_)",
			username,
		)
	}

	return nil
}

var PasswordValidator = func(password string) error {
	e := passwordvalidator.GetEntropy(password)

	if e < PasswordEntropy {
		//goland:noinspection GoErrorStringFormat
		return fmt.Errorf("That isn't a very secure password. Try something stronger.")
	}

	return nil
}

type system struct{}

func (e system) Name() string {
	return "account"
}

func (e system) Component() string {
	return "_"
}

func (e system) Match(_ string, _ interface{}) bool { return true }

func (e system) ComponentAdded(_ string, _ string, _ interface{}) error { return nil }

func (e system) ComponentUpdated(_ string, _ string, _ interface{}, _ interface{}) error { return nil }

func (e system) ComponentRemoved(_ string, _ string) error { return nil }

func (e system) MatchingComponentAdded(_ string, _ string, _ interface{}) error { return nil }

func (e system) MatchingComponentUpdated(_ string, _ string, _ interface{}, _ interface{}) error {
	return nil
}

func (e system) MatchingComponentRemoved(_ string, _ string) error { return nil }

var System = system{}

func ValidateEmail(email string) error {
	email = strings.ToLower(email)
	_, err := mail.ParseAddress(email)

	if err != nil {
		return fmt.Errorf("'%s' is not a valid email address.", email)
	}

	return nil
}

func ValidateUsername(username string) error {
	return UsernameValidator(username)
}

func ValidatePassword(username string, email string, password string) error {
	if password == username {
		//goland:noinspection ALL
		return fmt.Errorf("Your password cannot be the same as your username.")
	}

	if password == email {
		//goland:noinspection ALL
		return fmt.Errorf("Your password cannot be the same as your email address.")
	}

	return PasswordValidator(password)
}

// CompareAccountCredentials validates the account credentials returning the account id.
func CompareAccountCredentials(args Credentials) (string, error) {
	r, err := data_sources.FindOne("accounts", map[string]interface{}{"username": args.Username})

	if err != nil {
		return "", err
	}

	if r == nil {
		return "", errors.AccountNotFoundError{}
	}

	err = bcrypt.CompareHashAndPassword([]byte(r.Record["hashedPassword"].(string)), []byte(args.Password))

	if err != nil {
		return "", errors.AccountNotFoundError{}
	}

	return r.Id, nil
}
