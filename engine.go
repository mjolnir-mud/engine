package engine

import (
	"context"
	"fmt"

	redis2 "github.com/go-redis/redis/v9"
	"github.com/mjolnir-mud/engine/internal/nats"
	"github.com/mjolnir-mud/engine/internal/plugin_registry"
	"github.com/mjolnir-mud/engine/internal/pubsub"
	"github.com/mjolnir-mud/engine/internal/redis"
	"github.com/mjolnir-mud/engine/pkg/plugin"
	nats2 "github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type engine struct {
	name        string
	baseCommand *cobra.Command
}

// Subscription represents a pubsub to an event.
type Subscription interface {
	Stop()
}

var e = &engine{}

func Start(name string) {
	viper.SetEnvPrefix("MJOLNIR")
	err := viper.BindEnv("env")

	if err != nil {
		panic(err)
	}

	viper.SetDefault("env", "development")

	e.name = name

	e.baseCommand = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("manage the %s Mjolnir Game", name),
		Long:  fmt.Sprintf("manage the %s Mjolnir Game", name),
	}

	logger.Info().Str("plugin", "engine").Msgf("initializing engine for game %s", name)
	redis.Start()

	Redis = redis.GetClient()

	nats.Start()
	plugin_registry.StartPlugins()

	err = e.baseCommand.Execute()

	if err != nil {
		panic(err)
	}
}

func Stop() {
	logger.Info().Str("plugin", "engine").Msg("shutting down engine")
	nats.Stop()
	redis.Stop()
}

// Subscribe subscribes to an event on the message bus. It accepts a topic, an event constructor, and a callback
// function. The event constructor is used to create an event object that is passed to the callback function. The event
// should be a pointer to a struct that can be umarshalled from JSON. The callback function is called when a message is
// received on the topic. When the pubsub is no longer needed, it should be stopped by calling the `Stop` method.
func Subscribe(topic string, event func() interface{}, callback func(interface{})) Subscription {
	return pubsub.Subscribe(topic, event, callback)
}

// PSubscribe subscribes to a pattern on the message bus. It accepts a topic, an event constructor, and a callback
// function. The event constructor is used to create an event object that is passed to the callback function. The event
// should be a pointer to a struct that can be umarshalled from JSON. The callback function is called when a message is
// received on the topic. When the pubsub is no longer needed, it should be stopped by calling the `Stop` method.
// Unlike Subscribe, PSubscribe will match any topic that matches the pattern. For more information, see the Redis
// documentation on patterns [https://redis.io/commands/psubscribe](https://redis.io/commands/psubscribe).
func PSubscribe(topic string, event func() interface{}, callback func(interface{})) Subscription {
	return pubsub.PSubscribe(topic, event, callback)
}

// Publish publishes a message to the message bus. It accepts a topic and a message payload. The message payload should
// be a pointer to a struct that can be marshalled to JSON.
func Publish(topic string, payload interface{}) error {
	return Redis.Publish(context.Background(), topic, payload).Err()
}

// SetEnv sets the environment for the engine.
func SetEnv(env string) {
	viper.Set("env", env)
}

// RegisterPlugin registers a plugin with the engine.
func RegisterPlugin(plugin plugin.Plugin) {
	plugin_registry.Register(plugin)
}

func PublishEvent(event string, data interface{}) error {
	return nats.PublishEvent(event, data)
}

func SubscribeToEvent(event string, handler nats2.Handler) (*nats2.Subscription, error) {
	return nats.SubscribeToEvent(event, handler)
}

func AddCLICommand(command *cobra.Command) {
	e.baseCommand.AddCommand(command)
}

var Redis *redis2.Client

var logger = log.With().Str("plugin", "engine").Logger()
