package logger

import (
	"github.com/mjolnir-mud/engine/logger"
)

var Instance = logger.Instance.
	With().
	Str("plugin", "yaml_data_source").
	Logger()
