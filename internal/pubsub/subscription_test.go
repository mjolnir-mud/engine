package pubsub

import (
	"context"
	"testing"

	"github.com/mjolnir-mud/engine/internal/redis"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type testEvent struct {
	Payload string
}

func createTestEvent() interface{} {
	return &testEvent{}
}

func setup() {
	viper.Set("env", "test")
	redis.Start()
}

func teardown() {
	redis.Stop()
}

func TestSubscribe(t *testing.T) {
	setup()
	defer teardown()
	ch := make(chan interface{})

	s := Subscribe("test", createTestEvent, func(payload interface{}) {
		ch <- payload
	})

	redis.GetClient().Publish(context.Background(), "test", `{"Payload": "test"}`)

	v := <-ch

	assert.Equal(t, "test", v.(*testEvent).Payload)

	s.Stop()
}

func TestPSubscribe(t *testing.T) {
	setup()
	defer teardown()
	ch := make(chan interface{})

	s := PSubscribe("test:*", createTestEvent, func(payload interface{}) {
		ch <- payload
	})

	redis.GetClient().Publish(context.Background(), "test:foo*", `{"Payload": "test"}`)

	v := <-ch

	assert.Equal(t, "test", v.(*testEvent).Payload)

	s.Stop()
}
