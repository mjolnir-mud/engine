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
	err := RedisSet("test:get", "test", 0)
	assert.Nil(t, err)

	err = redisClient.Get(context.Background(), "test:get").Err()

	assert.Nil(t, err)
}

func TestRedisGet(t *testing.T) {
	setupTest()
	defer teardownTest()
	err := RedisSet("test:set", "test", 0)
	assert.Nil(t, err)

	var val string

	err = RedisGet("test:set", &val)

	assert.Nil(t, err)
	assert.Equal(t, "test", val)
}
