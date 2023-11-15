package errors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	caos_errs "github.com/zitadel/zitadel/v2/internal/errors"
)

func TestInternalError(t *testing.T) {
	var err interface{}
	err = new(caos_errs.InternalError)
	_, ok := err.(caos_errs.Internal)
	assert.True(t, ok)
}

func TestThrowInternalf(t *testing.T) {
	err := caos_errs.ThrowInternalf(nil, "id", "msg")
	_, ok := err.(*caos_errs.InternalError)
	assert.True(t, ok)
}

func TestIsInternal(t *testing.T) {
	err := caos_errs.ThrowInternal(nil, "id", "msg")
	ok := caos_errs.IsInternal(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = caos_errs.IsInternal(err)
	assert.False(t, ok)
}
