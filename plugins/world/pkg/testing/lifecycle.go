package testing

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/world"
)

func Setup() {
	engine.RegisterPlugin(world.Plugin)
	testing.Setup()
}

func Teardown() {
	testing.Teardown()
}
