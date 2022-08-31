package testing

import (
	"github.com/mjolnir-mud/engine"
)

func Setup() chan bool {
	engine.SetEnv("test")
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
