package errors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	caos_errs "github.com/zitadel/zitadel/v2/internal/errors"
)

func TestAlreadyExistsError(t *testing.T) {
	var alreadyExistsError interface{}
	alreadyExistsError = new(caos_errs.AlreadyExistsError)
	_, ok := alreadyExistsError.(caos_errs.AlreadyExists)
	assert.True(t, ok)
}

func TestThrowAlreadyExistsf(t *testing.T) {
	err := caos_errs.ThrowAlreadyExistsf(nil, "id", "msg")

	_, ok := err.(*caos_errs.AlreadyExistsError)
	assert.True(t, ok)
}

func TestIsErrorAlreadyExists(t *testing.T) {
	err := caos_errs.ThrowAlreadyExists(nil, "id", "msg")
	ok := caos_errs.IsErrorAlreadyExists(err)
	assert.True(t, ok)

	err = errors.New("Already Exists!")
	ok = caos_errs.IsErrorAlreadyExists(err)
	assert.False(t, ok)
}
