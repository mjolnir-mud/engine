package entity_type

// EntityType is a constructor for entities. The Create function should return a map of the components for the entity
// that when passed to the `ecs.AddEntity` function will create the entity by adding the returned key value pairs
// as components.
type EntityType interface {
	// Name returns the name of the entity type.
	Name() string

	// New returns a map of the components for the entity that can then be added to the game instance.
	New(args map[string]interface{}) map[string]interface{}

	// Validate returns an error if the args are invalid for the entity type.
	Validate(args map[string]interface{}) error
}
