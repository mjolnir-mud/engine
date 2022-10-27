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

type pluginRecord struct {
	plugin Plugin
	failed bool
}

type pluginRegistry struct {
	plugins []*pluginRecord
	engine  *Engine
	logger  zerolog.Logger
}

func newPluginRegistry(engine *Engine) *pluginRegistry {
	return &pluginRegistry{
		plugins: make([]*pluginRecord, 0),
		engine:  engine,
		logger:  engine.logger.With().Str("component", "plugin_registry").Logger(),
	}
}

func (r *pluginRegistry) start() {
	for _, record := range r.plugins {
		r.logger.Info().Str("name", record.plugin.Name()).Msg("starting plugin")
		err := record.plugin.Start(r.engine)
		if err != nil {
			r.logger.Error().Str("name", record.plugin.Name()).Err(err).Msg("failed to start plugin")
			record.failed = true
		}
	}
}

func (r *pluginRegistry) stop() {
	for _, record := range r.plugins {
		if record.failed {
			continue
		}

		r.logger.Info().Str("name", record.plugin.Name()).Msg("stopping plugin")
		err := record.plugin.Stop(r.engine)
		if err != nil {
			r.logger.Error().Str("name", record.plugin.Name()).Err(err).Msg("failed to stop plugin")
		}
	}
}

func (r *pluginRegistry) register(plugin Plugin) {
	if r.hasPlugin(plugin.Name()) {
		r.logger.Warn().Str("plugin", plugin.Name()).Msg("plugin already registered")
		return
	}

	r.logger.Debug().Str("name", plugin.Name()).Msg("registering plugin")

	err := plugin.Init(r.engine)

	if err != nil {
		r.logger.Error().Str("name", plugin.Name()).Err(err).Msg("failed to initialize plugin")
		panic(err)
	}

	r.plugins = append(r.plugins, &pluginRecord{
		plugin: plugin,
		failed: false,
	})
}

func (r *pluginRegistry) hasPlugin(name string) bool {
	for _, record := range r.plugins {
		if record.plugin.Name() == name {
			return true
		}
	}

	return false
}

// RegisterPlugin registers a plugin with the engine.
func (e *Engine) RegisterPlugin(plugin Plugin) {
	e.pluginRegistry.register(plugin)
}
