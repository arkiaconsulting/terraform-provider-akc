package client

import (
	"fmt"
)

var (
	// KVNotFoundError The given App Configuration key-value was not found
	KVNotFoundError = AppConfigClientError{Message: "KV not found"}
	// UnexpectedError An unexpected error has occurred
	UnexpectedError = AppConfigClientError{Message: "Unexpected error"}
)

// AppConfigClientError Main type for AppConfigClient errors
type AppConfigClientError struct {
	Message string
	Info    string
	Err     error
}

func (err AppConfigClientError) Error() string {
	if err.Err != nil {
		return fmt.Sprintf("%s (%s): %v", err.Message, err.Info, err.Err.Error())
	}
	return fmt.Sprintf("%s (%s)", err.Message, err.Info)
}
func (err AppConfigClientError) wrap(inner error) error {
	return AppConfigClientError{Message: err.Message, Err: inner}
}
func (err AppConfigClientError) with(info string) error {
	return AppConfigClientError{Message: err.Message, Info: info}
}
func (err AppConfigClientError) Unwrap() error {
	return err.Err
}

// Is Wether the given error is an AppConfigClientError
func (err AppConfigClientError) Is(target error) bool {
	t, ok := target.(AppConfigClientError)
	if !ok {
		return false
	}

	return t.Error() == err.Message
}
