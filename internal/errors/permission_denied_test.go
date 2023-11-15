package errors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	caos_errs "github.com/zitadel/zitadel/v2/internal/errors"
)

func TestPermissionDeniedError(t *testing.T) {
	var err interface{}
	err = new(caos_errs.PermissionDeniedError)
	_, ok := err.(caos_errs.PermissionDenied)
	assert.True(t, ok)
}

func TestThrowPermissionDeniedf(t *testing.T) {
	err := caos_errs.ThrowPermissionDeniedf(nil, "id", "msg")
	_, ok := err.(*caos_errs.PermissionDeniedError)
	assert.True(t, ok)
}

func TestIsPermissionDenied(t *testing.T) {
	err := caos_errs.ThrowPermissionDenied(nil, "id", "msg")
	ok := caos_errs.IsPermissionDenied(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = caos_errs.IsPermissionDenied(err)
	assert.False(t, ok)
}
