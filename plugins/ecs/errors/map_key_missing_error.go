package errors

import "fmt"

type MapKeyMissingError struct {
	ID   string
	Name string
	Key  string
}

func (e MapKeyMissingError) Error() string {
	return fmt.Sprintf(
		"map key missing for entity %s, component %s, key %s",
		e.ID,
		e.Name,
		e.Key,
	)
}
