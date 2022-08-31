package testing

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/ecs"
)

func Setup() {
	engine.RegisterPlugin(ecs.Plugin)
	engine.RegisterOnEnvStartCallback("test", func() {
		_ = engine.RedisFlushAll()
	})
}
