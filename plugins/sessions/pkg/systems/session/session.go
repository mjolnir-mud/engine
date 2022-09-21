/*
 * Copyright (c) 2022 eightfivefour llc. All rights reserved.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package session

import (
	"fmt"
	"strings"

	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/ecs/pkg/errors"
	"github.com/mjolnir-mud/engine/plugins/sessions/internal/registry"
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

// GetIntFromFlash gets a string value from the store, under the given key. If the session does not exist, an error
// is returned.
func GetIntFromFlash(id, key string) (int, error) {
	return ecs.GetIntFromMapComponent(id, "flash", key)
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

// SendLine sends a line to the session. If the session does not exist, an error is returned.
func SendLine(id, line string) error {
	return registry.SendLine(id, line)
}

// SendLines sends a slice of lines to the session. If the session does not exist, an error is returned.
func SendLines(id string, lines []string) error {
	for _, line := range lines {
		if err := SendLine(id, line); err != nil {
			return err
		}
	}

	return nil
}

// SendLinef sends a formatted line to the session. If the session does not exist, an error is returned.
func SendLinef(id, format string, args ...interface{}) error {
	return SendLine(id, fmt.Sprintf(format, args...))
}

// Send breaks up a string into lines and calles SendLines. If the session does not exist, an error is returned.
func Send(id, text string) error {
	return SendLines(id, strings.Split(text, "\n"))
}

// Sendf breaks up a string into lines and calles SendLines. If the session does not exist, an error is returned.
func Sendf(id, format string, args ...interface{}) error {
	return Send(id, fmt.Sprintf(format, args...))
}

// Render renders a template to the session. If the session does not exist, an error is returned. If the template does
// not exist, an error is returned.
func Render(id, template string, data interface{}) error {
	content, err := templates.RenderTemplate(template, data)

	if err != nil {
		return err
	}

	return Send(id, content)

}
