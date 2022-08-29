package logger

import "github.com/rs/zerolog/log"

var Instance = log.
	With().
	Str("plugin", "ecs").
	Logger()
