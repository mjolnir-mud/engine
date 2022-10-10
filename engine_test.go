package engine

import (
	"github.com/google/uuid"
	"github.com/mjolnir-mud/engine/internal/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	e := createEngineInstance()

	assert.NotNil(t, e)
}

func createEngineInstance() *Engine {
	puid, _ := uuid.NewRandom()
	prefix := puid.String()

	engine := New(&Configuration{
		Redis: &redis.Configuration{
			Host: "localhost",
			Port: 6379,
			DB:   1,
		},
		InstanceId:        prefix,
		DefaultController: "test",
	})

	_ = engine.Start()
	return engine
}
