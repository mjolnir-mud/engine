package plugin_registry

import (
	"testing"

	"github.com/mjolnir-mud/engine/test"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setup() {
	viper.Set("env", "test")
}

func tearDown() {
}

func TestRegister(t *testing.T) {
	setup()
	defer tearDown()
	Register(test.CreateTestPlugin())
	assert.Equal(t, len(plugins), 1)
}
