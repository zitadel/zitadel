package postgres

import (
	"context"
	"reflect"
	"testing"

	"github.com/zitadel/zitadel/internal/v2/database"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

func Test_filterQuery(t *testing.T) {
	type args struct {
		filter *eventstore.Filter
	}
	type want struct {
		query string
		args  []any
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "only instance",
			args: args{
				filter: eventstore.NewFilter(
					context.Background(),
					eventstore.FilterInstances("instance"),
				),
			},
			want: want{
				query: `SELECT created_at, event_type, "sequence", "position", payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = ANY($1) ORDER BY position, in_tx_order`,
				args: []any{
					[]string{"", "instance"},
				},
			},
		},
		{
			name: "descending",
			args: args{
				filter: eventstore.NewFilter(
					context.Background(),
					eventstore.FilterDescending(),
				),
			},
			want: want{
				query: `SELECT created_at, event_type, "sequence", "position", payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 ORDER BY position DESC, in_tx_order DESC`,
				args: []any{
					"",
				},
			},
		},
		{
			name: "event query",
			args: args{
				filter: eventstore.NewFilter(
					context.Background(),
					eventstore.FilterEventQuery(
						eventstore.FilterAggregateTypes("agg"),
					),
				),
			},
			want: want{
				query: `SELECT created_at, event_type, "sequence", "position", payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 ORDER BY position, in_tx_order`,
				args: []any{
					"",
					"agg",
				},
			},
		},
		{
			name: "multiple event queries",
			args: args{
				filter: eventstore.NewFilter(
					context.Background(),
					eventstore.FilterEventQuery(
						eventstore.FilterAggregateTypes("agg1"),
					),
					eventstore.FilterEventQuery(
						eventstore.FilterAggregateTypes("agg2", "agg3"),
					),
				),
			},
			want: want{
				query: `SELECT created_at, event_type, "sequence", "position", payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $2 OR aggregate_type = ANY($3)) ORDER BY position, in_tx_order`,
				args: []any{
					"",
					"agg1",
					[]string{"agg2", "agg3"},
				},
			},
		},
		{
			name: "multiple event queries and multiple filter in queries",
			args: args{
				filter: eventstore.NewFilter(
					context.Background(),
					eventstore.FilterEventQuery(
						eventstore.FilterAggregateTypes("agg1"),
						eventstore.FilterAggregateIDs("1", "2"),
					),
					eventstore.FilterEventQuery(
						eventstore.FilterAggregateTypes("agg2", "agg3"),
						eventstore.FilterAggregateIDs("3"),
					),
				),
			},
			want: want{
				query: `SELECT created_at, event_type, "sequence", "position", payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND ((aggregate_type = $2 AND aggregate_id = ANY($3)) OR (aggregate_type = ANY($4) AND aggregate_id = $5)) ORDER BY position, in_tx_order`,
				args: []any{
					"",
					"agg1",
					[]string{"1", "2"},
					[]string{"agg2", "agg3"},
					"3",
				},
			},
		},
		// {
		// 	name: "",
		// 	args: args{
		// 		filter: eventstore.NewFilter(
		// 			context.Background(),
		// 		),
		// 	},
		// 	want: want{
		// 		query: ``,
		// 		args:  []any{},
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stmt database.Statement
			filterQuery(&stmt, tt.args.filter)

			if got := stmt.String(); got != tt.want.query {
				t.Errorf("unexpected query:\nwant: %q\n got: %q", tt.want.query, got)
			}
			if len(stmt.Args()) != len(tt.want.args) {
				t.Errorf("unexpected length of args, want: %d got %d", len(tt.want.args), len(stmt.Args()))
				return
			}
			for i, got := range stmt.Args() {
				if !reflect.DeepEqual(got, tt.want.args[i]) {
					t.Errorf("unexpected arg at position %d: want %v got %v", i, tt.want.args[i], got)
				}
			}
		})
	}
}
