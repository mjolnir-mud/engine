package controllers

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/plugins/controllers/internal/logger"
	"github.com/mjolnir-mud/plugins/controllers/pkg/systems"
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
	})

	engine.RegisterOnServiceStopCallback("world", func() {
		logger.Instance.Info().Msg("stopping controllers")
	})

	return nil
}
