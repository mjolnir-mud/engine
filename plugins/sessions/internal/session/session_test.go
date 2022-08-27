package session

//
//import (
//	"testing"
//
//	"github.com/mjolnir-mud/engine"
//	testing2 "github.com/mjolnir-mud/engine/pkg/testing"
//	"github.com/mjolnir-mud/engine/plugins/world/pkg/events"
//	"github.com/stretchr/testify/assert"
//)
//
//type lineReceived struct {
//	Id   string
//	Line string
//}
//
//func TestNew(t *testing.T) {
//	s := New(&events.AssertSessionEvent{
//		UUID: "test",
//	})
//
//	assert.NotNil(t, s)
//	assert.Equal(t, "test", s.id)
//	assert.NotNil(t, s.logger)
//}
//
//func TestSession_Start(t *testing.T) {
//	started := testing2.Setup()
//	defer testing2.Teardown()
//
//	<-started
//	sessionHandlerCalled := make(chan string)
//	lineHandlerCalled := make(chan lineReceived)
//
//	RegisterSessionStartedHandler(func(id string) error {
//		go func() { sessionHandlerCalled <- id }()
//		return nil
//	})
//
//	RegisterLineHandler(func(id string, line string) error {
//		go func() { lineHandlerCalled <- lineReceived{id, line} }()
//		return nil
//	})
//
//	s := New(&events.AssertSessionEvent{
//		UUID: "test",
//	})
//
//	s.Start()
//	defer s.Stop()
//
//	id := <-sessionHandlerCalled
//
//	assert.Equal(t, "test", id)
//
//	err := engine.Publish(events.InputEvent{}, "test", "test")
//
//	assert.NoError(t, err)
//
//	lh := <-lineHandlerCalled
//
//	assert.Equal(t, "test", lh.Id)
//	assert.Equal(t, "test", lh.Line)
//}
//
//func TestSession_Stop(t *testing.T) {
//	started := testing2.Setup()
//	defer testing2.Teardown()
//
//	<-started
//	sessionHandlerCalled := make(chan string)
//	stopEventCalled := make(chan string)
//
//	s := New(&events.AssertSessionEvent{
//		UUID: "test",
//	})
//
//	s.Start()
//
//	RegisterSessionStoppedHandler(func(id string) error {
//		go func() { sessionHandlerCalled <- id }()
//		return nil
//	})
//
//	engine.Subscribe(events.SessionStoppedEvent{}, "test", func(e interface{}) {
//		event := e.(*events.SessionStoppedEvent)
//
//		go func() { stopEventCalled <- event.UUID }()
//	})
//
//	s.Stop()
//
//	id := <-sessionHandlerCalled
//
//	assert.Equal(t, "test", id)
//
//	uuid := <-stopEventCalled
//
//	assert.Equal(t, "test", uuid)
//}
