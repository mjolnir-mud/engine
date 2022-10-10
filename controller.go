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

import "github.com/mjolnir-mud/engine/uid"

// ControllerContext is a context that is passed to a controller.
type ControllerContext struct {
	SessionId *uid.UID
}

func newControllerContext(sessionId *uid.UID) *ControllerContext {
	return &ControllerContext{
		SessionId: sessionId,
	}
}

// Controller is the data_source for a session controller. A session controller handles interactions from the player
// session to the game world
type Controller interface {
	// Name returns the name of the controller. If multiple controllers of the same name are registered with the world
	// then the last one registered will be used. This enables the developer to override specific controllers with their
	// own implementation.
	Name() string

	// Start is called when the controller is set.
	Start(context *ControllerContext) error

	// Resume called when the world restarts, causing the portal to reset-assert the session.
	Resume(context *ControllerContext) error

	// Stop is called when the controller is unset.
	Stop(context *ControllerContext) error

	// HandleInput is called when the player sends input to the world.
	HandleInput(context *ControllerContext, input string) error
}
