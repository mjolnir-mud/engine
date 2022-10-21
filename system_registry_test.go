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

type componentAddedCall struct {
	EntityId uid.UID
	Key      string
	Value    interface{}
}

type componentUpdatedCall struct {
	EntityId uid.UID
	Key      string
	OldValue interface{}
	NewValue interface{}
}

type componentRemovedCall struct {
	EntityId uid.UID
	Key      string
	Value    interface{}
}

type fakeSystem struct {
	ComponentAddedCalled   chan componentAddedCall
	ComponentUpdatedCalled chan componentUpdatedCall
	ComponentRemovedCalled chan componentRemovedCall
}

func (s fakeSystem) Name() string {
	return "testing"
}

func (s fakeSystem) Component() string {
	return "testComponent"
}

func (s fakeSystem) Match(_ string, _ interface{}) bool {
	return true
}

func (s fakeSystem) ComponentAdded(entityId uid.UID, key string, value interface{}) error {
	go func() {
		s.ComponentAddedCalled <- componentAddedCall{
			EntityId: entityId,
			Key:      key,
			Value:    value,
		}
	}()

	return nil
}

func (s fakeSystem) ComponentUpdated(entityId uid.UID, key string, oldValue interface{}, newValue interface{}) error {
	go func() {
		s.ComponentUpdatedCalled <- componentUpdatedCall{
			EntityId: entityId,
			Key:      key,
			OldValue: oldValue,
			NewValue: newValue,
		}
	}()

	return nil
}

func (s fakeSystem) ComponentRemoved(entityId uid.UID, key string, value interface{}) error {
	go func() {
		s.ComponentRemovedCalled <- componentRemovedCall{
			EntityId: entityId,
			Key:      key,
			Value:    value,
		}
	}()

	return nil
}

func (s fakeSystem) MatchingComponentAdded(entityId uid.UID, value interface{}) error {
	go func() {
		s.ComponentAddedCalled <- componentAddedCall{
			EntityId: entityId,
			Key:      s.Component(),
			Value:    value,
		}
	}()

	return nil
}

func (s fakeSystem) MatchingComponentUpdated(entityId uid.UID, oldValue interface{}, newValue interface{}) error {
	go func() {
		s.ComponentUpdatedCalled <- componentUpdatedCall{
			EntityId: entityId,
			Key:      s.Component(),
			OldValue: oldValue,
			NewValue: newValue,
		}
	}()

	return nil
}

func (s fakeSystem) MatchingComponentRemoved(entityId uid.UID, value interface{}) error {
	go func() {
		s.ComponentRemovedCalled <- componentRemovedCall{
			EntityId: entityId,
			Key:      s.Component(),
			Value:    value,
		}
	}()

	return nil
}

func TestSystemRegistry_Register(t *testing.T) {
	engine := createEngineInstance()
	engine.Start("test")

	defer engine.Stop()

	engine.RegisterSystem(fakeSystem{
		ComponentAddedCalled:   make(chan componentAddedCall),
		ComponentUpdatedCalled: make(chan componentUpdatedCall),
		ComponentRemovedCalled: make(chan componentRemovedCall),
	})

	assert.Equal(t, 1, len(engine.systemRegistry.systems))
}

func TestSystemRegistry_ComponentAdded(t *testing.T) {
	engine := createEngineInstance()
	engine.Start("test")
	defer engine.Stop()

	fs := fakeSystem{
		ComponentAddedCalled:   make(chan componentAddedCall),
		ComponentUpdatedCalled: make(chan componentUpdatedCall),
		ComponentRemovedCalled: make(chan componentRemovedCall),
	}

	engine.RegisterSystem(fs)

	id, err := engine.AddEntity(fakeEntity{
		Value: "test",
	})

	assert.Nil(t, err)
	err = engine.AddComponent(id, fs.Component(), "test")

	assert.Nil(t, err)

	v := <-fs.ComponentAddedCalled

	assert.Equal(t, id, v.EntityId)

}
