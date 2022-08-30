package session

import (
	"fmt"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/ecs/pkg/errors"
	"github.com/mjolnir-mud/engine/plugins/templates"
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
func Start(_ string) error {
	return nil
}

// SetStringInStore sets a string value in the store, under the given key. If the session does not exist, an error is
// returned.
func SetStringInStore(id, key string, value string) error {
	return ecs.AddOrUpdateStringInMapComponent(id, "store", key, value)
}

// SendLine sends a line to the session's connection. If the session does not exist, an error is returned.
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

// SendLineF sends a line to the session's connection. If the session does not exist, an error is returned.
func SendLineF(id string, format string, args ...interface{}) error {
	return SendLine(id, fmt.Sprintf(format, args...))
}

// HandleInput passes the input to the session controller. If the session does not exist, an error is returned.
func HandleInput(id, input string) error {
	c, err := GetController(id)

	if err != nil {
		return err
	}

	return c.HandleInput(id, input)
}

// SetIntInStore sets an int value in the store, under the given key. If the session does not exist, an error is
// returned.
func SetIntInStore(id, key string, value int) error {
	return ecs.AddOrUpdateIntInMapComponent(id, "store", key, value)
}

// SetIntInFlash sets an int value in the flash, under the given key. If the session does not exist, an error is
// returned.
func SetIntInFlash(id, key string, value int) error {
	return ecs.AddOrUpdateIntInMapComponent(id, "flash", key, value)
}

// SetStringInFlash sets a string value in the flash, under the given key. If the session does not exist, an error is
// returned.
func SetStringInFlash(id, key string, value string) error {
	return ecs.AddOrUpdateStringInMapComponent(id, "flash", key, value)
}

// GetIntFromStore gets an int value from the store, under the given key. If the session does not exist, an error is
// returned.
func GetIntFromStore(id, key string) (int, error) {
	return ecs.GetIntFromMapComponent(id, "store", key)
}

// GetIntFromFlash gets a string value from the store, under the given key. If the session does not exist, an error
// is returned.
func GetIntFromFlash(id, key string) (int, error) {
	return ecs.GetIntFromMapComponent(id, "flash", key)
}

// GetStringFromStore gets a string value from the store, under the given key. If the session does not exist, an error
// is returned.
func GetStringFromStore(id, key string) (string, error) {
	return ecs.GetStringFromMapComponent(id, "store", key)
}

// GetStringFromFlash gets a string value from the store, under the given key. If the session does not exist, an error
// is returned.
func GetStringFromFlash(id, key string) (string, error) {
	return ecs.GetStringFromMapComponent(id, "flash", key)
}

// GetIntFromFlashWithDefault gets a string value from the store, under the given key. If the session does not exist,`
// an error is returned. If the key is not found, the default value is returned.
func GetIntFromFlashWithDefault(id, key string, defaultValue int) (int, error) {
	i, err := GetIntFromFlash(id, key)

	if err != nil {
		switch err.(type) {
		case errors.MapKeyNotFoundError:
			return defaultValue, nil
		default:
			return 0, err
		}
	}

	return i, nil
}

// RenderTemplate renders the template  with the given name and data. If the session does not exist, an error is
// returned.
func RenderTemplate(id, template string, data interface{}) error {
	rendered, err := templates.RenderTemplate(template, data)

	if err != nil {
		return err
	}

	return SendLine(id, rendered)
}

var System = session{}
