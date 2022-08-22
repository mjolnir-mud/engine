package account

import (
	"github.com/mjolnir-mud/engine/plugins/accounts/pkg/errors"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	"golang.org/x/crypto/bcrypt"
)

type accountType struct{}

type Credentials struct {
	Username string
	Password string
}

func (a accountType) Name() string {
	return "account"
}

func (a accountType) Create(args map[string]interface{}) map[string]interface{} {
	return args
}

// ValidateAccount validates the account credentials returning the account id.
func ValidateAccount(args Credentials) (string, error) {
	id, r, err := data_sources.FindOne("accounts", map[string]interface{}{"username": args.Username})

	if err != nil {
		return "", err
	}

	if r == nil {
		return "", errors.AccountNotFoundError{}
	}

	err = bcrypt.CompareHashAndPassword([]byte(r["password"].(string)), []byte(args.Password))

	if err != nil {
		return "", errors.AccountNotFoundError{}
	}

	return id, nil
}

var Type = accountType{}