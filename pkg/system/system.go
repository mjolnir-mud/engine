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

package system

import (
	"context"
)

// System is the interface that all systems must implement. They are reactive to changes in components that match
// the systems `Component` method.
type System interface {
	Name() string

	// Component is the type of the component that this system is responsible for.
	Component() string

	// EntityAdded is called when an entity has been added which has a component that matches the systems `Component`
	EntityAdded(ctx context.Context) context.Context

	// EntityUpdated is called when an entity has been updated which has a component that matches the systems `Component`
	EntityUpdated(ctx context.Context) context.Context

	// EntityRemoved is called when an entity has been removed which has a component that matches the systems `Component`
	EntityRemoved(ctx context.Context) context.Context

	// ComponentAdded is called when a matching component has been added to an entity. This is also called when an entity
	// is added which has a matching component.
	ComponentAdded(ctx context.Context) context.Context

	// ComponentUpdated is called when a matching component has been updated on an entity. It is not called when an entity
	// did not have a matching component and one is added.
	ComponentUpdated(ctx context.Context) context.Context

	// ComponentRemoved is called when a matching component has been removed from an entity. It is also called when an
	// entity is removed which had a matching component.
	ComponentRemoved(ctx context.Context) context.Context
}
