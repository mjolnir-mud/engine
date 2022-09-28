package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayerDisconnectedEvent_Topic(t *testing.T) {
	e := PlayerDisconnectedEvent{Id: "testing"}
	assert.Equal(t, "testing.disconnected", e.Topic())
}
