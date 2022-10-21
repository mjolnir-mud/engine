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

import (
	"github.com/mjolnir-engine/engine/uid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestMongoEntity struct {
	Id    uid.UID `bson:"_id,omitempty"`
	Value string
}

func setup() (DataSource, *Engine) {
	engine := createEngineInstance()
	ds := NewMongoDataSource("test", engine)

	engine.RegisterDataSource(ds)

	engine.Start("test")

	return ds, engine
}

func TestMongoDataSource_Save(t *testing.T) {
	ds, engine := setup()
	defer engine.Stop()

	entity := &TestMongoEntity{
		Value: "test",
	}

	id, err := ds.Save(&entity)

	assert.Nil(t, err)
	assert.NotNil(t, id)

	entity = &TestMongoEntity{}

	err = ds.FindOne(map[string]interface{}{"_id": id}, entity)

	assert.Nil(t, err)
	assert.Equal(t, "test", entity.Value)
}
