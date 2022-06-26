package engine

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if os.Getenv("MJOLNIR_ENV") == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
}

func GetLoggerForPlugin(plugin string) zerolog.Logger {
	return log.With().Str("plugin", plugin).Logger()
}
