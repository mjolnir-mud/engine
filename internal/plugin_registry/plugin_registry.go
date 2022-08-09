package plugin_registry

import (
	"github.com/mjolnir-mud/engine/internal/logger"
	"github.com/mjolnir-mud/engine/pkg/plugin"
)

var plugins = []plugin.Plugin{}
var log = logger.Instance.With().Str("service", "pluginRegistry").Logger()

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

func Start() {
	log.Info().Msg("starting plugins")

	for _, p := range plugins {
		log.Info().Msgf("initializing plugin %s", p.Name())
		err := p.Start()

		if err != nil {
			log.Fatal().Err(err).Msg("error initializing plugin")
			panic(err)
		}
	}
}

func Stop() {
	log.Info().Msg("stopping plugins")

	for _, p := range plugins {
		log.Info().Msgf("stopping plugin %s", p.Name())
		err := p.Stop()

		if err != nil {
			log.Fatal().Err(err).Msg("error stopping plugin")
			continue
		}
	}
}
