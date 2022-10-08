package registry

import (
	"github.com/mjolnir-mud/engine"
	events2 "github.com/mjolnir-mud/engine/plugins/sessions/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSessionHandler(t *testing.T) {
	setup()
	defer teardown()

	sh := NewSessionHandler("testing")

	assert.Equal(t, "testing", sh.Id)
}

func TestSessionHandler_Start(t *testing.T) {
	setup()
	defer teardown()

	ch := make(chan string)

	RegisterSessionStartedHandler(func(id string) error {
		go func() { ch <- id }()
		return nil
	})

	sh := NewSessionHandler("testing")
	sh.Start()
	defer sh.Stop()

	v := <-ch

	assert.Equal(t, "testing", v)
}

func TestSessionHandler_Stop(t *testing.T) {
	setup()
	defer teardown()

	ch := make(chan string)

	RegisterSessionStoppedHandler(func(id string) error {
		go func() { ch <- id }()
		return nil
	})

	sh := NewSessionHandler("testing")
	sh.Start()
	sh.Stop()

	v := <-ch

	assert.Equal(t, "testing", v)
	assert.Len(t, sessionHandlers, 0)
}

func TestSessionHandler_SendLine(t *testing.T) {
	setup()
	defer teardown()

	ch := make(chan string)

	sub := engine.Subscribe(events2.PlayerOutputEvent{Id: "testing"}, func(e engine.EventPayload) {
		ev := &events2.PlayerOutputEvent{}
		_ = e.Unmarshal(ev)

		go func() { ch <- ev.Line }()
	})

	defer sub.Stop()

	sh := NewSessionHandler("testing")
	sh.Start()
	defer sh.Stop()

	err := sh.SendLine("testing")

	assert.NoError(t, err)
	assert.Equal(t, "testing\r\n", <-ch)

}

func TestLineHandler(t *testing.T) {
	setup()
	defer teardown()

	ch := make(chan string)

	RegisterReceiveLineHandler(func(id string, line string) error {
		go func() { ch <- line }()
		return nil
	})

	sh := NewSessionHandler("testing")
	sh.Start()
	defer sh.Stop()

	err := engine.Publish(events2.PlayerInputEvent{Id: "testing", Line: "testing"})

	assert.NoError(t, err)

	v := <-ch

	assert.Equal(t, "testing", v)
}
