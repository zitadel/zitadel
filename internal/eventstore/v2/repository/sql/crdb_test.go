package sql

import (
	"context"
	"sync"
	"testing"

	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func TestCRDB_placeholder(t *testing.T) {
	type args struct {
		query string
	}
	type res struct {
		query string
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no placeholders",
			args: args{
				query: "SELECT * FROM eventstore.events",
			},
			res: res{
				query: "SELECT * FROM eventstore.events",
			},
		},
		{
			name: "one placeholder",
			args: args{
				query: "SELECT * FROM eventstore.events WHERE aggregate_type = ?",
			},
			res: res{
				query: "SELECT * FROM eventstore.events WHERE aggregate_type = $1",
			},
		},
		{
			name: "multiple placeholders",
			args: args{
				query: "SELECT * FROM eventstore.events WHERE aggregate_type = ? AND aggregate_id = ? LIMIT ?",
			},
			res: res{
				query: "SELECT * FROM eventstore.events WHERE aggregate_type = $1 AND aggregate_id = $2 LIMIT $3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{}
			if query := db.placeholder(tt.args.query); query != tt.res.query {
				t.Errorf("CRDB.placeholder() = %v, want %v", query, tt.res.query)
			}
		})
	}
}

func TestCRDB_operation(t *testing.T) {
	type res struct {
		op string
	}
	type args struct {
		operation repository.Operation
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no op",
			args: args{
				operation: repository.Operation(-1),
			},
			res: res{
				op: "",
			},
		},
		{
			name: "greater",
			args: args{
				operation: repository.OperationGreater,
			},
			res: res{
				op: ">",
			},
		},
		{
			name: "less",
			args: args{
				operation: repository.OperationLess,
			},
			res: res{
				op: "<",
			},
		},
		{
			name: "equals",
			args: args{
				operation: repository.OperationEquals,
			},
			res: res{
				op: "=",
			},
		},
		{
			name: "in",
			args: args{
				operation: repository.OperationIn,
			},
			res: res{
				op: "=",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{}
			if got := db.operation(tt.args.operation); got != tt.res.op {
				t.Errorf("CRDB.operation() = %v, want %v", got, tt.res.op)
			}
		})
	}
}

func TestCRDB_conditionFormat(t *testing.T) {
	type res struct {
		format string
	}
	type args struct {
		operation repository.Operation
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "default",
			args: args{
				operation: repository.OperationEquals,
			},
			res: res{
				format: "%s %s ?",
			},
		},
		{
			name: "in",
			args: args{
				operation: repository.OperationIn,
			},
			res: res{
				format: "%s %s ANY(?)",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{}
			if got := db.conditionFormat(tt.args.operation); got != tt.res.format {
				t.Errorf("CRDB.conditionFormat() = %v, want %v", got, tt.res.format)
			}
		})
	}
}

func TestCRDB_columnName(t *testing.T) {
	type res struct {
		name string
	}
	type args struct {
		field repository.Field
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "invalid field",
			args: args{
				field: repository.Field(-1),
			},
			res: res{
				name: "",
			},
		},
		{
			name: "aggregate id",
			args: args{
				field: repository.FieldAggregateID,
			},
			res: res{
				name: "aggregate_id",
			},
		},
		{
			name: "aggregate type",
			args: args{
				field: repository.FieldAggregateType,
			},
			res: res{
				name: "aggregate_type",
			},
		},
		{
			name: "editor service",
			args: args{
				field: repository.FieldEditorService,
			},
			res: res{
				name: "editor_service",
			},
		},
		{
			name: "editor user",
			args: args{
				field: repository.FieldEditorUser,
			},
			res: res{
				name: "editor_user",
			},
		},
		{
			name: "event type",
			args: args{
				field: repository.FieldEventType,
			},
			res: res{
				name: "event_type",
			},
		},
		{
			name: "latest sequence",
			args: args{
				field: repository.FieldSequence,
			},
			res: res{
				name: "event_sequence",
			},
		},
		{
			name: "resource owner",
			args: args{
				field: repository.FieldResourceOwner,
			},
			res: res{
				name: "resource_owner",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{}
			if got := db.columnName(tt.args.field); got != tt.res.name {
				t.Errorf("CRDB.operation() = %v, want %v", got, tt.res.name)
			}
		})
	}
}

