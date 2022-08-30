package systems

import (
	"github.com/mjolnir-mud/engine/plugins/controllers/internal/registry"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testController struct {
	StartCalled       chan string
	StopCalled        chan string
	HandleInputCalled chan []string
}

func (c testController) Name() string {
	return "test"
}

func (c testController) Start(entityId string) error {
	go func() { c.StartCalled <- entityId }()

	return nil
}

func (c testController) Resume(_ string) error {
	return nil
}

func (c testController) Stop(entityId string) error {
	go func() { c.StopCalled <- entityId }()

	return nil
}

func (c testController) HandleInput(entityId string, input string) error {
	go func() { c.HandleInputCalled <- []string{entityId, input} }()

	return nil
}

func setup() {
	tc = &testController{
		StartCalled:       make(chan string),
		StopCalled:        make(chan string),
		HandleInputCalled: make(chan []string),
	}

	registry.Start()
	registry.Register(tc)
}

func teardown() {
	registry.Stop()
}

var tc *testController

func TestControllerSystem_Name(t *testing.T) {
	assert.Equal(t, "controller", ControllerSystem.Name())
}

func TestControllerSystem_Component(t *testing.T) {
	assert.Equal(t, "controller", ControllerSystem.Component())
}

func TestControllerSystem_Match(t *testing.T) {
	assert.True(t, ControllerSystem.Match("controller", "test"))
}

func TestControllerSystem_MatchingComponentAdded(t *testing.T) {
	setup()
	defer teardown()

	assert.Nil(t, ControllerSystem.MatchingComponentAdded("test", "test"))
	entityId := <-tc.StartCalled

	assert.Equal(t, "test", entityId)
}

func TestControllerSystem_MatchingComponentUpdated(t *testing.T) {
	setup()
	defer teardown()

	assert.Nil(t, ControllerSystem.MatchingComponentUpdated("test", "test", "test"))
	entityId := <-tc.StopCalled

	assert.Equal(t, "test", entityId)

	entityId = <-tc.StartCalled

	assert.Equal(t, "test", entityId)
}

func TestGetController(t *testing.T) {
	setup()
	defer teardown()

	c, err := GetController("test")

	assert.Nil(t, err)
	assert.Equal(t, tc, c)
}

func TestHandleInput(t *testing.T) {
	setup()
	defer teardown()

	assert.Nil(t, HandleInput("test", "test"))

	called := <-tc.HandleInputCalled

	assert.Equal(t, "test", called[0])
	assert.Equal(t, "test", called[1])
}
