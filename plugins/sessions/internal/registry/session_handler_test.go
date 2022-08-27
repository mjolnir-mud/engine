package registry

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/pkg/event"
	"github.com/mjolnir-mud/engine/plugins/sessions/pkg/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSessionHandler(t *testing.T) {
	setup()
	defer teardown()

	sh := NewSessionHandler("test")

	assert.Equal(t, "test", sh.Id)
}

func TestSessionHandler_Start(t *testing.T) {
	setup()
	defer teardown()

	ch := make(chan string)

	RegisterSessionStartedHandler(func(id string) error {
		go func() { ch <- id }()
		return nil
	})

	sh := NewSessionHandler("test")
	sh.Start()
	defer sh.Stop()

	v := <-ch

	assert.Equal(t, "test", v)
}

func TestSessionHandler_Stop(t *testing.T) {
	setup()
	defer teardown()

	ch := make(chan string)

	RegisterSessionStoppedHandler(func(id string) error {
		go func() { ch <- id }()
		return nil
	})

	sh := NewSessionHandler("test")
	sh.Start()
	sh.Stop()

	v := <-ch

	assert.Equal(t, "test", v)
	assert.Len(t, sessionHandlers, 0)
}

func TestSessionHandler_SendLine(t *testing.T) {
	setup()
	defer teardown()

	ch := make(chan string)

	sub := engine.Subscribe(events.PlayerOutputEvent{Id: "test"}, func(e event.EventPayload) {
		ev := &events.PlayerOutputEvent{}
		_ = e.Unmarshal(ev)

		go func() { ch <- ev.Line }()
	})

	defer sub.Stop()

	sh := NewSessionHandler("test")
	sh.Start()
	defer sh.Stop()

	err := sh.SendLine("test")

	assert.NoError(t, err)
	assert.Equal(t, "test", <-ch)

}

func TestLineHandler(t *testing.T) {
	setup()
	defer teardown()

	ch := make(chan string)

	RegisterLineHandler(func(id string, line string) error {
		go func() { ch <- line }()
		return nil
	})

	sh := NewSessionHandler("test")
	sh.Start()
	defer sh.Stop()

	err := engine.Publish(events.PlayerInputEvent{Id: "test", Line: "test"})

	assert.NoError(t, err)

	v := <-ch

	assert.Equal(t, "test", v)
}
