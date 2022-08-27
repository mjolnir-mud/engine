package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInputEvent_Topic(t *testing.T) {
	e := PlayerInputEvent{Id: "test"}
	assert.Equal(t, "test.input", e.Topic())
}
