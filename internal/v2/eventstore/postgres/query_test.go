package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/v2/database"
	"github.com/zitadel/zitadel/internal/v2/database/mock"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

func Test_writeOrdering(t *testing.T) {
	type args struct {
		descending bool
	}
	tests := []struct {
		name string
		args args
		want wantQuery
	}{
		{
			name: "asc",
			args: args{
				descending: false,
			},
			want: wantQuery{
				query: " ORDER BY position, in_tx_order",
			},
		},
		{
			name: "desc",
			args: args{
				descending: true,
			},
			want: wantQuery{
				query: " ORDER BY position DESC, in_tx_order DESC",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stmt database.Statement
			writeOrdering(&stmt, tt.args.descending)
			assertQuery(t, &stmt, tt.want)
		})
	}
}

func Test_writeConditionsIfSet(t *testing.T) {
	type args struct {
		conditions []*condition
		sep        string
	}
	tests := []struct {
		name string
		args args
		want wantQuery
	}{
		{
			name: "no condition",
			args: args{
				conditions: []*condition{},
				sep:        " AND ",
			},
			want: wantQuery{
				query: "",
				args:  []any{},
			},
		},
		{
			name: "1 condition set",
			args: args{
				conditions: []*condition{
					{column: "column", condition: database.NewTextEqual("asdf")},
				},
				sep: " AND ",
			},
			want: wantQuery{
				query: "column = $1",
				args:  []any{"asdf"},
			},
		},
		{
			name: "multiple conditions set",
			args: args{
				conditions: []*condition{
					{column: "column1", condition: database.NewTextEqual("asdf")},
					{column: "column2", condition: database.NewNumberAtLeast(12)},
					{column: "column3", condition: database.NewNumberBetween(1, 100)},
				},
				sep: " AND ",
			},
			want: wantQuery{
				query: "column1 = $1 AND column2 >= $2 AND column3 >= $3 AND column3 <= $4",
				args:  []any{"asdf", 12, 1, 100},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stmt database.Statement
			writeConditions(&stmt, tt.args.conditions, tt.args.sep)
			assertQuery(t, &stmt, tt.want)
		})
	}
}

