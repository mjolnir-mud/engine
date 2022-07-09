package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluginRegistry_Register(t *testing.T) {
	plugins.Register(createTestPlugin())

	assert.Equal(t, plugins.Count(), 1)
}
