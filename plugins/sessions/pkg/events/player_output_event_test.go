package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayerOutputEvent_Topic(t *testing.T) {
	e := PlayerOutputEvent{Id: "testing"}
	assert.Equal(t, "testing.output", e.Topic())
}
