package common

import "fmt"

type BackendError struct {
	Identifier, Message string
	Code                int
	Err                 error
}

func (e *BackendError) Error() string {
	return fmt.Sprintf("Error %s: %s.\n%v", e.Identifier, e.Message, e.Err)
}

func NewBackendError(code int, identifier, message string, err error, a ...any) *BackendError {
	return &BackendError{
		Identifier: identifier,
		Code:       code,
		Message:    fmt.Sprintf(message, a...),
		Err:        err,
	}
}
