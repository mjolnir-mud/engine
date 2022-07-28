package main

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/telnet_portal"
	"github.com/mjolnir-mud/engine/plugins/templates"
	"github.com/mjolnir-mud/engine/plugins/world"
)

func main() {
	engine.RegisterPlugin(world.Plugin)
	engine.RegisterPlugin(telnet_portal.Plugin)
	engine.RegisterPlugin(templates.Plugin)

	world.RegisterLoadableDirectory("world")

	engine.Start("example")
}
