package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testPlugin struct{}

var initCalled = false

func (p *testPlugin) Name() string {
	return "test"
}

func (p *testPlugin) Init(state *State) error {
	initCalled = true

	return nil
}

func TestLoadPlugins(t *testing.T) {
	Init("test", []Plugin{&testPlugin{}})

	assert.True(t, initCalled)

	assert.NotNilf(t, state.baseCommand, "baseCommand should not be nil")
}
