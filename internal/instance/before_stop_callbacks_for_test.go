package instance

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterBeforeServiceStopCallback(t *testing.T) {
	setup()
	defer teardown()

	RegisterBeforeServiceStopCallback("testing", func() {})

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

	RegisterBeforeServiceStopCallbackForEnv("testing", "testing", func() {})

	assert.Equal(t, 1, len(beforeServiceStopCallbacksForEnv["testing"]["testing"]))
}

func TestRegisterBeforeStopCallbackForEnv(t *testing.T) {
	setup()
	defer teardown()

	RegisterBeforeStopCallbackForEnv("testing", func() {})

	assert.Equal(t, 1, len(beforeStopCallbacksForEnv["testing"]))
}
