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

package registry

import (
	controllers2 "github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/errors"
	errors2 "github.com/mjolnir-mud/engine/plugins/controllers/errors"
	"github.com/mjolnir-mud/engine/plugins/controllers/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/rs/zerolog"
)

var controllers map[string]controllers2.Controller
var log zerolog.Logger

func HandleInput(entityId string, line string) error {
	exists, err := ecs.EntityExists(entityId)

	if err != nil {
		return err
	}

	if !exists {
		return errors2.SessionNotFoundError{
			SessionId: entityId,
		}
	}

	cName, err := ecs.GetStringComponent(entityId, "controller")

	if err != nil {
		return err
	}

	c, err := Get(cName)

	if err != nil {
		return err
	}

	return c.HandleInput(entityId, line)
}

func Start() {
	log = logger.Instance.With().Str("service", "registry").Logger()
	controllers = make(map[string]controllers2.Controller, 0)
	log.Info().Msg("started")
}

func Stop() {
	log.Info().Msg("stopped")
}

func Register(c controllers2.Controller) {
	log.Info().Str("name", c.Name()).Msg("registering controller")
	controllers[c.Name()] = c
}

func Get(name string) (controllers2.Controller, error) {
	c, ok := controllers[name]

	if !ok {
		return nil, errors.ControllerNotFoundError{
			Name: name,
		}
	}

	return c, nil
}
