package errors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	caos_errs "github.com/zitadel/zitadel/v2/internal/errors"
)

func TestUnauthenticatedError(t *testing.T) {
	var err interface{}
	err = new(caos_errs.UnauthenticatedError)
	_, ok := err.(caos_errs.Unauthenticated)
	assert.True(t, ok)
}

func TestThrowUnauthenticatedf(t *testing.T) {
	err := caos_errs.ThrowUnauthenticatedf(nil, "id", "msg")
	_, ok := err.(*caos_errs.UnauthenticatedError)
	assert.True(t, ok)
}

func TestIsUnauthenticated(t *testing.T) {
	err := caos_errs.ThrowUnauthenticated(nil, "id", "msg")
	ok := caos_errs.IsUnauthenticated(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = caos_errs.IsUnauthenticated(err)
	assert.False(t, ok)
}
