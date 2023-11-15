package errors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	caos_errors "github.com/zitadel/zitadel/v2/internal/errors"
)

func TestContains(t *testing.T) {
	err := errors.New("hello world")
	world := caos_errors.Contains(err, "hello")
	assert.True(t, world)

	mars := caos_errors.Contains(err, "mars")
	assert.False(t, mars)
}
