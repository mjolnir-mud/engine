package game

import (
	"github.com/mjolnir-mud/engine/plugins/command_parser"
	session2 "github.com/mjolnir-mud/engine/plugins/world/pkg/systems/session"
)

type controller struct{}

var Controller = controller{}

func (c controller) Name() string {
	return "game"
}

func (C controller) Start(id string) error {
	return nil
}

func (C controller) Resume(id string) error {
	return nil
}

func (C controller) Stop(id string) error {
	return nil
}

func (C controller) HandleInput(id string, input string) error {
	// get the command sets from the store
	sets := session2.GetCommandSets(session2)
	command_parser.ParseCommand(sets, session2, input)

	return nil
}
