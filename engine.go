package engine

import (
	"context"

	redis2 "github.com/go-redis/redis/v9"
	"github.com/mjolnir-mud/engine/internal/instance"
	"github.com/mjolnir-mud/engine/internal/plugin_registry"
	"github.com/mjolnir-mud/engine/internal/pubsub"
	"github.com/mjolnir-mud/engine/internal/redis"
	"github.com/mjolnir-mud/engine/pkg/plugin"
	"github.com/spf13/viper"
)

// Subscription represents a pubsub to an event.
type Subscription interface {
	Stop()
}

// EnsureRegistered ensures that the plugin is registered with the engine. If the plugin is not registered, an the
// engine will panic. This should be used by plugins that need to ensure that another plugin is registered before
// they can start.
func EnsureRegistered(pluginName string) {
	plugin_registry.EnsureRegistered(pluginName)
}

// Ping pings the Redis server. This is a direct pass-through to the Redis client, simply setting the context.
func Ping() *redis2.StatusCmd {
	return Redis.Ping(context.Background())
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
	return pubsub.Publish(topic, payload)
}

// Initialize should be called by every games `main` function as the first step in the game's initialization. Once the
// engine is initialized, various plugins can be registered and the engine can be started.
func Initialize(name string) {
	instance.Initialize(name)
	Redis = redis.Client
}

// RegisterPlugin registers a plugin with the engine. Plugins need to be registered before the engine is started, but
// after the engine is initialized. Plugins should conform to the plugin interface.
func RegisterPlugin(plugin plugin.Plugin) {
	plugin_registry.Register(plugin)
}

// Start starts the engine. This should be called after all plugins are registered and all other initialization is
// complete.
func Start() {
	instance.Start()
}

// Stop	stops the engine.
func Stop() {
	instance.Stop()
}

// Subscribe subscribes to an event on the message bus. It accepts a topic, an event constructor, and a callback
// function. The event constructor is used to create an event object that is passed to the callback function. The event
// should be a pointer to a struct that can be umarshalled from JSON. The callback function is called when a message is
// received on the topic. When the pubsub is no longer needed, it should be stopped by calling the `Stop` method.
func Subscribe(topic string, event func() interface{}, callback func(interface{})) Subscription {
	return pubsub.Subscribe(topic, event, callback)
}

// SetEnv sets the environment for the engine. Mjolnir recognizes three different environments by default, development
// test, and production. The environment is set by setting the `MJOLNIR_ENV` environment variable.
func SetEnv(env string) {
	viper.Set("env", env)
}

// Redis returns a client to the Redis server. This can be used to interact with the Redis server directly.
var Redis *redis2.Client
