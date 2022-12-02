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
	"reflect"

	"github.com/fugufish/unstruct"
	"github.com/mjolnir-engine/engine/pkg/errors"
	"github.com/mjolnir-engine/engine/pkg/uid"
)

func getComponentMap(entity interface{}) map[string]interface{} {
	return unstruct.ToMapDeep(entity)
}

func getEntityId(entityOrId interface{}) (uid.UID, error) {
	if entityOrId == nil {
		return "", errors.EntityDataTypedError{}
	}

	switch entityOrId.(type) {
	case uid.UID:
		return entityOrId.(uid.UID), nil
	default:
		m := getComponentMap(entityOrId)

		if id, ok := m["Id"]; ok {
			return id.(uid.UID), nil
		}
	}

	return uid.New(), errors.EntityInvalidError{}
}

func mustGetEntityId(entityOrId interface{}) uid.UID {
	id, err := getEntityId(entityOrId)

	if err != nil {
		panic(err)
	}

	return id
}

func typeCheckStructPointer(entity interface{}) error {
	if entity == nil {
		return errors.EntityDataTypedError{
			Expected: "not nil",
		}
	}

	entityType := reflect.TypeOf(entity)

	if entityType.Kind() != reflect.Ptr {
		return errors.EntityDataTypedError{
			Expected: "pointer",
		}
	}

	if entityType.Elem().Kind() != reflect.Struct {
		return errors.EntityDataTypedError{
			Expected: "struct",
		}
	}

	return nil
}
