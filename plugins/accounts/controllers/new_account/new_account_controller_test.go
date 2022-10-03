package new_account

import (
	accountDataSource "github.com/mjolnir-mud/engine/plugins/accounts/data_sources/account"
	"github.com/mjolnir-mud/engine/plugins/accounts/entities/account"
	"github.com/mjolnir-mud/engine/plugins/accounts/templates"
	"github.com/mjolnir-mud/engine/plugins/controllers"
	dataSourcesTesting "github.com/mjolnir-mud/engine/plugins/data_sources/testing"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	ecsTesting "github.com/mjolnir-mud/engine/plugins/ecs/pkg/testing"
	mongoDataSourceTesting "github.com/mjolnir-mud/engine/plugins/mongo_data_source/testing"
	"github.com/mjolnir-mud/engine/plugins/sessions/systems/session"
	sessionsTesting "github.com/mjolnir-mud/engine/plugins/sessions/testing"
	"github.com/mjolnir-mud/engine/plugins/sessions/testing/helpers"
	templatesTesting "github.com/mjolnir-mud/engine/plugins/templates/pkg/testing"
	engineTesting "github.com/mjolnir-mud/engine/testing"
	"testing"

	"github.com/mjolnir-mud/engine"
	controllersTesting "github.com/mjolnir-mud/engine/plugins/controllers/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setup() {
	engineTesting.RegisterSetupCallback("accounts", func() {
		ecsTesting.Setup()
		templatesTesting.Setup()
		dataSourcesTesting.Setup()
		mongoDataSourceTesting.Setup()
		sessionsTesting.Setup()
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

			ecs.RegisterEntityType(account.EntityType)
			templates.RegisterAll()
		})
	})

	engineTesting.Setup("world")
}

func teardown() {
	_ = engine.RedisFlushAll()
	_ = data_sources.FindAndDelete("accounts", map[string]interface{}{"username": "testing-account"})
	_ = data_sources.FindAndDelete("accounts", map[string]interface{}{"username": "New_Random_Account"})
	controllersTesting.Teardown()
	templatesTesting.Teardown()
	mongoDataSourceTesting.Teardown()
	sessionsTesting.Teardown()
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

	id, receivedLine, sub, err := helpers.CreateSessionWithOutputSubscription()
	assert.NoError(t, err)

	defer sub.Stop()

	err = Controller.Start(id)

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter a username:\r\n", <-receivedLine)

	err = Controller.HandleInput(id, "New_Random_Account")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter an email:\r\n", <-receivedLine)

	err = Controller.HandleInput(id, "new_random_account@testing.com")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter a password:\r\n", <-receivedLine)

	err = Controller.HandleInput(id, "A VERY secure password with lots of entropy")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Confirm your password:\r\n", <-receivedLine)

	err = Controller.HandleInput(id, "A VERY secure password with lots of entropy")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	//assert.Equal(t, "Account created!", <-receivedLine)
}

func TestUsernameTooShort(t *testing.T) {
	setup()
	defer teardown()

	id, receivedLine, sub, err := helpers.CreateSessionWithOutputSubscription()

	defer sub.Stop()

	assert.NoError(t, err)

	err = helpers.CreateSessionWithId(id)

	assert.NoError(t, err)

	err = Controller.Start(id)

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter a username:\r\n", <-receivedLine)

	err = Controller.HandleInput(id, "New")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "'New' is not a valid username. It must be at least 4 characters long.\r\n", <-receivedLine)
}

func TestUsernameTooLong(t *testing.T) {
	setup()
	defer teardown()

	id, receivedLine, sub, err := helpers.CreateSessionWithOutputSubscription()

	defer sub.Stop()

	assert.NoError(t, err)

	err = helpers.CreateSessionWithId(id)

	assert.NoError(t, err)

	err = Controller.Start(id)

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter a username:\r\n", <-receivedLine)

	err = Controller.HandleInput(id, "New_Account_That_Is_Too_Long_For_Username")

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

	id, receivedLine, sub, err := helpers.CreateSessionWithOutputSubscription()

	defer sub.Stop()

	assert.NoError(t, err)

	err = helpers.CreateSessionWithId(id)

	assert.NoError(t, err)

	err = Controller.Start(id)

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter a username:\r\n", <-receivedLine)

	err = Controller.HandleInput(id, "New Accounts")

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

	id, receivedLine, sub, err := helpers.CreateSessionWithOutputSubscription()

	defer sub.Stop()

	assert.NoError(t, err)

	err = helpers.CreateSessionWithId(id)

	assert.NoError(t, err)

	err = session.SetIntInFlash(id, "step", EmailStep)

	assert.NoError(t, err)

	err = Controller.HandleInput(id, "Email is invalid")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "'email is invalid' is not a valid email address.\r\n", <-receivedLine)
}

func TestPasswordTooShort(t *testing.T) {
	setup()
	defer teardown()

	id, receivedLine, sub, err := helpers.CreateSessionWithOutputSubscription()

	defer sub.Stop()

	assert.NoError(t, err)

	err = helpers.CreateSessionWithId(id)

	assert.NoError(t, err)

	err = session.SetIntInFlash(id, "step", PasswordStep)

	assert.NoError(t, err)

	err = session.SetStringInFlash(id, "username", "New_Account")

	assert.NoError(t, err)

	err = session.SetStringInFlash(id, "email", "test_account@email.com")

	assert.NoError(t, err)

	err = Controller.HandleInput(id, "New")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "That isn't a very secure password. Try something stronger.\r\n", <-receivedLine)
}
