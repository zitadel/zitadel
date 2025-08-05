package zerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestErrorMethod(t *testing.T) {
	err := zerrors.ThrowError(nil, "id", "msg")
	expected := "ID=id Message=msg"
	assert.Equal(t, expected, err.Error())

	err = zerrors.ThrowError(err, "subID", "subMsg")
	subExptected := "ID=subID Message=subMsg Parent=(ID=id Message=msg)"
	assert.Equal(t, subExptected, err.Error())
}

func TestIsZitadelError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "zitadel error",
			err:  zerrors.ThrowInvalidArgument(nil, "id", "msg"),
			want: true,
		},
		{
			name: "other error",
			err:  errors.New("just a random error"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, zerrors.IsZitadelError(tt.err), "IsZitadelError(%v)", tt.err)
		})
	}
}
