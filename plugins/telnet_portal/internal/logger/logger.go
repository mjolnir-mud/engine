package logger

import "github.com/rs/zerolog/log"

var Logger = log.
	With().
	Str("plugin", "telnet_portal").
	Logger()
