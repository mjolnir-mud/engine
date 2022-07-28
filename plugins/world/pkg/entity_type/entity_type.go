package entity_type

// EntityType is a constructor for entities. The Create function should return a map of the components for the entity
// that when passed to the `entity_registry.Add` function will create the entity by adding its components.
type EntityType interface {
	// Name returns the name of the entity type.
	Name() string

	// Create returns a map of the components for the entity that when passed to the `entity_registry.Add` function. It
	// takes map of arguments that can be used to create the entity.
	Create(id string, args map[string]interface{}) map[string]interface{}
}
