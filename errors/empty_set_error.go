package errors

import "fmt"

// EmptySetError is an error type that is returned when an empty set is attempted to be added to an entity .
type EmptySetError struct {
	ID string
}

func (e EmptySetError) Error() string {
	return fmt.Sprintf("empty set cannot be added %s", e.ID)
}
