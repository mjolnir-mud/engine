package engine

// Config is the configuration for the Mjolnir MUD engine. This is passed to a new instance of the engine on creation.
type Config struct {
	// LogLevel is the log level to use for the engine. Valid values are "debug", "info", "warn", "error", "fatal",
	// and "panic". Defaults to "info". The log level can also be set via the MJOLNIR_LOG_LEVEL environment variable.
	// see [github.com/rs/zerolog/log.Level] for more information.
	LogLevel string
}
