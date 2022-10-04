/*
 * Copyright (c) 2022 eightfivefour llc. All rights reserved.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package engine

import (
	"github.com/alecthomas/kong"
	redis2 "github.com/go-redis/redis/v9"
	"github.com/mjolnir-mud/engine/config"
	"github.com/mjolnir-mud/engine/event"
	"github.com/mjolnir-mud/engine/internal/cli"
	"github.com/mjolnir-mud/engine/internal/instance"
	"github.com/mjolnir-mud/engine/internal/plugin_registry"
	"github.com/mjolnir-mud/engine/internal/redis"
	"github.com/mjolnir-mud/engine/plugin"
)

// Subscription represents a pubsub to an event.
type Subscription interface {
	Stop()
}

// Cli returns and execute CLI commands.
func Cli() {
	ctx := kong.Parse(&cli.StartCmd{})

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

// ConfigureForEnv configures the engine for the provided environment.
func ConfigureForEnv(env string, cb func(cfg *config.Configuration) *config.Configuration) {
	instance.ConfigureForEnv(env, cb)
}

// SetEnv sets the environment for the engine. The environment defaults to "dev".
func SetEnv(env string) {
	instance.SetEnv(env)
}

// GetEnv returns the environment for the engine.
func GetEnv() string {
	return instance.GetEnv()
}

// GetGameName returns the name of the game.
func GetGameName() string {
	return instance.GetGameName()
}

// Running returns a channel that is closed when the engine is stopped.
func Running() chan bool {
	return instance.Running
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
// // event, even though it uses the underlying Redis PubSub command.
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
// [Redis PubSub PSubscribe](https://redis.io/commands/psubscribe). It should not be confused with the `RedisPSubscribe`
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

// RegisterAfterStartCallbackForEnv registers a callback function that is called after the engine is started, but only
// if the engine is running in the provided environment.
func RegisterAfterStartCallbackForEnv(env string, callback func()) {
	instance.RegisterAfterStartCallbackForEnv(env, callback)
}

// RegisterAfterStopCallback registers a callback function that is called after the engine is stopped.
func RegisterAfterStopCallback(callback func()) {
	instance.RegisterAfterStopCallback(callback)
}

// RegisterAfterStopCallbackForEnv registers a callback function that is called after the engine is stopped, but only
// if the engine is running in the provided environment.
func RegisterAfterStopCallbackForEnv(env string, callback func()) {
	instance.RegisterAfterStopCallbackForEnv(env, callback)
}

// RegisterAfterServiceStartCallback registers a callback function that is called after the service is started.
func RegisterAfterServiceStartCallback(service string, callback func()) {
	instance.RegisterAfterServiceStartCallback(service, callback)
}

// RegisterBeforeStartCallback registers a callback function that is called before the engine is started.
func RegisterBeforeStartCallback(callback func()) {
	instance.RegisterBeforeStartCallback(callback)
}

// RegisterBeforeStartCallbackForEnv registers a callback function that is called before the engine is started, but only
// if the provided environment matches the current environment.
func RegisterBeforeStartCallbackForEnv(env string, callback func()) {
	instance.RegisterBeforeStartCallbackForEnv(env, callback)
}

// RegisterBeforeStopCallback registers a callback function that is called before the engine is stopped.
func RegisterBeforeStopCallback(callback func()) {
	instance.RegisterBeforeStopCallback(callback)
}

// RegisterBeforeStopCallbackForEnv registers a callback function that is called before the engine is stopped, but only
// if the provided environment matches the current environment.
func RegisterBeforeStopCallbackForEnv(env string, callback func()) {
	instance.RegisterBeforeStopCallbackForEnv(env, callback)
}

// RegisterBeforeServiceStartCallback registers a callback function that is called when the engine is started.
func RegisterBeforeServiceStartCallback(service string, callback func()) {
	instance.RegisterBeforeServiceStartCallback(service, callback)
}

// RegisterBeforeServiceStartCallbackForEnv registers a callback function that is called when the engine is started, but
// only if the provided environment matches the current environment.
func RegisterBeforeServiceStartCallbackForEnv(service string, env string, callback func()) {
	instance.RegisterBeforeServiceStartCallbackForEnv(service, env, callback)
}

// RegisterBeforeServiceStopCallback registers a callback function that is called when the engine is stopped.
func RegisterBeforeServiceStopCallback(service string, callback func()) {
	instance.RegisterBeforeServiceStopCallback(service, callback)
}

// RegisterBeforeServiceStopCallbackForEnv registers a callback function that is called when the engine is stopped, but
// only if the provided environment matches the current environment.
func RegisterBeforeServiceStopCallbackForEnv(service string, env string, callback func()) {
	instance.RegisterBeforeServiceStopCallbackForEnv(service, env, callback)
}

// RegisterPlugin registers a plugin with the engine. Plugins need to be registered before the engine is started, but
// after the engine is initialized. Plugins should conform to the plugin data_source.
func RegisterPlugin(plugin plugin.Plugin) {
	plugin_registry.Register(plugin)
}

// Initialize is the first function called when the engine is started. It initializes the engine, and should be called
// before any other engine functions.
func Initialize(gameName string, env string) {
	instance.Initialize(gameName, env)
}

// Start starts the engine. It should be called after the engine is initialized, and after all plugins are registered.
func Start() {
	instance.Start()
}

// Stop	stops the engine.
func Stop() {
	instance.Stop()
}

// Subscribe subscribes to an event on the message bus. It accepts an event and a callback function that will be
// called when the event is published. The callback function will be passed an EventPayload, with which `Unmarshal` can
// be called to unmarshal the payload into a Go struct. PSubscribe uses
// [Redis PubSub PSubscribe](https://redis.io/commands/subscribe). It should not be confused with the `RedisSubscribe`
// function, which simply subscribes to a channel on Redis, without wiring any of the underlying event handling.
func Subscribe(e event.Event, cb func(event.EventPayload)) Subscription {
	return redis.NewSubscription(e, cb)
}
