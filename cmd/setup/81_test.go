package setup

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/migration"
)

func TestCommandsToEventsClockTimestamp_String(t *testing.T) {
	assert.Equal(t, migration.EventstorePositionClockTimestampStep, new(StampEventPositionAtInsert).String())
	assert.Equal(t, "81_eventstore_commands_to_events_clock_timestamp", new(CommandsToEventsClockTimestamp).String())
	assert.NotEqual(t, new(StampEventPositionAtInsert).String(), new(CommandsToEventsClockTimestamp).String())
}

func TestCommandsToEventsClockTimestampSQL(t *testing.T) {
	assert.Contains(t, stampEventPositionAtInsert, "clock_timestamp()")
	assert.Contains(t, stampEventPositionAtInsert, "CREATE OR REPLACE FUNCTION eventstore.commands_to_events")
}

func TestSetupExecutionSliceOmits77(t *testing.T) {
	src, err := os.ReadFile("setup.go")
	require.NoError(t, err)

	re := regexp.MustCompile(`for _, step := range \[\]migration\.Migration\{([^}]+)\}`)
	matches := re.FindAllSubmatch(src, -1)
	require.NotEmpty(t, matches)

	var found bool
	for _, m := range matches {
		slice := string(m[1])
		if !regexp.MustCompile(`s81CommandsToEventsClockTimestamp`).MatchString(slice) {
			continue
		}
		found = true
		assert.NotContains(t, slice, "s77StampEventPositionAtInsert")
		assert.Contains(t, slice, "s81CommandsToEventsClockTimestamp")
		assert.Contains(t, slice, "s82CurrentStatesInTxOrder")
	}
	assert.True(t, found, "expected the one-shot migration slice to include 81 and 82")
}
