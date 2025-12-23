package errors

import (
	"errors"
	"fmt"
)

// Define standard sentinel errors for the domain.
var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("resource conflict")
	ErrInternal     = errors.New("internal system error")
	ErrUnauthorized = errors.New("unauthorized")
)

// Is reports whether any error in err's chain matches target.
// It is a wrapper for errors.Is.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// Wrap wraps an error with a message.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}
