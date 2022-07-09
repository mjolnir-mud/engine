package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTest() {
	connectToRedis()
}

func teardownTest() {
	disconnectFromRedis()
}

func TestRedisSet(t *testing.T) {
	setupTest()
	defer teardownTest()
	err := RedisSet("test:test", "test", 0)
	assert.Nil(t, err)

	err = redisClient.Get(context.Background(), "test:test").Err()

	assert.Nil(t, err)
}
