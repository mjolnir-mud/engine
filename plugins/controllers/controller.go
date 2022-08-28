package controllers

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/plugins/controllers/internal/logger"
)

type plugin struct{}

func (p *plugin) Name() string {
	return "controllers"
}

func (p *plugin) Registered() error {
	engine.RegisterOnServiceStartCallback("world", func() {
		logger.Instance.Info().Msg("started")
	})

	engine.RegisterOnServiceStopCallback("world", func() {
		logger.Instance.Info().Msg("stopping controllers")
	})

	return nil
}
