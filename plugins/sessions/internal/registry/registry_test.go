package registry

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/sessions/pkg/events"
	"github.com/stretchr/testify/assert"
	"testing"

	testing2 "github.com/mjolnir-mud/engine/pkg/testing"
)

func setup() {
	testing2.Setup()
	Start()
}

func teardown() {
	Stop()
	testing2.Teardown()
}

func TestPlayerConnected(t *testing.T) {
	setup()
	defer teardown()

	ch := make(chan bool)

	RegisterSessionStartedHandler(func(id string) error {
		go func() { ch <- true }()
		return nil
	})

	err := engine.Publish(events.PlayerConnectedEvent{})

	assert.NoError(t, err)

	<-ch

	assert.Len(t, sessionHandlers, 1)
}

func TestRegisterSessionStartedHandler(t *testing.T) {
	setup()
	defer teardown()

	RegisterSessionStartedHandler(func(id string) error {
		return nil
	})

	assert.Len(t, sessionStartedHandlers, 1)
}

func TestRegisterSessionStoppedHandler(t *testing.T) {
	setup()
	defer teardown()

	RegisterSessionStoppedHandler(func(id string) error {
		return nil
	})

	assert.Len(t, sessionStoppedHandlers, 1)
}

func TestRegisterLineHandler(t *testing.T) {
	setup()
	defer teardown()

	RegisterLineHandler(func(id string, line string) error {
		return nil
	})

	assert.Len(t, lineHandlers, 1)
}
