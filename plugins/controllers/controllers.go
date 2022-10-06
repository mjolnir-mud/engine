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

package controllers

import (
	"github.com/mjolnir-mud/engine/plugins/controllers/controller"
	"github.com/mjolnir-mud/engine/plugins/controllers/internal/plugin"
	"github.com/mjolnir-mud/engine/plugins/controllers/internal/registry"
	"github.com/mjolnir-mud/engine/plugins/ecs"
)

var Plugin = plugin.Plugin

// Set sets the controller for the provided entity
func Set(entityId string, controllerName string) error {
	return ecs.AddOrUpdateStringComponentToEntity(entityId, "controller", controllerName)
}

// Register registers a controller with the registry. If a controller with the same name already exists, it will be
// overwritten.
func Register(controller controller.Controller) {
	registry.Register(controller)
}
