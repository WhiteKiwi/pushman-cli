package cli

import (
	"context"
	"errors"
	"fmt"
)

type UsageError struct{ Err error }

func (e *UsageError) Error() string { return e.Err.Error() }
func (e *UsageError) Unwrap() error { return e.Err }

type AuthorizationError struct{ Err error }

func (e *AuthorizationError) Error() string { return e.Err.Error() }
func (e *AuthorizationError) Unwrap() error { return e.Err }

func usagef(format string, args ...any) error {
	return &UsageError{Err: fmt.Errorf(format, args...)}
}

func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	var usage *UsageError
	if errors.As(err, &usage) {
		return 2
	}
	var auth *AuthorizationError
	if errors.As(err, &auth) {
		return 4
	}
	if errors.Is(err, context.Canceled) {
		return 130
	}
	return 1
}
