package testing

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
)

func Setup() {
	engine.RegisterAfterServiceStartCallback("world", func() {
		_, _ = data_sources.CreateEntityWithId("accounts", "account", "test-account", map[string]interface{}{})
	})
}
