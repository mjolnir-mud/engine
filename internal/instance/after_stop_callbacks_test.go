package instance

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterAfterServiceStopCallback(t *testing.T) {
	setup()
	defer teardown()

	RegisterAfterServiceStopCallback("test", func() {})

	assert.Equal(t, 1, len(afterServiceStopCallbacks))
}

func TestRegisterAfterStopCallback(t *testing.T) {
	setup()
	defer teardown()

	RegisterAfterStopCallback(func() {})

	assert.Equal(t, 1, len(afterStopCallbacks))
}

func TestRegisterAfterServiceStopCallbackForEnv(t *testing.T) {
	setup()
	defer teardown()

	RegisterAfterServiceStopCallbackForEnv("test", "test", func() {})

	assert.Equal(t, 1, len(afterServiceStopCallbacksForEnv["test"]["test"]))
}

func TestRegisterAfterStopCallbackForEnv(t *testing.T) {
	setup()
	defer teardown()

	RegisterAfterStopCallbackForEnv("test", func() {})

	assert.Equal(t, 1, len(afterStopCallbacksForEnv["test"]))
}
