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

package engine

import "github.com/mjolnir-mud/engine/uid"

// System is the interface that all systems must implement. They are reactive to changes in components that match
// the systems `Component` method.
type System interface {
	Name() string

	// Component is the type of the component that this system is responsible for.
	Component() string

	// Match is a check to see if the system should be run for the given key and value.
	Match(key string, value interface{}) bool

	// ComponentAdded is called when any component is added to an entity, including those the system is interested in.
	ComponentAdded(entityId *uid.UID, component string, value interface{}) error

	// ComponentUpdated is called when a component is updated on an entity, including those the system is interested in.
	ComponentUpdated(entityId *uid.UID, component string, oldValue interface{}, newValue interface{}) error

	// ComponentRemoved is called when a component is removed from an entity, including those the system is interested
	// in.
	ComponentRemoved(entityId *uid.UID, component string, value interface{}) error

	// MatchingComponentAdded is called when a component is added to an entity, but only if the system is interested
	// in the component.
	MatchingComponentAdded(entityId *uid.UID, value interface{}) error

	// MatchingComponentUpdated is called when a component is updated on an entity, but only if the system is interested
	// in the component.
	MatchingComponentUpdated(entityId *uid.UID, oldValue interface{}, newValue interface{}) error

	// MatchingComponentRemoved is called when a component is removed from an entity, but only if the system is
	// interested in the component.
	MatchingComponentRemoved(entityId *uid.UID, value interface{}) error
}