func TestCRDB_Push_OneAggregate(t *testing.T) {
	type args struct {
		ctx    context.Context
		events []*repository.Event
	}
	type eventsRes struct {
		pushedEventsCount int
		aggType           repository.AggregateType
		aggID             []string
	}
	type res struct {
		wantErr   bool
		eventsRes eventsRes
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "push 1 event with check previous",
			args: args{
				ctx: context.Background(),
				events: []*repository.Event{
					generateEvent(t, "1", true, 0),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 1,
					aggID:             []string{"1"},
					aggType:           repository.AggregateType(t.Name()),
				}},
		},
		{
			name: "fail push 1 event with check previous wrong sequence",
			args: args{
				ctx: context.Background(),
				events: []*repository.Event{
					generateEvent(t, "2", true, 5),
				},
			},
			res: res{
				wantErr: true,
				eventsRes: eventsRes{
					pushedEventsCount: 0,
					aggID:             []string{"2"},
					aggType:           repository.AggregateType(t.Name()),
				},
			},
		},
		{
			name: "push 1 event without check previous",
			args: args{
				ctx: context.Background(),
				events: []*repository.Event{
					generateEvent(t, "3", false, 0),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 1,
					aggID:             []string{"3"},
					aggType:           repository.AggregateType(t.Name()),
				},
			},
		},
		{
			name: "push 1 event without check previous wrong sequence",
			args: args{
				ctx: context.Background(),
				events: []*repository.Event{
					generateEvent(t, "4", false, 5),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 1,
					aggID:             []string{"4"},
					aggType:           repository.AggregateType(t.Name()),
				},
			},
		},
		{
			name: "fail on push two events on agg without linking",
			args: args{
				ctx: context.Background(),
				events: []*repository.Event{
					generateEvent(t, "5", true, 0),
					generateEvent(t, "5", true, 0),
				},
			},
			res: res{
				wantErr: true,
				eventsRes: eventsRes{
					pushedEventsCount: 0,
					aggID:             []string{"5"},
					aggType:           repository.AggregateType(t.Name()),
				},
			},
		},
		{
			name: "push two events on agg with linking",
			args: args{
				ctx: context.Background(),
				events: linkEvents(
					generateEvent(t, "6", true, 0),
					generateEvent(t, "6", true, 0),
				),
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 2,
					aggID:             []string{"6"},
					aggType:           repository.AggregateType(t.Name()),
				},
			},
		},
		{
			name: "push two events on agg with linking without check previous",
			args: args{
				ctx: context.Background(),
				events: linkEvents(
					generateEvent(t, "7", false, 0),
					generateEvent(t, "7", false, 0),
				),
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 2,
					aggID:             []string{"7"},
					aggType:           repository.AggregateType(t.Name()),
				},
			},
		},
		{
			name: "push two events on agg with linking mixed check previous",
			args: args{
				ctx: context.Background(),
				events: linkEvents(
					generateEvent(t, "8", false, 0),
					generateEvent(t, "8", true, 0),
					generateEvent(t, "8", false, 0),
					generateEvent(t, "8", true, 0),
					generateEvent(t, "8", true, 0),
				),
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 5,
					aggID:             []string{"8"},
					aggType:           repository.AggregateType(t.Name()),
				},
			},
		},
		{
			name: "failed push because context canceled",
			args: args{
				ctx: canceledCtx(),
				events: []*repository.Event{
					generateEvent(t, "9", true, 0),
				},
			},
			res: res{
				wantErr: true,
				eventsRes: eventsRes{
					pushedEventsCount: 0,
					aggID:             []string{"9"},
					aggType:           repository.AggregateType(t.Name()),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{
				client: testCRDBClient,
			}
			if err := db.Push(tt.args.ctx, tt.args.events...); (err != nil) != tt.res.wantErr {
				t.Errorf("CRDB.Push() error = %v, wantErr %v", err, tt.res.wantErr)
			}

			countRow := testCRDBClient.QueryRow("SELECT COUNT(*) FROM eventstore.events where aggregate_type = $1 AND aggregate_id = ANY($2)", tt.res.eventsRes.aggType, pq.Array(tt.res.eventsRes.aggID))
			var count int
			err := countRow.Scan(&count)
			if err != nil {
				t.Error("unable to query inserted rows: ", err)
				return
			}
			if count != tt.res.eventsRes.pushedEventsCount {
				t.Errorf("expected push count %d got %d", tt.res.eventsRes.pushedEventsCount, count)
			}
		})
	}
}

