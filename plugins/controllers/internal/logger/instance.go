package logger

import (
	"github.com/mjolnir-mud/engine/logger"
	"github.com/rs/zerolog"
)

var Instance zerolog.Logger

func Start() {
	Instance = logger.Instance.
		With().
		Str("plugin", "controllers").
		Logger()
}
