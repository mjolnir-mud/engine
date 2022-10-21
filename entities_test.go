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
	engineEvents "github.com/mjolnir-engine/engine/events"
	"github.com/mjolnir-engine/engine/uid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakeEntity struct {
	Value string
}

func TestEngine_AddComponent(t *testing.T) {
	e := createEngineInstance()
	e.Start("test")
	defer e.Stop()

	id := uid.New()

	err := e.AddEntityWithId(id, fakeEntity{
		Value: "test",
	})

	assert.Nil(t, err)

	receivedMsg := make(chan EventMessage)

	e.Subscribe(engineEvents.ComponentAddedEvent{EntityId: id, Name: "Value2"}, func(event EventMessage) {
		receivedMsg <- event
	})

	err = e.AddComponent(id, "Value2", "test2")

	assert.Nil(t, err)

	msg := <-receivedMsg
	ca := engineEvents.ComponentAddedEvent{}

	err = msg.Unmarshal(&ca)

	assert.Nil(t, err)
	assert.Equal(t, id, ca.EntityId)
	assert.Equal(t, "Value2", ca.Name)

	exists, err := e.HasComponent(id, "Value2")

	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestEngine_AddEntity(t *testing.T) {
	e := createEngineInstance()
	e.Start("test")
	defer e.Stop()

	id, err := e.AddEntity(fakeEntity{
		Value: "test",
	})

	assert.Nil(t, err)
	assert.NotNil(t, id)
}

func TestEngine_AddEntityWithId(t *testing.T) {
	e := createEngineInstance()
	e.Start("test")
	defer e.Stop()

	id := uid.New()

	entityAdded := make(chan EventMessage, 0)
	componentAdded := make(chan EventMessage, 0)

	e.Subscribe(engineEvents.EntityAddedEvent{
		Id: id,
	}, func(event EventMessage) {
		go func() { entityAdded <- event }()
	})

	e.Subscribe(engineEvents.ComponentAddedEvent{EntityId: id, Name: "Value"}, func(event EventMessage) {
		go func() { componentAdded <- event }()
	})

	err := e.AddEntityWithId(id, fakeEntity{
		Value: "test",
	})

	assert.Nil(t, err)

	assert.NotNil(t, <-entityAdded)
	assert.NotNil(t, <-componentAdded)

	exists, err := e.redis.Do(
		context.Background(),
		e.redis.B().Exists().Key(string(id)).Build(),
	).AsBool()

	assert.NoError(t, err)
	assert.True(t, exists)

	err = e.AddEntityWithId(id, fakeEntity{
		Value: "test",
	})

	assert.Error(t, err, "entity with id 123 already exists")
}

func TestEngine_HasComponent(t *testing.T) {
	e := createEngineInstance()
	e.Start("test")
	defer e.Stop()

	id, err := e.AddEntity(fakeEntity{
		Value: "test",
	})

	assert.Nil(t, err)

	exists, err := e.HasComponent(id, "Value")

	assert.Nil(t, err)
	assert.True(t, exists)

	exists, err = e.HasComponent(id, "Value2")

	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestEngine_HasEntity(t *testing.T) {
	e := createEngineInstance()
	e.Start("test")
	defer e.Stop()

	added := make(chan EventMessage, 0)

	e.PSubscribe(engineEvents.ComponentAddedEvent{}, func(event EventMessage) {
		go func() { added <- event }()
	})

	id, err := e.AddEntity(fakeEntity{
		Value: "test",
	})

	assert.Nil(t, err)

	<-added

	exists, err := e.HasEntity(id)

	assert.Nil(t, err)
	assert.True(t, exists)

	exists, err = e.HasEntity(uid.New())

	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestEngine_UpdateComponent(t *testing.T) {
	e := createEngineInstance()
	e.Start("test")
	defer e.Stop()

	id, err := e.AddEntity(fakeEntity{
		Value: "test",
	})

	assert.Nil(t, err)

	receivedMsg := make(chan EventMessage)

	e.Subscribe(engineEvents.ComponentUpdatedEvent{EntityId: id, Name: "Value"}, func(event EventMessage) {
		go func() { receivedMsg <- event }()
	})

	err = e.UpdateComponent(id, "Value", "test2")

	assert.Nil(t, err)

	msg := <-receivedMsg
	ca := engineEvents.ComponentUpdatedEvent{}

	err = msg.Unmarshal(&ca)

	assert.Nil(t, err)
	assert.Equal(t, id, ca.EntityId)
	assert.Equal(t, "Value", ca.Name)
	assert.Equal(t, "test2", ca.Value)
	assert.Equal(t, "test", ca.PreviousValue)

	exists, err := e.HasComponent(id, "Value")

	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestEngine_GetComponent(t *testing.T) {
	e := createEngineInstance()
	e.Start("test")
	defer e.Stop()

	id, err := e.AddEntity(&fakeEntity{
		Value: "test",
	})

	assert.Nil(t, err)

	var component string

	err = e.GetComponent(id, "Value", &component)

	assert.Nil(t, err)
	assert.Equal(t, "test", component)
}

func TestEngine_RemoveComponent(t *testing.T) {
	e := createEngineInstance()
	e.Start("test")
	defer e.Stop()

	id, err := e.AddEntity(&fakeEntity{
		Value: "test",
	})

	assert.Nil(t, err)

	receivedMsg := make(chan EventMessage)

	e.Subscribe(engineEvents.ComponentRemovedEvent{EntityId: id, Name: "Value"}, func(event EventMessage) {
		receivedMsg <- event
	})

	err = e.RemoveComponent(id, "Value", "")

	assert.Nil(t, err)

	msg := <-receivedMsg
	ca := engineEvents.ComponentRemovedEvent{}

	err = msg.Unmarshal(&ca)

	assert.Nil(t, err)
	assert.Equal(t, id, ca.EntityId)
	assert.Equal(t, "Value", ca.Name)
	assert.Equal(t, "test", ca.Value)

	exists, err := e.HasComponent(id, "Value")

	assert.Nil(t, err)
	assert.False(t, exists)
}
