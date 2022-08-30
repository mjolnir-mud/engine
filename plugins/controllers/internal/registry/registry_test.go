package registry

import (
	"github.com/mjolnir-mud/engine"
	engineTesting "github.com/mjolnir-mud/engine/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/sessions"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() {
	engine.RegisterPlugin(ecs.Plugin)
	engine.RegisterPlugin(sessions.Plugin)

	engineTesting.Setup()
	Start()
}

func teardown() {
	Stop()
	engineTesting.Teardown()
}

type testController struct {
	HandleInputCalled chan []string
}

func (c testController) Name() string {
	return "test"
}

func (c testController) Start(_ string) error {
	return nil
}

func (c testController) Resume(_ string) error {
	return nil
}

func (c testController) Stop(_ string) error {
	return nil
}

func (c testController) HandleInput(_ string, _ string) error {
	return nil
}

func TestStart(t *testing.T) {
	setup()
	defer teardown()

	Start()

	assert.NotNil(t, controllers)
}

func TestRegister(t *testing.T) {
	setup()
	defer teardown()

	Register(testController{})

	assert.Len(t, controllers, 1)
}

func TestGet(t *testing.T) {
	setup()
	defer teardown()

	Register(testController{})

	c, err := Get("test")

	assert.Nil(t, err)
	assert.NotNil(t, c)
}

func TestHandleInput(t *testing.T) {
	setup()
	defer teardown()

	tc := &testController{
		HandleInputCalled: make(chan []string, 1),
	}

	Register(tc)

	err := HandleInput("test", "test")

	assert.Nil(t, err)

	res := <-tc.HandleInputCalled

	assert.Equal(t, []string{"test", "test"}, res)
}
