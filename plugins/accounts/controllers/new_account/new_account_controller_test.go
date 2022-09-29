package new_account

import (
	account2 "github.com/mjolnir-mud/engine/plugins/accounts/data_sources/account"
	"github.com/mjolnir-mud/engine/plugins/accounts/entities/account"
	"github.com/mjolnir-mud/engine/plugins/accounts/templates"
	dataSourcesTesting "github.com/mjolnir-mud/engine/plugins/data_sources/testing"
	mongoDataSourceTesting "github.com/mjolnir-mud/engine/plugins/mongo_data_source/testing"
	"github.com/mjolnir-mud/engine/plugins/sessions/systems/session"
	testing2 "github.com/mjolnir-mud/engine/plugins/sessions/testing"
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
		dataSourcesTesting.Setup()
		mongoDataSourceTesting.Setup()
		testing2.Setup()
		controllersTesting.Setup()
		engine.RegisterBeforeServiceStartCallback("world", func() {
			data_sources.Register(account2.Create())
		})

		engine.RegisterAfterServiceStartCallback("world", func() {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
			_ = data_sources.FindAndDelete("accounts", map[string]interface{}{"username": "testing-account"})

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
	_ = engine.RedisFlushAll()
	_ = data_sources.FindAndDelete("accounts", map[string]interface{}{"username": "testing-account"})
	_ = data_sources.FindAndDelete("accounts", map[string]interface{}{"username": "New_Random_Account"})
	engineTesting.Teardown()
}

func TestController_Name(t *testing.T) {
	setup()
	defer teardown()

	assert.Equal(t, Controller.Name(), "new_account")
}

func TestSignupHappyPath(t *testing.T) {
	setup()
	defer teardown()

	receivedLine, sub := testing2.CreateReceiveOutputSubscription()

	defer sub.Stop()

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	err = testing2.RegisterSession("sess")

	assert.NoError(t, err)

	assert.NoError(t, err)

	err = Controller.Start("sess")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter a username:\r\n", <-receivedLine)

	err = Controller.HandleInput("sess", "New_Random_Account")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter an email:\r\n", <-receivedLine)

	err = Controller.HandleInput("sess", "new_random_account@testing.com")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter a password:\r\n", <-receivedLine)

	err = Controller.HandleInput("sess", "A VERY secure password with lots of entropy")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Confirm your password:\r\n", <-receivedLine)

	err = Controller.HandleInput("sess", "A VERY secure password with lots of entropy")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	//assert.Equal(t, "Account created!", <-receivedLine)
}

func TestUsernameTooShort(t *testing.T) {
	setup()
	defer teardown()

	receivedLine, sub := testing2.CreateReceiveOutputSubscription()

	defer sub.Stop()

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = testing2.RegisterSession("sess")

	assert.NoError(t, err)

	err = Controller.Start("sess")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter a username:\r\n", <-receivedLine)

	err = Controller.HandleInput("sess", "New")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "'New' is not a valid username. It must be at least 4 characters long.\r\n", <-receivedLine)
}

func TestUsernameTooLong(t *testing.T) {
	setup()
	defer teardown()

	receivedLine, sub := testing2.CreateReceiveOutputSubscription()

	defer sub.Stop()

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = testing2.RegisterSession("sess")

	assert.NoError(t, err)

	err = Controller.Start("sess")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter a username:\r\n", <-receivedLine)

	err = Controller.HandleInput("sess", "New_Account_That_Is_Too_Long_For_Username")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(
		t,
		"'New_Account_That_Is_Too_Long_For_Username' is not a valid username. It must be at most 20 characters long.\r\n",
		<-receivedLine,
	)
}

func TestUsernameContainsInvalidCharacters(t *testing.T) {
	setup()
	defer teardown()

	receivedLine, sub := testing2.CreateReceiveOutputSubscription()

	defer sub.Stop()

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = testing2.RegisterSession("sess")

	assert.NoError(t, err)

	err = Controller.Start("sess")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter a username:\r\n", <-receivedLine)

	err = Controller.HandleInput("sess", "New Accounts")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(
		t,
		"'New Accounts' is not a valid username. It must contain only alpha-numeric characters, dashes (-) or underscores (_)\r\n",
		<-receivedLine,
	)
}

func TestInvalidEmail(t *testing.T) {
	setup()
	defer teardown()

	receivedLine, sub := testing2.CreateReceiveOutputSubscription()

	defer sub.Stop()

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = testing2.RegisterSession("sess")

	assert.NoError(t, err)

	err = session.SetIntInFlash("sess", "step", EmailStep)

	assert.NoError(t, err)

	err = Controller.HandleInput("sess", "Email is invalid")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "'email is invalid' is not a valid email address.\r\n", <-receivedLine)
}

func TestPasswordTooShort(t *testing.T) {
	setup()
	defer teardown()

	receivedLine, sub := testing2.CreateReceiveOutputSubscription()

	defer sub.Stop()

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = testing2.RegisterSession("sess")

	assert.NoError(t, err)

	err = session.SetIntInFlash("sess", "step", PasswordStep)

	assert.NoError(t, err)

	err = session.SetStringInFlash("sess", "username", "New_Account")

	assert.NoError(t, err)

	err = session.SetStringInFlash("sess", "email", "test_account@email.com")

	assert.NoError(t, err)

	err = Controller.HandleInput("sess", "New")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "That isn't a very secure password. Try something stronger.\r\n", <-receivedLine)
}
