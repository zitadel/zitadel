package errors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

func TestResourceExhaustedError(t *testing.T) {
	var err interface{} = new(caos_errs.ResourceExhaustedError)
	_, ok := err.(caos_errs.ResourceExhausted)
	assert.True(t, ok)
}

func TestThrowResourceExhaustedf(t *testing.T) {
	err := caos_errs.ThrowResourceExhaustedf(nil, "id", "msg")
	// TODO: refactor errors package
	//nolint:errorlint
	_, ok := err.(*caos_errs.ResourceExhaustedError)
	assert.True(t, ok)
}

func TestIsResourceExhausted(t *testing.T) {
	err := caos_errs.ThrowResourceExhausted(nil, "id", "msg")
	ok := caos_errs.IsResourceExhausted(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = caos_errs.IsResourceExhausted(err)
	assert.False(t, ok)
}
