package logger

import "github.com/rs/zerolog/log"

var Logger = log.
	With().
	Str("plugin", "world").
	Logger()
