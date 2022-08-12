package plugin_registry

import (
	"github.com/mjolnir-mud/engine/pkg/logger"
	"github.com/mjolnir-mud/engine/pkg/plugin"
)

var plugins = []plugin.Plugin{}
var log = logger.Instance.With().Str("service", "pluginRegistry").Logger()

func EnsureRegistered(pluginName string) {
	for _, p := range plugins {
		if p.Name() == pluginName {
			return
		}
	}

	log.Fatal().Str("plugin", pluginName).Msg("plugin not registered")
	panic("plugin not registered")
}

func Register(p plugin.Plugin) {
	log.Info().Str("plugin", p.Name()).Msg("registering plugin")
	plugins = append(plugins, p)

	log.Debug().Str("plugin", p.Name()).Msg("calling registered callback")
	err := p.Registered()

	if err != nil {
		log.Fatal().Err(err).Msg("error calling registered callback")
		panic(err)
	}
}

func Start() {}

func Stop() {}
