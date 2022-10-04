package logger

import (
	"github.com/mjolnir-mud/engine/logger"
)

var Intsance = logger.Instance.With().
	Str("plugin", "telnet_portal").
	Logger()
