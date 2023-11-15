package errors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	caos_errs "github.com/zitadel/zitadel/v2/internal/errors"
)

func TestInvalidArgumentError(t *testing.T) {
	var invalidArgumentError interface{}
	invalidArgumentError = new(caos_errs.InvalidArgumentError)
	_, ok := invalidArgumentError.(caos_errs.InvalidArgument)
	assert.True(t, ok)
}

func TestThrowInvalidArgumentf(t *testing.T) {
	err := caos_errs.ThrowInvalidArgumentf(nil, "id", "msg")
	_, ok := err.(*caos_errs.InvalidArgumentError)
	assert.True(t, ok)
}

func TestIsErrorInvalidArgument(t *testing.T) {
	err := caos_errs.ThrowInvalidArgument(nil, "id", "msg")
	ok := caos_errs.IsErrorInvalidArgument(err)
	assert.True(t, ok)

	err = errors.New("I am invalid!")
	ok = caos_errs.IsErrorInvalidArgument(err)
	assert.False(t, ok)
}
