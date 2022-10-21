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
	"context"
	"fmt"
	"github.com/mjolnir-engine/engine/uid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func teardown(ds DataSource, engine *Engine) {
	c, err := mongo.NewClient(options.
		Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%d", engine.config.Mongo.Host, engine.config.Mongo.Port)),
	)

	if err != nil {
		panic(err)
	}

	err = c.Connect(context.Background())

	if err != nil {
		panic(err)
	}

	err = c.Database(engine.config.Mongo.Database).Drop(context.Background())

	if err != nil {
		panic(err)
	}

	engine.Stop()
}

func TestMongoDataSource_Name(t *testing.T) {
	ds, e := setup()
	defer teardown(ds, e)

	assert.Equal(t, "test", ds.Name())
}

func TestMongoDataSource_Save(t *testing.T) {
	ds, engine := setup()
	defer teardown(ds, engine)

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

func TestMongoDataSource_Find(t *testing.T) {
	ds, engine := setup()
	defer teardown(ds, engine)

	entity := &TestMongoEntity{
		Value: "test",
	}

	_, err := ds.Save(&entity)

	assert.Nil(t, err)

	var results []TestMongoEntity

	err = ds.Find(map[string]interface{}{"value": "test"}, &results)

	assert.Nil(t, err)
	assert.Len(t, results, 1)
}
