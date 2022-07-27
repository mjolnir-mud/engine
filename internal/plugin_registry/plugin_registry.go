package plugin_registry

import (
	"github.com/mjolnir-mud/engine/pkg/plugin"
	"github.com/rs/zerolog/log"
)

type pluginRegistry struct {
	plugins []plugin.Plugin
}

func (r *pluginRegistry) Register(p plugin.Plugin) {
	logger.Info().Str("plugin", p.Name()).Msg("registering plugin")

	r.plugins = append(r.plugins, p)
}

func (r *pluginRegistry) Count() int {
	return len(r.plugins)
}

func (r *pluginRegistry) StartPlugins() {
	logger.Info().Msg("initializing plugins")

	for _, p := range r.plugins {
		logger.Info().Msgf("initializing plugin %s", p.Name())
		err := p.Start()

		if err != nil {
			logger.Fatal().Err(err).Msg("error initializing plugin")
		}
	}
}

var plugins = &pluginRegistry{
	plugins: []plugin.Plugin{},
}

var logger = log.
	With().
	Str("plugin", "engine").
	Str("service", "plugin_registry").
	Logger()

func Register(p plugin.Plugin) {
	plugins.Register(p)
}

func StartPlugins() {
	plugins.StartPlugins()
}
