package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayerDisconnectedEvent_Topic(t *testing.T) {
	e := PlayerDisconnectedEvent{}
	assert.Equal(t, "%s.disconnected", e.Topic("id"))
}

func TestPlayerDisconnectedEvent_Payload(t *testing.T) {
	e := &PlayerDisconnectedEvent{}
	assert.Equal(t, &PlayerDisconnectedEvent{}, e.Payload())

	e = &PlayerDisconnectedEvent{Id: "id"}
	assert.Equal(t, e, PlayerDisconnectedEvent{}.Payload("id"))
}
