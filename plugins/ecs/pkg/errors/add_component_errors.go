package errors

import (
	"fmt"
	"strings"
)

// AddComponentErrors is a collection of errors that occurred while adding components to an entity.
type AddComponentErrors struct {
	Errors []error
}

func (e AddComponentErrors) Error() string {
	errorStrings := make([]string, len(e.Errors))

	for i, err := range e.Errors {
		errorStrings[i] = err.Error()
	}

	return fmt.Sprintf(
		"%d errors occurred while adding components to an entity: %s", len(e.Errors),
		strings.Join(errorStrings, ", "),
	)
}
