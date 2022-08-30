package errors

import "fmt"

type SessionNotFoundError struct {
	SessionId string
}

func (e SessionNotFoundError) Error() string {
	return fmt.Sprintf("session not found: %s", e.SessionId)
}
