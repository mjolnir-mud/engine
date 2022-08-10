package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndStart(t *testing.T) {
	CreateAndStart()

	assert.NotNil(t, Client)

	err := Client.Ping(context.Background()).Err()

	assert.Nil(t, err)

	Stop()
}
