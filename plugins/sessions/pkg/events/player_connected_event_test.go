package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayerConnectedEvent_Payload(t *testing.T) {
	e := PlayerConnectedEvent{}
	assert.Equal(t, PlayerConnectedEvent{}, e.Payload())
}

func TestPlayerConnectedEvent_Topic(t *testing.T) {
	e := PlayerConnectedEvent{}
	assert.Equal(t, "player_connected", e.Topic())
}
