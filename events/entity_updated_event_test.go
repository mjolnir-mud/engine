package events

import (
	"fmt"
	"testing"

	"github.com/mjolnir-engine/engine/uid"
	"github.com/stretchr/testify/assert"
)

func TestEntityUpdatedEvent_Topic(t *testing.T) {
	event := EntityUpdatedEvent{
		Id: uid.New(),
	}

	assert.Equal(t, fmt.Sprintf("entity.%s.updated", event.Id), event.Topic())
}

func TestEntityUpdatedEvent_AllTopics(t *testing.T) {
	event := EntityUpdatedEvent{
		Id: uid.New(),
	}

	assert.Equal(t, "entity.*.updated", event.AllTopics())
}
