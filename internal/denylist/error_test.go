package denylist

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// nilUnsafeError is an error whose Unwrap method must not be called on a nil receiver,
// like most error types.
type nilUnsafeError struct {
	parent error
}

func (e *nilUnsafeError) Error() string { return "nil unsafe" }

func (e *nilUnsafeError) Unwrap() error { return e.parent }

// zeroValueTargetError is an error whose zero value returns its nil parent
// as a non-nil error interface from Unwrap.
// Callers use such zero values as target of [errors.Is] to check for the type.
type zeroValueTargetError struct {
	parent *nilUnsafeError
}

func (e *zeroValueTargetError) Error() string { return "zero value target" }

func (e *zeroValueTargetError) Unwrap() error { return e.parent }

func TestAddressDeniedError_Is(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		target error
		want   bool
	}{
		{
			name:   "same deniedBy",
			target: NewAddressDeniedError("localhost"),
			want:   true,
		},
		{
			name:   "different deniedBy",
			target: NewAddressDeniedError("127.0.0.0/8"),
			want:   false,
		},
		{
			name:   "other error type",
			target: errors.New("address is denied by 'localhost'"),
			want:   false,
		},
		{
			name:   "target is not unwrapped",
			target: new(zeroValueTargetError),
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := fmt.Errorf("dial: %w", NewAddressDeniedError("localhost"))
			assert.Equal(t, tt.want, errors.Is(err, tt.target))
		})
	}
}
