package testing

import (
	"github.com/mjolnir-mud/engine"
	accountDataSource "github.com/mjolnir-mud/engine/plugins/accounts/data_sources/account"
	"github.com/mjolnir-mud/engine/plugins/accounts/internal/plugin"
	"github.com/mjolnir-mud/engine/plugins/accounts/testing/helpers"
	"github.com/mjolnir-mud/engine/plugins/controllers"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	dataSourcesTesting "github.com/mjolnir-mud/engine/plugins/data_sources/testing"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	ecsTesting "github.com/mjolnir-mud/engine/plugins/ecs/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/mongo_data_source"
	mongoDataSourceTesting "github.com/mjolnir-mud/engine/plugins/mongo_data_source/testing"
	"github.com/mjolnir-mud/engine/plugins/sessions"
	sessionsTesting "github.com/mjolnir-mud/engine/plugins/sessions/testing"
	"github.com/mjolnir-mud/engine/plugins/templates"
	templatesTesting "github.com/mjolnir-mud/engine/plugins/templates/pkg/testing"
	engineTesting "github.com/mjolnir-mud/engine/testing"
)

func Setup() {
	engineTesting.RegisterSetupCallback("accounts", func() {
		engine.RegisterPlugin(ecs.Plugin)
		engine.RegisterPlugin(data_sources.Plugin)
		engine.RegisterPlugin(mongo_data_source.Plugin)
		engine.RegisterPlugin(templates.Plugin)
		engine.RegisterPlugin(sessions.Plugin)
		engine.RegisterPlugin(controllers.Plugin)
		engine.RegisterPlugin(plugin.Plugin)

		ecsTesting.Setup()
		dataSourcesTesting.Setup()
		mongoDataSourceTesting.Setup()
		sessionsTesting.Setup()
		templatesTesting.Setup()

		engine.RegisterAfterServiceStartCallback("world", func() {
			helpers.CreateDefaultAccount()
			data_sources.Register(accountDataSource.Create())
		})
	})
}
