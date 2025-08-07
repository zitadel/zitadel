package zerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestInternalError(t *testing.T) {
	var err interface{} = new(zerrors.InternalError)
	_, ok := err.(zerrors.Internal)
	assert.True(t, ok)
}

func TestThrowInternalf(t *testing.T) {
	err := zerrors.ThrowInternalf(nil, "id", "msg")
	//nolint:errorlint
	_, ok := err.(*zerrors.InternalError)
	assert.True(t, ok)
}

func TestIsInternal(t *testing.T) {
	err := zerrors.ThrowInternal(nil, "id", "msg")
	ok := zerrors.IsInternal(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = zerrors.IsInternal(err)
	assert.False(t, ok)
}
