package engine

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Start("test")

	assert.Equal(t, "development", viper.GetString("env"))
}
