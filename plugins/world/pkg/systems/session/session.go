package session

import (
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/ecs/pkg/errors"
	"github.com/mjolnir-mud/engine/plugins/world/internal/controller_registry"
	session2 "github.com/mjolnir-mud/engine/plugins/world/internal/session"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/controller"
)

type session struct{}

func (e session) Name() string {
	return "expiration"
}

func (e session) Component() string {
	return "session"
}

func (e session) Match(key string, value interface{}) bool {
	// return false if value is not a string
	if _, ok := value.(string); !ok {
		return false
	}

	// return true if the value is session and the key is _type
	return key == "_type" && value.(string) == "session"
}

func (e session) ComponentAdded(_ string, _ string, _ interface{}) error { return nil }

func (e session) ComponentUpdated(_ string, _ string, _ interface{}, _ interface{}) error { return nil }

func (e session) ComponentRemoved(_ string, _ string) error { return nil }

func (e session) MatchingComponentAdded(_ string, _ string, _ interface{}) error { return nil }

func (e session) MatchingComponentUpdated(_ string, _ string, _ interface{}, _ interface{}) error {
	return nil
}

func (e session) MatchingComponentRemoved(_ string, _ string) error { return nil }

// Start starts the session with the given id. This calls the start method of the testController. If the session does not
// exist, an error is returned.
func Start(id string) error {
	c, err := GetController(id)

	if err != nil {
		return err
	}

	return c.Start(id)
}

// GetController returns the testController for the session. If the session does not exist, an error will be returned.
// If the controller is not found, an error will be returned.
func GetController(name string) (controller.Controller, error) {
	c, err := ecs.GetStringFromMapComponent(name, "store", "controller")

	if err != nil {
		return nil, err
	}

	return controller_registry.Get(c)
}

// SetController sets the testController for the session. All input is passed through the session controller. If the
// session does not exist, an error is returned.
func SetController(id, controller string) error {
	return SetStringInStore(id, "controller", controller)
}

// SetStringInStore sets a string value in the store, under the given key. If the session does not exist, an error is
// returned.
func SetStringInStore(id, key string, value string) error {
	return ecs.AddOrUpdateStringInMapComponent(id, "store", key, value)
}

// SendLine sends a line to the sessions connection. If the session does not exist, an error is returned.
func SendLine(id, line string) error {
	exists, err := ecs.EntityExists(id)

	if err != nil {
		return err
	}

	if !exists {
		return errors.EntityNotFoundError{
			ID: id,
		}
	}

	session2.SendLine(id, line)

	return nil
}

// HandleInput passes the input to the session controller. If the session does not exist, an error is returned.
func HandleInput(id, input string) error {
	c, err := GetController(id)

	if err != nil {
		return err
	}

	return c.HandleInput(id, input)
}

var System = session{}
