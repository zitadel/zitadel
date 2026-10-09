package setup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStampEventPositionAtInsert_ExecuteDoesNotRun(t *testing.T) {
	err := new(StampEventPositionAtInsert).Execute(context.Background(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "historical")
}
