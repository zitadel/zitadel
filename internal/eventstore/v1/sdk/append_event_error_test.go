package sdk

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendEventError(t *testing.T) {
	var err interface{}
	err = new(appendEventError)
	_, ok := err.(*appendEventError)
	assert.True(t, ok)
}

func TestThrowAppendEventErrorf(t *testing.T) {
	err := ThrowAggregaterf(nil, "id", "msg")
	_, ok := err.(*appendEventError)
	assert.True(t, ok)
}

func TestIsAppendEventError(t *testing.T) {
	err := ThrowAppendEventError(nil, "id", "msg")
	ok := IsAppendEventError(err)
	assert.True(t, ok)

	err = errors.New("i am found")
	ok = IsAppendEventError(err)
	assert.False(t, ok)
}
