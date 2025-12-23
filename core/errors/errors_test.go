package errors

import (
	"errors"
	"testing"
)

func TestIs(t *testing.T) {
	base := errors.New("base")
	wrapped := Wrap(base, "wrapped")

	if !Is(wrapped, base) {
		t.Error("expected wrapped error to match base")
	}
	if Is(wrapped, errors.New("other")) {
		t.Error("expected wrapped error not to match other")
	}
}

func TestWrap(t *testing.T) {
	t.Run("Wrap non-nil", func(t *testing.T) {
		base := errors.New("base")
		wrapped := Wrap(base, "context")
		if wrapped == nil {
			t.Fatal("expected non-nil error")
		}
		if wrapped.Error() != "context: base" {
			t.Errorf("expected 'context: base', got %q", wrapped.Error())
		}
	})

	t.Run("Wrap nil", func(t *testing.T) {
		if Wrap(nil, "msg") != nil {
			t.Error("expected nil when wrapping nil error")
		}
	})
}
