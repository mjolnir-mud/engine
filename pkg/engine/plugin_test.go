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
	state := newState([]Plugin{&testPlugin{}})

	if err := LoadPlugins(state); err != nil {
		t.Error("Expected LoadPlugins to return nil error")
	}

	if !initCalled {
		t.Error("Expected Init to be called")
	}
}