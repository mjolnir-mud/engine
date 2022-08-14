package errors

import "fmt"

type SetValueTypeError struct {
	ID       string
	Name     string
	Expected string
	Actual   string
}

func (e SetValueTypeError) Error() string {
	return fmt.Sprintf(
		"set value type mismatch for entity %s, component %s, expected %s, actual %s",
		e.ID,
		e.Name,
		e.Expected,
		e.Actual,
	)
}
