package systems

import (
	"github.com/mjolnir-mud/plugins/controllers/internal/registry"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testController struct {
	StartCalled chan string
	StopCalled  chan string
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

func (c testController) HandleInput(_ string, _ string) error {
	return nil
}

func setup() {
	tc = &testController{
		StartCalled: make(chan string),
		StopCalled:  make(chan string),
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
