package templates

import (
	"github.com/mjolnir-mud/engine"
	engineTesting "github.com/mjolnir-mud/engine/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/templates"
)

func setup() {
	engineTesting.Setup("world", func() {
		engine.RegisterPlugin(templates.Plugin)

		engine.RegisterAfterServiceStartCallback("world", func() {
			RegisterAll()
		})
	})
}
