package plugin_registry

import (
	"github.com/mjolnir-mud/engine/testing/fakes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setup() {
	viper.Set("env", "testing")
}

func tearDown() {
}

func TestRegister(t *testing.T) {
	setup()
	defer tearDown()

	Start()
	defer Stop()

	Register(fakes.CreateFakePlugin())
	assert.Equal(t, len(plugins), 1)
}
