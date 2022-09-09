package instance

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterBeforeServiceStartCallback(t *testing.T) {
	setup()
	defer teardown()

	RegisterBeforeServiceStartCallback("test", func() {})

	assert.Equal(t, 1, len(beforeServiceStartCallbacks))
}

func TestRegisterBeforeStartCallback(t *testing.T) {
	setup()
	defer teardown()

	RegisterBeforeStartCallback(func() {})

	assert.Equal(t, 1, len(beforeStartCallbacks))
}

func TestRegisterBeforeServiceStartCallbackForEnv(t *testing.T) {
	setup()
	defer teardown()

	RegisterBeforeServiceStartCallbackForEnv("test", "test", func() {})

	assert.Equal(t, 1, len(beforeServiceStartCallbacksForEnv["test"]["test"]))
}

func TestRegisterBeforeStartCallbackForEnv(t *testing.T) {
	setup()
	defer teardown()

	RegisterBeforeStartCallbackForEnv("test", func() {})

	assert.Equal(t, 1, len(beforeStartCallbacksForEnv["test"]))
}
