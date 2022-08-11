package engine

import (
	redis2 "github.com/go-redis/redis/v9"
	"github.com/mjolnir-mud/engine/internal/instance"
	"github.com/mjolnir-mud/engine/internal/plugin_registry"
	"github.com/mjolnir-mud/engine/internal/redis"
	"github.com/mjolnir-mud/engine/pkg/event"
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
func Ping() error {
	return redis.Ping()
}

// PSubscribe subscribes to an event on the message bus. It accepts an event and an arbitrary number of arguments that
// will be passed to the event's topic constructor. The last argument should be a callback function otherwise the engine
// will panic. The callback function is called when a message is received on the topic, with a payload decoded against
// the event's Payload constructor. PSubscribe accepts topic patterns, and will subscribe to all matching topics. See
// the [Redis documentation](https://redis.io/topics/pubsub) for more information.
func PSubscribe(e event.Event, args ...interface{}) Subscription {
	return redis.PSubscribe(e, args...)
}

// Publish publishes an event on the message bus. It accepts an event and an arbitrary number of arguments that will be
// passed to the event's topic and payload constructors.
func Publish(e event.Event, args ...interface{}) error {
	return redis.Publish(e, args...)
}

// RegisterAfterStartCallback registers a callback function that is called after the engine is started.
func RegisterAfterStartCallback(callback func()) {
	instance.RegisterAfterStartCallback(callback)
}

// RegisterAfterStopCallback registers a callback function that is called after the engine is stopped.
func RegisterAfterStopCallback(callback func()) {
	instance.RegisterAfterStopCallback(callback)
}

// RegisterBeforeStartCallback registers a callback function that is called before the engine is started.
func RegisterBeforeStartCallback(callback func()) {
	instance.RegisterBeforeStartCallback(callback)
}

// RegisterBeforeStopCallback registers a callback function that is called before the engine is stopped.
func RegisterBeforeStopCallback(callback func()) {
	instance.RegisterBeforeStopCallback(callback)
}

// RegisterPlugin registers a plugin with the engine. Plugins need to be registered before the engine is started, but
// after the engine is initialized. Plugins should conform to the plugin interface.
func RegisterPlugin(plugin plugin.Plugin) {
	plugin_registry.Register(plugin)
}

// Start starts with the provided game name.
func Start(gameName string) {
	instance.Start(gameName)
}

// Stop	stops the engine.
func Stop() {
	instance.Stop()
}

// Subscribe subscribes to an event on the message bus. It accepts an event and an arbitrary number of arguments that
// will be passed to the event's topic constructor. The last argument should be a callback function otherwise the engine
// will panic. The callback function is called when a message is received on the topic, with a payload decoded against
// the event's Payload constructor. If it is wanted to subscribe against a pattern, the `PSubscribe` method should be
// used instead. See the [Redis documentation](https://redis.io/topics/pubsub) for more information.
func Subscribe(e event.Event, args ...interface{}) Subscription {
	return redis.Subscribe(e, args...)
}

// SetEnv sets the environment for the engine. Mjolnir recognizes three different environments by default, development
// test, and production. The environment is set by setting the `MJOLNIR_ENV` environment variable.
func SetEnv(env string) {
	viper.Set("env", env)
}

// Redis returns a client to the Redis server. This can be used to interact with the Redis server directly.
var Redis *redis2.Client
