package redis

import (
	"fmt"
	"github.com/mjolnir-mud/engine/pkg/event"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type testEvent struct {
	Value string
}

func (e testEvent) Topic() string {

	return fmt.Sprintf("test.%s", e.Value)
}

func setup() {
	viper.Set("env", "test")
	Start()
}

func teardown() {
	Stop()
}

func TestSubscribe(t *testing.T) {
	setup()
	defer teardown()
	ch := make(chan interface{})

	s := NewSubscription(&testEvent{Value: "test"}, func(payload event.EventPayload) {
		p := &testEvent{}

		_ = payload.Unmarshal(p)

		ch <- p
	})

	err := Publish(&testEvent{
		Value: "test",
	})

	assert.Nil(t, err)

	v := <-ch

	assert.Equal(t, &testEvent{
		Value: "test",
	}, v)

	s.Stop()
}

func TestPSubscribe(t *testing.T) {
	setup()
	defer teardown()
	ch := make(chan interface{})

	s := NewPatternSubscription(&testEvent{Value: "*"}, func(payload event.EventPayload) {
		e := &testEvent{}
		_ = payload.Unmarshal(e)

		ch <- e
	})

	err := Publish(&testEvent{
		Value: "testing",
	})

	assert.Nil(t, err)

	v := <-ch

	assert.Equal(t, &testEvent{
		Value: "testing",
	}, v)

	s.Stop()
}
