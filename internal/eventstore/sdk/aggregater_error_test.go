package sdk

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregaterError(t *testing.T) {
	var err interface{}
	err = new(aggregaterError)
	_, ok := err.(AggregaterError)
	assert.True(t, ok)
}

func TestThrowAggregaterf(t *testing.T) {
	err := ThrowAggregaterf(nil, "id", "msg")
	_, ok := err.(*aggregaterError)
	assert.True(t, ok)
}

func TestIsAggregater(t *testing.T) {
	err := ThrowAggregater(nil, "id", "msg")
	ok := IsAggregater(err)
	assert.True(t, ok)

	err = errors.New("i am found")
	ok = IsAggregater(err)
	assert.False(t, ok)
}
