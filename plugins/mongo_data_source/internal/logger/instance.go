package logger

import "github.com/mjolnir-mud/engine/pkg/logger"

var Instance = logger.Instance.
	With().
	Str("plugin", "mongo_data_source").
	Logger()
