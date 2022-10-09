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
	engineErrors "github.com/mjolnir-mud/engine/errors"
	"github.com/rs/zerolog"
)

type controllerRegistry struct {
	controllers map[string]Controller
	logger      zerolog.Logger
	engine      *Engine
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

// RegisterController registers a controller with the engine. Controllers must implement the `Controller` interface. If
// a controller with the same name is already registered, it will be overwritten.
func (e *Engine) RegisterController(controller Controller) {
	e.controllerRegistry.Register(controller)
}

// GetController returns a controller by name. If the controller is not found, an error is returned.
func (e *Engine) GetController(name string) (Controller, error) {
	if controller, ok := e.controllerRegistry.controllers[name]; ok {
		return controller, nil
	}
	return nil, engineErrors.ControllerNotFoundError{
		Name: name,
	}
}
