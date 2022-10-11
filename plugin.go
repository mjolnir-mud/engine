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

// Plugin is the interface that must be implemented by a Mjolnir plugin.
type Plugin interface {
	// Name returns the name of the plugin. Plugin names must be unique. If two plugins are registered with the same
	// name, an error will be returned.
	Name() string

	// Start is called when the plugin is started. This is where the plugin should initialize itself and start any
	// of its own components. This method is called after the rest of the Mjolnir engine has been started. It will be
	// passed the `Engine` instance that it is running in as a parameter.
	Start(engine *Engine) error

	// Stop is called when the plugin is stopped. This is where the plugin should stop any of its own components. This
	// method should not return until all components have been stopped. It will be passed the `Engine` instance that
	// it is running in as a parameter.
	Stop(engine *Engine) error
}
