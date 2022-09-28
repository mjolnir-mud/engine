package instance

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterAfterServiceStopCallback(t *testing.T) {
	setup()
	defer teardown()

	RegisterAfterServiceStopCallback("testing", func() {})

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

	RegisterAfterServiceStopCallbackForEnv("testing", "testing", func() {})

	assert.Equal(t, 1, len(afterServiceStopCallbacksForEnv["testing"]["testing"]))
}

func TestRegisterAfterStopCallbackForEnv(t *testing.T) {
	setup()
	defer teardown()

	RegisterAfterStopCallbackForEnv("testing", func() {})

	assert.Equal(t, 1, len(afterStopCallbacksForEnv["testing"]))
}
