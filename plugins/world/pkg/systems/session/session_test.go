package session

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/internal/redis"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/world/internal/controller_registry"
	session3 "github.com/mjolnir-mud/engine/plugins/world/internal/session"
	session2 "github.com/mjolnir-mud/engine/plugins/world/pkg/entities/session"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/events"
	testing2 "github.com/mjolnir-mud/engine/plugins/world/pkg/testing"
	"github.com/stretchr/testify/assert"
	"testing"
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

func (l testController) Start(id string) error {
	//_ = promptLoginUsername(session)

	return nil
}

func (l testController) Resume(id string) error {
	return nil
}

func (l testController) Stop(id string) error {
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

func (a altController) Start(id string) error {
	return nil
}

func (a altController) Resume(id string) error {
	return nil
}

func (a altController) Stop(id string) error {
	return nil
}

func (a altController) HandleInput(id string, input string) error {
	return nil
}

func setup() {
	ecs.RegisterEntityType(session2.Type)
	controller_registry.Register(testController{})
	controller_registry.Register(altController{})
	controller_registry.Start()
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

	sub := engine.Subscribe(events.SendLineEvent{}, "test", func(e interface{}) {
		go func() { ch <- e.(*events.SendLineEvent).Line }()
	})

	defer sub.Stop()

	err := Start("test")

	assert.NoError(t, err)

	err = SendLine("test", "test")

	assert.NoError(t, err)

	line := <-ch

	assert.Equal(t, "test", line)
}
