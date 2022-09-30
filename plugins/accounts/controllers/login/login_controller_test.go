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

package login

import (
	accountDataSource "github.com/mjolnir-mud/engine/plugins/accounts/data_sources/account"
	"github.com/mjolnir-mud/engine/plugins/accounts/entities/account"
	"github.com/mjolnir-mud/engine/plugins/accounts/templates"
	"github.com/mjolnir-mud/engine/plugins/controllers"
	testing2 "github.com/mjolnir-mud/engine/plugins/data_sources/testing"
	testing3 "github.com/mjolnir-mud/engine/plugins/mongo_data_source/testing"
	"github.com/mjolnir-mud/engine/plugins/sessions/systems/session"
	testing4 "github.com/mjolnir-mud/engine/plugins/sessions/testing"
	"github.com/mjolnir-mud/engine/plugins/sessions/testing/helpers"
	"testing"

	"github.com/mjolnir-mud/engine"
	engineTesting "github.com/mjolnir-mud/engine/pkg/testing"
	controllersTesting "github.com/mjolnir-mud/engine/plugins/controllers/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	ecsTesting "github.com/mjolnir-mud/engine/plugins/ecs/pkg/testing"
	templatesTesting "github.com/mjolnir-mud/engine/plugins/templates/pkg/testing"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setup() {
	engineTesting.Setup("world", func() {
		ecsTesting.Setup()
		templatesTesting.Setup()
		testing2.Setup()
		testing3.Setup()
		testing4.Setup()
		controllersTesting.Setup()

		engine.RegisterBeforeServiceStartCallback("world", func() {
			data_sources.Register(accountDataSource.Create())
		})

		engine.RegisterAfterServiceStartCallback("world", func() {
			controllers.Register(controllersTesting.CreateMockController("new_account"))
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

			deleted := make(chan interface{})

			go func() {
				deleted <- data_sources.FindAndDelete("accounts", map[string]interface{}{"username": "testing-account"})
			}()

			<-deleted

			err := data_sources.SaveWithId(
				"accounts",
				"testing-account",
				map[string]interface{}{
					"username":       "testing-account",
					"hashedPassword": string(hashedPassword),
					"__metadata": map[string]interface{}{
						"entityType": "account",
						"collection": "accounts",
					},
				})

			if err != nil {
				panic(err)
			}
		})

		ecs.RegisterEntityType(account.EntityType)
	})

	templates.RegisterAll()

}

func teardown() {
	ecsTesting.Teardown()
	templatesTesting.Teardown()
	testing3.Teardown()
	testing2.Teardown()
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

	receivedLine, sub := helpers.CreateSessionWithOutputSubscription()

	defer sub.Stop()

	c := controller{}

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = helpers.RegisterSessionWithId("sess")

	assert.NoError(t, err)

	err = c.Start("sess")

	assert.Equal(t, nil, err)

	assert.NoError(t, err)

	line := <-receivedLine

	assert.Equal(t, "Enter your username, or type '\u001B[1mcreate\u001B[0m' to create a new account:\r\n", line)
}

func TestControllerHandlesInvalidLogin(t *testing.T) {
	setup()
	defer teardown()

	receivedLine, sub := helpers.CreateSessionWithOutputSubscription()

	defer sub.Stop()

	c := controller{}

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = helpers.RegisterSessionWithId("sess")

	assert.NoError(t, err)

	err = c.HandleInput("sess", "testing")

	assert.NoError(t, err)

	i, err := session.GetIntFromFlash("sess", "step")

	assert.NoError(t, err)
	assert.Equal(t, 2, i)

	s, err := session.GetStringFromFlash("sess", "username")

	assert.NoError(t, err)
	assert.Equal(t, "testing", s)

	line := <-receivedLine

	assert.Equal(t, "Enter your password:\r\n", line)

	err = c.HandleInput("sess", "testing")

	line = <-receivedLine

	assert.Equal(t, "An account with that username and password combination was not found.\r\n", line)
}

func TestControllerHandlesValidLogin(t *testing.T) {
	setup()
	defer teardown()

	receivedLine, sub := helpers.CreateSessionWithOutputSubscription()

	defer sub.Stop()
	c := controller{}

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = helpers.RegisterSessionWithId("sess")

	assert.NoError(t, err)

	err = c.HandleInput("sess", "testing-account")

	assert.NoError(t, err)
	assert.Equal(t, "Enter your password:\r\n", <-receivedLine)
	err = c.HandleInput("sess", "password")

	assert.NoError(t, err)

	s, err := ecs.GetStringComponent("sess", "accountId")

	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func TestControllerHandleUsernameCreate(t *testing.T) {
	setup()
	defer teardown()

	c := controller{}

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = helpers.RegisterSessionWithId("sess")

	assert.NoError(t, err)

	err = c.HandleInput("sess", "create")

	assert.NoError(t, err)

	i, err := ecs.GetStringComponent("sess", "controller")

	assert.NoError(t, err)
	assert.Equal(t, "new_account", i)
}
