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

package templates

import (
	"github.com/mjolnir-mud/engine/plugins/templates/internal/plugin"
	"github.com/mjolnir-mud/engine/plugins/templates/internal/template_registry"
	"github.com/mjolnir-mud/engine/plugins/templates/internal/theme_registry"
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/template"
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/theme"
)

// RegisterTheme registers a theme with the theme registry.
func RegisterTheme(t theme.Theme) {
	theme_registry.Register(t)
}

// RegisterTemplate registers a template with the template registry.
func RegisterTemplate(t template.Template) {
	template_registry.Register(t)
}

// GetTheme returns a theme with the given name. If the theme is not found, an error is returned.
func GetTheme(name string) (theme.Theme, error) {
	return theme_registry.GetTheme(name)
}

// RenderTemplate renders a template with the given name passing the given data to the template. If the template is not
// found, an error is returned.
func RenderTemplate(name string, ctx interface{}) (string, error) {
	return template_registry.Render(name, ctx)
}

var Plugin = plugin.Plugin