func Test_writeEventFilter(t *testing.T) {
	now := time.Now()
	type args struct {
		filter *eventstore.EventFilter
	}
	tests := []struct {
		name string
		args args
		want wantQuery
	}{
		{
			name: "no filters",
			args: args{
				filter: &eventstore.EventFilter{},
			},
			want: wantQuery{
				query: "",
				args:  []any{},
			},
		},
		{
			name: "event_type",
			args: args{
				filter: eventstore.NewEventFilter(
					eventstore.SetEventType("user.added"),
				),
			},
			want: wantQuery{
				query: "event_type = $1",
				args:  []any{"user.added"},
			},
		},
		{
			name: "created_at",
			args: args{
				filter: eventstore.NewEventFilter(
					eventstore.EventCreatedAtEquals(now),
				),
			},
			want: wantQuery{
				query: "created_at = $1",
				args:  []any{now},
			},
		},
		{
			name: "created_at between",
			args: args{
				filter: eventstore.NewEventFilter(
					eventstore.EventCreatedAtBetween(now, now.Add(time.Second)),
				),
			},
			want: wantQuery{
				query: "created_at >= $1 AND created_at <= $2",
				args:  []any{now, now.Add(time.Second)},
			},
		},
		{
			name: "sequence",
			args: args{
				filter: eventstore.NewEventFilter(
					eventstore.EventSequenceEquals(100),
				),
			},
			want: wantQuery{
				query: "sequence = $1",
				args:  []any{uint32(100)},
			},
		},
		{
			name: "sequence between",
			args: args{
				filter: eventstore.NewEventFilter(
					eventstore.EventSequenceBetween(0, 10),
				),
			},
			want: wantQuery{
				query: "sequence >= $1 AND sequence <= $2",
				args:  []any{uint32(0), uint32(10)},
			},
		},
		{
			name: "revision",
			args: args{
				filter: eventstore.NewEventFilter(
					eventstore.EventRevisionAtLeast(2),
				),
			},
			want: wantQuery{
				query: "revision >= $1",
				args:  []any{uint16(2)},
			},
		},
		{
			name: "creator",
			args: args{
				filter: eventstore.NewEventFilter(
					eventstore.EventCreatorsEqual("user-123"),
				),
			},
			want: wantQuery{
				query: "creator = $1",
				args:  []any{"user-123"},
			},
		},
		{
			name: "all",
			args: args{
				filter: eventstore.NewEventFilter(
					eventstore.SetEventType("user.added"),
					eventstore.EventCreatedAtAtLeast(now),
					eventstore.EventSequenceGreater(10),
					eventstore.EventRevisionEquals(1),
					eventstore.EventCreatorsEqual("user-123"),
				),
			},
			want: wantQuery{
				query: "(event_type = $1 AND created_at >= $2 AND sequence > $3 AND revision = $4 AND creator = $5)",
				args:  []any{"user.added", now, uint32(10), uint16(1), "user-123"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stmt database.Statement
			writeEventFilter(&stmt, tt.args.filter)
			assertQuery(t, &stmt, tt.want)
		})
	}
}

func Test_writeEventFilters(t *testing.T) {
	type args struct {
		filters []*eventstore.EventFilter
	}
	tests := []struct {
		name string
		args args
		want wantQuery
	}{
		{
			name: "no filters",
			args: args{},
			want: wantQuery{
				query: "",
				args:  []any{},
			},
		},
		{
			name: "1 filter",
			args: args{
				filters: []*eventstore.EventFilter{
					eventstore.NewEventFilter(
						eventstore.SetEventType("user.added"),
					),
				},
			},
			want: wantQuery{
				query: " AND event_type = $1",
				args:  []any{"user.added"},
			},
		},
		{
			name: "multiple filters",
			args: args{
				filters: []*eventstore.EventFilter{
					eventstore.NewEventFilter(
						eventstore.SetEventType("user.added"),
					),
					eventstore.NewEventFilter(
						eventstore.SetEventType("org.added"),
						eventstore.EventSequenceGreater(4),
					),
				},
			},
			want: wantQuery{
				query: " AND (event_type = $1 OR (event_type = $2 AND sequence > $3))",
				args:  []any{"user.added", "org.added", uint32(4)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stmt database.Statement
			writeEventFilters(&stmt, tt.args.filters)
			assertQuery(t, &stmt, tt.want)
		})
	}
}

func Test_writeAggregateFilter(t *testing.T) {
	type args struct {
		filter *eventstore.AggregateFilter
	}
	tests := []struct {
		name string
		args args
		want wantQuery
	}{
		{
			name: "minimal",
			args: args{
				filter: eventstore.NewAggregateFilter(
					"user",
				),
			},
			want: wantQuery{
				query: "aggregate_type = $1",
				args:  []any{"user"},
			},
		},
		{
			name: "all on aggregate",
			args: args{
				filter: eventstore.NewAggregateFilter(
					"user",
					eventstore.SetAggregateID("234"),
				),
			},
			want: wantQuery{
				query: "(aggregate_type = $1 AND aggregate_id = $2)",
				args:  []any{"user", "234"},
			},
		},
		{
			name: "1 event filter minimal aggregate",
			args: args{
				filter: eventstore.NewAggregateFilter(
					"user",
					eventstore.AppendEvent(
						eventstore.SetEventType("user.added"),
					),
				),
			},
			want: wantQuery{
				query: "(aggregate_type = $1 AND event_type = $2)",
				args:  []any{"user", "user.added"},
			},
		},
		{
			name: "1 event filter all aggregate",
			args: args{
				filter: eventstore.NewAggregateFilter(
					"user",
					eventstore.SetAggregateID("123"),
					eventstore.AppendEvent(
						eventstore.SetEventType("user.added"),
					),
				),
			},
			want: wantQuery{
				query: "(aggregate_type = $1 AND aggregate_id = $2 AND event_type = $3)",
				args:  []any{"user", "123", "user.added"},
			},
		},
		{
			name: "1 event filter with multiple conditions all aggregate",
			args: args{
				filter: eventstore.NewAggregateFilter(
					"user",
					eventstore.SetAggregateID("123"),
					eventstore.AppendEvent(
						eventstore.SetEventType("user.added"),
						eventstore.EventSequenceGreater(1),
					),
				),
			},
			want: wantQuery{
				query: "(aggregate_type = $1 AND aggregate_id = $2 AND (event_type = $3 AND sequence > $4))",
				args:  []any{"user", "123", "user.added", uint32(1)},
			},
		},
		{
			name: "2 event filters all aggregate",
			args: args{
				filter: eventstore.NewAggregateFilter(
					"user",
					eventstore.SetAggregateID("123"),
					eventstore.AppendEvent(
						eventstore.SetEventType("user.added"),
					),
					eventstore.AppendEvent(
						eventstore.EventSequenceGreater(1),
					),
				),
			},
			want: wantQuery{
				query: "(aggregate_type = $1 AND aggregate_id = $2 AND (event_type = $3 OR sequence > $4))",
				args:  []any{"user", "123", "user.added", uint32(1)},
			},
		},
		{
			name: "2 event filters with multiple conditions all aggregate",
			args: args{
				filter: eventstore.NewAggregateFilter(
					"user",
					eventstore.SetAggregateID("123"),
					eventstore.AppendEvents(
						eventstore.NewEventFilter(
							eventstore.SetEventType("user.added"),
							eventstore.EventSequenceGreater(1),
						),
					),
					eventstore.AppendEvent(
						eventstore.SetEventType("user.changed"),
						eventstore.EventSequenceGreater(4),
					),
				),
			},
			want: wantQuery{
				query: "(aggregate_type = $1 AND aggregate_id = $2 AND ((event_type = $3 AND sequence > $4) OR (event_type = $5 AND sequence > $6)))",
				args:  []any{"user", "123", "user.added", uint32(1), "user.changed", uint32(4)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stmt database.Statement
			writeAggregateFilter(&stmt, tt.args.filter)
			assertQuery(t, &stmt, tt.want)
		})
	}
}

func Test_writeAggregateFilters(t *testing.T) {
	type args struct {
		filters []*eventstore.AggregateFilter
	}
	tests := []struct {
		name string
		args args
		want wantQuery
	}{
		{
			name: "no filters",
			args: args{},
			want: wantQuery{
				query: "",
				args:  []any{},
			},
		},
		{
			name: "1 filter",
			args: args{
				filters: []*eventstore.AggregateFilter{
					eventstore.NewAggregateFilter("user"),
				},
			},
			want: wantQuery{
				query: " AND aggregate_type = $1",
				args:  []any{"user"},
			},
		},
		{
			name: "multiple filters",
			args: args{
				filters: []*eventstore.AggregateFilter{
					eventstore.NewAggregateFilter("user"),
					eventstore.NewAggregateFilter("org",
						eventstore.AppendEvent(
							eventstore.SetEventType("org.added"),
						),
					),
				},
			},
			want: wantQuery{
				query: " AND (aggregate_type = $1 OR (aggregate_type = $2 AND event_type = $3))",
				args:  []any{"user", "org", "org.added"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stmt database.Statement
			writeAggregateFilters(&stmt, tt.args.filters)
			assertQuery(t, &stmt, tt.want)
		})
	}
}

func Test_writeFilter(t *testing.T) {
	type args struct {
		filter *eventstore.Filter
	}
	tests := []struct {
		name string
		args args
		want wantQuery
	}{
		{
			name: "empty filters",
			args: args{
				filter: eventstore.NewFilter(),
			},
			want: wantQuery{
				query: " WHERE instance_id = $1 ORDER BY position, in_tx_order",
				args:  []any{"i1"},
			},
		},
		{
			name: "descending",
			args: args{
				filter: eventstore.NewFilter(
					eventstore.FilterPagination(
						eventstore.Descending(),
					),
				),
			},
			want: wantQuery{
				query: " WHERE instance_id = $1 ORDER BY position DESC, in_tx_order DESC",
				args:  []any{"i1"},
			},
		},
		{
			name: "database pagination",
			args: args{
				filter: eventstore.NewFilter(
					eventstore.FilterPagination(
						eventstore.Limit(10),
						eventstore.Offset(3),
					),
				),
			},
			want: wantQuery{
				query: " WHERE instance_id = $1 ORDER BY position, in_tx_order LIMIT $2 OFFSET $3",
				args:  []any{"i1", uint32(10), uint32(3)},
			},
		},
		{
			name: "position pagination",
			args: args{
				filter: eventstore.NewFilter(
					eventstore.FilterPagination(
						eventstore.PositionGreater(decimal.NewFromFloat(123.4), 0),
					),
				),
			},
			want: wantQuery{
				query: " WHERE instance_id = $1 AND position > $2 ORDER BY position, in_tx_order",
				args:  []any{"i1", decimal.NewFromFloat(123.4)},
			},
		},
		{
			name: "position pagination between",
			args: args{
				filter: eventstore.NewFilter(
					eventstore.FilterPagination(
						// 	eventstore.PositionGreater(decimal.NewFromFloat(123.4), 0),
						// 	eventstore.PositionLess(125.4, 10),
						eventstore.PositionBetween(
							&eventstore.GlobalPosition{Position: decimal.NewFromFloat(123.4)},
							&eventstore.GlobalPosition{Position: decimal.NewFromFloat(125.4), InPositionOrder: 10},
						),
					),
				),
			},
			want: wantQuery{
				query: " WHERE instance_id = $1 AND ((position = $2 AND in_tx_order < $3) OR position < $4) AND position > $5 ORDER BY position, in_tx_order",
				args:  []any{"i1", decimal.NewFromFloat(125.4), uint32(10), decimal.NewFromFloat(125.4), decimal.NewFromFloat(123.4)},
				// TODO: (adlerhurst) would require some refactoring to reuse existing args
				// query: " WHERE instance_id = $1 AND position > $2 AND ((position = $3 AND in_tx_order < $4) OR position < $3) ORDER BY position, in_tx_order",
				// args:  []any{"i1", 123.4, 125.4, uint32(10)},
			},
		},
		{
			name: "position and inPositionOrder pagination",
			args: args{
				filter: eventstore.NewFilter(
					eventstore.FilterPagination(
						eventstore.PositionGreater(decimal.NewFromFloat(123.4), 12),
					),
				),
			},
			want: wantQuery{
				query: " WHERE instance_id = $1 AND ((position = $2 AND in_tx_order > $3) OR position > $4) ORDER BY position, in_tx_order",
				args:  []any{"i1", decimal.NewFromFloat(123.4), uint32(12), decimal.NewFromFloat(123.4)},
			},
		},
		{
			name: "pagination",
			args: args{
				filter: eventstore.NewFilter(
					eventstore.FilterPagination(
						eventstore.Limit(10),
						eventstore.Offset(3),
						eventstore.PositionGreater(decimal.NewFromFloat(123.4), 12),
					),
				),
			},
			want: wantQuery{
				query: " WHERE instance_id = $1 AND ((position = $2 AND in_tx_order > $3) OR position > $4) ORDER BY position, in_tx_order LIMIT $5 OFFSET $6",
				args:  []any{"i1", decimal.NewFromFloat(123.4), uint32(12), decimal.NewFromFloat(123.4), uint32(10), uint32(3)},
			},
		},
		{
			name: "aggregate and pagination",
			args: args{
				filter: eventstore.NewFilter(
					eventstore.FilterPagination(
						eventstore.Limit(10),
						eventstore.Offset(3),
						eventstore.PositionGreater(decimal.NewFromFloat(123.4), 12),
					),
					eventstore.AppendAggregateFilter("user"),
				),
			},
			want: wantQuery{
				query: " WHERE instance_id = $1 AND aggregate_type = $2 AND ((position = $3 AND in_tx_order > $4) OR position > $5) ORDER BY position, in_tx_order LIMIT $6 OFFSET $7",
				args:  []any{"i1", "user", decimal.NewFromFloat(123.4), uint32(12), decimal.NewFromFloat(123.4), uint32(10), uint32(3)},
			},
		},
		{
			name: "aggregates and pagination",
			args: args{
				filter: eventstore.NewFilter(
					eventstore.FilterPagination(
						eventstore.Limit(10),
						eventstore.Offset(3),
						eventstore.PositionGreater(decimal.NewFromFloat(123.4), 12),
					),
					eventstore.AppendAggregateFilter("user"),
					eventstore.AppendAggregateFilter(
						"org",
						eventstore.SetAggregateID("o1"),
					),
				),
			},
			want: wantQuery{
				query: " WHERE instance_id = $1 AND (aggregate_type = $2 OR (aggregate_type = $3 AND aggregate_id = $4)) AND ((position = $5 AND in_tx_order > $6) OR position > $7) ORDER BY position, in_tx_order LIMIT $8 OFFSET $9",
				args:  []any{"i1", "user", "org", "o1", decimal.NewFromFloat(123.4), uint32(12), decimal.NewFromFloat(123.4), uint32(10), uint32(3)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stmt database.Statement
			eventstore.NewQuery("i1", nil, eventstore.AppendFilters(tt.args.filter))

			writeFilter(&stmt, tt.args.filter)
			assertQuery(t, &stmt, tt.want)
		})
	}
}

func Test_writeQuery(t *testing.T) {
	type args struct {
		query *eventstore.Query
	}
	tests := []struct {
		name string
		args args
		want wantQuery
	}{
		{
			name: "empty filter",
			args: args{
				query: eventstore.NewQuery(
					"i1",
					nil,
					eventstore.AppendFilters(
						eventstore.NewFilter(),
					),
				),
			},
			want: wantQuery{
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 ORDER BY position, in_tx_order)) sub ORDER BY position, in_tx_order`,
				args:  []any{"i1"},
			},
		},
		{
			name: "1 filter",
			args: args{
				query: eventstore.NewQuery(
					"i1",
					nil,
					eventstore.AppendFilters(
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter(
								"user",
								eventstore.AggregateIDs("a", "b"),
							),
						),
					),
				),
			},
			want: wantQuery{
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $2 AND aggregate_id = ANY($3)) ORDER BY position, in_tx_order)) sub ORDER BY position, in_tx_order`,
				args:  []any{"i1", "user", []string{"a", "b"}},
			},
		},
		{
			name: "multiple filters",
			args: args{
				query: eventstore.NewQuery(
					"i1",
					nil,
					eventstore.AppendFilters(
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter(
								"user",
								eventstore.AggregateIDs("a", "b"),
							),
						),
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter(
								"org",
								eventstore.AppendEvent(
									eventstore.SetEventType("org.added"),
								),
							),
						),
					),
				),
			},
			want: wantQuery{
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $2 AND aggregate_id = ANY($3)) ORDER BY position, in_tx_order) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $4 AND (aggregate_type = $5 AND event_type = $6) ORDER BY position, in_tx_order)) sub ORDER BY position, in_tx_order`,
				args:  []any{"i1", "user", []string{"a", "b"}, "i1", "org", "org.added"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stmt database.Statement
			writeQuery(&stmt, tt.args.query)
			assertQuery(t, &stmt, tt.want)
		})
	}
}

func Test_writeQueryUse_examples(t *testing.T) {
	type args struct {
		query *eventstore.Query
	}
	tests := []struct {
		name string
		args args
		want wantQuery
	}{
		{
			name: "aggregate type",
			args: args{
				query: eventstore.NewQuery(
					"instance",
					nil,
					eventstore.AppendFilters(
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter("aggregate"),
						),
					),
				),
			},
			want: wantQuery{
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 ORDER BY position, in_tx_order)) sub ORDER BY position, in_tx_order`,
				args: []any{
					"instance",
					"aggregate",
				},
			},
		},
		{
			name: "descending",
			args: args{
				query: eventstore.NewQuery(
					"instance",
					nil,
					eventstore.QueryPagination(
						eventstore.Descending(),
					),
					eventstore.AppendFilter(
						eventstore.AppendAggregateFilter("aggregate"),
					),
				),
			},
			want: wantQuery{
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 ORDER BY position DESC, in_tx_order DESC)) sub ORDER BY position DESC, in_tx_order DESC`,
				args: []any{
					"instance",
					"aggregate",
				},
			},
		},
		{
			name: "multiple aggregates",
			args: args{
				query: eventstore.NewQuery(
					"instance",
					nil,
					eventstore.AppendFilters(
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter("agg1"),
						),
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter("agg2"),
							eventstore.AppendAggregateFilter("agg3"),
						),
					),
				),
			},
			want: wantQuery{
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 ORDER BY position, in_tx_order) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $3 AND (aggregate_type = $4 OR aggregate_type = $5) ORDER BY position, in_tx_order)) sub ORDER BY position, in_tx_order`,
				args: []any{
					"instance",
					"agg1",
					"instance",
					"agg2",
					"agg3",
				},
			},
		},
		{
			name: "multiple aggregates with ids",
			args: args{
				query: eventstore.NewQuery(
					"instance",
					nil,
					eventstore.AppendFilters(
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter("agg1", eventstore.SetAggregateID("id")),
						),
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter("agg2", eventstore.SetAggregateID("id2")),
							eventstore.AppendAggregateFilter("agg3"),
						),
					),
				),
			},
			want: wantQuery{
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $2 AND aggregate_id = $3) ORDER BY position, in_tx_order) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $4 AND ((aggregate_type = $5 AND aggregate_id = $6) OR aggregate_type = $7) ORDER BY position, in_tx_order)) sub ORDER BY position, in_tx_order`,
				args: []any{
					"instance",
					"agg1",
					"id",
					"instance",
					"agg2",
					"id2",
					"agg3",
				},
			},
		},
		{
			name: "multiple event queries and multiple filter in queries",
			args: args{
				query: eventstore.NewQuery(
					"instance",
					nil,
					eventstore.AppendFilter(
						eventstore.AppendAggregateFilter(
							"agg1",
							eventstore.AggregateIDs("1", "2"),
						),
						eventstore.AppendAggregateFilter(
							"agg2",
							eventstore.SetAggregateID("3"),
						),
						eventstore.AppendAggregateFilter(
							"agg3",
							eventstore.SetAggregateID("3"),
						),
					),
				),
			},
			want: wantQuery{
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND ((aggregate_type = $2 AND aggregate_id = ANY($3)) OR (aggregate_type = $4 AND aggregate_id = $5) OR (aggregate_type = $6 AND aggregate_id = $7)) ORDER BY position, in_tx_order)) sub ORDER BY position, in_tx_order`,
				args: []any{
					"instance",
					"agg1",
					[]string{"1", "2"},
					"agg2",
					"3",
					"agg3",
					"3",
				},
			},
		},
		{
			name: "milestones",
			args: args{
				query: eventstore.NewQuery(
					"instance",
					nil,
					eventstore.AppendFilters(
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter(
								"instance",
								eventstore.AppendEvent(
									eventstore.SetEventType("instance.added"),
								),
							),
							eventstore.FilterPagination(
								eventstore.Limit(1),
							),
						),
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter(
								"instance",
								eventstore.AppendEvent(
									eventstore.SetEventType("instance.removed"),
								),
							),
							eventstore.FilterPagination(
								eventstore.Limit(1),
							),
						),
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter(
								"instance",
								eventstore.AppendEvent(
									eventstore.SetEventType("instance.domain.primary.set"),
									eventstore.EventCreatorsNotContains("", "SYSTEM"),
								),
							),
							eventstore.FilterPagination(
								eventstore.Limit(1),
							),
						),
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter(
								"project",
								eventstore.AppendEvent(
									eventstore.SetEventType("project.added"),
									eventstore.EventCreatorsNotContains("", "SYSTEM"),
								),
							),
							eventstore.FilterPagination(
								eventstore.Limit(1),
							),
						),
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter(
								"project",
								eventstore.AppendEvent(
									eventstore.EventCreatorsNotContains("", "SYSTEM"),
									eventstore.SetEventType("project.application.added"),
								),
							),
							eventstore.FilterPagination(
								eventstore.Limit(1),
							),
						),
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter(
								"user",
								eventstore.AppendEvent(
									eventstore.SetEventType("user.token.added"),
								),
							),
							eventstore.FilterPagination(
								// used because we need to check for first login and an app which is not console
								eventstore.PositionGreater(decimal.NewFromInt(12), 4),
							),
						),
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter(
								"instance",
								eventstore.AppendEvent(
									eventstore.SetEventTypes(
										"instance.idp.config.added",
										"instance.idp.oauth.added",
										"instance.idp.oidc.added",
										"instance.idp.jwt.added",
										"instance.idp.azure.added",
										"instance.idp.github.added",
										"instance.idp.github.enterprise.added",
										"instance.idp.gitlab.added",
										"instance.idp.gitlab.selfhosted.added",
										"instance.idp.google.added",
										"instance.idp.ldap.added",
										"instance.idp.config.apple.added",
										"instance.idp.saml.added",
									),
								),
							),
							eventstore.AppendAggregateFilter(
								"org",
								eventstore.AppendEvent(
									eventstore.SetEventTypes(
										"org.idp.config.added",
										"org.idp.oauth.added",
										"org.idp.oidc.added",
										"org.idp.jwt.added",
										"org.idp.azure.added",
										"org.idp.github.added",
										"org.idp.github.enterprise.added",
										"org.idp.gitlab.added",
										"org.idp.gitlab.selfhosted.added",
										"org.idp.google.added",
										"org.idp.ldap.added",
										"org.idp.config.apple.added",
										"org.idp.saml.added",
									),
								),
							),
							eventstore.FilterPagination(
								eventstore.Limit(1),
							),
						),
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter(
								"instance",
								eventstore.AppendEvent(
									eventstore.SetEventType("instance.login.policy.idp.added"),
								),
							),
							eventstore.AppendAggregateFilter(
								"org",
								eventstore.AppendEvent(
									eventstore.SetEventType("org.login.policy.idp.added"),
								),
							),
							eventstore.FilterPagination(
								eventstore.Limit(1),
							),
						),
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter(
								"instance",
								eventstore.AppendEvent(
									eventstore.SetEventType("instance.smtp.config.added"),
									eventstore.EventCreatorsNotContains("", "SYSTEM", "<SYSTEM-USER>"),
								),
							),
							eventstore.FilterPagination(
								eventstore.Limit(1),
							),
						),
					),
				),
			},
			want: wantQuery{
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $2 AND event_type = $3) ORDER BY position, in_tx_order LIMIT $4) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $5 AND (aggregate_type = $6 AND event_type = $7) ORDER BY position, in_tx_order LIMIT $8) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $9 AND (aggregate_type = $10 AND (event_type = $11 AND NOT(creator = ANY($12)))) ORDER BY position, in_tx_order LIMIT $13) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $14 AND (aggregate_type = $15 AND (event_type = $16 AND NOT(creator = ANY($17)))) ORDER BY position, in_tx_order LIMIT $18) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $19 AND (aggregate_type = $20 AND (event_type = $21 AND NOT(creator = ANY($22)))) ORDER BY position, in_tx_order LIMIT $23) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $24 AND (aggregate_type = $25 AND event_type = $26) AND ((position = $27 AND in_tx_order > $28) OR position > $29) ORDER BY position, in_tx_order) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $30 AND ((aggregate_type = $31 AND event_type = ANY($32)) OR (aggregate_type = $33 AND event_type = ANY($34))) ORDER BY position, in_tx_order LIMIT $35) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $36 AND ((aggregate_type = $37 AND event_type = $38) OR (aggregate_type = $39 AND event_type = $40)) ORDER BY position, in_tx_order LIMIT $41) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $42 AND (aggregate_type = $43 AND (event_type = $44 AND NOT(creator = ANY($45)))) ORDER BY position, in_tx_order LIMIT $46)) sub ORDER BY position, in_tx_order`,
				args: []any{
					"instance",
					"instance",
					"instance.added",
					uint32(1),
					"instance",
					"instance",
					"instance.removed",
					uint32(1),
					"instance",
					"instance",
					"instance.domain.primary.set",
					[]string{"", "SYSTEM"},
					uint32(1),
					"instance",
					"project",
					"project.added",
					[]string{"", "SYSTEM"},
					uint32(1),
					"instance",
					"project",
					"project.application.added",
					[]string{"", "SYSTEM"},
					uint32(1),
					"instance",
					"user",
					"user.token.added",
					decimal.NewFromInt(12),
					uint32(4),
					decimal.NewFromInt(12),
					"instance",
					"instance",
					[]string{"instance.idp.config.added", "instance.idp.oauth.added", "instance.idp.oidc.added", "instance.idp.jwt.added", "instance.idp.azure.added", "instance.idp.github.added", "instance.idp.github.enterprise.added", "instance.idp.gitlab.added", "instance.idp.gitlab.selfhosted.added", "instance.idp.google.added", "instance.idp.ldap.added", "instance.idp.config.apple.added", "instance.idp.saml.added"},
					"org",
					[]string{"org.idp.config.added", "org.idp.oauth.added", "org.idp.oidc.added", "org.idp.jwt.added", "org.idp.azure.added", "org.idp.github.added", "org.idp.github.enterprise.added", "org.idp.gitlab.added", "org.idp.gitlab.selfhosted.added", "org.idp.google.added", "org.idp.ldap.added", "org.idp.config.apple.added", "org.idp.saml.added"},
					uint32(1),
					"instance",
					"instance",
					"instance.login.policy.idp.added",
					"org",
					"org.login.policy.idp.added",
					uint32(1),
					"instance",
					"instance",
					"instance.smtp.config.added",
					[]string{"", "SYSTEM", "<SYSTEM-USER>"},
					uint32(1),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stmt database.Statement
			writeQuery(&stmt, tt.args.query)
			assertQuery(t, &stmt, tt.want)
		})
	}
}

type wantQuery struct {
	query string
	args  []any
}

func assertQuery(t *testing.T, stmt *database.Statement, want wantQuery) bool {
	t.Helper()
	ok := true

	defer func() {
		if !ok {
			t.Logf("generated statement: %s\n", stmt.Debug())
		}
	}()

	got := stmt.String()
	if got != want.query {
		t.Errorf("unexpected query:\n want: %q\n  got: %q", want.query, got)
		ok = false
	}

	if len(want.args) != len(stmt.Args()) {
		t.Errorf("unexpected length of args, want: %d got: %d", len(want.args), len(stmt.Args()))
		return false
	}

	for i, arg := range want.args {
		if !reflect.DeepEqual(arg, stmt.Args()[i]) {
			t.Errorf("unexpected arg at %d, want %v got: %v", i, arg, stmt.Args()[i])
			ok = false
		}
	}

	return ok
}

var _ eventstore.Reducer = (*testReducer)(nil)

type testReducer struct {
	expectedReduces int
	reduceCount     int
	shouldErr       bool
}

// Reduce implements eventstore.Reducer.
func (r *testReducer) Reduce(events ...*eventstore.StorageEvent) error {
	if r == nil {
		return nil
	}
	r.reduceCount++
	if r.shouldErr {
		return errReduce
	}
	return nil
}

func (r *testReducer) assert(t *testing.T) {
	if r.expectedReduces == r.reduceCount {
		return
	}

	t.Errorf("unexpected reduces, want %d, got %d", r.expectedReduces, r.reduceCount)
}

func Test_executeQuery(t *testing.T) {
	type args struct {
		values  [][]driver.Value
		reducer *testReducer
	}
	type want struct {
		eventCount int
		assertErr  func(t *testing.T, err error) bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no result",
			args: args{
				values:  [][]driver.Value{},
				reducer: &testReducer{},
			},
			want: want{
				eventCount: 0,
				assertErr: func(t *testing.T, err error) bool {
					is := errors.Is(err, nil)
					if !is {
						t.Errorf("no error expected got: %v", err)
					}
					return is
				},
			},
		},
		{
			name: "1 event without payload",
			args: args{
				values: [][]driver.Value{
					{
						time.Now(),
						"event.type",
						uint32(23),
						decimal.NewFromInt(123).String(),
						uint32(0),
						nil,
						"gigi",
						"owner",
						"instance",
						"aggregate.type",
						"aggregate.id",
						uint16(1),
					},
				},
				reducer: &testReducer{
					expectedReduces: 1,
				},
			},
			want: want{
				eventCount: 1,
				assertErr: func(t *testing.T, err error) bool {
					is := errors.Is(err, nil)
					if !is {
						t.Errorf("no error expected got: %v", err)
					}
					return is
				},
			},
		},
		{
			name: "1 event with payload",
			args: args{
				values: [][]driver.Value{
					{
						time.Now(),
						"event.type",
						uint32(23),
						decimal.NewFromInt(123).String(),
						uint32(0),
						[]byte(`{"name": "gigi"}`),
						"gigi",
						"owner",
						"instance",
						"aggregate.type",
						"aggregate.id",
						uint16(1),
					},
				},
				reducer: &testReducer{
					expectedReduces: 1,
				},
			},
			want: want{
				eventCount: 1,
				assertErr: func(t *testing.T, err error) bool {
					is := errors.Is(err, nil)
					if !is {
						t.Errorf("no error expected got: %v", err)
					}
					return is
				},
			},
		},
		{
			name: "multiple events",
			args: args{
				values: [][]driver.Value{
					{
						time.Now(),
						"event.type",
						uint32(23),
						decimal.NewFromInt(123).String(),
						uint32(0),
						nil,
						"gigi",
						"owner",
						"instance",
						"aggregate.type",
						"aggregate.id",
						uint16(1),
					},
					{
						time.Now(),
						"event.type",
						uint32(24),
						decimal.NewFromInt(124).String(),
						uint32(0),
						[]byte(`{"name": "gigi"}`),
						"gigi",
						"owner",
						"instance",
						"aggregate.type",
						"aggregate.id",
						uint16(1),
					},
				},
				reducer: &testReducer{
					expectedReduces: 2,
				},
			},
			want: want{
				eventCount: 2,
				assertErr: func(t *testing.T, err error) bool {
					is := errors.Is(err, nil)
					if !is {
						t.Errorf("no error expected got: %v", err)
					}
					return is
				},
			},
		},
		{
			name: "reduce error",
			args: args{
				values: [][]driver.Value{
					{
						time.Now(),
						"event.type",
						uint32(23),
						decimal.NewFromInt(123).String(),
						uint32(0),
						nil,
						"gigi",
						"owner",
						"instance",
						"aggregate.type",
						"aggregate.id",
						uint16(1),
					},
					{
						time.Now(),
						"event.type",
						uint32(24),
						decimal.NewFromInt(124).String(),
						uint32(0),
						[]byte(`{"name": "gigi"}`),
						"gigi",
						"owner",
						"instance",
						"aggregate.type",
						"aggregate.id",
						uint16(1),
					},
				},
				reducer: &testReducer{
					expectedReduces: 1,
					shouldErr:       true,
				},
			},
			want: want{
				eventCount: 1,
				assertErr: func(t *testing.T, err error) bool {
					is := errors.Is(err, errReduce)
					if !is {
						t.Errorf("no error expected got: %v", err)
					}
					return is
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := mock.NewSQLMock(t,
				mock.ExpectQuery(
					"",
					mock.WithQueryResult(
						[]string{"created_at", "event_type", "sequence", "position", "in_tx_order", "payload", "creator", "owner", "instance_id", "aggregate_type", "aggregate_id", "revision"},
						tt.args.values,
					),
				),
			)
			gotEventCount, err := executeQuery(context.Background(), mockDB.DB, &database.Statement{}, tt.args.reducer)
			tt.want.assertErr(t, err)
			if gotEventCount != tt.want.eventCount {
				t.Errorf("executeQuery() = %v, want %v", gotEventCount, tt.want.eventCount)
			}
		})
	}
}
