package zerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestPreconditionFailedError(t *testing.T) {
	var err interface{} = new(zerrors.PreconditionFailedError)
	_, ok := err.(zerrors.PreconditionFailed)
	assert.True(t, ok)
}

func TestThrowPreconditionFailedf(t *testing.T) {
	err := zerrors.ThrowPreconditionFailedf(nil, "id", "msg")
	//nolint:errorlint
	_, ok := err.(*zerrors.PreconditionFailedError)
	assert.True(t, ok)
}

func TestIsPreconditionFailed(t *testing.T) {
	err := zerrors.ThrowPreconditionFailed(nil, "id", "msg")
	ok := zerrors.IsPreconditionFailed(err)
	assert.True(t, ok)

	err = errors.New("Precondition failed!")
	ok = zerrors.IsPreconditionFailed(err)
	assert.False(t, ok)
}
