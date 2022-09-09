package testing

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/ecs/internal/plugin"
)

func Setup() {
	engine.RegisterPlugin(plugin.Plugin)
	engine.RegisterOnEnvStartCallback("test", func() {
		_ = engine.RedisFlushAll()
		plugin.RegisterEntityType(TestEntityType{})
		plugin.RegisterSystem(TestSystem{})
	})
}
