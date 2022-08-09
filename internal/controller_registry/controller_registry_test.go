package controller_registry

import (
	"testing"

	"github.com/mjolnir-mud/engine/pkg/reactor"
	"github.com/stretchr/testify/assert"
)

type testControllerForRegistry struct{}

func (l testControllerForRegistry) Name() string {
	return "login"
}

func (l testControllerForRegistry) Start(session reactor.Session) error {
	return nil
}

func (l testControllerForRegistry) Resume(session reactor.Session) error {
	return nil
}

func (l testControllerForRegistry) Stop(session reactor.Session) error {
	return nil
}

func (l testControllerForRegistry) HandleInput(session reactor.Session, input string) error {
	return nil
}

func TestControllerRegistry_RegisterAndGet(t *testing.T) {
	ControllerRegistry.Register(testControllerForRegistry{})

	assert.NotNil(t, ControllerRegistry.Get("login"))
}
