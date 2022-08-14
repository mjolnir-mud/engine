package controller_registry

import (
	"github.com/mjolnir-mud/engine/plugins/world/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/controller"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/errors"
	"github.com/rs/zerolog"
)

var registry map[string]controller.Controller
var log zerolog.Logger

func Start() {
	log = logger.Instance.
		With().
		Str("service", "controller_registry").
		Logger()

	registry = make(map[string]controller.Controller)
}

func Register(c controller.Controller) {
	log.Info().
		Str("controller", c.Name()).
		Msg("registering controller")

	registry[c.Name()] = c
}

func Get(name string) (controller.Controller, error) {
	c, ok := registry[name]

	if !ok {
		return nil, errors.ControllerNotFoundError{
			Controller: name,
		}
	}

	return c, nil
}
