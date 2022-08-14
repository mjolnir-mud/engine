package errors

import "fmt"

// ComponentTypeMismatchError is an error type that is returned when the component type does not match the entity type.
type ComponentTypeMismatchError struct {
	ID       string
	Name     string
	Expected string
	Actual   string
}

func (e ComponentTypeMismatchError) Error() string {
	return fmt.Sprintf(
		"component type mismatch for entity %s, component %s, expected %s, actual %s",
		e.ID,
		e.Name,
		e.Expected,
		e.Actual,
	)
}
