package zerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestUnauthenticatedError(t *testing.T) {
	var err interface{} = new(zerrors.UnauthenticatedError)
	_, ok := err.(zerrors.Unauthenticated)
	assert.True(t, ok)
}

func TestThrowUnauthenticatedf(t *testing.T) {
	err := zerrors.ThrowUnauthenticatedf(nil, "id", "msg")
	//nolint:errorlint
	_, ok := err.(*zerrors.UnauthenticatedError)
	assert.True(t, ok)
}

func TestIsUnauthenticated(t *testing.T) {
	err := zerrors.ThrowUnauthenticated(nil, "id", "msg")
	ok := zerrors.IsUnauthenticated(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = zerrors.IsUnauthenticated(err)
	assert.False(t, ok)
}
