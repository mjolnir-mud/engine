package errors

import "fmt"

// ComponentNotFoundError is returned when a component is not found.
type ComponentNotFoundError struct {
	EntityId string
	Name     string
}

func (e ComponentNotFoundError) Error() string {
	return fmt.Sprintf("component %s not found for entity %s", e.Name, e.EntityId)
}
