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

package plugin

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/controllers/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/controllers/internal/registry"
	"github.com/mjolnir-mud/engine/plugins/controllers/pkg/systems"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/sessions"
)

type plugin struct{}

var Plugin = &plugin{}

func (p *plugin) Name() string {
	return "controllers"
}

func (p *plugin) Registered() error {
	engine.EnsureRegistered(ecs.Plugin.Name())
	engine.EnsureRegistered(sessions.Plugin.Name())

	engine.RegisterBeforeServiceStartCallback("world", func() {
		logger.Start()
		registry.Start()

		ecs.RegisterSystem(systems.ControllerSystem)

		sessions.RegisterLineHandler(func(entityId string, line string) error {
			return registry.HandleInput(entityId, line)
		})

		logger.Instance.Info().Msg("started")
	})

	engine.RegisterBeforeServiceStopCallback("world", func() {
		logger.Instance.Info().Msg("stopping")
		registry.Stop()
	})

	return nil
}
