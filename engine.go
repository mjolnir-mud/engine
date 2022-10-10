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
	"github.com/rs/zerolog"
	"github.com/rueian/rueidis"
)

// Engine is an instance of the Mjolnir game engine.
type Engine struct {
	redis                rueidis.Client
	instanceId           string
	controllerRegistry   *controllerRegistry
	systemRegistry       *systemRegistry
	sessionRegistry      *sessionRegistry
	subscriptionRegistry *subscriptionRegistry
	config               *Configuration

	// Logger is the logger for the engine. Plugins can use this logger to create their own tagged loggers. See
	// [zerolog](https://github.com/rs/zerolog) for more information.
	Logger zerolog.Logger
}

// New creates a new instance of the Mjolnir game engine. If the redis connection fails, an error is returned.
func New(config *Configuration) *Engine {
	e := &Engine{
		instanceId: config.InstanceId,
		config:     config,
		Logger:     newLogger(config.Log),
	}

	e.systemRegistry = newSystemRegistry(e)
	e.sessionRegistry = newSessionRegistry(e)
	e.subscriptionRegistry = newSubscriptionRegistry(e)
	e.controllerRegistry = newControllerRegistry(e)

	return e
}

// RegisterSystem registers a system with the engine. System should implement the `System` interface.
func (e *Engine) RegisterSystem(system System) {
	e.systemRegistry.Register(system)
}

// Start starts the Mjolnir game engine. If the redis connection fails, an error is returned.
func (e *Engine) Start() error {
	e.Logger.Info().Msg("starting engine")
	redisClient, err := newRedisClient(e)

	if err != nil {
		e.Logger.Fatal().Err(err).Msg("failed to connect to redis")
		return err
	}
	e.redis = redisClient

	e.systemRegistry.Start()
	e.sessionRegistry.Start()
	e.controllerRegistry.Start()

	return nil
}

// Stop stops the Mjolnir game engine.
func (e *Engine) Stop() {
	e.Logger.Info().Msg("stopping engine")
	e.systemRegistry.Stop()
}
