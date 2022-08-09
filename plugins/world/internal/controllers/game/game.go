package game

import (
	"github.com/mjolnir-mud/engine/pkg/reactor"
	"github.com/mjolnir-mud/engine/plugins/command_parser"
	session2 "github.com/mjolnir-mud/engine/plugins/world/internal/systems/session"
)

type controller struct{}

var Controller = controller{}

func (c controller) Name() string {
	return "game"
}

func (C controller) Start(session reactor.Session) error {
	return nil
}

func (C controller) Resume(session reactor.Session) error {
	return nil
}

func (C controller) Stop(session reactor.Session) error {
	return nil
}

func (C controller) HandleInput(session reactor.Session, input string) error {
	// get the command sets from the store
	sets := session2.GetCommandSets(session)
	command_parser.ParseCommand(sets, session, input)

	return nil
}
