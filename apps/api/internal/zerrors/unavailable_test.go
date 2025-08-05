package zerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestUnavailableError(t *testing.T) {
	var err interface{} = new(zerrors.UnavailableError)
	_, ok := err.(zerrors.Unavailable)
	assert.True(t, ok)
}

func TestThrowUnavailablef(t *testing.T) {
	err := zerrors.ThrowUnavailablef(nil, "id", "msg")
	//nolint:errorlint
	_, ok := err.(*zerrors.UnavailableError)
	assert.True(t, ok)
}

func TestIsUnavailable(t *testing.T) {
	err := zerrors.ThrowUnavailable(nil, "id", "msg")
	ok := zerrors.IsUnavailable(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = zerrors.IsUnavailable(err)
	assert.False(t, ok)
}
