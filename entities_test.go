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
	engineEvents "github.com/mjolnir-mud/engine/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func teardown(e *Engine) {
	_ = e.FlushEntities()
}

type fakeEntity struct {
	Value string
}

func TestEngine_AddComponent(t *testing.T) {
	e := createEngineInstance()
	defer teardown(e)

	err := e.AddEntityWithId("123", fakeEntity{
		Value: "test",
	})

	assert.Nil(t, err)

	receivedMsg := make(chan EventMessage)

	e.Subscribe(engineEvents.ComponentAddedEvent{EntityId: "123", Name: "Value"}, func(event EventMessage) {
		receivedMsg <- event
	})

	err = e.AddComponent("123", "Value2", "test2")

	assert.Nil(t, err)

	msg := <-receivedMsg
	ca := engineEvents.ComponentAddedEvent{}

	err = msg.Unmarshal(&ca)

	assert.Nil(t, err)
	assert.Equal(t, "123", ca.EntityId)
	assert.Equal(t, "Value2", ca.Name)

	exists, err := e.HasComponent("123", "Value2")

	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestEngine_AddEntity(t *testing.T) {
	e := createEngineInstance()
	defer teardown(e)

	id, err := e.AddEntity(fakeEntity{
		Value: "test",
	})

	assert.Nil(t, err)
	assert.NotNil(t, id)
}

func TestEngine_AddEntityWithId(t *testing.T) {
	e := createEngineInstance()
	defer teardown(e)

	receivedMsg := make(chan EventMessage)

	e.Subscribe(engineEvents.EntityAddedEvent{}, func(event EventMessage) {
		receivedMsg <- event
	})

	e.Subscribe(engineEvents.ComponentAddedEvent{EntityId: "123", Name: "Value"}, func(event EventMessage) {
		receivedMsg <- event
	})

	err := e.AddEntityWithId("123", fakeEntity{
		Value: "test",
	})

	msg := <-receivedMsg
	ea := engineEvents.EntityAddedEvent{}

	err = msg.Unmarshal(&ea)

	assert.Nil(t, err)
	assert.Equal(t, "123", ea.Id)

	msg = <-receivedMsg
	ca := engineEvents.ComponentAddedEvent{}

	err = msg.Unmarshal(&ca)

	assert.Nil(t, err)
	assert.Equal(t, "123", ca.EntityId)
	assert.Equal(t, "Value", ca.Name)

	exists, err := e.redis.Do(
		context.Background(),
		e.redis.B().Exists().Key(e.stringToKey("123")).Build(),
	).AsBool()

	assert.NoError(t, err)
	assert.True(t, exists)

	err = e.AddEntityWithId("123", fakeEntity{
		Value: "test",
	})

	assert.Error(t, err, "entity with id 123 already exists")
}

func TestEngine_HasComponent(t *testing.T) {
	e := createEngineInstance()
	defer teardown(e)

	err := e.AddEntityWithId("123", fakeEntity{
		Value: "test",
	})

	assert.Nil(t, err)

	exists, err := e.HasComponent("123", "Value")

	assert.Nil(t, err)
	assert.True(t, exists)

	exists, err = e.HasComponent("123", "Value2")

	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestEngine_HasEntity(t *testing.T) {
	e := createEngineInstance()
	defer teardown(e)

	err := e.AddEntityWithId("123", fakeEntity{
		Value: "test",
	})

	assert.Nil(t, err)

	exists, err := e.HasEntity("123")

	assert.Nil(t, err)
	assert.True(t, exists)

	exists, err = e.HasEntity("1234")

	assert.Nil(t, err)
	assert.False(t, exists)
}
