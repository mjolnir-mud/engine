package systems

import (
	"github.com/mjolnir-mud/plugins/controllers/internal/registry"
	"github.com/mjolnir-mud/plugins/controllers/pkg/controller"
)

type controllerSystem struct{}

func (s controllerSystem) Name() string {
	return "controller"
}

func (s controllerSystem) Component() string {
	return "controller"
}

func (s controllerSystem) Match(_ string, _ interface{}) bool {
	return true
}

func (s controllerSystem) ComponentAdded(_ string, _ string, _ interface{}) error { return nil }

func (s controllerSystem) ComponentUpdated(_ string, _ string, _ interface{}, _ interface{}) error {
	return nil
}

func (s controllerSystem) ComponentRemoved(_ string, _ string) error { return nil }

func (s controllerSystem) MatchingComponentAdded(entityId string, value interface{}) error {
	c, err := GetController(value.(string))

	if err != nil {
		return err
	}

	return c.Start(entityId)
}

func (s controllerSystem) MatchingComponentUpdated(entityId string, oldValue interface{}, newValue interface{}) error {
	oldController, err := GetController(oldValue.(string))

	if err != nil {
		return err
	}

	newController, err := GetController(newValue.(string))

	if err != nil {
		return err
	}

	err = oldController.Stop(entityId)

	if err != nil {
		return err
	}

	return newController.Start(entityId)
}

func (s controllerSystem) MatchingComponentRemoved(_ string) error {
	return nil
}

var ControllerSystem = controllerSystem{}

// GetController returns the Name for the session. If the session does not exist, an error will be returned.
func GetController(name string) (controller.Controller, error) {
	return registry.Get(name)
}

func HandleInput(entityId string, input string) error {
	c, err := GetController(entityId)

	if err != nil {
		return err
	}

	return c.HandleInput(entityId, input)
}
