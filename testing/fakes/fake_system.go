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

package fakes

import "github.com/mjolnir-mud/engine/uid"

type ComponentAddedCall struct {
	EntityId *uid.UID
	Key      string
	Value    interface{}
}

type ComponentUpdatedCall struct {
	EntityId *uid.UID
	Key      string
	OldValue interface{}
	NewValue interface{}
}

type ComponentRemovedCall struct {
	EntityId *uid.UID
	Key      string
	Value    interface{}
}

type FakeSystem struct {
	ComponentAddedCalled   chan ComponentAddedCall
	ComponentUpdatedCalled chan ComponentUpdatedCall
	ComponentRemovedCalled chan ComponentRemovedCall
}

func NewFakeSystem() *FakeSystem {
	return &FakeSystem{
		ComponentAddedCalled:   make(chan ComponentAddedCall),
		ComponentUpdatedCalled: make(chan ComponentUpdatedCall),
		ComponentRemovedCalled: make(chan ComponentRemovedCall),
	}
}

func (s FakeSystem) Name() string {
	return "testing"
}

func (s FakeSystem) Component() string {
	return "testComponent"
}

func (s FakeSystem) Match(_ string, _ interface{}) bool {
	return true
}

func (s FakeSystem) WorldStarted() {}

func (s FakeSystem) ComponentAdded(entityId *uid.UID, key string, value interface{}) error {
	go func() {
		s.ComponentAddedCalled <- ComponentAddedCall{
			EntityId: entityId,
			Key:      key,
			Value:    value,
		}
	}()

	return nil
}

func (s FakeSystem) ComponentUpdated(entityId *uid.UID, key string, oldValue interface{}, newValue interface{}) error {
	go func() {
		s.ComponentUpdatedCalled <- ComponentUpdatedCall{
			EntityId: entityId,
			Key:      key,
			OldValue: oldValue,
			NewValue: newValue,
		}
	}()

	return nil
}

func (s FakeSystem) ComponentRemoved(entityId *uid.UID, key string, value interface{}) error {
	go func() {
		s.ComponentRemovedCalled <- ComponentRemovedCall{
			EntityId: entityId,
			Key:      key,
			Value:    value,
		}
	}()

	return nil
}

func (s FakeSystem) MatchingComponentAdded(entityId *uid.UID, value interface{}) error {
	go func() {
		s.ComponentAddedCalled <- ComponentAddedCall{
			EntityId: entityId,
			Key:      s.Component(),
			Value:    value,
		}
	}()

	return nil
}

func (s FakeSystem) MatchingComponentUpdated(entityId *uid.UID, oldValue interface{}, newValue interface{}) error {
	go func() {
		s.ComponentUpdatedCalled <- ComponentUpdatedCall{
			EntityId: entityId,
			Key:      s.Component(),
			OldValue: oldValue,
			NewValue: newValue,
		}
	}()

	return nil
}

func (s FakeSystem) MatchingComponentRemoved(entityId *uid.UID, value interface{}) error {
	go func() {
		s.ComponentRemovedCalled <- ComponentRemovedCall{
			EntityId: entityId,
			Key:      s.Component(),
			Value:    value,
		}
	}()

	return nil
}
