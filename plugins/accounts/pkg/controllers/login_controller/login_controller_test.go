package login_controller

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/pkg/testing"
	templates2 "github.com/mjolnir-mud/engine/plugins/accounts/internal/templates"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/templates"
	"github.com/mjolnir-mud/engine/plugins/world"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/events"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/systems/session"
	"github.com/stretchr/testify/assert"
	testing2 "testing"
)

func setup() {
	engine.RegisterPlugin(ecs.Plugin)
	engine.RegisterPlugin(templates.Plugin)
	engine.RegisterPlugin(world.Plugin)

	testing.Setup()
	templates.RegisterTemplate(templates2.PromptUsernameTemplate)
	templates.RegisterTemplate(templates2.PromptPasswordTemplate)
	_ = engine.RedisFlushAll()
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

func TestControllerHandlesUsername(t *testing2.T) {
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
}

func TestControllerHandleUsernameCreate(t *testing2.T) {
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

	err = c.HandleInput("sess", "create")

	assert.NoError(t, err)

	i, err := session.GetIntFromFlash("sess", "step")

	assert.NoError(t, err)
	assert.Equal(t, 1, i)
}
