package handler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOverwrite_usesPassedFilterOffsetIsCursor(t *testing.T) {
	cfg := Config{
		BulkLimit:             100,
		FailureCountUntilSkip: 5,
		TransactionDuration:   time.Second,
		Handlers: map[string]*ConfigOverwrites{
			"User": {MinimumCycleDuration: 2 * time.Minute},
		},
	}

	got := cfg.overwrite("User", true)
	assert.True(t, got.FilterOffsetIsCursor)
	assert.Equal(t, 2*time.Minute, got.RequeueEvery)

	got = cfg.overwrite("missing", false)
	assert.False(t, got.FilterOffsetIsCursor)
	assert.Equal(t, 3*time.Minute, got.RequeueEvery)
}
