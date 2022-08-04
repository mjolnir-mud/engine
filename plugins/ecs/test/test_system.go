package test

import (
	"fmt"
	"reflect"
)

type call struct {
	args map[string]interface{}
}

type TestSystem struct {
	calls map[string][]call
}

func (s TestSystem) Name() string {
	return "test"
}

func (s TestSystem) Component() string {
	return "testComponent"
}

func (s TestSystem) Match(key string, value interface{}) bool {
	return true
}

func (s TestSystem) WorldStarted() {}

func (s TestSystem) ComponentAdded(entityId string, key string, value interface{}) error {
	s.appendCall("ComponentAdded", map[string]interface{}{
		"entityId": entityId,
		"key":      key,
		"value":    value,
	})

	return nil
}

func (s TestSystem) ComponentUpdated(entityId string, key string, oldValue interface{}, newValue interface{}) error {
	s.appendCall("ComponentUpdated", map[string]interface{}{
		"entityId": entityId,
		"key":      key,
		"oldValue": oldValue,
		"newValue": newValue,
	})

	return nil
}

func (s TestSystem) ComponentRemoved(entityId string, key string) error {
	s.appendCall("ComponentRemoved", map[string]interface{}{
		"entityId": entityId,
		"key":      key,
	})

	return nil
}

func (s TestSystem) MatchingComponentAdded(entityId string, key string, value interface{}) error {
	s.appendCall("MatchingComponentAdded", map[string]interface{}{
		"entityId": entityId,
		"key":      key,
		"value":    value,
	})

	return nil
}

func (s TestSystem) MatchingComponentUpdated(entityId string, key string, oldValue interface{}, newValue interface{}) error {
	s.appendCall("MatchingComponentUpdated", map[string]interface{}{
		"entityId": entityId,
		"key":      key,
		"oldValue": oldValue,
		"newValue": newValue,
	})

	return nil
}

func (s TestSystem) MatchingComponentRemoved(entityId string, key string) error {
	s.appendCall("MatchingComponentRemoved", map[string]interface{}{
		"entityId": entityId,
		"key":      key,
	})

	return nil
}

func (s TestSystem) HasBeenCalledWith(function string, args map[string]interface{}) bool {
	if _, ok := s.calls[function]; !ok {
		return false
	}

	for _, call := range s.calls[function] {
		if reflect.DeepEqual(call.args, args) {
			return true
		}
	}

	return false
}

func (s TestSystem) appendCall(function string, args map[string]interface{}) {
	if _, ok := s.calls[function]; !ok {
		s.calls[function] = []call{}
	}

	s.calls[function] = append(s.calls[function], call{args: args})

	fmt.Printf("%s(%v)\n", function, args)
}

func NewTestSystem() TestSystem {
	return TestSystem{
		calls: make(map[string][]call),
	}
}
