package testing

import "github.com/mjolnir-mud/engine"

func Teardown() {
	engine.RegisterOnEnvStopCallback("test", func() {
		_ = engine.RedisFlushAll()
	})
}
