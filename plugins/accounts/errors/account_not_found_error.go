package errors

type AccountNotFoundError struct{}

func (e AccountNotFoundError) Error() string {
	return "account not found"
}
