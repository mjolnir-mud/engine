package errors

type ControllerNotFoundError struct {
	Controller string
}

func (e ControllerNotFoundError) Error() string {
	return "controller not found: " + e.Controller
}
