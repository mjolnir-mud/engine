package errors

import "fmt"

// EntityTypeNotRegisteredError is called when an entity type is not registered.
type EntityTypeNotRegisteredError struct {
	Type string
}

func (e EntityTypeNotRegisteredError) Error() string {
	return fmt.Sprintf("entity type %s is not registered", e.Type)
}
