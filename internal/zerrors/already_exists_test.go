package zerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestAlreadyExistsError(t *testing.T) {
	var alreadyExistsError interface{} = new(zerrors.AlreadyExistsError)
	_, ok := alreadyExistsError.(zerrors.AlreadyExists)
	assert.True(t, ok)
}

func TestThrowAlreadyExistsf(t *testing.T) {
	err := zerrors.ThrowAlreadyExistsf(nil, "id", "msg")
	//nolint:errorlint
	_, ok := err.(*zerrors.AlreadyExistsError)
	assert.True(t, ok)
}

func TestIsErrorAlreadyExists(t *testing.T) {
	err := zerrors.ThrowAlreadyExists(nil, "id", "msg")
	ok := zerrors.IsErrorAlreadyExists(err)
	assert.True(t, ok)

	err = errors.New("Already Exists!")
	ok = zerrors.IsErrorAlreadyExists(err)
	assert.False(t, ok)
}
