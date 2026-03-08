package apperr

import (
	"errors"
	"testing"
)

func TestWrapPreservesWrappedError(t *testing.T) {
	t.Parallel()

	cause := errors.New("boom")
	err := Wrap(cause, Internal, "internal.failure", "request failed")

	if !errors.Is(err, cause) {
		t.Fatalf("expected wrapped error to match cause")
	}

	if _, ok := As(err); !ok {
		t.Fatal("expected wrapped app error")
	}
}

func TestStatusMapping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		kind Kind
		want int
	}{
		{name: "invalid_argument", kind: InvalidArgument, want: 400},
		{name: "not_found", kind: NotFound, want: 404},
		{name: "conflict", kind: Conflict, want: 409},
		{name: "unauthorized", kind: Unauthorized, want: 401},
		{name: "forbidden", kind: Forbidden, want: 403},
		{name: "unavailable", kind: Unavailable, want: 503},
		{name: "internal", kind: Internal, want: 500},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := New(tt.kind, Code("test.code"), "failed")
			if got := Status(err); got != tt.want {
				t.Fatalf("expected status %d, got %d", tt.want, got)
			}
		})
	}
}

func TestHelpersMatchKindAndCode(t *testing.T) {
	t.Parallel()

	err := New(NotFound, "users.not_found", "user not found")

	if !HasKind(err, NotFound) {
		t.Fatal("expected kind match")
	}

	if !HasCode(err, "users.not_found") {
		t.Fatal("expected code match")
	}
}
