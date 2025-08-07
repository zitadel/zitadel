package zerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestContains(t *testing.T) {
	err := errors.New("hello world")
	world := zerrors.Contains(err, "hello")
	assert.True(t, world)

	mars := zerrors.Contains(err, "mars")
	assert.False(t, mars)
}
