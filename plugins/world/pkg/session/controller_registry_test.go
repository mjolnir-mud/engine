package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testControllerForRegistry struct{}

func (l testControllerForRegistry) Name() string {
	return "login"
}

func (l testControllerForRegistry) Start(session Session) error {
	return nil
}

func (l testControllerForRegistry) Resume(session Session) error {
	return nil
}

func (l testControllerForRegistry) Stop(session Session) error {
	return nil
}

func (l testControllerForRegistry) HandleInput(session Session, input string) error {
	return nil
}

func TestControllerRegistry_RegisterAndGet(t *testing.T) {
	ControllerRegistry.Register(testControllerForRegistry{})

	assert.NotNil(t, ControllerRegistry.Get("login"))
}
