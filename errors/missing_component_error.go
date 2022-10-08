package errors

import "fmt"

// MissingComponentError is an error type that is returned when the entity does not have the given component.
type MissingComponentError struct {
	ID   string
	Name string
}

func (e MissingComponentError) Error() string {
	return fmt.Sprintf("entity %s does not have component %s", e.ID, e.Name)
}
