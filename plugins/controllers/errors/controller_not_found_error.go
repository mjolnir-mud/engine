package errors

type ControllerNotFoundError struct {
	Name string
}

func (e ControllerNotFoundError) Error() string {
	return "controller not found: " + e.Name
}
