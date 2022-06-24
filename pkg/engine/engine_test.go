package engine

import "testing"

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
	Init([]Plugin{&testPlugin{}})

	if !initCalled {
		t.Error("Expected Init to be called")
	}
}
