package registry

import (
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/plugins/controllers/internal/logger"
	"github.com/mjolnir-mud/plugins/controllers/pkg/controller"
	"github.com/mjolnir-mud/plugins/controllers/pkg/errors"
	"github.com/rs/zerolog"
)

var controllers map[string]controller.Controller
var log zerolog.Logger

func HandleInput(entityId string, line string) error {
	exists, err := ecs.EntityExists(entityId)

	if err != nil {
		return err
	}

	if !exists {
		return errors.SessionNotFoundError{
			SessionId: entityId,
		}
	}

	cName, err := ecs.GetStringComponent(entityId, "controller")

	if err != nil {
		return err
	}

	c, err := Get(cName)

	if err != nil {
		return err
	}

	return c.HandleInput(entityId, line)
}

func Start() {
	log = logger.Instance.With().Str("service", "registry").Logger()
	controllers = make(map[string]controller.Controller, 0)
	log.Info().Msg("started")
}

func Stop() {}

func Register(c controller.Controller) {
	log.Info().Str("name", c.Name()).Msg("registering controller")
	controllers[c.Name()] = c
}

func Get(name string) (controller.Controller, error) {
	c, ok := controllers[name]

	if !ok {
		return nil, errors.ControllerNotFoundError{
			Name: name,
		}
	}

	return c, nil
}
