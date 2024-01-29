package engine

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"os"
)

// Engine is the core Mjolnir MUD engine, providing all the basic functionality for a Mjolnir based game.
type Engine struct {
	instanceId uuid.UUID
	logger     *zerolog.Logger
	config     *Config
	plugins    map[string]Plugin
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

	return &Engine{instanceId, &logger, config, make(map[string]Plugin)}
}

// NewContext creates a new context to be used by various components of the engine.
func (engine *Engine) NewContext() context.Context {
	return context.WithValue(context.Background(), "engineInstanceId", engine.instanceId)
}
