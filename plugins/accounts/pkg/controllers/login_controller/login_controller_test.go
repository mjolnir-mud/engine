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

package login_controller

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/pkg/event"
	"github.com/mjolnir-mud/engine/pkg/redis"
	engineTesting "github.com/mjolnir-mud/engine/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/accounts/internal/data_source"
	"github.com/mjolnir-mud/engine/plugins/accounts/internal/templates"
	"github.com/mjolnir-mud/engine/plugins/accounts/pkg/entities/account"
	controllersTesting "github.com/mjolnir-mud/engine/plugins/controllers/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	dataSourcesTesting "github.com/mjolnir-mud/engine/plugins/data_sources/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	ecsTesting "github.com/mjolnir-mud/engine/plugins/ecs/pkg/testing"
	mongoDataSourceTesting "github.com/mjolnir-mud/engine/plugins/mongo_data_source/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/sessions/pkg/events"
	"github.com/mjolnir-mud/engine/plugins/sessions/pkg/systems/session"
	templatesTesting "github.com/mjolnir-mud/engine/plugins/templates/pkg/testing"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func setup() {
	engineTesting.Setup(func() {
		ecsTesting.Setup()
		templatesTesting.Setup()
		dataSourcesTesting.Setup()
		mongoDataSourceTesting.Setup()
		controllersTesting.Setup()

		templates.RegisterAll()
		ecs.RegisterEntityType(account.Type)
		data_sources.Register(data_source.Create())
	})

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	_ = data_sources.FindAndDelete("accounts", map[string]interface{}{"username": "testaccount"})

	err := data_sources.Save(
		"accounts",
		"testaccount",
		map[string]interface{}{
			"username": "testaccount",
			"password": string(hashedPassword),
			"__metadata": map[string]interface{}{
				"entityType": "account",
				"collection": "accounts",
			},
		})

	if err != nil {
		panic(err)
	}
}

func teardown() {
	ecsTesting.Teardown()
	templatesTesting.Teardown()
	mongoDataSourceTesting.Teardown()
	dataSourcesTesting.Teardown()
	engineTesting.Teardown()
}

func TestController_Name(t *testing.T) {
	setup()
	defer teardown()

	c := controller{}

	assert.Equal(t, "login", c.Name())
}

func TestController_Start(t *testing.T) {
	setup()
	defer teardown()

	receivedLine := make(chan string)

	sub := createReceiveLineSubscription()

	defer sub.Stop()

	c := controller{}

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = c.Start("sess")

	assert.Equal(t, nil, err)

	assert.NoError(t, err)

	line := <-receivedLine

	assert.Equal(t, "Enter your username, or type '\u001B[1mcreate\u001B[0m' to create a new account:", line)
}

func TestControllerHandlesInvalidLogin(t *testing.T) {
	setup()
	defer teardown()

	receivedLine := make(chan string)

	sub := createReceiveLineSubscription()

	defer sub.Stop()

	c := controller{}

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = c.HandleInput("sess", "test")

	assert.NoError(t, err)

	i, err := session.GetIntFromFlash("sess", "step")

	assert.NoError(t, err)
	assert.Equal(t, 2, i)

	s, err := session.GetStringFromFlash("sess", "username")

	assert.NoError(t, err)
	assert.Equal(t, "test", s)

	line := <-receivedLine

	assert.Equal(t, "Enter your password:", line)

	err = c.HandleInput("sess", "test")

	line = <-receivedLine

	assert.Equal(t, "An account with that username and password combination was not found.", line)
}

func TestControllerHandlesValidLogin(t *testing.T) {
	setup()
	defer teardown()

	receivedLine := make(chan string)

	sub := createReceiveLineSubscription()

	defer sub.Stop()
	c := controller{}

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = c.HandleInput("sess", "testaccount")

	assert.NoError(t, err)
	assert.Equal(t, "Enter your password:", <-receivedLine)
	err = c.HandleInput("sess", "password")

	assert.NoError(t, err)

	s, err := session.GetStringFromStore("sess", "accountId")

	assert.NoError(t, err)
	assert.Equal(t, "sess", s)

	controller, err := session.GetStringFromStore("sess", "controller")

	assert.NoError(t, err)
	assert.Equal(t, "game", controller)
}

func TestControllerHandleUsernameCreate(t *testing.T) {
	setup()
	defer teardown()

	c := controller{}

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = c.HandleInput("sess", "create")

	assert.NoError(t, err)

	i, err := session.GetStringFromStore("sess", "controller")

	assert.NoError(t, err)
	assert.Equal(t, "new_account", i)
}

func createReceiveLineSubscription() redis.Subscription {
	receivedLine := make(chan string)

	return engine.Subscribe(events.PlayerOutputEvent{}, func(e event.EventPayload) {
		go func() {
			poe := &events.PlayerOutputEvent{}

			err := e.Unmarshal(poe)

			if err != nil {
				panic(err)
			}

			receivedLine <- poe.Line
		}()
	})
}
