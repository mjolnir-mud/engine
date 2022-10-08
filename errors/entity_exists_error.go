package errors

import "fmt"

// EntityExistsError is an error type that is returned when an entity with the given id already exists.
type EntityExistsError struct {
	Id string
}

func (e EntityExistsError) Error() string {
	return fmt.Sprintf("entity with id %s already exists", e.Id)
}
