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

package systems

import (
	"github.com/mjolnir-mud/engine/plugins/controllers/internal/registry"
	"github.com/mjolnir-mud/engine/plugins/controllers/pkg/controller"
	"github.com/mjolnir-mud/engine/plugins/ecs"
)

type controllerSystem struct{}

func (s controllerSystem) Name() string {
	return "controller"
}

func (s controllerSystem) Component() string {
	return "controller"
}

func (s controllerSystem) Match(_ string, _ interface{}) bool {
	return true
}

func (s controllerSystem) ComponentAdded(_ string, _ string, _ interface{}) error { return nil }

func (s controllerSystem) ComponentUpdated(_ string, _ string, _ interface{}, _ interface{}) error {
	return nil
}

func (s controllerSystem) ComponentRemoved(_ string, _ string) error { return nil }

func (s controllerSystem) MatchingComponentAdded(entityId string, value interface{}) error {
	c, err := GetController(value.(string))

	if err != nil {
		return err
	}

	return c.Start(entityId)
}

func (s controllerSystem) MatchingComponentUpdated(entityId string, oldValue interface{}, newValue interface{}) error {
	oldController, err := GetController(oldValue.(string))

	if err != nil {
		return err
	}

	newController, err := GetController(newValue.(string))

	if err != nil {
		return err
	}

	err = oldController.Stop(entityId)

	if err != nil {
		return err
	}

	return newController.Start(entityId)
}

func (s controllerSystem) MatchingComponentRemoved(_ string) error {
	return nil
}

var ControllerSystem = controllerSystem{}

// GetController returns the Name for the session. If the session does not exist, an error will be returned.
func GetController(name string) (controller.Controller, error) {
	return registry.Get(name)
}

func SetController(entityId string, name string) error {
	_, err := GetController(name)

	if err != nil {
		return err
	}

	return ecs.AddOrUpdateStringComponentToEntity(entityId, "controller", name)
}

func HandleInput(entityId string, input string) error {
	c, err := GetController(entityId)

	if err != nil {
		return err
	}

	return c.HandleInput(entityId, input)
}
