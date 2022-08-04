package system

type System interface {
	Name() string

	// Component is the type of the component that this system is responsible for.
	Component() string

	// Match is a check to see if the system should be run for the given key and value.
	Match(key string, value interface{}) bool

	// WorldStarted is called when the world is started.
	WorldStarted()

	// ComponentAdded is called when any component is added to an entity, including those the system is interested in.
	ComponentAdded(entityId string, component string, value interface{}) error

	// ComponentUpdated is called when a component is updated on an entity, including those the system is interested in.
	ComponentUpdated(entityId string, component string, oldValue interface{}, newValue interface{}) error

	// ComponentRemoved is called when a component is removed from an entity, including those the system is interested
	// in.
	ComponentRemoved(entityId string, component string) error

	// MatchingComponentAdded is called when a component is added to an entity, but only if the system is interested
	// in the component.
	MatchingComponentAdded(entityId string, component string, value interface{}) error

	// MatchingComponentUpdated is called when a component is updated on an entity, but only if the system is interested
	// in the component.
	MatchingComponentUpdated(entityId string, component string, oldValue interface{}, newValue interface{}) error

	// MatchingComponentRemoved is called when a component is removed from an entity, but only if the system is
	// interested in the component.
	MatchingComponentRemoved(entityId string, component string) error
}
