package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInputEvent_Topic(t *testing.T) {
	e := PlayerInputEvent{Id: "testing"}
	assert.Equal(t, "testing.input", e.Topic())
}
