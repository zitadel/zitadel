package zerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestDeadlineExceededError(t *testing.T) {
	var err interface{} = new(zerrors.DeadlineExceededError)
	_, ok := err.(zerrors.DeadlineExceeded)
	assert.True(t, ok)
}

func TestThrowDeadlineExceededf(t *testing.T) {
	err := zerrors.ThrowDeadlineExceededf(nil, "id", "msg")
	//nolint:errorlint
	_, ok := err.(*zerrors.DeadlineExceededError)
	assert.True(t, ok)
}

func TestIsDeadlineExceeded(t *testing.T) {
	err := zerrors.ThrowDeadlineExceeded(nil, "id", "msg")
	ok := zerrors.IsDeadlineExceeded(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = zerrors.IsDeadlineExceeded(err)
	assert.False(t, ok)
}
