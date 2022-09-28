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
	engineTesting "github.com/mjolnir-mud/engine/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/controllers/internal/registry"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	ecsTesting "github.com/mjolnir-mud/engine/plugins/ecs/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/ecs/pkg/testing/fakes"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testController struct {
	StartCalled       chan string
	StopCalled        chan string
	HandleInputCalled chan []string
}

func (c testController) Name() string {
	return "testing"
}

func (c testController) Start(entityId string) error {
	go func() {
		c.StartCalled <- entityId
	}()

	return nil
}

func (c testController) Resume(_ string) error {
	return nil
}

func (c testController) Stop(entityId string) error {
	go func() { c.StopCalled <- entityId }()

	return nil
}

func (c testController) HandleInput(entityId string, input string) error {
	go func() { c.HandleInputCalled <- []string{entityId, input} }()

	return nil
}

func setup() {
	tc = &testController{
		StartCalled:       make(chan string),
		StopCalled:        make(chan string),
		HandleInputCalled: make(chan []string),
	}

	engineTesting.Setup("world", func() {
		registry.Start()
		registry.Register(tc)
		ecsTesting.Setup()
	})

	ecs.RegisterEntityType(fakes.FakeEntityType{})
	ecs.RegisterSystem(ControllerSystem)

}

func teardown() {
	registry.Stop()
	ecsTesting.Teardown()
	engineTesting.Teardown()
}

var tc *testController

func TestControllerSystem_Name(t *testing.T) {
	assert.Equal(t, "controller", ControllerSystem.Name())
}

func TestControllerSystem_Component(t *testing.T) {
	assert.Equal(t, "controller", ControllerSystem.Component())
}

func TestControllerSystem_Match(t *testing.T) {
	assert.True(t, ControllerSystem.Match("controller", "testing"))
}

func TestControllerSystem_MatchingComponentAdded(t *testing.T) {
	setup()
	defer teardown()

	assert.Nil(t, ControllerSystem.MatchingComponentAdded("testing", "testing"))

	entityId := <-tc.StartCalled

	assert.Equal(t, "testing", entityId)
}

func TestControllerSystem_MatchingComponentUpdated(t *testing.T) {
	setup()
	defer teardown()

	assert.Nil(t, ControllerSystem.MatchingComponentUpdated("testing", "testing", "testing"))
	entityId := <-tc.StopCalled

	assert.Equal(t, "testing", entityId)

	entityId = <-tc.StartCalled

	assert.Equal(t, "testing", entityId)
}

func TestGetController(t *testing.T) {
	setup()
	defer teardown()

	c, err := GetController("testing")

	assert.Nil(t, err)
	assert.Equal(t, tc, c)
}

func TestSetController(t *testing.T) {
	setup()
	defer teardown()

	err := ecs.AddEntityWithID("testing", "testing", map[string]interface{}{})

	assert.Nil(t, err)

	assert.Nil(t, SetController("testing", "testing"))

	called := <-tc.StartCalled

	assert.Equal(t, "testing", called)
}

func TestHandleInput(t *testing.T) {
	setup()
	defer teardown()

	assert.Nil(t, HandleInput("testing", "testing"))

	called := <-tc.HandleInputCalled

	assert.Equal(t, "testing", called[0])
	assert.Equal(t, "testing", called[1])
}
