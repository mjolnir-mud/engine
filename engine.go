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
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rueian/rueidis"
	"os"
)

// Engine is an instance of the Mjolnir game engine.
type Engine struct {
	redis                rueidis.Client
	instanceId           string
	controllerRegistry   *controllerRegistry
	systemRegistry       *systemRegistry
	sessionRegistry      *sessionRegistry
	subscriptionRegistry *subscriptionRegistry
	dataSourceRegistry   *dataSourceRegistry
	pluginRegistry       *pluginRegistry
	config               *Configuration
	services             []string
	service              string

	// logger is the logger for the engine. Plugins can use this logger to create their own tagged loggers. See
	// [zerolog](https://github.com/rs/zerolog) for more information.
	logger zerolog.Logger
}

// New creates a new instance of the Mjolnir game engine. If the redis connection fails, an error is returned.
func New(config *Configuration) *Engine {
	if config.Environment == "" {
		env := os.Getenv("MJOLNIR_ENV")

		if env == "" {
			config.Environment = "development"
		} else {
			config.Environment = env
		}
	}

	if config.Mongo == nil {
		config.Mongo = &MongoConfiguration{
			Host:     "localhost",
			Port:     27017,
			Database: fmt.Sprintf("mjolnir_%s", config.Environment),
		}
	} else {
		if config.Mongo.Database == "" {
			config.Mongo.Database = fmt.Sprintf("mjolnir_%s", config.Environment)
		}

		if config.Mongo.Host == "" {
			config.Mongo.Host = "localhost"
		}

		if config.Mongo.Port == 0 {
			config.Mongo.Port = 27017
		}
	}

	e := &Engine{
		instanceId: config.InstanceId,
		config:     config,
		logger:     newLogger(config),
	}

	e.systemRegistry = newSystemRegistry(e)
	e.sessionRegistry = newSessionRegistry(e)
	e.subscriptionRegistry = newSubscriptionRegistry(e)
	e.controllerRegistry = newControllerRegistry(e)
	e.dataSourceRegistry = newDataSourceRegistry(e)
	e.pluginRegistry = newPluginRegistry(e)

	return e
}

// GetEnv returns the current environment.
func (e *Engine) GetEnv() string {
	return e.config.Environment
}

// GetService returns the current service.
func (e *Engine) GetService() string {
	return e.service
}

// RegisterService registers a service with the engine. Services individual processes that can be started and stopped
// independently of each other. The engine has a single default service called `engine`.
func (e *Engine) RegisterService(service string) {
	e.logger.Info().Str("service", service).Msg("registering service")
	e.services = append(e.services, service)
}

// RegisterSystem registers a system with the engine. System should implement the `System` interface.
func (e *Engine) RegisterSystem(system System) {
	e.systemRegistry.register(system)
}

// Start starts the Mjolnir game engine. If the redis connection fails, an error is returned.
func (e *Engine) Start(service string) {
	e.RegisterService("engine")

	if !e.hasService(service) {
		e.logger.Fatal().Str("service", service).Msg("unknown service")
		panic("unknown service")
	}

	e.service = service

	e.logger.Info().Msg("starting engine")
	redisClient, err := newRedisClient(e)

	if err != nil {
		e.logger.Fatal().Err(err).Msg("failed to connect to redis")
		panic(err)
	}

	e.redis = redisClient

	e.systemRegistry.start()
	e.sessionRegistry.start()
	e.controllerRegistry.start()
	e.dataSourceRegistry.start()

	// starts all the registered plugins
	e.pluginRegistry.start()

	e.service = service
}

// Stop stops the Mjolnir game engine.
func (e *Engine) Stop() {
	e.logger.Info().Msg("stopping engine")
	e.systemRegistry.stop()
	e.sessionRegistry.stop()
	e.controllerRegistry.stop()
	e.dataSourceRegistry.stop()

	// stops all the registered plugins
	e.pluginRegistry.stop()
}

func (e *Engine) hasService(name string) bool {
	for _, service := range e.services {
		if service == name {
			return true
		}
	}
	return false
}
