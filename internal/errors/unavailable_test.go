package errors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

func TestUnavailableError(t *testing.T) {
	var err interface{}
	err = new(caos_errs.UnavailableError)
	_, ok := err.(caos_errs.Unavailable)
	assert.True(t, ok)
}

func TestThrowUnavailablef(t *testing.T) {
	err := caos_errs.ThrowUnavailablef(nil, "id", "msg")
	_, ok := err.(*caos_errs.UnavailableError)
	assert.True(t, ok)
}

func TestIsUnavailable(t *testing.T) {
	err := caos_errs.ThrowUnavailable(nil, "id", "msg")
	ok := caos_errs.IsUnavailable(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = caos_errs.IsUnavailable(err)
	assert.False(t, ok)
}
