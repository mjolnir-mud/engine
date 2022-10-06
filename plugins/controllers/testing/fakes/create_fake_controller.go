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

package fakes

import (
	"github.com/mjolnir-mud/engine/plugins/controllers/controller"
	"github.com/mjolnir-mud/engine/plugins/sessions/systems/session"
)

type fakeController struct {
	ControllerName    string
	HandleInputCalled chan []string
}

func (c fakeController) Name() string {
	return c.ControllerName
}

func (c fakeController) Start(id string) error {
	return session.SendLine(id, "testing")
}

func (c fakeController) Resume(_ string) error {
	return nil
}

func (c fakeController) Stop(_ string) error {
	return nil
}

func (c fakeController) HandleInput(_ string, _ string) error {
	go func() { c.HandleInputCalled <- []string{"testing", "testing"} }()

	return nil
}

func CreateFakeController(name string) controller.Controller {
	return fakeController{
		ControllerName:    name,
		HandleInputCalled: make(chan []string),
	}
}
