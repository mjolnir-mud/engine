package world

import (
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/world/internal/session"
)

type world struct {
	stop chan bool
}

func (w world) Name() string {
	return "world"
}

func (w world) Registered() error {
	session.RegisterSessionStartedHandler(func(id string) error {
		err := ecs.AddEntityWithID(id, "session", map[string]interface{}{})

		if err != nil {
			return err
		}

		return nil
	})

	return nil
}

var Plugin = world{
	stop: make(chan bool),
}
