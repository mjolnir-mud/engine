package session

import (
	testing2 "testing"

	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/events"
	"github.com/stretchr/testify/assert"
)

type lineReceived struct {
	Id   string
	Line string
}

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

//func TestSendLine(t *testing2.T) {
//	started := testing.Setup()
//	defer testing.Teardown()
//
//	ch := make(chan string)
//
//	<-started
//
//	sub := engine.
//		Subscribe(events.SendLineTopic("test"), events.SendLine, func(event interface{}) {
//			ch <- true
//		})
//
//	defer sub.Stop()
//	defer StopRegistry()
//
//	SendLine("test", "test")
//	<-ch
//}
