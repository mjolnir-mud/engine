package engine

import (
	"fmt"
	"os"
	"time"

	"github.com/mjolnir-mud/engine/internal/nats"
	"github.com/mjolnir-mud/engine/internal/plugin_registry"
	"github.com/mjolnir-mud/engine/internal/redis"
	"github.com/mjolnir-mud/engine/pkg/plugin"
	nats2 "github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type engine struct {
	name        string
	baseCommand *cobra.Command
}

var e = &engine{}

func Start(name string) {
	viper.SetEnvPrefix("MJOLNIR")
	err := viper.BindEnv("env")

	if err != nil {
		panic(err)
	}

	viper.SetDefault("env", "development")

	if viper.GetString("env") == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	e.name = name

	e.baseCommand = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("manage the %s Mjolnir Game", name),
		Long:  fmt.Sprintf("manage the %s Mjolnir Game", name),
	}

	logger.Info().Str("plugin", "engine").Msgf("initializing engine for game %s", name)
	nats.Start()
	redis.Start()
	plugin_registry.InitPlugins()

	err = e.baseCommand.Execute()

	if err != nil {
		panic(err)
	}
}

func Stop() {
	logger.Info().Str("plugin", "engine").Msg("shutting down engine")
	redis.Stop()
	nats.Stop()
}

func RegisterPlugin(plugin plugin.Plugin) {
	plugin_registry.Register(plugin)
}

func RedisSet(key string, value interface{}, expiration time.Duration) error {
	return redis.Set(key, value, expiration)
}

func RedisGet(key string, value interface{}) error {
	return redis.Get(key, value)
}

func RedisKeyExists(key string) (bool, error) {
	return redis.KeyExists(key)
}

func RedisDelete(key string) error {
	return redis.Delete(key)
}

func RedisCount(key string) (int64, error) {
	return redis.Count(key)
}

func RedisDeleteAll(key string) error {
	return redis.DeleteAll(key)
}

func PublishEvent(event string, data interface{}) error {
	return nats.PublishEvent(event, data)
}

func SubscribeToEvent(event string, handler nats2.Handler) (*nats2.Subscription, error) {
	return nats.SubscribeToEvent(event, handler)
}

func AddCommand(command *cobra.Command) {
	e.baseCommand.AddCommand(command)
}

var logger = log.With().Str("plugin", "engine").Logger()
