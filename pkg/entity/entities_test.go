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
	//"github.com/mjolnir-engine/engine/pkg/redis"

	"context"
	"testing"

	"github.com/mjolnir-engine/engine/internal/redis"
	internalTesting "github.com/mjolnir-engine/engine/internal/testing"
	"github.com/mjolnir-engine/engine/pkg/uid"
)

type fakeEntity struct {
	Id    uid.UID
	Value string
}

type updatedFakeEntity struct {
	Id         uid.UID
	Value      string
	OtherValue string
}

func setup(t *testing.T) context.Context {
	ctx := internalTesting.Setup(t)
	ctx = redis.WithRedis(ctx)

	return ctx
}

func teardown(t *testing.T, ctx context.Context) {
	internalTesting.Teardown(t, ctx)
}

//func TestEngine_AddEntity_WithoutId(t *testing.T) {
//	e := engine.createEngineInstance()
//	e.Start("test")
//	defer e.Stop()
//
//	fe := &fakeEntity{
//		Value: "test",
//	}
//
//	err := e.AddEntity(fe)
//
//	assert.Nil(t, err)
//	assert.NotEqual(t, uid.UID(""), fe.Id)
//
//	exists, err := e.HasEntity(fe.Id)
//
//	assert.Nil(t, err)
//	assert.True(t, exists)
//}
//
//func TestEngine_AddEntity_AlreadyExists(t *testing.T) {
//	e := engine.createEngineInstance()
//	e.Start("test")
//	defer e.Stop()
//
//	fe := &fakeEntity{
//		Value: "test",
//	}
//
//	err := e.AddEntity(fe)
//
//	assert.Nil(t, err)
//
//	err = e.AddEntity(&fakeEntity{
//		Id:    fe.Id,
//		Value: "test",
//	})
//
//	assert.NotNil(t, err)
//}
//
//func TestEngine_AddEntity_NotStruct(t *testing.T) {
//	e := engine.createEngineInstance()
//	e.Start("test")
//	defer e.Stop()
//
//	err := e.AddEntity("test")
//
//	assert.NotNil(t, err)
//}
//
//func TestEngine_AddEntity_Nil(t *testing.T) {
//	e := engine.createEngineInstance()
//	e.Start("test")
//	defer e.Stop()
//
//	err := e.AddEntity(nil)
//
//	assert.NotNil(t, err)
//}
//
//func TestEngine_HasComponent(t *testing.T) {
//	e := engine.createEngineInstance()
//	e.Start("test")
//	defer e.Stop()
//
//	fe := &fakeEntity{
//		Value: "test",
//	}
//
//	err := e.AddEntity(fe)
//
//	assert.Nil(t, err)
//
//	exists, err := e.HasComponent(fe.Id, "Value")
//
//	assert.Nil(t, err)
//	assert.True(t, exists)
//
//	exists, err = e.HasComponent(fe.Id, "Value2")
//
//	assert.Nil(t, err)
//	assert.False(t, exists)
//}
//
//func TestEngine_HasEntity(t *testing.T) {
//	e := engine.createEngineInstance()
//	e.Start("test")
//	defer e.Stop()
//
//	added := make(chan redis.EventMessage, 0)
//
//	e.PSubscribe(events.ComponentAddedEvent{}, func(event redis.EventMessage) {
//		go func() { added <- event }()
//	})
//
//	fe := &fakeEntity{
//		Value: "test",
//	}
//
//	err := e.AddEntity(fe)
//
//	assert.Nil(t, err)
//
//	<-added
//
//	exists, err := e.HasEntity(fe.Id)
//
//	assert.Nil(t, err)
//	assert.True(t, exists)
//
//	exists, err = e.HasEntity(uid.New())
//
//	assert.Nil(t, err)
//	assert.False(t, exists)
//}
//
//func TestEngine_UpdateEntity(t *testing.T) {
//	e := engine.createEngineInstance()
//	e.Start("test")
//	defer e.Stop()
//
//	fe := &fakeEntity{
//		Value: "test",
//	}
//
//	err := e.AddEntity(fe)
//
//	assert.Nil(t, err)
//
//	receivedMsg := make(chan redis.EventMessage)
//
//	e.Subscribe(events.ComponentUpdatedEvent{EntityId: fe.Id, Name: "Value"}, func(event redis.EventMessage) {
//		go func() { receivedMsg <- event }()
//	})
//
//	fe.Value = "test2"
//
//	err = e.Update(fe)
//
//	assert.Nil(t, err)
//
//	msg := <-receivedMsg
//	ca := events.ComponentUpdatedEvent{}
//
//	err = msg.Unmarshal(&ca)
//
//	assert.Nil(t, err)
//	assert.Equal(t, fe.Id, ca.EntityId)
//	assert.Equal(t, "Value", ca.Name)
//	assert.Equal(t, "test2", ca.Value)
//	assert.Equal(t, "test", ca.PreviousValue)
//}
//
//func TestEngine_GetComponent(t *testing.T) {
//	e := engine.createEngineInstance()
//	e.Start("test")
//	defer e.Stop()
//
//	fe := &fakeEntity{
//		Value: "test",
//	}
//
//	err := e.AddEntity(fe)
//
//	assert.Nil(t, err)
//
//	var component string
//
//	err = e.GetComponent(fe.Id, "Value", &component)
//
//	assert.Nil(t, err)
//	assert.Equal(t, "test", component)
//}
