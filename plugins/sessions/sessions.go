package sessions

import (
	"github.com/mjolnir-mud/engine/plugins/sessions/internal/plugin"
	"github.com/mjolnir-mud/engine/plugins/sessions/internal/registry"
)

// RegisterSessionStartedHandler registers a handler that is called when a session is started.
func RegisterSessionStartedHandler(h func(id string) error) {
	registry.RegisterSessionStartedHandler(h)
}

// RegisterSessionStoppedHandler registers a handler that is called when a session is stopped.
func RegisterSessionStoppedHandler(h func(id string) error) {
	registry.RegisterSessionStoppedHandler(h)
}

// RegisterSessionLineHandler registers a handler that is called when a line is received from a session.
func RegisterLineHandler(h func(id string, line string) error) {
	registry.RegisterLineHandler(h)
}

// StopSessionRegistry stops the session registry. This should only be called non-portal services.
func StopSessionRegistry() {
	registry.Stop()
}

// StartSessionRegistry starts the session registry. This should only be called non-portal services.
func StartSessionRegistry() {
	registry.Start()
}

var Plugin = plugin.Plugin