func TestCRDB_Push_MultipleAggregate(t *testing.T) {
	type args struct {
		events []*repository.Event
	}
	type eventsRes struct {
		pushedEventsCount int
		aggType           []repository.AggregateType
		aggID             []string
	}
	type res struct {
		wantErr   bool
		eventsRes eventsRes
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "push two aggregates both check previous",
			args: args{
				events: []*repository.Event{
					generateEvent(t, "100", true, 0),
					generateEvent(t, "101", true, 0),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 2,
					aggID:             []string{"100", "101"},
					aggType:           []repository.AggregateType{repository.AggregateType(t.Name())},
				},
			},
		},
		{
			name: "push two aggregates both check previous multiple events",
			args: args{
				events: combineEventLists(
					linkEvents(
						generateEvent(t, "102", true, 0),
						generateEvent(t, "102", true, 0),
					),
					linkEvents(
						generateEvent(t, "103", true, 0),
						generateEvent(t, "103", true, 0),
					),
				),
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 4,
					aggID:             []string{"102", "103"},
					aggType:           []repository.AggregateType{repository.AggregateType(t.Name())},
				},
			},
		},
		{
			name: "fail push linked events of different aggregates",
			args: args{
				events: linkEvents(
					generateEvent(t, "104", false, 0),
					generateEvent(t, "105", false, 0),
				),
			},
			res: res{
				wantErr: true,
				eventsRes: eventsRes{
					pushedEventsCount: 0,
					aggID:             []string{"104", "105"},
					aggType:           []repository.AggregateType{repository.AggregateType(t.Name())},
				},
			},
		},
		{
			name: "push two aggregates mixed check previous multiple events",
			args: args{
				events: combineEventLists(
					linkEvents(
						generateEvent(t, "106", true, 0),
						generateEvent(t, "106", false, 0),
						generateEvent(t, "106", false, 0),
						generateEvent(t, "106", true, 0),
					),
					linkEvents(
						generateEvent(t, "107", false, 0),
						generateEvent(t, "107", true, 0),
						generateEvent(t, "107", false, 0),
						generateEvent(t, "107", true, 0),
					),
					linkEvents(
						generateEvent(t, "108", true, 0),
						generateEvent(t, "108", false, 0),
						generateEvent(t, "108", false, 0),
						generateEvent(t, "108", true, 0),
					),
				),
			},
		},
		{
			name: "failed push same aggregate in two transactions",
			args: args{
				events: combineEventLists(
					linkEvents(
						generateEvent(t, "109", true, 0),
					),
					linkEvents(
						generateEvent(t, "109", true, 0),
					),
				),
			},
			res: res{
				wantErr: true,
				eventsRes: eventsRes{
					pushedEventsCount: 0,
					aggID:             []string{"109"},
					aggType:           []repository.AggregateType{repository.AggregateType(t.Name())},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{
				client: testCRDBClient,
			}
			if err := db.Push(context.Background(), tt.args.events...); (err != nil) != tt.res.wantErr {
				t.Errorf("CRDB.Push() error = %v, wantErr %v", err, tt.res.wantErr)
			}

			countRow := testCRDBClient.QueryRow("SELECT COUNT(*) FROM eventstore.events where aggregate_type = ANY($1) AND aggregate_id = ANY($2)", pq.Array(tt.res.eventsRes.aggType), pq.Array(tt.res.eventsRes.aggID))
			var count int
			err := countRow.Scan(&count)
			if err != nil {
				t.Error("unable to query inserted rows: ", err)
				return
			}
			if count != tt.res.eventsRes.pushedEventsCount {
				t.Errorf("expected push count %d got %d", tt.res.eventsRes.pushedEventsCount, count)
			}
		})
	}
}

