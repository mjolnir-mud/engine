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

package sessions

import (
	"github.com/mjolnir-mud/engine/plugins/sessions/internal/plugin"
	"github.com/mjolnir-mud/engine/plugins/sessions/internal/registry"
)

// RegisterSessionStartedHandler registers a handler that is called when a session is started.
func RegisterSessionStartedHandler(h func(id string) error) {
	registry.RegisterSessionStartedHandler(h)
}

// RegisterSessionStoppedHandler registers a handler that is called when a session is stopped.
func RegisterSessionStoppedHandler(h func(id string) error) {
	registry.RegisterSessionStoppedHandler(h)
}

// RegisterSessionLineHandler registers a handler that is called when a line is received from a session.
func RegisterReceiveLineHandler(h func(id string, line string) error) {
	registry.RegisterReceiveLineHandler(h)
}

func RegisterSendLineHandler(h func(id string, line string) error) {
	registry.RegisterSendLineHandler(h)
}

// StopSessionRegistry stops the session registry. This should only be called non-portal services.
func StopSessionRegistry() {
	registry.Stop()
}

// StartSessionRegistry starts the session registry. This should only be called non-portal services.
func StartSessionRegistry() {
	registry.Start()
}

var Plugin = plugin.Plugin
