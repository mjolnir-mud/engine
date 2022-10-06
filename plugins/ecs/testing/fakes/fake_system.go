package fakes

type ComponentAddedCall struct {
	EntityId string
	Key      string
	Value    interface{}
}

type ComponentUpdatedCall struct {
	EntityId string
	Key      string
	OldValue interface{}
	NewValue interface{}
}

type ComponentRemovedCall struct {
	EntityId string
	Key      string
}

type FakeSystem struct {
	ComponentAddedCalled   chan ComponentAddedCall
	ComponentUpdatedCalled chan ComponentUpdatedCall
	ComponentRemovedCalled chan ComponentRemovedCall
}

func NewFakeSystem() *FakeSystem {
	return &FakeSystem{
		ComponentAddedCalled:   make(chan ComponentAddedCall, 1),
		ComponentUpdatedCalled: make(chan ComponentUpdatedCall, 1),
		ComponentRemovedCalled: make(chan ComponentRemovedCall, 1),
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

func (s FakeSystem) ComponentAdded(entityId string, key string, value interface{}) error {
	go func() {
		s.ComponentAddedCalled <- ComponentAddedCall{
			EntityId: entityId,
			Key:      key,
			Value:    value,
		}
	}()

	return nil
}

func (s FakeSystem) ComponentUpdated(entityId string, key string, oldValue interface{}, newValue interface{}) error {
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

func (s FakeSystem) ComponentRemoved(entityId string, key string) error {
	go func() {
		s.ComponentRemovedCalled <- ComponentRemovedCall{
			EntityId: entityId,
			Key:      key,
		}
	}()

	return nil
}

func (s FakeSystem) MatchingComponentAdded(entityId string, value interface{}) error {
	go func() {
		s.ComponentAddedCalled <- ComponentAddedCall{
			EntityId: entityId,
			Key:      s.Component(),
			Value:    value,
		}
	}()

	return nil
}

func (s FakeSystem) MatchingComponentUpdated(entityId string, oldValue interface{}, newValue interface{}) error {
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

func (s FakeSystem) MatchingComponentRemoved(entityId string) error {
	go func() {
		s.ComponentRemovedCalled <- ComponentRemovedCall{
			EntityId: entityId,
			Key:      s.Component(),
		}
	}()

	return nil
}
