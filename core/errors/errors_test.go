package errors

import (
	"errors"
	"testing"
)

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		err      error
		expected string
	}{
		{ErrNotFound, "resource not found"},
		{ErrInvalidInput, "invalid input"},
		{ErrConflict, "resource conflict"},
		{ErrInternal, "internal system error"},
		{ErrUnauthorized, "unauthorized"},
	}

	for _, tt := range tests {
		if tt.err.Error() != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, tt.err.Error())
		}
	}
}

func TestIsComplicated(t *testing.T) {
	err := Wrap(ErrNotFound, "failed to fetch user")
	if !Is(err, ErrNotFound) {
		t.Error("expected wrapped error to be ErrNotFound")
	}

	deep := Wrap(err, "outer context")
	if !Is(deep, ErrNotFound) {
		t.Error("expected deeply wrapped error to be ErrNotFound")
	}

	if Is(deep, ErrInternal) {
		t.Error("did not expect ErrInternal")
	}
}

func TestWrapComplicated(t *testing.T) {
	if Wrap(nil, "something") != nil {
		t.Error("expected nil")
	}

	err := errors.New("original")
	wrapped := Wrap(err, "context")
	if !errors.Is(wrapped, err) {
		t.Error("expected wrapped error to wrap original")
	}
}
