package errors

import "fmt"

// MapValueTypeMismatchError is an error type that is returned when the map key does not match the entity type.
type MapValueTypeMismatchError struct {
	ID       string
	Name     string
	Key      string
	Expected string
	Actual   string
}

func (e MapValueTypeMismatchError) Error() string {
	return fmt.Sprintf(
		"map value type mismatch for entity %s, component %s, key %s, expected %s, actual %s",
		e.ID,
		e.Name,
		e.Key,
		e.Expected,
		e.Actual,
	)
}
