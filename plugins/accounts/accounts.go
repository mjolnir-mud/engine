package accounts

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/accounts/internal/data_source"
	templates2 "github.com/mjolnir-mud/engine/plugins/accounts/internal/templates"
	"github.com/mjolnir-mud/engine/plugins/accounts/pkg/controllers/login_controller"
	"github.com/mjolnir-mud/engine/plugins/accounts/pkg/controllers/new_account"
	"github.com/mjolnir-mud/engine/plugins/accounts/pkg/entities/account"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/mongo_data_source"
	"github.com/mjolnir-mud/engine/plugins/templates"
)

type plugin struct{}

func (p plugin) Name() string {
	return "accounts"
}

func (p plugin) Registered() error {
	engine.EnsureRegistered(data_sources.Plugin.Name())
	engine.EnsureRegistered(mongo_data_source.Plugin.Name())
	engine.EnsureRegistered(ecs.Plugin.Name())
	engine.EnsureRegistered(templates.Plugin.Name())

	engine.RegisterAfterServiceStartCallback("world", func() {
		templates2.RegisterAll()

		data_sources.Register(data_source.Create())
		ecs.RegisterEntityType(account.Type)
	})

	return nil
}

// RegisterUsernameValidator overwrites the default username validator. This is a function that simply returns an error
// whose message will be presented to the connected user when the error is present.
func RegisterUsernameValidator(validator func(username string) error) {
	new_account.UsernameValidator = validator
}

// RegisterPasswordValidator overwrites the default password validator. This is a function that simply returns an error
// whose message will be presented to the connected user when the error is present.
func RegisterPasswordValidator(validator func(password string) error) {
	new_account.PasswordValidator = validator
}

// RegisterAfterLoginCallback registers a callback that will be called after a successful login. This can be used to set
// up post login tasks such as setting a new controller. The callback will be called with the session id as an argument.
func RegisterAfterLoginCallback(callback func(id string) error) {
	login_controller.AfterLoginCallback = callback
}

// RegisterAfterCreateCallback registers a callback that will be called after a successful account creation. This
// can be used to set up post account creation tasks such as setting a new controller. The callback will be called with
// the session id as an argument, and the account id as an argument.
func RegisterAfterCreateCallback(callback func(id string, accountId string) error) {
	new_account.AfterCreateCallback = callback
}

var Plugin = plugin{}
