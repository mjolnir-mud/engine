package controllers

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/controllers/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/controllers/internal/registry"
	"github.com/mjolnir-mud/engine/plugins/controllers/pkg/systems"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/sessions"
)

type plugin struct{}

func (p *plugin) Name() string {
	return "controllers"
}

func (p *plugin) Registered() error {
	engine.EnsureRegistered(ecs.Plugin.Name())

	engine.RegisterOnServiceStartCallback("world", func() {
		logger.Instance.Info().Msg("started")

		ecs.RegisterSystem(systems.ControllerSystem)

		sessions.RegisterLineHandler(func(entityId string, line string) error {
			return registry.HandleInput(entityId, line)
		})

	})

	engine.RegisterOnServiceStopCallback("world", func() {
		logger.Instance.Info().Msg("stopping controllers")
	})

	return nil
}

var Plugin = &plugin{}
