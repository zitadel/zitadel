package zerrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorMethod(t *testing.T) {
	err := ThrowError(nil, "id", "msg")
	expected := "ID=id Message=msg"
	assert.Equal(t, expected, err.Error())

	err = ThrowError(err, "subID", "subMsg")
	subExptected := "ID=subID Message=subMsg Parent=(ID=id Message=msg)"
	assert.Equal(t, subExptected, err.Error())
}

func TestZitadelError_Is(t *testing.T) {
	parent := errors.New("parent error")
	target := CreateZitadelError(KindAborted, parent, "id", "message", 0)
	tests := []struct {
		name string // description of this test case
		err  error
		want bool
	}{
		{
			name: "wrong type",
			err:  errors.New("some other error"),
			want: false,
		},
		{
			name: "different kind",
			err:  CreateZitadelError(KindNotFound, parent, "id", "message", 0),
			want: false,
		},
		{
			name: "different id",
			err:  CreateZitadelError(KindAborted, parent, "otherID", "message", 0),
			want: false,
		},
		{
			name: "different message",
			err:  CreateZitadelError(KindAborted, parent, "id", "other message", 0),
			want: false,
		},
		{
			name: "different parent",
			err:  CreateZitadelError(KindAborted, errors.New("other parent"), "id", "message", 0),
			want: false,
		},
		{
			name: "same error",
			err:  CreateZitadelError(KindAborted, parent, "id", "message", 0),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errors.Is(tt.err, target)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestZitadelError_Is_WrappedTarget(t *testing.T) {
	parent := errors.New("parent error")
	err := CreateZitadelError(KindAborted, parent, "id", "message", 0)
	tests := []struct {
		name   string // description of this test case
		target error
		want   bool
	}{
		{
			name:   "unwrapped target",
			target: CreateZitadelError(KindAborted, parent, "id", "message", 0),
			want:   true,
		},
		{
			name:   "wrapped target",
			target: fmt.Errorf("wrapped: %w", CreateZitadelError(KindAborted, parent, "id", "message", 0)),
			want:   true,
		},
		{
			name:   "wrapped target of different kind",
			target: fmt.Errorf("wrapped: %w", CreateZitadelError(KindNotFound, parent, "id", "message", 0)),
			want:   false,
		},
		{
			name:   "wrapped target without zitadel error",
			target: fmt.Errorf("wrapped: %w", errors.New("some other error")),
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errors.Is(err, tt.target)
			assert.Equal(t, tt.want, got)
		})
	}
}
