package templates

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/templates"
	engineTesting "github.com/mjolnir-mud/engine/testing"
)

func setup() {
	engineTesting.RegisterSetupCallback("templates", func() {
		engine.RegisterPlugin(templates.Plugin)

		engine.RegisterAfterServiceStartCallback("world", func() {
			RegisterAll()
		})
	})

	engineTesting.Setup("world")
}
