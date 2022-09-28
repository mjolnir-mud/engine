package errors

import "fmt"

type EntityNotFoundError struct {
	ID string
}

func (e EntityNotFoundError) Error() string {
	return fmt.Sprintf("entity %s not found", e.ID)
}
