package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	Start()

	err := client.Ping(context.Background()).Err()
	assert.Nil(t, err)
}

func TestStop(t *testing.T) {
	Start()
	Stop()

	err := client.Ping(context.Background()).Err()

	assert.NotNil(t, err)
}
