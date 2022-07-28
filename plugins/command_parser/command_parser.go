package command_parser

import (
	"github.com/mjolnir-mud/engine/plugins/command_parser/internal/command_registry"
	"github.com/mjolnir-mud/engine/plugins/command_parser/pkg/command_set"
)

type plugin struct{}

func (p plugin) Name() string {
	return "command_parser"
}

func (p plugin) Start() error {
	return nil
}

func RegisterCommandSet(set command_set.CommandSet) {
	command_registry.RegisterCommandSet(set)
}

func ParseCommand(sets []string, sess command_registry.Session, input string) {
	command_registry.ParseCommand(sets, sess, input)
}

var Plugin = plugin{}
