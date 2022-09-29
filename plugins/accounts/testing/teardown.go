package testing

import (
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	"github.com/mjolnir-mud/engine/plugins/ecs"
)

func Teardown() {
	_ = ecs.RemoveEntity("test-account")
	_ = data_sources.Delete("accounts", "test-account")
}
