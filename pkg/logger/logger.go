package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var Instance zerolog.Logger

func Start() {
	if viper.GetString("env") == "production" {
		Instance = zerolog.New(os.Stdout).
			With().
			Timestamp().
			Logger().
			Level(zerolog.InfoLevel)

	} else {
		// set up with console writer
		Instance = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
			With().
			Timestamp().
			Logger().
			Level(zerolog.TraceLevel)
	}
}
