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
		data_sources.Register(data_source.Accounts)
		ecs.RegisterEntityType(entities.Account)
	})

	return nil
}

func RegisterUsernameValidator(validator func(username string) (bool, error)) {
	new_acccount.UsernameValidator = validator
}

func RegisterPasswordValidator(validator func(password string) (bool, error)) {
	new_acccount.PasswordValidator = validator
}

var Plugin = plugin{}
