package testing

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/accounts/testing/helpers"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/mongo_data_source"
	"github.com/mjolnir-mud/engine/plugins/templates"
)

func Setup() {
	engine.RegisterPlugin(ecs.Plugin)
	engine.RegisterPlugin(data_sources.Plugin)
	engine.RegisterPlugin(mongo_data_source.Plugin)
	engine.RegisterPlugin(templates.Plugin)

	engine.RegisterAfterServiceStartCallback("world", func() {
		helpers.CreateDefaultAccount()
	})
}
