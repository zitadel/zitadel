package errors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

func TestPreconditionFailedError(t *testing.T) {
	var err interface{}
	err = new(caos_errs.PreconditionFailedError)
	_, ok := err.(caos_errs.PreconditionFailed)
	assert.True(t, ok)
}

func TestThrowPreconditionFailedf(t *testing.T) {
	err := caos_errs.ThrowPreconditionFailedf(nil, "id", "msg")
	_, ok := err.(*caos_errs.PreconditionFailedError)
	assert.True(t, ok)
}

func TestIsPreconditionFailed(t *testing.T) {
	err := caos_errs.ThrowPreconditionFailed(nil, "id", "msg")
	ok := caos_errs.IsPreconditionFailed(err)
	assert.True(t, ok)

	err = errors.New("Precondition failed!")
	ok = caos_errs.IsPreconditionFailed(err)
	assert.False(t, ok)
}
