package command_sets

import (
	"github.com/mjolnir-mud/engine/plugins/world/internal/commands"
)

type base struct {
	Say  commands.Say  `cmd`
	Look commands.Look `cmd`
}

func (c *base) Name() string {
	return "base"
}

var Base = &base{}