func TestCRDB_Push_Parallel(t *testing.T) {
	type args struct {
		events [][]*repository.Event
	}
	type eventsRes struct {
		pushedEventsCount int
		aggTypes          []repository.AggregateType
		aggIDs            []string
	}
	type res struct {
		errCount  int
		eventsRes eventsRes
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "clients push different aggregates",
			args: args{
				events: [][]*repository.Event{
					linkEvents(
						generateEvent(t, "200", false, 0),
						generateEvent(t, "200", true, 0),
						generateEvent(t, "200", false, 0),
					),
					linkEvents(
						generateEvent(t, "201", false, 0),
						generateEvent(t, "201", true, 0),
						generateEvent(t, "201", false, 0),
					),
					combineEventLists(
						linkEvents(
							generateEvent(t, "202", false, 0),
						),
						linkEvents(
							generateEvent(t, "203", true, 0),
							generateEvent(t, "203", false, 0),
						),
					),
				},
			},
			res: res{
				errCount: 0,
				eventsRes: eventsRes{
					aggIDs:            []string{"200", "201", "202", "203"},
					pushedEventsCount: 9,
					aggTypes:          []repository.AggregateType{repository.AggregateType(t.Name())},
				},
			},
		},
		{
			name: "clients push same aggregates no check previous",
			args: args{
				events: [][]*repository.Event{
					linkEvents(
						generateEvent(t, "204", false, 0),
						generateEvent(t, "204", false, 0),
					),
					linkEvents(
						generateEvent(t, "204", false, 0),
						generateEvent(t, "204", false, 0),
					),
					combineEventLists(
						linkEvents(
							generateEvent(t, "205", false, 0),
							generateEvent(t, "205", false, 0),
							generateEvent(t, "205", false, 0),
						),
						linkEvents(
							generateEvent(t, "206", false, 0),
							generateEvent(t, "206", false, 0),
							generateEvent(t, "206", false, 0),
						),
					),
					combineEventLists(
						linkEvents(
							generateEvent(t, "204", false, 0),
						),
						linkEvents(
							generateEvent(t, "205", false, 0),
							generateEvent(t, "205", false, 0),
						),
						linkEvents(
							generateEvent(t, "206", false, 0),
						),
					),
				},
			},
			res: res{
				errCount: 0,
				eventsRes: eventsRes{
					aggIDs:            []string{"204", "205", "206"},
					pushedEventsCount: 14,
					aggTypes:          []repository.AggregateType{repository.AggregateType(t.Name())},
				},
			},
		},
		{
			name: "clients push different aggregates one with check previous",
			args: args{
				events: [][]*repository.Event{
					linkEvents(
						generateEvent(t, "207", false, 0),
						generateEvent(t, "207", false, 0),
						generateEvent(t, "207", false, 0),
						generateEvent(t, "207", false, 0),
						generateEvent(t, "207", false, 0),
						generateEvent(t, "207", false, 0),
					),
					linkEvents(
						generateEvent(t, "208", true, 0),
						generateEvent(t, "208", true, 0),
						generateEvent(t, "208", true, 0),
						generateEvent(t, "208", true, 0),
						generateEvent(t, "208", true, 0),
					),
				},
			},
			res: res{
				errCount: 0,
				eventsRes: eventsRes{
					aggIDs:            []string{"207", "208"},
					pushedEventsCount: 11,
					aggTypes:          []repository.AggregateType{repository.AggregateType(t.Name())},
				},
			},
		},
		{
			name: "clients push different aggregates all with check previous on first event fail",
			args: args{
				events: [][]*repository.Event{
					linkEvents(
						generateEventWithData(t, "210", true, 0, []byte(`{ "transaction": 1 }`)),
						generateEventWithData(t, "210", false, 0, []byte(`{ "transaction": 1.1 }`)),
					),
					linkEvents(
						generateEventWithData(t, "210", true, 0, []byte(`{ "transaction": 2 }`)),
						generateEventWithData(t, "210", false, 0, []byte(`{ "transaction": 2.1 }`)),
					),
					linkEvents(
						generateEventWithData(t, "210", true, 0, []byte(`{ "transaction": 3 }`)),
						generateEventWithData(t, "210", false, 0, []byte(`{ "transaction": 30.1 }`)),
					),
				},
			},
			res: res{
				errCount: 2,
				eventsRes: eventsRes{
					aggIDs:            []string{"210"},
					pushedEventsCount: 2,
					aggTypes:          []repository.AggregateType{repository.AggregateType(t.Name())},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{
				client: testCRDBClient,
			}
			wg := sync.WaitGroup{}

			errs := make([]error, 0, tt.res.errCount)
			errsMu := sync.Mutex{}
			for _, events := range tt.args.events {
				wg.Add(1)
				go func(events []*repository.Event) {
					err := db.Push(context.Background(), events...)
					if err != nil {
						errsMu.Lock()
						errs = append(errs, err)
						errsMu.Unlock()
					}

					wg.Done()
				}(events)
			}
			wg.Wait()

			if len(errs) != tt.res.errCount {
				t.Errorf("CRDB.Push() error count = %d, wanted err count %d, errs: %v", len(errs), tt.res.errCount, errs)
			}

			rows, err := testCRDBClient.Query("SELECT event_data FROM eventstore.events where aggregate_type = ANY($1) AND aggregate_id = ANY($2) order by event_sequence", pq.Array(tt.res.eventsRes.aggTypes), pq.Array(tt.res.eventsRes.aggIDs))
			if err != nil {
				t.Error("unable to query inserted rows: ", err)
				return
			}
			var count int

			for rows.Next() {
				count++
				data := make(Data, 0)

				err := rows.Scan(&data)
				if err != nil {
					t.Error("unable to query inserted rows: ", err)
					return
				}
				t.Logf("inserted data: %v", string(data))
			}
			if count != tt.res.eventsRes.pushedEventsCount {
				t.Errorf("expected push count %d got %d", tt.res.eventsRes.pushedEventsCount, count)
			}
		})
	}
}

