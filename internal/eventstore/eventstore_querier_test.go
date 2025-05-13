package eventstore_test

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func TestCRDB_Filter(t *testing.T) {
	type args struct {
		searchQuery *eventstore.SearchQueryBuilder
	}
	type fields struct {
		existingEvents []eventstore.Command
	}
	type res struct {
		eventCount int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		res     res
		wantErr bool
	}{
		{
			name: "aggregate type filter no events",
			args: args{
				searchQuery: eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
					AddQuery().
					AggregateTypes("not found").
					Builder(),
			},
			fields: fields{
				existingEvents: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "300"),
					generateCommand(eventstore.AggregateType(t.Name()), "300"),
					generateCommand(eventstore.AggregateType(t.Name()), "300"),
				},
			},
			res: res{
				eventCount: 0,
			},
			wantErr: false,
		},
		{
			name: "aggregate type and id filter events found",
			args: args{
				searchQuery: eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
					AddQuery().
					AggregateTypes(eventstore.AggregateType(t.Name())).
					AggregateIDs("303").
					Builder(),
			},
			fields: fields{
				existingEvents: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "303"),
					generateCommand(eventstore.AggregateType(t.Name()), "303"),
					generateCommand(eventstore.AggregateType(t.Name()), "303"),
					generateCommand(eventstore.AggregateType(t.Name()), "305"),
				},
			},
			res: res{
				eventCount: 3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		for querierName, querier := range queriers {
			t.Run(querierName+"/"+tt.name, func(t *testing.T) {
				t.Cleanup(cleanupEventstore(clients[querierName]))

				db := eventstore.NewEventstore(
					&eventstore.Config{
						Querier: querier,
						Pusher:  pushers["v3(inmemory)"],
					},
				)

				// setup initial data for query
				if _, err := db.Push(context.Background(), tt.fields.existingEvents...); err != nil {
					t.Errorf("error in setup = %v", err)
					return
				}

				events, err := db.Filter(context.Background(), tt.args.searchQuery)
				if (err != nil) != tt.wantErr {
					t.Errorf("CRDB.query() error = %v, wantErr %v", err, tt.wantErr)
				}

				if len(events) != tt.res.eventCount {
					t.Errorf("CRDB.query() expected event count: %d got %d", tt.res.eventCount, len(events))
				}
			})
		}
	}
}

func TestCRDB_LatestPosition(t *testing.T) {
	type args struct {
		searchQuery *eventstore.SearchQueryBuilder
	}
	type fields struct {
		existingEvents []eventstore.Command
	}
	type res struct {
		position decimal.Decimal
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		res     res
		wantErr bool
	}{
		{
			name: "aggregate type filter no sequence",
			args: args{
				searchQuery: eventstore.NewSearchQueryBuilder(eventstore.ColumnsMaxPosition).
					AddQuery().
					AggregateTypes("not found").
					Builder(),
			},
			fields: fields{
				existingEvents: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "400"),
					generateCommand(eventstore.AggregateType(t.Name()), "400"),
					generateCommand(eventstore.AggregateType(t.Name()), "400"),
				},
			},
			wantErr: false,
		},
		{
			name: "aggregate type filter sequence",
			args: args{
				searchQuery: eventstore.NewSearchQueryBuilder(eventstore.ColumnsMaxPosition).
					AddQuery().
					AggregateTypes(eventstore.AggregateType(t.Name())).
					Builder(),
			},
			fields: fields{
				existingEvents: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "401"),
					generateCommand(eventstore.AggregateType(t.Name()), "401"),
					generateCommand(eventstore.AggregateType(t.Name()), "401"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		for querierName, querier := range queriers {
			t.Run(querierName+"/"+tt.name, func(t *testing.T) {
				t.Cleanup(cleanupEventstore(clients[querierName]))

				db := eventstore.NewEventstore(
					&eventstore.Config{
						Querier: querier,
						Pusher:  pushers["v3(inmemory)"],
					},
				)

				// setup initial data for query
				_, err := db.Push(context.Background(), tt.fields.existingEvents...)
				if err != nil {
					t.Errorf("error in setup = %v", err)
					return
				}

				position, err := db.LatestPosition(context.Background(), tt.args.searchQuery)
				if (err != nil) != tt.wantErr {
					t.Errorf("CRDB.query() error = %v, wantErr %v", err, tt.wantErr)
				}
				if tt.res.position.GreaterThan(position) {
					t.Errorf("CRDB.query() expected sequence: %v got %v", tt.res.position, position)
				}
			})
		}
	}
}
