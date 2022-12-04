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

package entity

import (
	"context"

	"github.com/mjolnir-engine/engine/pkg/redis"
	"github.com/mjolnir-engine/engine/pkg/uid"
	"github.com/rueian/rueidis"
)

// Record is the record that is stored in Redis for an entity. It is stored as JSON using ReJSON.
type Record struct {
	// Id is the unique identifier for the entity.
	Id uid.UID

	// Version is the version of the entity. This is incremented every time a component is added, removed or updated.
	Version int64

	// Entity is the entity itself.
	Entity interface{}
}

//

//

//	func componentPath(componentName string) string {
//		return fmt.Sprintf(".Entity.%s", componentName)
//	}
func getEntityRecord(ctx context.Context, id uid.UID) Record {
	var record Record

	redis.UnmarshallResult(
		redis.MustExecuteCommands(ctx, rueidis.Commands{
			redis.GetClient(ctx).B().JsonGet().Key(string(id)).Build(),
		})[0],
		&record,
	)

	return record
}
