package instance

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterBeforeServiceStopCallback(t *testing.T) {
	setup()
	defer teardown()

	RegisterBeforeServiceStopCallback("test", func() {})

	assert.Equal(t, 1, len(beforeServiceStopCallbacks))
}

func TestRegisterBeforeStopCallback(t *testing.T) {
	setup()
	defer teardown()

	RegisterBeforeStopCallback(func() {})

	assert.Equal(t, 1, len(beforeStopCallbacks))
}

func TestRegisterBeforeServiceStopCallbackForEnv(t *testing.T) {
	setup()
	defer teardown()

	RegisterBeforeServiceStopCallbackForEnv("test", "test", func() {})

	assert.Equal(t, 1, len(beforeServiceStopCallbacksForEnv["test"]["test"]))
}

func TestRegisterBeforeStopCallbackForEnv(t *testing.T) {
	setup()
	defer teardown()

	RegisterBeforeStopCallbackForEnv("test", func() {})

	assert.Equal(t, 1, len(beforeStopCallbacksForEnv["test"]))
}