func TestCRDB_query_events(t *testing.T) {
	type args struct {
		searchQuery *repository.SearchQuery
	}
	type fields struct {
		existingEvents []*repository.Event
	}
	type res struct {
		events []*repository.Event
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
				searchQuery: &repository.SearchQuery{},
			},
			fields: fields{
				existingEvents: []*repository.Event{},
			},
			res: res{
				events: []*repository.Event{},
			},
			wantErr: false,
		},
		// {
		// 	name: "aggregate type filter events found",
		// 	args: args{
		// 		searchQuery: &repository.SearchQuery{},
		// 	},
		// 	fields: fields{
		// 		existingEvents: []*repository.Event{},
		// 	},
		// 	res: res{
		// 		events: []*repository.Event{},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "aggregate type and id filter events found",
		// 	args: args{
		// 		searchQuery: &repository.SearchQuery{},
		// 	},
		// 	fields: fields{
		// 		existingEvents: []*repository.Event{},
		// 	},
		// 	res: res{
		// 		events: []*repository.Event{},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "sequence filter events found",
		// 	args: args{
		// 		searchQuery: &repository.SearchQuery{},
		// 	},
		// 	fields: fields{
		// 		existingEvents: []*repository.Event{},
		// 	},
		// 	res: res{
		// 		events: []*repository.Event{},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "resource owner filter events found",
		// 	args: args{
		// 		searchQuery: &repository.SearchQuery{},
		// 	},
		// 	fields: fields{
		// 		existingEvents: []*repository.Event{},
		// 	},
		// 	res: res{
		// 		events: []*repository.Event{},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "editor service filter events found",
		// 	args: args{
		// 		searchQuery: &repository.SearchQuery{},
		// 	},
		// 	fields: fields{
		// 		existingEvents: []*repository.Event{},
		// 	},
		// 	res: res{
		// 		events: []*repository.Event{},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "editor user filter events found",
		// 	args: args{
		// 		searchQuery: &repository.SearchQuery{},
		// 	},
		// 	fields: fields{
		// 		existingEvents: []*repository.Event{},
		// 	},
		// 	res: res{
		// 		events: []*repository.Event{},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "event type filter events found",
		// 	args: args{
		// 		searchQuery: &repository.SearchQuery{},
		// 	},
		// 	fields: fields{
		// 		existingEvents: []*repository.Event{},
		// 	},
		// 	res: res{
		// 		events: []*repository.Event{},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "no filter events found",
		// 	args: args{
		// 		searchQuery: &repository.SearchQuery{},
		// 	},
		// 	fields: fields{
		// 		existingEvents: []*repository.Event{},
		// 	},
		// 	res: res{
		// 		events: []*repository.Event{},
		// 	},
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{
				client: testCRDBClient,
			}

			// setup initial data for query
			if err := db.Push(context.Background(), tt.fields.existingEvents...); err != nil {
				t.Errorf("error in setup = %v", err)
				return
			}

			events := []*repository.Event{}
			if err := db.query(tt.args.searchQuery, &events); (err != nil) != tt.wantErr {
				t.Errorf("CRDB.query() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func canceledCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func combineEventLists(lists ...[]*repository.Event) []*repository.Event {
	combined := make([]*repository.Event, 0)
	for _, list := range lists {
		combined = append(combined, list...)
	}
	return combined
}

func linkEvents(events ...*repository.Event) []*repository.Event {
	for i := 1; i < len(events); i++ {
		events[i].PreviousEvent = events[i-1]
	}
	return events
}

func generateEventForAggregate(aggregateType repository.AggregateType, aggregateID string, checkPrevious bool, previousSeq uint64) *repository.Event {
	return &repository.Event{
		AggregateID:           aggregateID,
		AggregateType:         aggregateType,
		CheckPreviousSequence: checkPrevious,
		EditorService:         "svc",
		EditorUser:            "user",
		PreviousEvent:         nil,
		PreviousSequence:      previousSeq,
		ResourceOwner:         "ro",
		Type:                  "test.created",
		Version:               "v1",
	}
}

func generateEvent(t *testing.T, aggregateID string, checkPrevious bool, previousSeq uint64) *repository.Event {
	t.Helper()
	return &repository.Event{
		AggregateID:           aggregateID,
		AggregateType:         repository.AggregateType(t.Name()),
		CheckPreviousSequence: checkPrevious,
		EditorService:         "svc",
		EditorUser:            "user",
		PreviousEvent:         nil,
		PreviousSequence:      previousSeq,
		ResourceOwner:         "ro",
		Type:                  "test.created",
		Version:               "v1",
	}
}

func generateEventWithData(t *testing.T, aggregateID string, checkPrevious bool, previousSeq uint64, data []byte) *repository.Event {
	t.Helper()
	return &repository.Event{
		AggregateID:           aggregateID,
		AggregateType:         repository.AggregateType(t.Name()),
		CheckPreviousSequence: checkPrevious,
		EditorService:         "svc",
		EditorUser:            "user",
		PreviousEvent:         nil,
		PreviousSequence:      previousSeq,
		ResourceOwner:         "ro",
		Type:                  "test.created",
		Version:               "v1",
		Data:                  data,
	}
}
