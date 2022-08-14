package redis

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type testEvent struct {
	Value string
}

func (e testEvent) Topic(args ...interface{}) string {
	id := args[0].(string)

	return fmt.Sprintf("test.%s", id)
}

func (e testEvent) Payload(args ...interface{}) interface{} {
	if len(args) > 0 {
		id := args[0].(string)

		return &testEvent{
			Value: id,
		}
	}

	return &testEvent{}
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

	s := NewSubscription(&testEvent{}, "test", func(payload interface{}) {
		ch <- payload
	})

	err := Publish(&testEvent{}, "test")

	assert.Nil(t, err)

	v := <-ch

	assert.Equal(t, &testEvent{
		Value: "test",
	}, v.(*testEvent))

	s.Stop()
}

func TestPSubscribe(t *testing.T) {
	setup()
	defer teardown()
	ch := make(chan interface{})

	s := NewPatternSubscription(&testEvent{}, "*", func(payload interface{}) {
		ch <- payload
	})

	err := Publish(&testEvent{}, "testing")

	assert.Nil(t, err)

	v := <-ch

	assert.Equal(t, &testEvent{
		Value: "testing",
	}, v.(*testEvent))

	s.Stop()
}
