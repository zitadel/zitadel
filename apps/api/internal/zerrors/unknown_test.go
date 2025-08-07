package zerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestUnknownError(t *testing.T) {
	var err interface{} = new(zerrors.UnknownError)
	_, ok := err.(zerrors.Unknown)
	assert.True(t, ok)
}

func TestThrowUnknownf(t *testing.T) {
	err := zerrors.ThrowUnknownf(nil, "id", "msg")
	//nolint:errorlint
	_, ok := err.(*zerrors.UnknownError)
	assert.True(t, ok)
}

func TestIsUnknown(t *testing.T) {
	err := zerrors.ThrowUnknown(nil, "id", "msg")
	ok := zerrors.IsUnknown(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = zerrors.IsUnknown(err)
	assert.False(t, ok)
}
