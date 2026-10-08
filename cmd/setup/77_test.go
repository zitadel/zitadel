package setup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStampEventPositionAtInsertSQL(t *testing.T) {
	assert.NotContains(t, stampEventPositionAtInsert, "last_owner")
	assert.Contains(t, stampEventPositionAtInsert, "i INTEGER")
	assert.Contains(t, stampEventPositionAtInsert, "clock_timestamp()")
	assert.Contains(t, stampEventPositionAtInsert, "VOLATILE PARALLEL SAFE")
}

func TestStampEventPositionAtInsert_ExecuteDoesNotRun(t *testing.T) {
	err := new(StampEventPositionAtInsert).Execute(context.Background(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "historical")
}
