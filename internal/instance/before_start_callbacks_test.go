package instance

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterBeforeServiceStartCallback(t *testing.T) {
	setup()
	defer teardown()

	RegisterBeforeServiceStartCallback("testing", func() {})

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

	RegisterBeforeServiceStartCallbackForEnv("testing", "testing", func() {})

	assert.Equal(t, 1, len(beforeServiceStartCallbacksForEnv["testing"]["testing"]))
}

func TestRegisterBeforeStartCallbackForEnv(t *testing.T) {
	setup()
	defer teardown()

	RegisterBeforeStartCallbackForEnv("testing", func() {})

	assert.Equal(t, 1, len(beforeStartCallbacksForEnv["testing"]))
}
