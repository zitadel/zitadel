package sql

import (
	"context"
	"database/sql/driver"
	"regexp"
	"strings"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

func Test_query_event_sort_key(t *testing.T) {
	cursor := eventstore.EventSortKey{
		Position:      decimal.NewFromFloat(1.5),
		InTxOrder:     5,
		AggregateType: "user",
		AggregateID:   "agg-a",
		Sequence:      5,
	}
	tuple := `(` + eventSortKeySQL + `) > (`

	tests := []struct {
		name      string
		query     *eventstore.SearchQueryBuilder
		sql       string
		args      []driver.Value
		awaitLock bool
	}{
		{
			name: "after sort key with await cap",
			query: eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
				InstanceID("instanceID").
				OrderAsc().
				AwaitOpenTransactions().
				Limit(200).
				AfterEventSortKey(cursor).
				AddQuery().
				AggregateTypes("user").
				EventTypes("user.human.added").
				Builder(),
			sql: `SELECT created_at, event_type, "sequence", "position", payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision, in_tx_order FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 AND event_type = $3 AND (` + eventSortKeySQL + `) > ($4, $5, $6, $7, $8) AND "position" <= EXTRACT(EPOCH FROM now()) ORDER BY ` + eventSortKeySQL + ` LIMIT $9`,
			args: []driver.Value{
				"instanceID",
				eventstore.AggregateType("user"),
				eventstore.EventType("user.human.added"),
				cursor.Position,
				cursor.InTxOrder,
				cursor.AggregateType,
				cursor.AggregateID,
				cursor.Sequence,
				uint64(200),
			},
			awaitLock: true,
		},
		{
			name: "after sort key without await",
			query: eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
				InstanceID("instanceID").
				OrderAsc().
				Limit(200).
				AfterEventSortKey(cursor).
				AddQuery().
				AggregateTypes("user").
				Builder(),
			sql: `SELECT created_at, event_type, "sequence", "position", payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision, in_tx_order FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 AND (` + eventSortKeySQL + `) > ($3, $4, $5, $6, $7) ORDER BY ` + eventSortKeySQL + ` LIMIT $8`,
			args: []driver.Value{
				"instanceID",
				eventstore.AggregateType("user"),
				cursor.Position,
				cursor.InTxOrder,
				cursor.AggregateType,
				cursor.AggregateID,
				cursor.Sequence,
				uint64(200),
			},
		},
	}

	client := NewPostgres(&database.DB{Database: new(testDB)})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if strings.Contains(tt.sql, "OFFSET") {
				t.Fatal("sort-key resume must not use OFFSET")
			}
			if !strings.Contains(tt.sql, tuple) {
				t.Fatalf("expected sort-key tuple %s in SQL", tuple)
			}

			mock := newMockClient(t)
			if tt.awaitLock {
				mock = mock.expectExec(regexp.QuoteMeta(
					`select pg_advisory_lock('eventstore.events2'::REGCLASS::OID::INTEGER, hashtext($1)), pg_advisory_unlock('eventstore.events2'::REGCLASS::OID::INTEGER, hashtext($1))`),
					[]driver.Value{"instanceID"},
				)
			}
			mock = mock.expectQuery(regexp.QuoteMeta(tt.sql), tt.args)
			client.DB.DB = mock.client

			err := query(context.Background(), client, tt.query, &[]*repository.Event{}, false)
			if err != nil {
				t.Fatalf("query() error = %v", err)
			}
			if err := mock.mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all expectations met: %v", err)
			}
		})
	}
}
