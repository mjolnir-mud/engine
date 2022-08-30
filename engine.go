package engine

import (
	redis2 "github.com/go-redis/redis/v9"
	"github.com/mjolnir-mud/engine/internal/instance"
	"github.com/mjolnir-mud/engine/internal/plugin_registry"
	"github.com/mjolnir-mud/engine/internal/redis"
	"github.com/mjolnir-mud/engine/pkg/event"
	"github.com/mjolnir-mud/engine/pkg/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Subscription represents a pubsub to an event.
type Subscription interface {
	Stop()
}

// EnsureRegistered ensures that the plugin is registered with the engine. If the plugin is not registered, and the
// engine will panic. This should be used by plugins that need to ensure that another plugin is registered before
// they can start.
func EnsureRegistered(pluginName string) {
	plugin_registry.EnsureRegistered(pluginName)
}

// RedisPing pings the Redis server. This is a direct pass-through to the Redis client, simply setting the context.
func RedisPing() error {
	return redis.Ping()
}

// RedisFlushAll flushes all Redis databases.
func RedisFlushAll() error {
	return redis.FlushAll()
}

// RedisGet returns the value for the provided key from the Redis database.
func RedisGet(key string) *redis2.StringCmd {
	return redis.Get(key)
}

// RedisSet sets the value for the provided key in the Redis database.
func RedisSet(key string, value interface{}) *redis2.StatusCmd {
	return redis.Set(key, value)
}

// RedisExists returns if the provided key exists in the Redis database.
func RedisExists(key string) *redis2.IntCmd {
	return redis.Exists(key)
}

// RedisDel deletes the provided key from the Redis database.
func RedisDel(key string) *redis2.IntCmd {
	return redis.Del(key)
}

// RedisHGet returns the value for the provided map key and map key from the Redis database.
func RedisHGet(key string, mapKey string) *redis2.StringCmd {
	return redis.HGet(key, mapKey)
}

// RedisHGetAll returns the values for the provided map key from the Redis database.
func RedisHGetAll(key string) *redis2.MapStringStringCmd {
	return redis.HGetAll(key)
}

// RedisHExists returns true if the provided map key exists in the provided map key.
func RedisHExists(key string, mapKey string) *redis2.BoolCmd {
	return redis.HExists(key, mapKey)
}

// RedisHSet sets the value for the provided map key and map key in the Redis database.
func RedisHSet(key string, mapKey string, value interface{}) *redis2.IntCmd {
	return redis.HSet(key, mapKey, value)
}

// RedisHMSet sets the values for the provided map key from the Redis database.
func RedisHMSet(key string, values map[string]interface{}) *redis2.BoolCmd {
	return redis.HMSet(key, values)
}

// RedisKeys returns the keys for the provided pattern from the Redis database.
func RedisKeys(pattern string) *redis2.StringSliceCmd {
	return redis.Keys(pattern)
}

// RedisSMembers returns the members of the set for the provided key.
func RedisSMembers(key string) *redis2.StringSliceCmd {
	return redis.SMembers(key)
}

// RedisSIsMember returns true if the provided member is a member of the set for the provided key.
func RedisSIsMember(key string, value interface{}) *redis2.BoolCmd {
	return redis.SIsMember(key, value)
}

// RedisSAdd adds the provided members to the set for the provided key.
func RedisSAdd(key string, value interface{}) *redis2.IntCmd {
	return redis.SAdd(key, value)
}

// RedisSRem removes the provided members from the set for the provided key.
func RedisSRem(key string, value interface{}) *redis2.IntCmd {
	return redis.SRem(key, value)
}

// RedisSubscribe subscribes to a channel on Redis. This should not be confused with PSubscribe, which subscribes to an
//// event, even though it uses the underlying Redis PubSub command.
func RedisSubscribe(channels ...string) *redis2.PubSub {
	return redis.Subscribe(channels...)
}

// RedisPSubscribe pattern subscribes to a channel on Redis. This should not be confused with PSubscribe, which
// subscribes to an event, even though it uses the underlying Redis PubSub command.
func RedisPSubscribe(channels ...string) *redis2.PubSub {
	return redis.PSubscribe(channels...)
}

// PSubscribe subscribes to an event on the message bus. It accepts an event and a callback function that will be
// called when the event is published. The callback function will be passed an EventPayload, with which `Unmarshal` can
// be called to unmarshal the payload into a Go struct. PSubscribe uses
//[Redis PubSub PSubscribe](https://redis.io/commands/psubscribe). It should not be confused with the `RedisPSubscribe`
// function, which simply subscribes to a channel on Redis, without wiring any of the underlying event handling.
func PSubscribe(e event.Event, cb func(payload event.EventPayload)) Subscription {
	return redis.NewPatternSubscription(e, cb)
}

// Publish publishes an event on the message bus. It accepts an event and an arbitrary number of arguments that will be
// passed to the event's topic and payload constructors.
func Publish(e interface{}) error {
	return redis.Publish(e)
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

// RegisterOnServiceStartCallback registers a callback function that is called when the engine is started.
func RegisterOnServiceStartCallback(service string, callback func()) {
	instance.RegisterOnServiceStartCallback(service, callback)
}

// RegisterOnServiceStopCallback registers a callback function that is called when the engine is stopped.
func RegisterOnServiceStopCallback(service string, callback func()) {
	instance.RegisterOnServiceStopCallback(service, callback)
}

// RegisterCLICommand registers a CLI command with the engine. The command will be available in the CLI when calling
// the compiled binary.
func RegisterCLICommand(command *cobra.Command) {
	instance.RegisterCLICommand(command)
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

// Subscribe subscribes to an event on the message bus. It accepts an event and a callback function that will be
// called when the event is published. The callback function will be passed an EventPayload, with which `Unmarshal` can
// be called to unmarshal the payload into a Go struct. PSubscribe uses
//[Redis PubSub PSubscribe](https://redis.io/commands/subscribe). It should not be confused with the `RedisSubscribe`
// function, which simply subscribes to a channel on Redis, without wiring any of the underlying event handling.
func Subscribe(e event.Event, cb func(event.EventPayload)) Subscription {
	return redis.NewSubscription(e, cb)
}

// SetEnv sets the environment for the engine. Mjolnir recognizes three different environments by default, development
// test, and production. The environment is set by setting the `MJOLNIR_ENV` environment variable.
func SetEnv(env string) {
	viper.Set("env", env)
}
