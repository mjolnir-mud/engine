package engine

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"os"
)

// Engine is the core Mjolnir MUD engine, providing all the basic functionality for a Mjolnir based game.
type Engine struct {
	instanceId    uuid.UUID
	logger        *zerolog.Logger
	pluginManager *pluginManager
	config        *Config
}

// New creates a new instance of the Mjolnir MUD engine.
func New(config *Config) *Engine {
	logLevel, err := zerolog.ParseLevel(config.LogLevel)

	if err != nil {
		panic(err)
	}

	instanceId := uuid.New()
	logger := zerolog.New(os.Stdout).
		Level(logLevel).
		With().
		Timestamp().
		Str("instanceId", instanceId.String()).
		Str("service", "engine").
		Logger()

	return &Engine{instanceId, &logger, newPluginManager(), config}
}
