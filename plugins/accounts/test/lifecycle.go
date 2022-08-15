package test

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/templates"
)

func Setup() {
	engine.RegisterPlugin(templates.Plugin)
	testing.Setup()
}

func Teardown() {
	testing.Teardown()
}
