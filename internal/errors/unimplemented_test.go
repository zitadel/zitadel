package errors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

func TestUnimplementedError(t *testing.T) {
	var unimplementedError interface{}
	unimplementedError = new(caos_errs.UnimplementedError)
	_, ok := unimplementedError.(caos_errs.Unimplemented)
	assert.True(t, ok)
}

func TestThrowUnimplementedf(t *testing.T) {
	err := caos_errs.ThrowUnimplementedf(nil, "id", "msg")
	_, ok := err.(*caos_errs.UnimplementedError)
	assert.True(t, ok)
}

func TestIsUnimplemented(t *testing.T) {
	err := caos_errs.ThrowUnimplemented(nil, "id", "msg")
	ok := caos_errs.IsUnimplemented(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = caos_errs.IsUnimplemented(err)
	assert.False(t, ok)
}
