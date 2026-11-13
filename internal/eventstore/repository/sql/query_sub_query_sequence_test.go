package sql

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

// Test_query_sub_query_sequence_greater_postgres reads the password state of a user together with
// only the password checks after the last reset, as done by the command side password check.
// The events have historical positions which disagree with their sequence order
// (positions were the transaction start time before setup step 77),
// so both the cutoff and the order must be based on the aggregate sequence.
func Test_query_sub_query_sequence_greater_postgres(t *testing.T) {
	ctx := t.Context()
	instanceID := "sequence-greater-" + t.Name()
	client := &Postgres{DB: &database.DB{DB: testClient, Database: new(testDB)}}

	t.Cleanup(func() {
		_, err := testClient.ExecContext(context.Background(), "DELETE FROM eventstore.events2 WHERE instance_id = $1", instanceID)
		assert.NoError(t, err)
	})

	type event struct {
		eventType eventstore.EventType
		position  int64
	}
	// index+1 is the sequence
	events := []event{
		{"user.human.added", 100},
		{"user.human.password.changed", 108},
		// newer password change with an older position: replaying by position would restore the superseded password
		{"user.human.password.changed", 102},
		// last reset of the failed count
		{"user.human.password.check.succeeded", 110},
		// failure after the reset with an older position: a position based cutoff would miss it
		{"user.human.password.check.failed", 105},
		{"user.human.password.check.failed", 120},
	}
	for i, e := range events {
		_, err := testClient.ExecContext(ctx,
			`INSERT INTO eventstore.events2 (instance_id, aggregate_type, aggregate_id, event_type, "sequence", revision, created_at, payload, creator, "owner", "position", in_tx_order)
			VALUES ($1, 'user', 'user1', $2, $3, 1, now(), NULL, 'creator', 'owner', $4, 1)`,
			instanceID, e.eventType, i+1, decimal.NewFromInt(e.position),
		)
		require.NoError(t, err)
	}

	filter := func(t *testing.T, builder *eventstore.SearchQueryBuilder) []eventstore.Event {
		t.Helper()
		var events []eventstore.Event
		err := query(ctx, client, builder.InstanceID(instanceID), eventstore.Reducer(func(event eventstore.Event) error {
			events = append(events, event)
			return nil
		}), false)
		require.NoError(t, err)
		return events
	}

	// last reset of the failed password check count
	resets := filter(t, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		OrderDesc().
		Limit(1).
		AddQuery().
		AggregateTypes("user").
		AggregateIDs("user1").
		EventTypes("user.human.added", "user.human.password.changed", "user.human.password.check.succeeded").
		Builder(),
	)
	require.Len(t, resets, 1)
	lastReset := resets[0].Sequence()
	assert.Equal(t, uint64(4), lastReset)

	// password state and only the password checks after the last reset
	got := filter(t, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes("user").
		AggregateIDs("user1").
		EventTypes("user.human.added", "user.human.password.changed").
		Builder().
		AddQuery().
		AggregateTypes("user").
		AggregateIDs("user1").
		EventTypes("user.human.password.check.succeeded", "user.human.password.check.failed").
		SequenceGreater(lastReset).
		Builder(),
	)
	gotSequences := make([]uint64, len(got))
	for i, event := range got {
		gotSequences[i] = event.Sequence()
	}
	// replayed in sequence order, including the failure with the older position
	assert.Equal(t, []uint64{1, 2, 3, 5, 6}, gotSequences)
}
