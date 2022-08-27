package session

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/internal/redis"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	session2 "github.com/mjolnir-mud/engine/plugins/sessions/pkg/entities/session"
	"github.com/mjolnir-mud/engine/plugins/sessions/pkg/events"
	"github.com/mjolnir-mud/engine/plugins/templates"
	testing2 "github.com/mjolnir-mud/engine/plugins/world/pkg/testing"
	"github.com/stretchr/testify/assert"
	"testing"
)

var ch = make(chan inputArgs)

type inputArgs struct {
	Id    string
	Input string
}

type testTemplate struct{}

func (t testTemplate) Render(args interface{}) (string, error) {
	return "test", nil
}

func (t testTemplate) Name() string {
	return "test"
}

func (t testTemplate) Style() string {
	return "default"
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
	ecs.RegisterEntityType(session2.Type)
	controller_registry.Register(testController{})
	controller_registry.Register(altController{})
	controller_registry.Start()
	engine.RegisterPlugin(templates.Plugin)
	testing2.Setup()
	session3.StartRegistry()
	_ = engine.RedisFlushAll()

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
	_ = redis.FlushAll()
	session3.StopRegistry()
	testing2.Teardown()
}

func TestStart(t *testing.T) {
	setup()
	defer teardown()

	err := Start("test")

	assert.NoError(t, err)
}

func TestGetController(t *testing.T) {
	setup()
	defer teardown()

	err := Start("test")

	assert.NoError(t, err)

	c, err := GetController("test")

	assert.NoError(t, err)

	assert.Equal(t, "login", c.Name())
}

func TestSetController(t *testing.T) {
	setup()
	defer teardown()

	err := SetController("test", "alt")

	assert.NoError(t, err)

	c, err := GetController("test")

	assert.NoError(t, err)
	assert.Equal(t, "alt", c.Name())
}

func TestHandleInput(t *testing.T) {
	setup()
	defer teardown()

	err := Start("test")

	assert.NoError(t, err)

	err = HandleInput("test", "test")

	assert.NoError(t, err)

	i := <-ch
	assert.Equal(t, "test", i.Input)
	assert.Equal(t, "test", i.Id)
}

func TestSendLine(t *testing.T) {
	setup()
	defer teardown()
	ch := make(chan string)

	sub := engine.Subscribe(events.PlayerOutputEvent{}, "test", func(e interface{}) {
		go func() { ch <- e.(*events.PlayerOutputEvent).Line }()
	})

	defer sub.Stop()

	err := Start("test")

	assert.NoError(t, err)

	err = SendLine("test", "test")

	assert.NoError(t, err)

	line := <-ch

	assert.Equal(t, "test", line)
}

func TestGetIntFromStore(t *testing.T) {
	setup()
	defer teardown()

	err := Start("test")

	assert.NoError(t, err)

	err = SetIntInStore("test", "test", 1)

	assert.NoError(t, err)

	i, err := GetIntFromStore("test", "test")

	assert.NoError(t, err)
	assert.Equal(t, 1, i)
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

func TestGetStringFromStore(t *testing.T) {
	setup()
	defer teardown()

	err := Start("test")

	assert.NoError(t, err)

	err = SetStringInStore("test", "test", "test")

	assert.NoError(t, err)

	s, err := GetStringFromStore("test", "test")

	assert.NoError(t, err)
	assert.Equal(t, "test", s)
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

func TestRenderTemplate(t *testing.T) {
	setup()
	defer teardown()

	ch := make(chan string)

	templates.RegisterTemplate(testTemplate{})

	sub := engine.Subscribe(events.PlayerOutputEvent{}, "test", func(e interface{}) {
		go func() { ch <- e.(*events.PlayerOutputEvent).Line }()
	})

	defer sub.Stop()

	err := Start("test")

	assert.NoError(t, err)

	err = RenderTemplate("test", "test", "test")

	assert.NoError(t, err)

	line := <-ch

	assert.Equal(t, "test", line)
}

func TestSendLineF(t *testing.T) {
	setup()
	defer teardown()

	ch := make(chan string)

	sub := engine.Subscribe(events.PlayerOutputEvent{}, "test", func(e interface{}) {
		go func() { ch <- e.(*events.PlayerOutputEvent).Line }()
	})

	defer sub.Stop()

	err := Start("test")

	assert.NoError(t, err)

	err = SendLineF("test", "test%s", "test")

	assert.NoError(t, err)

	line := <-ch

	assert.Equal(t, "testtest", line)
}
