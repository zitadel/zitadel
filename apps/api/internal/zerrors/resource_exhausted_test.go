package zerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestResourceExhaustedError(t *testing.T) {
	var err interface{} = new(zerrors.ResourceExhaustedError)
	_, ok := err.(zerrors.ResourceExhausted)
	assert.True(t, ok)
}

func TestThrowResourceExhaustedf(t *testing.T) {
	err := zerrors.ThrowResourceExhaustedf(nil, "id", "msg")
	//nolint:errorlint
	_, ok := err.(*zerrors.ResourceExhaustedError)
	assert.True(t, ok)
}

func TestIsResourceExhausted(t *testing.T) {
	err := zerrors.ThrowResourceExhausted(nil, "id", "msg")
	ok := zerrors.IsResourceExhausted(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = zerrors.IsResourceExhausted(err)
	assert.False(t, ok)
}
