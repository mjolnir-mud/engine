package errors

import "fmt"

// MapHasKeyError is an error type that is returned when the map already has the given key.
type MapHasKeyError struct {
	ID   string
	Name string
	Key  string
}

func (e MapHasKeyError) Error() string {
	return fmt.Sprintf("map %s for entity %s already has key %s", e.Name, e.ID, e.Key)
}
