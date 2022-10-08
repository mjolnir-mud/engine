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

package instance

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/mjolnir-mud/engine/internal/plugin_registry"
	engineRedis "github.com/mjolnir-mud/engine/internal/redis"
	"github.com/rs/zerolog"
)

var environment string
var gameName string
var cfg *engine.Configuration
var Running chan bool

func SetEnv(n string) {
	environment = n
}

func GetEnv() string {
	return environment
}

func GetGameName() string {
	return gameName
}

func IsRunning() bool {
	return engineRedis.Ping() == nil
}

func Initialize(name string, env string) {
	environment = env
	gameName = name

	initializeConfigurations()

	initializeBeforeStartCallbacks()
	initializeAfterStartCallbacks()
	initializeBeforeStopCallbacks()
	initializeAfterStopCallbacks()

	plugin_registry.Initialize()
	logger.Start()
}

func Start() {
	log = logger.Instance.With().Str("component", "engine").Logger()
	callBeforeStartCallbacks()

	env := GetEnv()

	callBeforeStartCallbacksForEnv(env)

	cfg = callConfigureForEnv(env)
	engineRedis.Start(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Db)

	redisReady := make(chan bool)

	go func() {
		for {
			if engineRedis.Ping() == nil {
				redisReady <- true
				break
			}
		}
	}()

	<-redisReady

	plugin_registry.Start()

	callAfterStartCallbacks()
	callAfterStartCallbacksForEnv(env)
}

func StopService(service string) {
	log.Info().Str("service", service).Msg("stopping service")
	callBeforeServiceStopCallbacks(service)
	callBeforeServiceStopCallbacksForEnv(service, GetEnv())

	callAfterServiceStopCallbacks(service)
	callAfterServiceStopCallbacksForEnv(service, GetEnv())
}

func StartService(service string) {
	log.Info().Str("service", service).Msg("starting service")

	Running = make(chan bool)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		StopService(service)
		Stop()
		Running <- true
	}()

	callBeforeServiceStartCallbacks(service)
	callBeforeServiceStartCallbacksForEnv(service, GetEnv())

	callAfterServiceStartCallbacks(service)
	callAfterServiceStartCallbacksForEnv(service, GetEnv())
}

func Stop() {
	log.Info().Msg("stopping engine")
	callBeforeStopCallbacks()

	callBeforeStopCallbacksForEnv(GetEnv())

	plugin_registry.Stop()
	engineRedis.Stop()

	callAfterStopCallbacks()
	callAfterStopCallbacksForEnv(GetEnv())
}

var log zerolog.Logger
