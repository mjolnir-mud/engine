package testing

import (
	"github.com/mjolnir-mud/engine"
	"github.com/spf13/viper"
)

func Setup(beforeStart func()) {
	viper.Set("env", "test")
	engine.Initialize("test")
	beforeStart()
	engine.Start()
}

func Teardown() {
	engine.Stop()
}
