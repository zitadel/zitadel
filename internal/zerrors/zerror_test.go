package zerrors

import (
	"errors"
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
	target := CreateZitadelError(KindAborted, parent, "id", "message")
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
			err:  CreateZitadelError(KindNotFound, parent, "id", "message"),
			want: false,
		},
		{
			name: "different id",
			err:  CreateZitadelError(KindAborted, parent, "otherID", "message"),
			want: false,
		},
		{
			name: "different message",
			err:  CreateZitadelError(KindAborted, parent, "id", "other message"),
			want: false,
		},
		{
			name: "different parent",
			err:  CreateZitadelError(KindAborted, errors.New("other parent"), "id", "message"),
			want: false,
		},
		{
			name: "same error",
			err:  CreateZitadelError(KindAborted, parent, "id", "message"),
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

func Test_newCaller(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		skip int
		want caller
	}{
		{
			name: "skip 0",
			skip: 0,
			want: caller{
				Func: "github.com/zitadel/zitadel/internal/zerrors.Test_newCaller.func1",
				File: "internal/zerrors/zerror_test.go",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newCaller(tt.skip)
			assert.Equal(t, tt.want.Func, got.Func)
			assert.Contains(t, got.File, tt.want.File)
			assert.NotZero(t, got.Line)
		})
	}
}
