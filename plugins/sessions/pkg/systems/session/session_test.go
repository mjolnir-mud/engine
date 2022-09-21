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

package session

import (
	"testing"

	engineTesting "github.com/mjolnir-mud/engine/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	ecsTesting "github.com/mjolnir-mud/engine/plugins/ecs/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/sessions/internal/registry"
	sessionEntity "github.com/mjolnir-mud/engine/plugins/sessions/pkg/entities/session"
	"github.com/stretchr/testify/assert"
)

var ch = make(chan inputArgs)

type inputArgs struct {
	Id    string
	Input string
}

// testController is the login testController, responsible handling user logins.
type testController struct{}

func (l testController) Name() string {
	return "login"
}

func (l testController) Start(_ string) error {
	return nil
}

func (l testController) Resume(_ string) error {
	return nil
}

func (l testController) Stop(_ string) error {
	return nil
}

func (l testController) HandleInput(id string, input string) error {
	go func() { ch <- inputArgs{Id: id, Input: input} }()

	return nil
}

type altController struct{}

func (a altController) Name() string {
	return "alt"
}

func (a altController) Start(_ string) error {
	return nil
}

func (a altController) Resume(_ string) error {
	return nil
}

func (a altController) Stop(_ string) error {
	return nil
}

func (a altController) HandleInput(_ string, _ string) error {
	return nil
}

func setup() {
	engineTesting.Setup(func() {
		ecsTesting.Setup()

		ecs.RegisterEntityType(sessionEntity.Type)
	})

	registry.Start()

	ent, err := ecs.CreateEntity("session", map[string]interface{}{})

	if err != nil {
		panic(err)
	}

	err = ecs.AddEntityWithID("session", "test", ent)

	if err != nil {
		panic(err)
	}
}

func teardown() {
	ecsTesting.Teardown()
	registry.Stop()
	engineTesting.Teardown()
}

func TestStart(t *testing.T) {
	setup()
	defer teardown()

	err := Start("test")

	assert.NoError(t, err)
}

func TestGetIntFromFlash(t *testing.T) {
	setup()
	defer teardown()

	err := Start("test")

	assert.NoError(t, err)

	err = SetIntInFlash("test", "test", 1)

	assert.NoError(t, err)

	i, err := GetIntFromFlash("test", "test")

	assert.NoError(t, err)
	assert.Equal(t, 1, i)
}

func TestGetIntFromFlashWithDefault(t *testing.T) {
	setup()
	defer teardown()

	err := Start("test")

	assert.NoError(t, err)

	err = SetIntInFlash("test", "test", 1)

	assert.NoError(t, err)

	i, err := GetIntFromFlashWithDefault("test", "test", 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, i)

	i, err = GetIntFromFlashWithDefault("test", "test2", 2)

	assert.NoError(t, err)
	assert.Equal(t, 2, i)

	i, err = GetIntFromFlashWithDefault("test3", "test2", 3)

	assert.Error(t, err)
}

func TestGetStringFromFlash(t *testing.T) {
	setup()
	defer teardown()

	err := Start("test")

	assert.NoError(t, err)

	err = SetStringInFlash("test", "test", "test")
	assert.NoError(t, err)

	s, err := GetStringFromFlash("test", "test")

	assert.NoError(t, err)
	assert.Equal(t, "test", s)
}
