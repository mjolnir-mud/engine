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
	engineErrors "github.com/mjolnir-engine/engine/errors"
	engineEvents "github.com/mjolnir-engine/engine/events"

	"github.com/mjolnir-engine/engine/uid"
	"github.com/rs/zerolog"
)

type controllerRegistry struct {
	controllers                map[string]Controller
	logger                     zerolog.Logger
	engine                     *Engine
	sessionStartedSubscription *uid.UID
}

func newControllerRegistry(engine *Engine) *controllerRegistry {
	return &controllerRegistry{
		controllers: make(map[string]Controller),
		logger:      engine.Logger.With().Str("component", "controller-registry").Logger(),
		engine:      engine,
	}
}

func (c *controllerRegistry) Register(controller Controller) {
	c.logger.Info().Str("controller", controller.Name()).Msg("registering controller")
	c.controllers[controller.Name()] = controller
}

func (c *controllerRegistry) Start() {
	c.logger.Info().Msg("starting")

	id := c.engine.PSubscribe(engineEvents.SessionStartedEvent{}, func(message EventMessage) {
		event := engineEvents.SessionStartedEvent{}
		err := message.Unmarshal(&event)

		if err != nil {
			c.logger.Error().Err(err).Msg("error unmarshalling event")
			return
		}

		controller, err := c.engine.GetSessionController(event.Id)

		if err != nil {
			c.logger.
				Error().
				Err(err).
				Str("sessionId", event.Id.String()).
				Msg("error getting session controller")
			return
		}

		context := &ControllerContext{
			SessionId: event.Id,
		}

		err = controller.Start(context)

		if err != nil {
			c.logger.Error().Err(err).Msg("error starting controller")
			return
		}
	})

	c.sessionStartedSubscription = id
}

func (r *controllerRegistry) Stop() {
	r.logger.Info().Msg("stopping")
	r.engine.Unsubscribe(r.sessionStartedSubscription)
}

func (c *controllerRegistry) Get(name string) (Controller, error) {
	controller, ok := c.controllers[name]
	if !ok {
		return nil, engineErrors.ControllerNotFoundError{Name: name}
	}
	return controller, nil
}

// RegisterController registers a controller with the engine. Controllers must implement the `Controller` interface. If
// a controller with the same name is already registered, it will be overwritten.
func (e *Engine) RegisterController(controller Controller) {
	e.controllerRegistry.Register(controller)
}

// GetController returns a controller by name. If the controller is not found, an error is returned.
func (e *Engine) GetController(name string) (Controller, error) {
	return e.controllerRegistry.Get(name)
}
