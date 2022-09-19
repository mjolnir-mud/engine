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
	"github.com/mjolnir-mud/engine/plugins/templates/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/templates/internal/template_registry"
	"github.com/mjolnir-mud/engine/plugins/templates/internal/theme_registry"
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/default_theme"
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/template"
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/theme"
)

type templatePlugin struct {
	themes    map[string]theme.Theme
	templates map[string]template.Template
}

func (p templatePlugin) Name() string {
	return "templates"
}

func (p templatePlugin) Registered() error {
	engine.RegisterAfterServiceStartCallback("world", func() {
		logger.Start()
		theme_registry.Start()
		template_registry.Start()

		theme_registry.Register(default_theme.Theme)

	})

	engine.RegisterBeforeServiceStopCallback("world", func() {
		theme_registry.Stop()
		template_registry.Stop()
	})

	return nil
}

var Plugin = &templatePlugin{}