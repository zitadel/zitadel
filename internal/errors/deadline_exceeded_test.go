package errors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	caos_errs "github.com/zitadel/zitadel/v2/internal/errors"
)

func TestDeadlineExceededError(t *testing.T) {
	var err interface{}
	err = new(caos_errs.DeadlineExceededError)
	_, ok := err.(caos_errs.DeadlineExceeded)
	assert.True(t, ok)
}

func TestThrowDeadlineExceededf(t *testing.T) {
	err := caos_errs.ThrowDeadlineExceededf(nil, "id", "msg")
	_, ok := err.(*caos_errs.DeadlineExceededError)
	assert.True(t, ok)
}

func TestIsDeadlineExceeded(t *testing.T) {
	err := caos_errs.ThrowDeadlineExceeded(nil, "id", "msg")
	ok := caos_errs.IsDeadlineExceeded(err)
	assert.True(t, ok)

	err = errors.New("I am found!")
	ok = caos_errs.IsDeadlineExceeded(err)
	assert.False(t, ok)
}
