package logger

import (
	"github.com/mjolnir-mud/engine/pkg/logger"
	"github.com/rs/zerolog"
)

var Instance zerolog.Logger

func Start() {
	Instance = logger.Instance.
		With().
		Str("plugin", "mongo_data_source").
		Logger()
}
