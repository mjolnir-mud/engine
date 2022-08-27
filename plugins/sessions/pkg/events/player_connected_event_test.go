package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayerConnectedEvent_Topic(t *testing.T) {
	e := PlayerConnectedEvent{}
	assert.Equal(t, "player_connected", e.Topic())
}
