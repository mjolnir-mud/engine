package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	Start()

	assert.NotNil(t, client)

	err := Ping()

	assert.Nil(t, err)

	Stop()
}
