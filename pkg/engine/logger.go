package engine

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

}

func GetLoggerForPlugin(plugin string) zerolog.Logger {
	return log.With().Str("plugin", plugin).Logger()
}
