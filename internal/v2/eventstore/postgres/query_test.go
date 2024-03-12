package postgres

import (
	"reflect"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/v2/database"
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
					eventstore.EventType("user.added"),
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
					eventstore.EventCreatedAt(
						database.NewNumberEquals(now),
					),
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
					eventstore.EventCreatedAt(
						database.NewNumberBetween(now, now.Add(time.Second)),
					),
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
					eventstore.EventSequence(
						database.NewNumberEquals(100),
					),
				),
			},
			want: wantQuery{
				query: "sequence = $1",
				args:  []any{100},
			},
		},
		{
			name: "sequence between",
			args: args{
				filter: eventstore.NewEventFilter(
					eventstore.EventSequence(
						database.NewNumberBetween(0, 10),
					),
				),
			},
			want: wantQuery{
				query: "sequence >= $1 AND sequence <= $2",
				args:  []any{0, 10},
			},
		},
		{
			name: "revision",
			args: args{
				filter: eventstore.NewEventFilter(
					eventstore.EventRevision(
						database.NewNumberAtLeast(2),
					),
				),
			},
			want: wantQuery{
				query: "revision >= $1",
				args:  []any{2},
			},
		},
		{
			name: "creator",
			args: args{
				filter: eventstore.NewEventFilter(
					eventstore.EventCreator("user-123"),
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
					eventstore.EventType("user.added"),
					eventstore.EventCreatedAt(database.NewNumberAtLeast(now)),
					eventstore.EventSequence(database.NewNumberGreater(10)),
					eventstore.EventRevision(database.NewNumberEquals(1)),
					eventstore.EventCreator("user-123"),
				),
			},
			want: wantQuery{
				query: "(event_type = $1 AND created_at >= $2 AND sequence > $3 AND revision = $4 AND creator = $5)",
				args:  []any{"user.added", now, 10, 1, "user-123"},
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
						eventstore.EventType("user.added"),
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
						eventstore.EventType("user.added"),
					),
					eventstore.NewEventFilter(
						eventstore.EventType("org.added"),
						eventstore.EventSequence(database.NewNumberGreater(4)),
					),
				},
			},
			want: wantQuery{
				query: " AND (event_type = $1 OR (event_type = $2 AND sequence > $3))",
				args:  []any{"user.added", "org.added", 4},
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
					eventstore.AggregateID("234"),
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
						eventstore.EventType("user.added"),
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
					eventstore.AggregateID("123"),
					eventstore.AppendEvent(
						eventstore.EventType("user.added"),
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
					eventstore.AggregateID("123"),
					eventstore.AppendEvent(
						eventstore.EventType("user.added"),
						eventstore.EventSequence(database.NewNumberGreater(1)),
					),
				),
			},
			want: wantQuery{
				query: "(aggregate_type = $1 AND aggregate_id = $2 AND (event_type = $3 AND sequence > $4))",
				args:  []any{"user", "123", "user.added", 1},
			},
		},
		{
			name: "2 event filters all aggregate",
			args: args{
				filter: eventstore.NewAggregateFilter(
					"user",
					eventstore.AggregateID("123"),
					eventstore.AppendEvent(
						eventstore.EventType("user.added"),
					),
					eventstore.AppendEvent(
						eventstore.EventSequence(database.NewNumberGreater(1)),
					),
				),
			},
			want: wantQuery{
				query: "(aggregate_type = $1 AND aggregate_id = $2 AND (event_type = $3 OR sequence > $4))",
				args:  []any{"user", "123", "user.added", 1},
			},
		},
		{
			name: "2 event filters with multiple conditions all aggregate",
			args: args{
				filter: eventstore.NewAggregateFilter(
					"user",
					eventstore.AggregateID("123"),
					eventstore.AppendEvents(
						eventstore.NewEventFilter(
							eventstore.EventType("user.added"),
							eventstore.EventSequence(database.NewNumberGreater(1)),
						),
					),
					eventstore.AppendEvent(
						eventstore.EventType("user.changed"),
						eventstore.EventSequence(database.NewNumberGreater(4)),
					),
				),
			},
			want: wantQuery{
				query: "(aggregate_type = $1 AND aggregate_id = $2 AND ((event_type = $3 AND sequence > $4) OR (event_type = $5 AND sequence > $6)))",
				args:  []any{"user", "123", "user.added", 1, "user.changed", 4},
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
							eventstore.EventType("org.added"),
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
						eventstore.PositionGreater(123.4, 0),
					),
				),
			},
			want: wantQuery{
				query: " WHERE instance_id = $1 AND position > $2 ORDER BY position, in_tx_order",
				args:  []any{"i1", 123.4},
			},
		},
		{
			name: "position and inPositionOrder pagination",
			args: args{
				filter: eventstore.NewFilter(
					eventstore.FilterPagination(
						eventstore.PositionGreater(123.4, 12),
					),
				),
			},
			want: wantQuery{
				query: " WHERE instance_id = $1 AND ((position = $2 AND in_tx_order > $3) OR position > $4) ORDER BY position, in_tx_order",
				args:  []any{"i1", 123.4, uint32(12), 123.4},
			},
		},
		{
			name: "pagination",
			args: args{
				filter: eventstore.NewFilter(
					eventstore.FilterPagination(
						eventstore.Limit(10),
						eventstore.Offset(3),
						eventstore.PositionGreater(123.4, 12),
					),
				),
			},
			want: wantQuery{
				query: " WHERE instance_id = $1 AND ((position = $2 AND in_tx_order > $3) OR position > $4) ORDER BY position, in_tx_order LIMIT $5 OFFSET $6",
				args:  []any{"i1", 123.4, uint32(12), 123.4, uint32(10), uint32(3)},
			},
		},
		{
			name: "aggregate and pagination",
			args: args{
				filter: eventstore.NewFilter(
					eventstore.FilterPagination(
						eventstore.Limit(10),
						eventstore.Offset(3),
						eventstore.PositionGreater(123.4, 12),
					),
					eventstore.AppendAggregateFilter("user"),
				),
			},
			want: wantQuery{
				query: " WHERE instance_id = $1 AND aggregate_type = $2 AND ((position = $3 AND in_tx_order > $4) OR position > $5) ORDER BY position, in_tx_order LIMIT $6 OFFSET $7",
				args:  []any{"i1", "user", 123.4, uint32(12), 123.4, uint32(10), uint32(3)},
			},
		},
		{
			name: "aggregates and pagination",
			args: args{
				filter: eventstore.NewFilter(
					eventstore.FilterPagination(
						eventstore.Limit(10),
						eventstore.Offset(3),
						eventstore.PositionGreater(123.4, 12),
					),
					eventstore.AppendAggregateFilter("user"),
					eventstore.AppendAggregateFilter(
						"org",
						eventstore.AggregateID("o1"),
					),
				),
			},
			want: wantQuery{
				query: " WHERE instance_id = $1 AND (aggregate_type = $2 OR (aggregate_type = $3 AND aggregate_id = $4)) AND ((position = $5 AND in_tx_order > $6) OR position > $7) ORDER BY position, in_tx_order LIMIT $8 OFFSET $9",
				args:  []any{"i1", "user", "org", "o1", 123.4, uint32(12), 123.4, uint32(10), uint32(3)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stmt database.Statement
			stmt.SetNamedArg(instancePlaceholder, "i1")
			// ensure a parent is set on the filter
			eventstore.NewQuery("instance", eventstore.AppendFilters(tt.args.filter))

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
					eventstore.AppendFilters(
						eventstore.NewFilter(),
					),
				),
			},
			want: wantQuery{
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 ORDER BY position, in_tx_order)) ORDER BY position, in_tx_order`,
				args:  []any{"i1"},
			},
		},
		{
			name: "1 filter",
			args: args{
				query: eventstore.NewQuery(
					"i1",
					eventstore.AppendFilters(
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter(
								"user",
								eventstore.AggregateIDList(database.NewListContains("a", "b")),
							),
						),
					),
				),
			},
			want: wantQuery{
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $2 AND aggregate_id = ANY($3)) ORDER BY position, in_tx_order)) ORDER BY position, in_tx_order`,
				args:  []any{"i1", "user", []string{"a", "b"}},
			},
		},
		{
			name: "multiple filters",
			args: args{
				query: eventstore.NewQuery(
					"i1",
					eventstore.AppendFilters(
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter(
								"user",
								eventstore.AggregateIDList(database.NewListContains("a", "b")),
							),
						),
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter(
								"org",
								eventstore.AppendEvent(
									eventstore.EventType("org.added"),
								),
							),
						),
					),
				),
			},
			want: wantQuery{
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $2 AND aggregate_id = ANY($3)) ORDER BY position, in_tx_order) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $4 AND event_type = $5) ORDER BY position, in_tx_order)) ORDER BY position, in_tx_order`,
				args:  []any{"i1", "user", []string{"a", "b"}, "org", "org.added"},
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
					eventstore.AppendFilters(
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter("aggregate"),
						),
					),
				),
			},
			want: wantQuery{
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 ORDER BY position, in_tx_order)) ORDER BY position, in_tx_order`,
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
					eventstore.QueryPagination(
						eventstore.Descending(),
					),
					eventstore.AppendFilter(
						eventstore.AppendAggregateFilter("aggregate"),
					),
				),
			},
			want: wantQuery{
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 ORDER BY position DESC, in_tx_order DESC)) ORDER BY position DESC, in_tx_order DESC`,
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
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 ORDER BY position, in_tx_order) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $3 OR aggregate_type = $4) ORDER BY position, in_tx_order)) ORDER BY position, in_tx_order`,
				args: []any{
					"instance",
					"agg1",
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
					eventstore.AppendFilters(
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter("agg1", eventstore.AggregateID("id")),
						),
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter("agg2", eventstore.AggregateID("id2")),
							eventstore.AppendAggregateFilter("agg3"),
						),
					),
				),
			},
			want: wantQuery{
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $2 AND aggregate_id = $3) ORDER BY position, in_tx_order) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND ((aggregate_type = $4 AND aggregate_id = $5) OR aggregate_type = $6) ORDER BY position, in_tx_order)) ORDER BY position, in_tx_order`,
				args: []any{
					"instance",
					"agg1",
					"id",
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
					eventstore.AppendFilter(
						eventstore.AppendAggregateFilter(
							"agg1",
							eventstore.AggregateIDList(database.NewListContains("1", "2")),
						),
						eventstore.AppendAggregateFilter(
							"agg2",
							eventstore.AggregateID("3"),
						),
						eventstore.AppendAggregateFilter(
							"agg3",
							eventstore.AggregateID("3"),
						),
					),
				),
			},
			want: wantQuery{
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND ((aggregate_type = $2 AND aggregate_id = ANY($3)) OR (aggregate_type = $4 AND aggregate_id = $5) OR (aggregate_type = $6 AND aggregate_id = $7)) ORDER BY position, in_tx_order)) ORDER BY position, in_tx_order`,
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
					eventstore.AppendFilters(

						eventstore.NewFilter(
							eventstore.AppendAggregateFilter(
								"instance",
								eventstore.AppendEvent(
									eventstore.EventType("instance.added"),
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
									eventstore.EventType("instance.removed"),
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
									eventstore.EventType("instance.domain.primary.set"),
									eventstore.EventCreatorList(database.NewListNotContains("", "SYSTEM")),
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
									eventstore.EventType("project.added"),
									eventstore.EventCreatorList(database.NewListNotContains("", "SYSTEM")),
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
									eventstore.EventCreatorList(database.NewListNotContains("", "SYSTEM")),
									eventstore.EventType("project.application.added"),
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
									eventstore.EventType("user.token.added"),
								),
							),
							eventstore.FilterPagination(
								// used because we need to check for first login and an app which is not console
								eventstore.PositionGreater(12, 4),
							),
						),
						eventstore.NewFilter(
							eventstore.AppendAggregateFilter(
								"instance",
								eventstore.AppendEvent(
									eventstore.EventType("instance.idp.config.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("instance.idp.oauth.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("instance.idp.oidc.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("instance.idp.jwt.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("instance.idp.azure.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("instance.idp.github.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("instance.idp.github.enterprise.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("instance.idp.gitlab.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("instance.idp.gitlab.selfhosted.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("instance.idp.google.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("instance.idp.ldap.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("instance.idp.config.apple.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("instance.idp.saml.added"),
								),
							),
							eventstore.AppendAggregateFilter(
								"org",
								eventstore.AppendEvent(
									eventstore.EventType("org.idp.config.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("org.idp.oauth.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("org.idp.oidc.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("org.idp.jwt.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("org.idp.azure.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("org.idp.github.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("org.idp.github.enterprise.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("org.idp.gitlab.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("org.idp.gitlab.selfhosted.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("org.idp.google.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("org.idp.ldap.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("org.idp.config.apple.added"),
								),
								eventstore.AppendEvent(
									eventstore.EventType("org.idp.saml.added"),
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
									eventstore.EventType("instance.login.policy.idp.added"),
								),
							),
							eventstore.AppendAggregateFilter(
								"org",
								eventstore.AppendEvent(
									eventstore.EventType("org.login.policy.idp.added"),
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
									eventstore.EventType("instance.smtp.config.added"),
									eventstore.EventCreatorList(database.NewListNotContains("", "SYSTEM", "<SYSTEM-USER>")),
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
				query: `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM ((SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $2 AND event_type = $3) ORDER BY position, in_tx_order LIMIT $4) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $5 AND event_type = $6) ORDER BY position, in_tx_order LIMIT $7) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $8 AND (event_type = $9 AND NOT(creator = ANY($10)))) ORDER BY position, in_tx_order LIMIT $11) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $12 AND (event_type = $13 AND NOT(creator = ANY($14)))) ORDER BY position, in_tx_order LIMIT $15) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $16 AND (event_type = $17 AND NOT(creator = ANY($18)))) ORDER BY position, in_tx_order LIMIT $19) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $20 AND event_type = $21) AND ((position = $22 AND in_tx_order > $23) OR position > $24) ORDER BY position, in_tx_order) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND ((aggregate_type = $25 AND (event_type = $26 OR event_type = $27 OR event_type = $28 OR event_type = $29 OR event_type = $30 OR event_type = $31 OR event_type = $32 OR event_type = $33 OR event_type = $34 OR event_type = $35 OR event_type = $36 OR event_type = $37 OR event_type = $38)) OR (aggregate_type = $39 AND (event_type = $40 OR event_type = $41 OR event_type = $42 OR event_type = $43 OR event_type = $44 OR event_type = $45 OR event_type = $46 OR event_type = $47 OR event_type = $48 OR event_type = $49 OR event_type = $50 OR event_type = $51 OR event_type = $52))) ORDER BY position, in_tx_order LIMIT $53) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND ((aggregate_type = $54 AND event_type = $55) OR (aggregate_type = $56 AND event_type = $57)) ORDER BY position, in_tx_order LIMIT $58) UNION ALL (SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $59 AND (event_type = $60 AND NOT(creator = ANY($61)))) ORDER BY position, in_tx_order LIMIT $62)) ORDER BY position, in_tx_order`,
				args: []any{
					"instance",
					"instance",
					"instance.added",
					uint32(1),
					"instance",
					"instance.removed",
					uint32(1),
					"instance",
					"instance.domain.primary.set",
					[]string{"", "SYSTEM"},
					uint32(1),
					"project",
					"project.added",
					[]string{"", "SYSTEM"},
					uint32(1),
					"project",
					"project.application.added",
					[]string{"", "SYSTEM"},
					uint32(1),
					"user",
					"user.token.added",
					float64(12),
					uint32(4),
					float64(12),
					"instance",
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
					"org",
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
					uint32(1),
					"instance",
					"instance.login.policy.idp.added",
					"org",
					"org.login.policy.idp.added",
					uint32(1),
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
