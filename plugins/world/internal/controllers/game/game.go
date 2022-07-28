package game

import (
	"github.com/mjolnir-mud/engine/plugins/world/internal/command_registry"
	session2 "github.com/mjolnir-mud/engine/plugins/world/internal/systems/session"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/session"
)

type controller struct{}

var Controller = controller{}

func (c controller) Name() string {
	return "game"
}

func (C controller) Start(session session.Session) error {
	return nil
}

func (C controller) Resume(session session.Session) error {
	return nil
}

func (C controller) Stop(session session.Session) error {
	return nil
}

func (C controller) HandleInput(session session.Session, input string) error {
	// get the command sets from the store
	sets := session2.GetCommandSets(session)
	command_registry.ParseCommand(sets, session, input)

	return nil
}
