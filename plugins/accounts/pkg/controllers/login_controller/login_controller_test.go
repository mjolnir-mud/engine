package login_controller

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/accounts/internal/data_source"
	templates2 "github.com/mjolnir-mud/engine/plugins/accounts/internal/templates"
	"github.com/mjolnir-mud/engine/plugins/accounts/pkg/entities/account"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/mongo_data_source"
	"github.com/mjolnir-mud/engine/plugins/sessions/pkg/events"
	"github.com/mjolnir-mud/engine/plugins/templates"
	"github.com/mjolnir-mud/engine/plugins/world"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/systems/session"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	testing2 "testing"
)

func setup() {
	engine.RegisterPlugin(ecs.Plugin)
	engine.RegisterPlugin(templates.Plugin)
	engine.RegisterPlugin(world.Plugin)
	engine.RegisterPlugin(data_sources.Plugin)
	engine.RegisterPlugin(mongo_data_source.Plugin)
	ecs.RegisterEntityType(account.Type)

	data_sources.Register(data_source.Create())

	testing.Setup()
	templates2.RegisterAll()
	_ = engine.RedisFlushAll()

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
	_ = engine.RedisFlushAll()
	testing.Teardown()
}

func TestController_Name(t *testing2.T) {
	setup()
	defer teardown()

	c := controller{}

	assert.Equal(t, "login", c.Name())
}

func TestController_Start(t *testing2.T) {
	setup()
	defer teardown()

	receivedLine := make(chan string)

	sub := engine.Subscribe(events.SendLineEvent{}, "sess", func(e interface{}) {
		go func() { receivedLine <- e.(*events.SendLineEvent).Line }()
	})

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

func TestControllerHandlesInvalidLogin(t *testing2.T) {
	setup()
	defer teardown()

	receivedLine := make(chan string)

	sub := engine.Subscribe(events.SendLineEvent{}, "sess", func(e interface{}) {
		go func() { receivedLine <- e.(*events.SendLineEvent).Line }()
	})

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

func TestControllerHandlesValidLogin(t *testing2.T) {
	setup()
	defer teardown()

	receivedLine := make(chan string)

	sub := engine.Subscribe(events.SendLineEvent{}, "sess", func(e interface{}) {
		go func() { receivedLine <- e.(*events.SendLineEvent).Line }()
	})

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

func TestControllerHandleUsernameCreate(t *testing2.T) {
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
