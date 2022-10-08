package errors

import "fmt"

// HasComponentError is an error type that is returned when the entity already has the given component.
type HasComponentError struct {
	ID   string
	Name string
}

func (e HasComponentError) Error() string {
	return fmt.Sprintf("entity %s already has component %s", e.ID, e.Name)
}
