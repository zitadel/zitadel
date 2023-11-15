package errors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	caos_errs "github.com/zitadel/zitadel/v2/internal/errors"
)

func TestNotFoundError(t *testing.T) {
	var notFoundError interface{}
	notFoundError = new(caos_errs.NotFoundError)
	_, ok := notFoundError.(caos_errs.NotFound)
	assert.True(t, ok)
}

func TestThrowNotFoundf(t *testing.T) {
	err := caos_errs.ThrowNotFoundf(nil, "id", "msg")
	_, ok := err.(*caos_errs.NotFoundError)
	assert.True(t, ok)
}

func TestIsNotFound(t *testing.T) {
	err := caos_errs.ThrowNotFound(nil, "id", "msg")
	ok := caos_errs.IsNotFound(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = caos_errs.IsNotFound(err)
	assert.False(t, ok)
}
