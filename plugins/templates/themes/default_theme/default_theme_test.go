package default_theme

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultTheme_Name(t *testing.T) {
	assert.Equal(t, "default", Theme.Name())
}
