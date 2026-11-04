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

// Test_query_sub_query_position_after_postgres reads the password state of a user together with
// only the password checks after the last reset, as done by the command side password check.
func Test_query_sub_query_position_after_postgres(t *testing.T) {
	ctx := t.Context()
	instanceID := "position-after-" + t.Name()
	client := &Postgres{DB: &database.DB{DB: testClient, Database: new(testDB)}}

	t.Cleanup(func() {
		_, err := testClient.ExecContext(context.Background(), "DELETE FROM eventstore.events2 WHERE instance_id = $1", instanceID)
		assert.NoError(t, err)
	})

	eventTypes := []eventstore.EventType{
		"user.human.added",
		"user.human.password.changed",
		"user.human.password.check.succeeded",
		"user.human.password.check.failed",
		"user.human.password.check.succeeded",
		"user.human.password.check.failed",
		"user.human.password.check.failed",
	}
	for i, eventType := range eventTypes {
		_, err := testClient.ExecContext(ctx,
			`INSERT INTO eventstore.events2 (instance_id, aggregate_type, aggregate_id, event_type, "sequence", revision, created_at, payload, creator, "owner", "position", in_tx_order)
			VALUES ($1, 'user', 'user1', $2, $3, 1, now(), NULL, 'creator', 'owner', $4, 1)`,
			instanceID, eventType, i+1, decimal.NewFromInt(int64(100+i)),
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
	lastReset := resets[0].Position()
	assert.True(t, lastReset.Equal(decimal.NewFromInt(104)), "got %s", lastReset)

	// password state and only the password checks after the last reset
	events := filter(t, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes("user").
		AggregateIDs("user1").
		EventTypes("user.human.added", "user.human.password.changed").
		Builder().
		AddQuery().
		AggregateTypes("user").
		AggregateIDs("user1").
		EventTypes("user.human.password.check.succeeded", "user.human.password.check.failed").
		PositionAfter(lastReset).
		Builder(),
	)
	gotTypes := make([]eventstore.EventType, len(events))
	for i, event := range events {
		gotTypes[i] = event.Type()
	}
	assert.Equal(t, []eventstore.EventType{
		"user.human.added",
		"user.human.password.changed",
		"user.human.password.check.failed",
		"user.human.password.check.failed",
	}, gotTypes)
}
