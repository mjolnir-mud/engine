package testing

import (
	"github.com/mjolnir-mud/engine"
	"github.com/spf13/viper"
)

func Setup() chan bool {
	viper.Set("env", "test")
	ch := make(chan bool)
	engine.RegisterAfterStartCallback(func() {
		go func() { ch <- true }()
	})
	engine.Start("test")

	return ch
}

func Teardown() {
	engine.Stop()
}
