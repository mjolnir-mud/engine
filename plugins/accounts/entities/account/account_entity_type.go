package account

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type accountType struct{}

func (a accountType) Name() string {
	return "account"
}

func (a accountType) Create(args map[string]interface{}) map[string]interface{} {
	password, ok := args["password"].(string)

	if ok {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if err != nil {
			panic(err)
		}

		args["hashedPassword"] = string(hashedPassword)
		delete(args, "password")
	}

	return args
}

func (a accountType) Validate(args map[string]interface{}) error {
	return fmt.Errorf("account entities cannot be added")
}

var EntityType = accountType{}
