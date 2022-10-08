package errors

import "fmt"

// ComponentExistsError is an error type that is returned when the entity already has the given component.
type ComponentExistsError struct {
	EntityId string
	Name     string
}

func (e ComponentExistsError) Error() string {
	return fmt.Sprintf("component %s exists on entity %s", e.EntityId, e.Name)
}
