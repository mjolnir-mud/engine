package plugin_registry

import (
	"github.com/mjolnir-mud/engine/pkg/plugin"
	"github.com/rs/zerolog/log"
)

type pluginRegistry struct {
	plugins []plugin.Plugin
}

func (r *pluginRegistry) Register(p plugin.Plugin) {
	pluginRegistryLogger.Info().Str("plugin", p.Name()).Msg("Registering plugin")

	r.plugins = append(r.plugins, p)
}

func (r *pluginRegistry) Count() int {
	return len(r.plugins)
}

var plugins = &pluginRegistry{
	plugins: []plugin.Plugin{},
}

var pluginRegistryLogger = log.
	With().
	Str("plugin", "engine").
	Str("service", "plugin_registry").
	Logger()

func Register(p plugin.Plugin) {
	plugins.Register(p)
}
