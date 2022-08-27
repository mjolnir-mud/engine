package sessions

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/sessions/internal/registry"
)

type plugin struct{}

func (p *plugin) Name() string {
	return "sessions"
}

func (p *plugin) Registered() error {
	engine.EnsureRegistered(ecs.Plugin.Name())

	return nil
}

// RegisterSessionStartedHandler registers a handler that is called when a session is started.
func RegisterSessionStartedHandler(h func(id string) error) {
	registry.RegisterSessionStartedHandler(h)
}

// RegisterSessionStoppedHandler registers a handler that is called when a session is stopped.
func RegisterSessionStoppedHandler(h func(id string) error) {
	registry.RegisterSessionStoppedHandler(h)
}

// StopSessionRegistry stops the session registry. This should only be called non-portal services.
func StopSessionRegistry() {
	registry.Stop()
}

// StartSessionRegistry starts the session registry. This should only be called non-portal services.
func StartSessionRegistry() {
	registry.Start()
}

var Plugin = &plugin{}
