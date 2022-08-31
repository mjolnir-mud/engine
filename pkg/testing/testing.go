package testing

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/internal/instance"
)

func Setup(services ...string) chan bool {
	engine.SetEnv("test")
	ch := make(chan bool)
	engine.RegisterAfterStartCallback(func() {
		go func() { ch <- true }()
	})
	engine.Start("test")

	if len(services) > 0 {
		for _, service := range services {
			instance.StartService(service)
		}
	}

	return ch
}

func Teardown() {
	engine.Stop()
}
