package accounts

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/accounts/internal/data_source"
	templates2 "github.com/mjolnir-mud/engine/plugins/accounts/internal/templates"
	"github.com/mjolnir-mud/engine/plugins/accounts/pkg/controllers/new_acccount"
	"github.com/mjolnir-mud/engine/plugins/accounts/pkg/entities"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/mongo_data_source"
	"github.com/mjolnir-mud/engine/plugins/templates"
)

type plugin struct{}

func (p *plugin) Name() string {
	return "accounts"
}

func (p *plugin) Registered() error {
	engine.EnsureRegistered(data_sources.Plugin.Name())
	engine.EnsureRegistered(mongo_data_source.Plugin.Name())
	engine.EnsureRegistered(ecs.Plugin.Name())
	engine.EnsureRegistered(templates.Plugin.Name())

	engine.RegisterBeforeStartCallback(func() {

		templates.RegisterTemplate(templates2.PromptUsernameTemplate)
		templates.RegisterTemplate(templates2.PromptPasswordTemplate)
		templates.RegisterTemplate(templates2.PromptEmailTemplate)
		templates.RegisterTemplate(templates2.PromptNewPasswordTemplate)
		templates.RegisterTemplate(templates2.PromptPasswordConfirmationTemplate)
		templates.RegisterTemplate(templates2.PromptNewUsernameTemplate)
		templates.RegisterTemplate(templates2.PromptNewEmailTemplate)
		templates.RegisterTemplate(templates2.InvalidEmailAddressTemplate)
		templates.RegisterTemplate(templates2.PasswordMatchFailTemplate)
		templates.RegisterTemplate(templates2.UsernameTakenTemplate)

		data_sources.Register(data_source.Accounts)
		ecs.RegisterEntityType(entities.Account)
	})

	return nil
}

// RegisterUsernameValidator overwrites the default username validator. This is a function that simply returns an error
// whose message will be presented to the connected user when the error is present.
func RegisterUsernameValidator(validator func(username string) error) {
	new_acccount.UsernameValidator = validator
}

// RegisterPasswordValidator overwrites the default password validator. This is a function that simply returns an error
// whose message will be presented to the connected user when the error is present.
func RegisterPasswordValidator(validator func(password string) error) {
	new_acccount.PasswordValidator = validator
}

var Plugin = plugin{}
