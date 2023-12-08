package zerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestPermissionDeniedError(t *testing.T) {
	var err interface{} = new(zerrors.PermissionDeniedError)
	_, ok := err.(zerrors.PermissionDenied)
	assert.True(t, ok)
}

func TestThrowPermissionDeniedf(t *testing.T) {
	err := zerrors.ThrowPermissionDeniedf(nil, "id", "msg")
	//nolint:errorlint
	_, ok := err.(*zerrors.PermissionDeniedError)
	assert.True(t, ok)
}

func TestIsPermissionDenied(t *testing.T) {
	err := zerrors.ThrowPermissionDenied(nil, "id", "msg")
	ok := zerrors.IsPermissionDenied(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = zerrors.IsPermissionDenied(err)
	assert.False(t, ok)
}
