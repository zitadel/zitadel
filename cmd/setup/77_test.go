package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStampEventPositionAtInsertSQL(t *testing.T) {
	assert.NotContains(t, stampEventPositionAtInsert, "last_owner")
	assert.Contains(t, stampEventPositionAtInsert, "i INTEGER")
	assert.Contains(t, stampEventPositionAtInsert, "clock_timestamp()")
	assert.Contains(t, stampEventPositionAtInsert, "VOLATILE PARALLEL SAFE")
}
