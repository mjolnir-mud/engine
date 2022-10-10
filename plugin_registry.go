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

import "github.com/rs/zerolog"

type pluginRegistry struct {
	plugins map[string]Plugin
	engine  *Engine
	logger  zerolog.Logger
}

func newPluginRegistry(engine *Engine) *pluginRegistry {
	return &pluginRegistry{
		plugins: make(map[string]Plugin),
		engine:  engine,
		logger:  engine.logger.With().Str("component", "plugin_registry").Logger(),
	}
}

func (r *pluginRegistry) start() {
}

func (r *pluginRegistry) stop() {
}

func (r *pluginRegistry) register(plugin Plugin) {
	r.logger.Debug().Str("name", plugin.Name()).Msg("registering plugin")
	r.plugins[plugin.Name()] = plugin
}

// RegisterPlugin registers a plugin with the engine.
func (r *Engine) RegisterPlugin(plugin Plugin) {
	r.pluginRegistry.register(plugin)
}
