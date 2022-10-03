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

package registry

import (
	"github.com/mjolnir-mud/engine/plugins/ecs"
	ecsTesting "github.com/mjolnir-mud/engine/plugins/ecs/pkg/testing"
	sessionsTesting "github.com/mjolnir-mud/engine/plugins/sessions/testing"
	engineTesting "github.com/mjolnir-mud/engine/testing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() {
	engineTesting.RegisterSetupCallback("controllers", func() {
		ecsTesting.Setup()
		sessionsTesting.Setup()
	})
	engineTesting.Setup("world")

	Start()
}

func teardown() {
	Stop()
	engineTesting.Teardown()
}

type testController struct {
	HandleInputCalled chan []string
}

func (c testController) Name() string {
	return "testing"
}

func (c testController) Start(_ string) error {
	return nil
}

func (c testController) Resume(_ string) error {
	return nil
}

func (c testController) Stop(_ string) error {
	return nil
}

func (c testController) HandleInput(_ string, _ string) error {
	go func() { c.HandleInputCalled <- []string{"testing", "testing"} }()

	return nil
}

func TestStart(t *testing.T) {
	setup()
	defer teardown()

	Start()

	assert.NotNil(t, controllers)
}

func TestRegister(t *testing.T) {
	setup()
	defer teardown()

	Register(testController{})

	assert.Len(t, controllers, 1)
}

func TestGet(t *testing.T) {
	setup()
	defer teardown()

	Register(testController{})

	c, err := Get("testing")

	assert.Nil(t, err)
	assert.NotNil(t, c)
}

func TestHandleInput(t *testing.T) {
	setup()
	defer teardown()

	tc := &testController{
		HandleInputCalled: make(chan []string, 1),
	}

	Register(tc)

	err := ecs.AddEntityWithID("session", "testing", map[string]interface{}{})

	assert.Nil(t, err)

	err = ecs.AddStringComponentToEntity("testing", "controller", "testing")

	assert.Nil(t, err)

	err = HandleInput("testing", "testing")

	assert.Nil(t, err)

	res := <-tc.HandleInputCalled

	assert.Equal(t, []string{"testing", "testing"}, res)
}
