package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/migration"
)

func TestCommandsToEventsClockTimestamp_String(t *testing.T) {
	assert.Equal(t, "81_eventstore_commands_to_events_clock_timestamp", new(CommandsToEventsClockTimestamp).String())
	assert.NotEqual(t, migration.EventstorePositionClockTimestampStep, new(CommandsToEventsClockTimestamp).String())
}

func TestCommandsToEventsClockTimestampSQL(t *testing.T) {
	assert.NotContains(t, commandsToEventsClockTimestampSQL, "last_owner")
	assert.Contains(t, commandsToEventsClockTimestampSQL, "i INTEGER")
	assert.Contains(t, commandsToEventsClockTimestampSQL, "clock_timestamp()")
	assert.Contains(t, commandsToEventsClockTimestampSQL, "VOLATILE PARALLEL SAFE")
	assert.Contains(t, commandsToEventsClockTimestampSQL, "CREATE OR REPLACE FUNCTION eventstore.commands_to_events")
}

func TestOneShotMigrationsOmits77(t *testing.T) {
	steps := &Steps{
		s81CommandsToEventsClockTimestamp: &CommandsToEventsClockTimestamp{},
		s82CurrentStatesInTxOrder:         &CurrentStatesInTxOrder{},
		s29FillFieldsForProjectGrant:      &FillFieldsForProjectGrant{},
	}
	names := make([]string, 0, len(oneShotMigrations(steps)))
	for _, step := range oneShotMigrations(steps) {
		if step == nil {
			continue
		}
		names = append(names, step.String())
	}

	assert.NotContains(t, names, migration.EventstorePositionClockTimestampStep)
	assert.Contains(t, names, "81_eventstore_commands_to_events_clock_timestamp")
	assert.Contains(t, names, "82_current_states_in_tx_order")
	assert.Contains(t, names, "29_init_fields_for_project_grant")

	idx82, idx29 := -1, -1
	for i, name := range names {
		switch name {
		case "82_current_states_in_tx_order":
			idx82 = i
		case "29_init_fields_for_project_grant":
			idx29 = i
		}
	}
	require.GreaterOrEqual(t, idx82, 0)
	require.GreaterOrEqual(t, idx29, 0)
	assert.Less(t, idx82, idx29, "82 must add in_tx_order before field handlers query current_states")
}
