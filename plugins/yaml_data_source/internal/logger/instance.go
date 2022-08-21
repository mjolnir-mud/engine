package logger

import "github.com/mjolnir-mud/engine/pkg/logger"

var Instance = logger.Instance.
	With().
	Str("plugin", "yaml_data_source").
	Logger()
