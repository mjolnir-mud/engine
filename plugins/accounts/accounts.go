package accounts

import (
	"github.com/mjolnir-mud/engine/plugins/accounts/controllers/login"
	"github.com/mjolnir-mud/engine/plugins/accounts/controllers/new_account"
	"github.com/mjolnir-mud/engine/plugins/accounts/internal/plugin"
	accountSystem "github.com/mjolnir-mud/engine/plugins/accounts/systems/account"
)

// RegisterUsernameValidator overwrites the default username validator. This is a function that simply returns an error
// whose message will be presented to the connected user when the error is present.
func RegisterUsernameValidator(validator func(username string) error) {
	accountSystem.UsernameValidator = validator
}

// RegisterPasswordValidator overwrites the default password validator. This is a function that simply returns an error
// whose message will be presented to the connected user when the error is present.
func RegisterPasswordValidator(validator func(password string) error) {
	accountSystem.PasswordValidator = validator
}

// RegisterAfterLoginCallback registers a callback that will be called after a successful login. This can be used to set
// up post login tasks such as setting a new controller. The callback will be called with the session id as an argument.
func RegisterAfterLoginCallback(callback func(id string) error) {
	login.AfterLoginCallback = callback
}

// RegisterAfterCreateCallback registers a callback that will be called after a successful account creation. This
// can be used to set up post account creation tasks such as setting a new controller. The callback will be called with
// the session id as an argument, and the account id as an argument.
func RegisterAfterCreateCallback(callback func(id string, accountId string) error) {
	new_account.AfterCreateCallback = callback
}

var Plugin = plugin.Plugin
