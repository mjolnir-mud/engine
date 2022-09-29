package testing

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
)

const TestAccountName = "test-account"

func Setup() {
	engine.RegisterAfterServiceStartCallback("world", func() {
		_, _ = data_sources.CreateEntityWithId(
			"accounts",
			"account",
			TestAccountName,
			map[string]interface{}{},
		)
	})
}
