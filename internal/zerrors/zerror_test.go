package zerrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func panicWith(v any) {
	panic(v)
}

// recoverFrom panics with v and returns the error ThrowRecovered creates for it.
func recoverFrom(v any) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = ThrowRecovered(rec)
		}
	}()
	panicWith(v)
	return nil
}

func TestThrowRecovered(t *testing.T) {
	prev := enableStackTrace.Load()
	t.Cleanup(func() { EnableStackTrace(prev) })

	parent := errors.New("parent error")
	tests := []struct {
		name        string
		stackTrace  bool
		panicValue  any
		wantIs      error
		wantContain string
	}{
		{
			name:        "error value",
			stackTrace:  true,
			panicValue:  parent,
			wantIs:      parent,
			wantContain: "parent error",
		},
		{
			name:        "stack trace disabled",
			stackTrace:  false,
			panicValue:  "boom",
			wantContain: "boom",
		},
		{
			name:        "zitadel error value",
			stackTrace:  false,
			panicValue:  CreateZitadelError(KindNotFound, nil, "id", "message", 0),
			wantContain: "ID=id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			EnableStackTrace(tt.stackTrace)
			err := recoverFrom(tt.panicValue)

			assert.True(t, IsInternal(err))
			zErr, ok := errors.AsType[*ZitadelError](err)
			require.True(t, ok)
			assert.Equal(t, IDRecover, zErr.ID)
			if tt.wantIs != nil {
				assert.ErrorIs(t, err, tt.wantIs)
			}
			assert.Contains(t, err.Error(), tt.wantContain)

			trace, ok := zErr.StackTrace()
			assert.True(t, ok)
			assert.Contains(t, string(trace), "panicWith")
		})
	}
}
