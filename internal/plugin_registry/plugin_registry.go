package plugin_registry

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/mjolnir-mud/engine/pkg/logger"
	"github.com/mjolnir-mud/engine/pkg/plugin"
	"github.com/rs/zerolog"
)

var pluginsNames mapset.Set[string]
var plugins []plugin.Plugin
var log zerolog.Logger

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
	if pluginsNames.Contains(p.Name()) {
		return
	}

	log.Info().Str("plugin", p.Name()).Msg("registering plugin")
	pluginsNames.Add(p.Name())
	plugins = append(plugins, p)

	log.Debug().Str("plugin", p.Name()).Msg("calling registered callback")
	err := p.Registered()

	if err != nil {
		log.Fatal().Err(err).Msg("error calling registered callback")
		panic(err)
	}
}

func Initialize() {
	pluginsNames = mapset.NewSet[string]()
	plugins = make([]plugin.Plugin, 0)
	log = logger.Instance.With().Str("component", "plugin_registry").Logger()
}

func Start() {
	log.Info().Msg("starting")
}

func Stop() {
	log.Info().Msg("stopping")
	for _, p := range plugins {
		log.Info().Str("plugin", p.Name()).Msg("stopping plugin")
	}
}
