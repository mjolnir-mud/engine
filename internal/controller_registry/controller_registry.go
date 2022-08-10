package controller_registry

import (
	"github.com/mjolnir-mud/engine/plugins/world/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/controller"
)

type controllerRegistry map[string]controller.Controller

// Start starts the controller registry.
func (cr *controllerRegistry) Start() {
}

// Get returns a controller from the registry.
func (cr *controllerRegistry) Get(name string) controller.Controller {
	log.Debug().Msgf("getting controller: %s", name)
	return ControllerRegistry[name]
}

// Register registers a controller with the registry.
func (cr *controllerRegistry) Register(controller controller.Controller) {
	log.Info().Msgf("registering controller: %s", controller.Name())
	ControllerRegistry[controller.Name()] = controller
}

var ControllerRegistry = controllerRegistry{}

var log = logger.Logger.
	With().
	Str("service", "controllerRegistry").
	Logger()
