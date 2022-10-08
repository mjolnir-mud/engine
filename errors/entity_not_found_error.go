package errors

import "fmt"

// EntityNotFoundError is an error type that is returned when an entity with the given id does not exist.
type EntityNotFoundError struct {
	Id string
}

func (e EntityNotFoundError) Error() string {
	return fmt.Sprintf("entity with id %s does not exist", e.Id)
}
