package logger

import "github.com/mjolnir-mud/engine/pkg/logger"

var Intsance = logger.Instance.With().
	Str("plugin", "telnet_portal").
	Logger()
