package pubsub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublish(t *testing.T) {
	setup()
	defer teardown()
	ch := make(chan interface{})

	sub := Subscribe("test", createTestEvent, func(payload interface{}) {
		ch <- payload
	})

	err := Publish("test", `{"Payload": "test"}`)

	assert.Nil(t, err)

	v := <-ch

	assert.Equal(t, "test", v.(*testEvent).Payload)

	sub.Stop()
}
