package engine

import "github.com/rs/zerolog/log"

type pluginRegistry struct {
	plugins []Plugin
}

func (r *pluginRegistry) Register(p Plugin) {
	pluginRegistryLogger.Info().Str("plugin", p.Name()).Msg("Registering plugin")

	r.plugins = append(r.plugins, p)
}

func (r *pluginRegistry) Count() int {
	return len(r.plugins)
}

var plugins = &pluginRegistry{
	plugins: []Plugin{},
}

var pluginRegistryLogger = log.
	With().
	Str("plugin", "engine").
	Str("service", "plugin_registry").
	Logger()
