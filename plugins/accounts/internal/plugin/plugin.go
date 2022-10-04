package plugin

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/accounts/controllers/login"
	"github.com/mjolnir-mud/engine/plugins/accounts/controllers/new_account"
	accountDataSource "github.com/mjolnir-mud/engine/plugins/accounts/data_sources/account"
	accountEntityType "github.com/mjolnir-mud/engine/plugins/accounts/entities/account"
	accountTemplates "github.com/mjolnir-mud/engine/plugins/accounts/templates"
	"github.com/mjolnir-mud/engine/plugins/controllers"
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

	// Ensure the data source gets registered before the data sources plugin starts all of its data sources.
	engine.RegisterBeforeServiceStartCallback("world", func() {
		data_sources.Register(accountDataSource.Create())
	})

	engine.RegisterAfterServiceStartCallback("world", func() {
		accountTemplates.RegisterAll()
		controllers.Register(login.Controller)
		controllers.Register(new_account.Controller)

		ecs.RegisterEntityType(accountEntityType.EntityType)
	})

	return nil
}

var Plugin = plugin{}
