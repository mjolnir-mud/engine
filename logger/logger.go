package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Instance zerolog.Logger

func Start() {
	Instance = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
		With().
		Timestamp().
		Logger().
		Level(zerolog.TraceLevel)
}
