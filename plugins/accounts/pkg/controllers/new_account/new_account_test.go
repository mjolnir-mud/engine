package new_account

import (
	"github.com/mjolnir-mud/engine"
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
	"github.com/mjolnir-mud/engine/plugins/sessions/pkg/systems/session"
	sessionsTesting "github.com/mjolnir-mud/engine/plugins/sessions/pkg/testing"
	templatesTesting "github.com/mjolnir-mud/engine/plugins/templates/pkg/testing"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func setup() {
	engineTesting.Setup("world", func() {
		ecsTesting.Setup()
		templatesTesting.Setup()
		dataSourcesTesting.Setup()
		mongoDataSourceTesting.Setup()
		sessionsTesting.Setup()
		controllersTesting.Setup()
		engine.RegisterBeforeServiceStartCallback("world", func() {
			data_sources.Register(data_source.Create())
		})

		engine.RegisterAfterServiceStartCallback("world", func() {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
			_ = data_sources.FindAndDelete("accounts", map[string]interface{}{"username": "testaccount"})

			err := data_sources.Save(
				"accounts",
				"testaccount",
				map[string]interface{}{
					"username":       "testaccount",
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

		ecs.RegisterEntityType(account.Type)
	})

	templates.RegisterAll()
}

func teardown() {
	_ = engine.RedisFlushAll()
	_ = data_sources.FindAndDelete("accounts", map[string]interface{}{"username": "testaccount"})
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

	receivedLine, sub := sessionsTesting.CreateReceiveOutputSubscription()

	defer sub.Stop()

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	err = sessionsTesting.RegisterSession("sess")

	assert.NoError(t, err)

	assert.NoError(t, err)

	err = Controller.Start("sess")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter a username:", <-receivedLine)

	err = Controller.HandleInput("sess", "New_Account")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter an email:", <-receivedLine)

	err = Controller.HandleInput("sess", "new_account@test.com")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter a password:", <-receivedLine)

	err = Controller.HandleInput("sess", "A VERY secure password with lots of entropy")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Confirm your password:", <-receivedLine)

	err = Controller.HandleInput("sess", "A VERY secure password with lots of entropy")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	//assert.Equal(t, "Account created!", <-receivedLine)
}

func TestUsernameTooShort(t *testing.T) {
	setup()
	defer teardown()

	receivedLine, sub := sessionsTesting.CreateReceiveOutputSubscription()

	defer sub.Stop()

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = sessionsTesting.RegisterSession("sess")

	assert.NoError(t, err)

	err = Controller.Start("sess")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter a username:", <-receivedLine)

	err = Controller.HandleInput("sess", "New")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "'New' is not a valid username. It must be at least 4 characters long.", <-receivedLine)
}

func TestUsernameTooLong(t *testing.T) {
	setup()
	defer teardown()

	receivedLine, sub := sessionsTesting.CreateReceiveOutputSubscription()

	defer sub.Stop()

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = sessionsTesting.RegisterSession("sess")

	assert.NoError(t, err)

	err = Controller.Start("sess")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter a username:", <-receivedLine)

	err = Controller.HandleInput("sess", "New_Account_That_Is_Too_Long_For_Username")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(
		t,
		"'New_Account_That_Is_Too_Long_For_Username' is not a valid username. It must be at most 20 characters long.",
		<-receivedLine,
	)
}

func TestUsernameContainsInvalidCharacters(t *testing.T) {
	setup()
	defer teardown()

	receivedLine, sub := sessionsTesting.CreateReceiveOutputSubscription()

	defer sub.Stop()

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = sessionsTesting.RegisterSession("sess")

	assert.NoError(t, err)

	err = Controller.Start("sess")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "Enter a username:", <-receivedLine)

	err = Controller.HandleInput("sess", "New Accounts")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(
		t,
		"'New Accounts' is not a valid username. It must contain only alpha-numeric characters, dashes (-) or underscores (_)",
		<-receivedLine,
	)
}

func TestInvalidEmail(t *testing.T) {
	setup()
	defer teardown()

	receivedLine, sub := sessionsTesting.CreateReceiveOutputSubscription()

	defer sub.Stop()

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = sessionsTesting.RegisterSession("sess")

	assert.NoError(t, err)

	err = session.SetIntInFlash("sess", "step", EmailStep)

	assert.NoError(t, err)

	err = Controller.HandleInput("sess", "Email is invalid")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "'email is invalid' is not a valid email address.", <-receivedLine)
}

func TestPasswordTooShort(t *testing.T) {
	setup()
	defer teardown()

	receivedLine, sub := sessionsTesting.CreateReceiveOutputSubscription()

	defer sub.Stop()

	err := ecs.AddEntityWithID("session", "sess", map[string]interface{}{})

	assert.NoError(t, err)

	err = sessionsTesting.RegisterSession("sess")

	assert.NoError(t, err)

	err = session.SetIntInFlash("sess", "step", PasswordStep)

	err = Controller.HandleInput("sess", "New")

	assert.Equal(t, nil, err)
	assert.NoError(t, err)
	assert.Equal(t, "That isn't a very secure password. Try something stronger.", <-receivedLine)
}
