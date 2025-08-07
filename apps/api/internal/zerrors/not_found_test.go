package zerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestNotFoundError(t *testing.T) {
	var notFoundError interface{} = new(zerrors.NotFoundError)
	_, ok := notFoundError.(zerrors.NotFound)
	assert.True(t, ok)
}

func TestThrowNotFoundf(t *testing.T) {
	err := zerrors.ThrowNotFoundf(nil, "id", "msg")
	//nolint:errorlint
	_, ok := err.(*zerrors.NotFoundError)
	assert.True(t, ok)
}

func TestIsNotFound(t *testing.T) {
	err := zerrors.ThrowNotFound(nil, "id", "msg")
	ok := zerrors.IsNotFound(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = zerrors.IsNotFound(err)
	assert.False(t, ok)
}
