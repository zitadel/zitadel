package sql

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

func Test_query_event_type_scans(t *testing.T) {
	cursor := eventstore.EventSortKey{
		Position:      decimal.NewFromFloat(1.5),
		InTxOrder:     5,
		InstanceID:    "instanceID",
		AggregateType: "user",
		AggregateID:   "agg-a",
		Sequence:      5,
	}
	scanSQL := func(positionCondition string) string {
		return `SELECT events.* FROM UNNEST($1::TEXT[], $2::TEXT[]) AS scans(scan_aggregate_type, scan_event_type)` +
			` CROSS JOIN LATERAL (SELECT ` + eventColumnsSQL + ` FROM eventstore.events2` +
			` WHERE instance_id = $3 AND aggregate_type = scans.scan_aggregate_type AND event_type = scans.scan_event_type` +
			positionCondition +
			` AND "position" <= EXTRACT(EPOCH FROM now()) ORDER BY ` + eventSortKeySQL + ` LIMIT $%d) AS events` +
			` ORDER BY ` + eventSortKeySQL + ` LIMIT $%d`
	}
	projectionQuery := func() *eventstore.SearchQueryBuilder {
		return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
			InstanceID("instanceID").
			OrderAsc().
			AwaitOpenTransactions().
			Limit(200).
			ScanEventTypesSeparately()
	}
	aggregateTypes := database.TextArray[eventstore.AggregateType]{"org", "user", "user"}
	eventTypes := database.TextArray[eventstore.EventType]{"org.removed", "user.human.added", "user.locked"}

	tests := []struct {
		name  string
		query *eventstore.SearchQueryBuilder
		sql   string
		args  []driver.Value
	}{
		{
			name: "after sort key, one scan per aggregate and event type",
			query: projectionQuery().
				AfterEventSortKey(cursor).
				AddQuery().AggregateTypes("user").EventTypes("user.locked", "user.human.added").Builder().
				AddQuery().AggregateTypes("org").EventTypes("org.removed").Builder(),
			sql: fmt.Sprintf(scanSQL(` AND (`+eventSortKeySQL+`) > ($4, $5, $6, $7, $8, $9)`), 10, 11),
			args: []driver.Value{
				aggregateTypes, eventTypes, "instanceID",
				cursor.Position, cursor.InTxOrder, cursor.InstanceID, cursor.AggregateType, cursor.AggregateID, cursor.Sequence,
				uint64(200), uint64(200),
			},
		},
		{
			name: "position at least, combinations listed twice are scanned once",
			query: projectionQuery().
				PositionAtLeast(decimal.NewFromInt(10)).
				AddQuery().AggregateTypes("user").EventTypes("user.locked", "user.human.added").Builder().
				AddQuery().AggregateTypes("org").EventTypes("org.removed").Builder().
				AddQuery().AggregateTypes("user").EventTypes("user.locked").Builder(),
			sql: fmt.Sprintf(scanSQL(` AND "position" >= $4`), 5, 6),
			args: []driver.Value{
				aggregateTypes, eventTypes, "instanceID",
				decimal.NewFromInt(10),
				uint64(200), uint64(200),
			},
		},
		{
			name: "sub query filtering by aggregate id is not scanned separately",
			query: projectionQuery().
				AfterEventSortKey(cursor).
				AddQuery().AggregateTypes("user").AggregateIDs("agg-a").EventTypes("user.locked").Builder(),
			sql: `SELECT ` + eventColumnsSQL + ` FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 AND aggregate_id = $3 AND event_type = $4 AND (` + eventSortKeySQL + `) > ($5, $6, $7, $8, $9, $10) AND "position" <= EXTRACT(EPOCH FROM now()) ORDER BY ` + eventSortKeySQL + ` LIMIT $11`,
			args: []driver.Value{
				"instanceID", eventstore.AggregateType("user"), "agg-a", eventstore.EventType("user.locked"),
				cursor.Position, cursor.InTxOrder, cursor.InstanceID, cursor.AggregateType, cursor.AggregateID, cursor.Sequence,
				uint64(200),
			},
		},
		{
			name: "descending order is not scanned separately",
			query: eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
				InstanceID("instanceID").
				OrderDesc().
				AwaitOpenTransactions().
				Limit(200).
				ScanEventTypesSeparately().
				AddQuery().AggregateTypes("org").EventTypes("org.removed").Builder(),
			sql: `SELECT ` + eventColumnsSQL + ` FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 AND event_type = $3 AND "position" <= EXTRACT(EPOCH FROM now()) ORDER BY (` + eventSortKeySQL + `) DESC LIMIT $4`,
			args: []driver.Value{
				"instanceID", eventstore.AggregateType("org"), eventstore.EventType("org.removed"),
				uint64(200),
			},
		},
	}

	client := NewPostgres(&database.DB{Database: new(testDB)})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockClient(t).
				expectExec(regexp.QuoteMeta(
					`select pg_advisory_lock('eventstore.events2'::REGCLASS::OID::INTEGER, hashtext($1)), pg_advisory_unlock('eventstore.events2'::REGCLASS::OID::INTEGER, hashtext($1))`),
					[]driver.Value{"instanceID"},
				).
				expectQuery("^"+regexp.QuoteMeta(tt.sql)+"$", tt.args)
			client.DB.DB = mock.client

			require.NoError(t, query(t.Context(), client, tt.query, &[]*repository.Event{}, false))
			assert.NoError(t, mock.mock.ExpectationsWereMet())
		})
	}
}

