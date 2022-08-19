package errors

type IDRequiredError struct{}

func (e IDRequiredError) Error() string {
	return "entity ID required"
}
