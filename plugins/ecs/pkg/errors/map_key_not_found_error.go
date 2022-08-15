package errors

import "fmt"

type MapKeyNotFoundError struct {
	ID   string
	Name string
	Key  string
}

func (e MapKeyNotFoundError) Error() string {
	return fmt.Sprintf("map %s for entity %s does not have key %s", e.Name, e.ID, e.Key)
}
