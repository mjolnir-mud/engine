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
	sessionEntity "github.com/mjolnir-mud/engine/plugins/sessions/entities/session"
	engineTesting "github.com/mjolnir-mud/engine/testing"
	"testing"

	"github.com/mjolnir-mud/engine/plugins/ecs"
	ecsTesting "github.com/mjolnir-mud/engine/plugins/ecs/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/sessions/internal/registry"
	"github.com/stretchr/testify/assert"
)

func setup() {
	engineTesting.RegisterSetupCallback("sessions", func() {
		ecsTesting.Setup()

		ecs.RegisterEntityType(sessionEntity.Type)
	})
	engineTesting.Setup("world")

	registry.Start()

	ent, err := ecs.NewEntity("session", map[string]interface{}{})

	if err != nil {
		panic(err)
	}

	err = ecs.AddEntityWithID("session", "testing", ent)

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

	err := Start("testing")

	assert.NoError(t, err)
}

func TestGetIntFromFlash(t *testing.T) {
	setup()
	defer teardown()

	err := Start("testing")

	assert.NoError(t, err)

	err = SetIntInFlash("testing", "testing", 1)

	assert.NoError(t, err)

	i, err := GetIntFromFlash("testing", "testing")

	assert.NoError(t, err)
	assert.Equal(t, 1, i)
}

func TestGetIntFromFlashWithDefault(t *testing.T) {
	setup()
	defer teardown()

	err := Start("testing")

	assert.NoError(t, err)

	err = SetIntInFlash("testing", "testing", 1)

	assert.NoError(t, err)

	i, err := GetIntFromFlashWithDefault("testing", "testing", 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, i)

	i, err = GetIntFromFlashWithDefault("testing", "test2", 2)

	assert.NoError(t, err)
	assert.Equal(t, 2, i)

	i, err = GetIntFromFlashWithDefault("test3", "test2", 3)

	assert.Error(t, err)
}

func TestGetStringFromFlash(t *testing.T) {
	setup()
	defer teardown()

	err := Start("testing")

	assert.NoError(t, err)

	err = SetStringInFlash("testing", "testing", "testing")
	assert.NoError(t, err)

	s, err := GetStringFromFlash("testing", "testing")

	assert.NoError(t, err)
	assert.Equal(t, "testing", s)
}

func TestGetAccountId(t *testing.T) {
	setup()
	defer teardown()

	err := Start("testing")

	assert.NoError(t, err)

	err = SetAccountId("testing", "testing")

	assert.NoError(t, err)

	i, err := GetAccountId("testing")

	assert.NoError(t, err)
	assert.Equal(t, "testing", i)
}
