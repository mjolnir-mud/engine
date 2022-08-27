package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayerDisconnectedEvent_Topic(t *testing.T) {
	e := PlayerDisconnectedEvent{Id: "test"}
	assert.Equal(t, "test.disconnected", e.Topic())
}