// Test_query_event_type_scans_postgres checks against a real database that reading every combination of
// aggregate type and event type separately returns exactly the events, in exactly the order,
// of the single query for all event types, from any cursor and for any limit.
// The events include positions shared by several transactions, whose in_tx_order restarts at 1,
// which is the situation where resuming from a cursor has skipped events before.
func Test_query_event_type_scans_postgres(t *testing.T) {
	ctx := t.Context()
	instanceID := "event-type-scans-" + t.Name()
	client := &Postgres{DB: &database.DB{DB: testClient, Database: new(testDB)}}

	t.Cleanup(func() {
		_, err := testClient.ExecContext(context.Background(), "DELETE FROM eventstore.events2 WHERE instance_id = $1", instanceID)
		assert.NoError(t, err)
	})

	type aggregateEvents struct {
		aggregateType eventstore.AggregateType
		eventTypes    []eventstore.EventType
	}
	// the last event type of every aggregate is written but not subscribed to
	aggregates := []aggregateEvents{
		{"user", []eventstore.EventType{"user.human.added", "user.locked", "user.human.changed"}},
		{"org", []eventstore.EventType{"org.removed", "org.added"}},
		{"session", []eventstore.EventType{"session.added", "session.token.set", "session.terminated"}},
	}
	sequences := make(map[string]uint64)
	for tx := range 90 {
		// every three transactions share a position, each numbering its events from 1
		position := decimal.NewFromInt(1_000_000).Add(decimal.NewFromInt(int64(tx / 3)))
		for inTxOrder := range tx%4 + 1 {
			agg := aggregates[(tx+inTxOrder)%len(aggregates)]
			eventType := agg.eventTypes[(tx*7+inTxOrder)%len(agg.eventTypes)]
			aggregateID := fmt.Sprintf("agg-%d", (tx*3+inTxOrder)%5)
			key := string(agg.aggregateType) + aggregateID
			sequences[key]++
			_, err := testClient.ExecContext(ctx,
				`INSERT INTO eventstore.events2 (instance_id, aggregate_type, aggregate_id, event_type, "sequence", revision, created_at, payload, creator, "owner", "position", in_tx_order)
				VALUES ($1, $2, $3, $4, $5, 1, now(), NULL, 'creator', 'owner', $6, $7)`,
				instanceID, agg.aggregateType, aggregateID, eventType, sequences[key], position, inTxOrder+1,
			)
			require.NoError(t, err)
		}
	}

	subscribe := func(builder *eventstore.SearchQueryBuilder) *eventstore.SearchQueryBuilder {
		for _, agg := range aggregates {
			builder = builder.AddQuery().AggregateTypes(agg.aggregateType).EventTypes(agg.eventTypes[:len(agg.eventTypes)-1]...).Builder()
		}
		return builder
	}
	filter := func(t *testing.T, scanSeparately bool, limit uint64, cursor *eventstore.EventSortKey, minPosition decimal.Decimal) []eventstore.EventSortKey {
		builder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
			InstanceID(instanceID).
			AwaitOpenTransactions().
			OrderAsc().
			Limit(limit)
		if scanSeparately {
			builder = builder.ScanEventTypesSeparately()
		}
		if cursor != nil {
			builder = builder.AfterEventSortKey(*cursor)
		}
		if !minPosition.IsZero() {
			builder = builder.PositionAtLeast(minPosition)
		}
		keys := []eventstore.EventSortKey{}
		err := query(ctx, client, subscribe(builder), eventstore.Reducer(func(event eventstore.Event) error {
			keys = append(keys, eventstore.EventSortKey{
				Position:      event.Position(),
				InTxOrder:     event.InTxOrder(),
				InstanceID:    event.Aggregate().InstanceID,
				AggregateType: event.Aggregate().Type,
				AggregateID:   event.Aggregate().ID,
				Sequence:      event.Sequence(),
			})
			return nil
		}), false)
		require.NoError(t, err)
		return keys
	}

	all := filter(t, false, 10_000, nil, decimal.Decimal{})
	require.Len(t, filter(t, true, 10_000, nil, decimal.Decimal{}), len(all))
	require.Greater(t, len(all), 100, "the generated events must include subscribed events")

	for _, limit := range []uint64{1, 2, 3, 7, 50} {
		t.Run(fmt.Sprintf("limit %d", limit), func(t *testing.T) {
			// from the start
			assert.Equal(t, filter(t, false, limit, nil, decimal.Decimal{}), filter(t, true, limit, nil, decimal.Decimal{}))
			// from after every event
			for i := range all {
				cursor := all[i]
				want := filter(t, false, limit, &cursor, decimal.Decimal{})
				require.Equal(t, all[i+1:min(len(all), i+1+int(limit))], want, "the single query must page through all events")
				assert.Equal(t, want, filter(t, true, limit, &cursor, decimal.Decimal{}), "cursor %d", i)
			}
			// from every position, inclusive
			for i := range all {
				assert.Equal(t,
					filter(t, false, limit, nil, all[i].Position),
					filter(t, true, limit, nil, all[i].Position),
					"position %s", all[i].Position,
				)
			}
		})
	}
}
