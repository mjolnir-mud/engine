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

package data_source

import (
	"context"

	"github.com/mjolnir-engine/engine/pkg/uid"
)

// DataSource is an interface that represents a persistent data store. Data sources are used to store and retrieve
// entities that may not be actively loaded into memory.
type DataSource interface {
	// All returns all entities within the data source.
	All(ctx context.Context, response interface{}) error

	// Count returns the number of entities within the data source based on the provided filter. The filter is a map,
	// the data source is responsible for translating that map into a filter to be used against its search.
	Count(ctx context.Context, query interface{}) (int64, error)

	// Delete deletes entities from the data source based on the provided filter.
	Delete(ctx context.Context, query interface{}) error

	// Find returns a list of entities from executing a search against a provided map. It returns a list of entities as
	// a map keyed by their ids.
	Find(ctx context.Context, query interface{}, result interface{}) error

	// Name returns the name of the data source. The name must be unique. Registering a data source with the same name
	// will replace the existing data source of the same name.
	Name() string

	// Save saves an entity to the data source.
	Save(ctx context.Context, entity interface{}) (uid.UID, error)
}
