package zerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestInvalidArgumentError(t *testing.T) {
	var invalidArgumentError interface{} = new(zerrors.InvalidArgumentError)
	_, ok := invalidArgumentError.(zerrors.InvalidArgument)
	assert.True(t, ok)
}

func TestThrowInvalidArgumentf(t *testing.T) {
	err := zerrors.ThrowInvalidArgumentf(nil, "id", "msg")
	//nolint:errorlint
	_, ok := err.(*zerrors.InvalidArgumentError)
	assert.True(t, ok)
}

func TestIsErrorInvalidArgument(t *testing.T) {
	err := zerrors.ThrowInvalidArgument(nil, "id", "msg")
	ok := zerrors.IsErrorInvalidArgument(err)
	assert.True(t, ok)

	err = errors.New("I am invalid!")
	ok = zerrors.IsErrorInvalidArgument(err)
	assert.False(t, ok)
}
