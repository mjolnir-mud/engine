package test

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

type TestSystem struct {
	ComponentAddedCalled   chan ComponentAddedCall
	ComponentUpdatedCalled chan ComponentUpdatedCall
}

func NewTestSystem() *TestSystem {
	return &TestSystem{
		ComponentAddedCalled: make(chan ComponentAddedCall, 1),
	}
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
	//s.appendCall("ComponentRemoved", map[string]interface{}{
	//	"entityId": entityId,
	//	"key":      key,
	//})

	return nil
}

func (s TestSystem) MatchingComponentAdded(entityId string, value interface{}) error {
	//s.appendCall("MatchingComponentAdded", map[string]interface{}{
	//	"entityId": entityId,
	//	"value":    value,
	//})

	return nil
}

func (s TestSystem) MatchingComponentUpdated(entityId string, oldValue interface{}, newValue interface{}) error {
	//s.appendCall("MatchingComponentUpdated", map[string]interface{}{
	//	"entityId": entityId,
	//	"oldValue": oldValue,
	//	"newValue": newValue,
	//})

	return nil
}

func (s TestSystem) MatchingComponentRemoved(entityId string) error {
	//s.appendCall("MatchingComponentRemoved", map[string]interface{}{
	//	"entityId": entityId,
	//})

	return nil
}
