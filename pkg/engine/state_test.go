package engine

import "testing"

func TestNewState(t *testing.T) {
	state := newState([]Plugin{})
	if state == nil {
		t.Error("Expected state to be non-nil")
	}
}
