package errors

import (
	"fmt"
	"github.com/mjolnir-engine/engine/uid"
)

// EntityExistsError is an error type that is returned when an entity with the given id already exists.
type EntityExistsError struct {
	Id uid.UID
}

func (e EntityExistsError) Error() string {
	return fmt.Sprintf("entity with id %s already exists", e.Id)
}
