package nats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testMessage struct {
	Test string `json:"test"`
}

func TestPublishAndSubscribeToEvent(t *testing.T) {
	Start()
	defer Stop()

	cont := make(chan bool)

	s, err := SubscribeToEvent("test", func(msg *testMessage) {
		assert.Equal(t, "test", msg.Test)
		cont <- true
	})

	defer s.Unsubscribe()

	assert.Nil(t, err)

	err = PublishEvent("test", &testMessage{Test: "test"})

	assert.Nil(t, err)

	<-cont
}
