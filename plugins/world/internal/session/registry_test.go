package session

import (
	testing2 "testing"

	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/events"
	"github.com/stretchr/testify/assert"
)

func TestStartRegistry(t *testing2.T) {
	started := testing.Setup()
	defer testing.Teardown()

	ch := make(chan bool)

	<-started

	sub := engine.
		Subscribe(events.SessionRegistryStartedEvent{}, func(_ interface{}) {
			ch <- true
		})

	defer sub.Stop()
	defer StopRegistry()

	StartRegistry()

	<-ch

	assert.NotNil(t, assertSessionSubscription)
	assert.NotNil(t, sessionStoppedSubscription)
}

func TestStopRegistry(t *testing2.T) {
	started := testing.Setup()
	defer testing.Teardown()

	ch := make(chan bool)

	<-started

	sub := engine.
		Subscribe(events.SessionRegistryStartedEvent{}, func(_ interface{}) {
			ch <- true
		})

	StartRegistry()
	StopRegistry()
	<-ch

	sub.Stop()
}

func TestRegisterLineHandler(t *testing2.T) {
	started := testing.Setup()
	defer testing.Teardown()

	ch := make(chan bool)

	<-started

	sub := engine.
		Subscribe(events.SessionRegistryStartedEvent{}, func(_ interface{}) {
			ch <- true
		})

	defer sub.Stop()
	defer StopRegistry()

	RegisterLineHandler(func(id string, line string) error {
		return nil
	})
	StartRegistry()
	<-ch

	assert.Equal(t, len(lineHandlers), 1)
}

func TestRegisterSessionStartedHandler(t *testing2.T) {
	started := testing.Setup()
	defer testing.Teardown()

	ch := make(chan bool)

	<-started

	sub := engine.
		Subscribe(events.SessionRegistryStartedEvent{},
			func(_ interface{}) {
				ch <- true
			},
		)

	defer sub.Stop()
	defer StopRegistry()

	RegisterSessionStartedHandler(func(id string) error {
		return nil
	})
	StartRegistry()
	<-ch

	assert.Equal(t, len(sessionStartedHandlers), 1)
}

func TestRegisterSessionStoppedHandler(t *testing2.T) {
	started := testing.Setup()
	defer testing.Teardown()

	ch := make(chan bool)

	<-started

	sub := engine.
		Subscribe(events.SessionRegistryStartedEvent{}, func(_ interface{}) {
			ch <- true
		})

	defer sub.Stop()
	defer StopRegistry()

	RegisterSessionStoppedHandler(func(id string) error {
		return nil
	})
	StartRegistry()
	<-ch

	assert.Equal(t, len(sessionStoppedHandlers), 1)
}

func TestSendLine(t *testing2.T) {
	started := testing.Setup()
	defer testing.Teardown()

	ch := make(chan string)

	<-started

	StartRegistry()

	sub := engine.
		Subscribe(events.SendLineEvent{}, "test", func(event interface{}) {
			e := event.(*events.SendLineEvent)
			ch <- e.Line
		})

	defer sub.Stop()
	defer StopRegistry()

	SendLine("test", "test")
	line := <-ch
	assert.Equal(t, "test", line)
}

func TestSendLineF(t *testing2.T) {
	started := testing.Setup()
	defer testing.Teardown()

	ch := make(chan string)

	<-started

	StartRegistry()

	sub := engine.
		Subscribe(events.SendLineEvent{}, "test", func(event interface{}) {
			e := event.(*events.SendLineEvent)
			ch <- e.Line
		})

	defer sub.Stop()
	defer StopRegistry()

	SendLineF("test", "test%s", "ing")
	line := <-ch
	assert.Equal(t, "testing", line)
}