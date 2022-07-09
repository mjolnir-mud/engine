package plugin_registry

import (
	"testing"

	"github.com/mjolnir-mud/engine/test"
	"github.com/stretchr/testify/assert"
)

func TestPluginRegistry_Register(t *testing.T) {
	plugins.Register(test.CreateTestPlugin())

	assert.Equal(t, plugins.Count(), 1)
}
