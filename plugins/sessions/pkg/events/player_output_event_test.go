package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayerOutputEvent_Topic(t *testing.T) {
	e := PlayerOutputEvent{Id: "test"}
	assert.Equal(t, "test.output", e.Topic())
}
