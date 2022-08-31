package testing

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

type TestSystem struct {
	ComponentAddedCalled   chan ComponentAddedCall
	ComponentUpdatedCalled chan ComponentUpdatedCall
	ComponentRemovedCalled chan ComponentRemovedCall
}

func NewTestSystem() *TestSystem {
	return &TestSystem{
		ComponentAddedCalled:   make(chan ComponentAddedCall, 1),
		ComponentUpdatedCalled: make(chan ComponentUpdatedCall, 1),
		ComponentRemovedCalled: make(chan ComponentRemovedCall, 1),
	}
}

func (s TestSystem) Name() string {
	return "testing"
}

func (s TestSystem) Component() string {
	return "testComponent"
}

func (s TestSystem) Match(key string, value interface{}) bool {
	return true
}

func (s TestSystem) WorldStarted() {}

func (s TestSystem) ComponentAdded(entityId string, key string, value interface{}) error {
	go func() {
		s.ComponentAddedCalled <- ComponentAddedCall{
			EntityId: entityId,
			Key:      key,
			Value:    value,
		}
	}()

	return nil
}

func (s TestSystem) ComponentUpdated(entityId string, key string, oldValue interface{}, newValue interface{}) error {
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

func (s TestSystem) ComponentRemoved(entityId string, key string) error {
	go func() {
		s.ComponentRemovedCalled <- ComponentRemovedCall{
			EntityId: entityId,
			Key:      key,
		}
	}()

	return nil
}

func (s TestSystem) MatchingComponentAdded(entityId string, value interface{}) error {
	go func() {
		s.ComponentAddedCalled <- ComponentAddedCall{
			EntityId: entityId,
			Key:      s.Component(),
			Value:    value,
		}
	}()

	return nil
}

func (s TestSystem) MatchingComponentUpdated(entityId string, oldValue interface{}, newValue interface{}) error {
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

func (s TestSystem) MatchingComponentRemoved(entityId string) error {
	go func() {
		s.ComponentRemovedCalled <- ComponentRemovedCall{
			EntityId: entityId,
			Key:      s.Component(),
		}
	}()

	return nil
}
