package errors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

func TestUnknownError(t *testing.T) {
	var err interface{}
	err = new(caos_errs.UnknownError)
	_, ok := err.(caos_errs.Unknown)
	assert.True(t, ok)
}

func TestThrowUnknownf(t *testing.T) {
	err := caos_errs.ThrowUnknownf(nil, "id", "msg")
	_, ok := err.(*caos_errs.UnknownError)
	assert.True(t, ok)
}

func TestIsUnknown(t *testing.T) {
	err := caos_errs.ThrowUnknown(nil, "id", "msg")
	ok := caos_errs.IsUnknown(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = caos_errs.IsUnknown(err)
	assert.False(t, ok)
}
