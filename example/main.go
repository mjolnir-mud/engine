package main

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/command_parser"
	"github.com/mjolnir-mud/engine/plugins/telnet_portal"
	"github.com/mjolnir-mud/engine/plugins/templates"
	"github.com/mjolnir-mud/engine/plugins/world"
)

func main() {
	engine.RegisterPlugin(command_parser.Plugin)
	engine.RegisterPlugin(world.Plugin)
	engine.RegisterPlugin(telnet_portal.Plugin)
	engine.RegisterPlugin(templates.Plugin)

	world.RegisterLoadableDirectory("world")

	engine.Start("example")
}
