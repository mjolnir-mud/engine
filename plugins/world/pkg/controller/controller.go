package controller

import "github.com/mjolnir-mud/engine/pkg/reactor"

// Controller is the interface for a session controller. A session controller handles interactions from the player
// session to the game world
type Controller interface {
	// Name returns the name of the controller. If multiple controllers of the same name are registered with the world
	// then the last one registered will be used. This enables the developer to override specific controllers with their
	// own implementation.
	Name() string

	// Start is called when the controller is set.
	Start(session reactor.Session) error

	// Resume called when the world restarts, causing the portal to reset-assert the session.
	Resume(session reactor.Session) error

	// Stop is called when the controller is unset.
	Stop(session reactor.Session) error

	// HandleInput is called when the player sends input to the world.
	HandleInput(session reactor.Session, input string) error
}
