package engine

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewEngine(t *testing.T) {
	// Create a new engine instance
	engine := New(&Config{
		LogLevel: "debug",
	})

	// Verify that the instanceId is a valid UUID
	_, err := uuid.Parse(engine.instanceId.String())
	if err != nil {
		t.Errorf("Expected a valid UUID for instanceId, got error: %v", err)
	}

	// Verify that the logger is not nil
	if engine.logger == nil {
		t.Errorf("Expected logger to be initialized, but it's nil")
	}

	// Verify that the pluginManager is not nil
	if engine.pluginManager == nil {
		t.Errorf("Expected pluginManager to be initialized, but it's nil")
	}

	// test panic if log level is invalid
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when passing invalid log level, but no panic occurred")
		}
	}()

	New(&Config{
		LogLevel: "invalid",
	})
}
