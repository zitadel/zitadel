package zerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestUnimplementedError(t *testing.T) {
	var unimplementedError interface{} = new(zerrors.UnimplementedError)
	_, ok := unimplementedError.(zerrors.Unimplemented)
	assert.True(t, ok)
}

func TestThrowUnimplementedf(t *testing.T) {
	err := zerrors.ThrowUnimplementedf(nil, "id", "msg")
	//nolint:errorlint
	_, ok := err.(*zerrors.UnimplementedError)
	assert.True(t, ok)
}

func TestIsUnimplemented(t *testing.T) {
	err := zerrors.ThrowUnimplemented(nil, "id", "msg")
	ok := zerrors.IsUnimplemented(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = zerrors.IsUnimplemented(err)
	assert.False(t, ok)
}
