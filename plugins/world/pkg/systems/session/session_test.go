package session

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/internal/redis"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/world/internal/controller_registry"
	session2 "github.com/mjolnir-mud/engine/plugins/world/pkg/entities/session"
	testing2 "github.com/mjolnir-mud/engine/plugins/world/pkg/testing"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	return nil
}

func setup() {
	ecs.RegisterEntityType(session2.Type)
	controller_registry.Register(testController{})
	controller_registry.Start()
	testing2.Setup()
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
	testing2.Teardown()
}

func TestStart(t *testing.T) {
	setup()
	defer teardown()

	err := Start("test")

	assert.NoError(t, err)
}
