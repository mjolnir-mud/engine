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

import (
	"fmt"
	"github.com/mjolnir-engine/engine/events"
	"github.com/mjolnir-engine/engine/uid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakeController struct {
	StartCalled chan *uid.UID
}

func (f fakeController) Name() string {
	return "test"
}

func (f fakeController) Start(c *ControllerContext) error {
	fmt.Print("starting")

	f.StartCalled <- c.SessionId

	return nil
}

func (f fakeController) Resume(_ *ControllerContext) error {
	return nil
}

func (f fakeController) Stop(_ *ControllerContext) error {
	return nil
}

func (f fakeController) HandleInput(_ *ControllerContext, _ string) error {
	return nil
}

func TestControllerRegistry_Register(t *testing.T) {
	engine := createEngineInstance()

	engine.RegisterController(fakeController{})

	assert.Len(t, engine.controllerRegistry.controllers, 1)
}

func TestEngine_GetController(t *testing.T) {
	engine := createEngineInstance()
	defer engine.Stop()

	engine.RegisterController(fakeController{})

	controller, err := engine.GetController("test")

	assert.Nil(t, err)
	assert.NotNil(t, controller)
}

func TestControllerRegistry_NewSession(t *testing.T) {
	engine := createEngineInstance()
	engine.Start("test")
	defer engine.Stop()

	fc := fakeController{
		StartCalled: make(chan *uid.UID),
	}

	engine.RegisterController(fc)

	id := uid.New()

	err := engine.Publish(events.SessionStartEvent{
		Id: id,
	})

	called := <-fc.StartCalled

	assert.Nil(t, err)
	assert.Equal(t, id, called)
}
