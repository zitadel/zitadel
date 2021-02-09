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
		ctx               context.Context
		events            []*repository.Event
		uniqueConstraints *repository.UniqueConstraint
		uniqueDataType    string
		uniqueDataField   string
	}
	type eventsRes struct {
		pushedEventsCount int
		uniqueCount       int
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
			name: "push 1 event",
			args: args{
				ctx: context.Background(),
				events: []*repository.Event{
					generateEvent(t, "1"),
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
			name: "push two events on agg",
			args: args{
				ctx: context.Background(),
				events: []*repository.Event{
					generateEvent(t, "6"),
					generateEvent(t, "6"),
				},
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
			name: "failed push because context canceled",
			args: args{
				ctx: canceledCtx(),
				events: []*repository.Event{
					generateEvent(t, "9"),
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
		{
			name: "push 1 event and add unique constraint",
			args: args{
				ctx: context.Background(),
				events: []*repository.Event{
					generateEvent(t, "10"),
				},
				uniqueConstraints: generateAddUniqueConstraint(t, "usernames", "field"),
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 1,
					uniqueCount:       1,
					aggID:             []string{"10"},
					aggType:           repository.AggregateType(t.Name()),
				}},
		},
		{
			name: "push 1 event and remove unique constraint",
			args: args{
				ctx: context.Background(),
				events: []*repository.Event{
					generateEvent(t, "11"),
				},
				uniqueConstraints: generateRemoveUniqueConstraint(t, "usernames", "testremove"),
				uniqueDataType:    "usernames",
				uniqueDataField:   "testremove",
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 1,
					uniqueCount:       0,
					aggID:             []string{"11"},
					aggType:           repository.AggregateType(t.Name()),
				}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{
				client: testCRDBClient,
			}
			if tt.args.uniqueDataType != "" && tt.args.uniqueDataField != "" {
				err := fillUniqueData(tt.args.uniqueDataType, tt.args.uniqueDataField)
				if err != nil {
					t.Error("unable to prefill insert unique data: ", err)
					return
				}
			}
			if err := db.Push(tt.args.ctx, tt.args.events, tt.args.uniqueConstraints); (err != nil) != tt.res.wantErr {
				t.Errorf("CRDB.Push() error = %v, wantErr %v", err, tt.res.wantErr)
			}

			countEventRow := testCRDBClient.QueryRow("SELECT COUNT(*) FROM eventstore.events where aggregate_type = $1 AND aggregate_id = ANY($2)", tt.res.eventsRes.aggType, pq.Array(tt.res.eventsRes.aggID))
			var eventCount int
			err := countEventRow.Scan(&eventCount)
			if err != nil {
				t.Error("unable to query inserted rows: ", err)
				return
			}
			if eventCount != tt.res.eventsRes.pushedEventsCount {
				t.Errorf("expected push count %d got %d", tt.res.eventsRes.pushedEventsCount, eventCount)
			}
			if tt.args.uniqueConstraints != nil {
				countUniqueRow := testCRDBClient.QueryRow("SELECT COUNT(*) FROM eventstore.unique_constraints where unique_type = $1 AND unique_field = $2", tt.args.uniqueConstraints.UniqueType, tt.args.uniqueConstraints.UniqueField)
				var uniqueCount int
				err := countUniqueRow.Scan(&uniqueCount)
				if err != nil {
					t.Error("unable to query inserted rows: ", err)
					return
				}
				if uniqueCount != tt.res.eventsRes.uniqueCount {
					t.Errorf("expected unique count %d got %d", tt.res.eventsRes.uniqueCount, uniqueCount)
				}
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
			name: "push two aggregates",
			args: args{
				events: []*repository.Event{
					generateEvent(t, "100"),
					generateEvent(t, "101"),
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
			name: "push two aggregates both multiple events",
			args: args{
				events: []*repository.Event{
					generateEvent(t, "102"),
					generateEvent(t, "102"),
					generateEvent(t, "103"),
					generateEvent(t, "103"),
				},
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
			name: "push two aggregates mixed multiple events",
			args: args{
				events: []*repository.Event{
					generateEvent(t, "106"),
					generateEvent(t, "106"),
					generateEvent(t, "106"),
					generateEvent(t, "106"),
					generateEvent(t, "107"),
					generateEvent(t, "107"),
					generateEvent(t, "107"),
					generateEvent(t, "107"),
					generateEvent(t, "108"),
					generateEvent(t, "108"),
					generateEvent(t, "108"),
					generateEvent(t, "108"),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 12,
					aggID:             []string{"106", "107", "108"},
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
			if err := db.Push(context.Background(), tt.args.events); (err != nil) != tt.res.wantErr {
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
					{
						generateEvent(t, "200"),
						generateEvent(t, "200"),
						generateEvent(t, "200"),
						generateEvent(t, "201"),
						generateEvent(t, "201"),
						generateEvent(t, "201"),
					},
					{
						generateEvent(t, "202"),
						generateEvent(t, "203"),
						generateEvent(t, "203"),
					},
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
			name: "clients push same aggregates",
			args: args{
				events: [][]*repository.Event{
					{
						generateEvent(t, "204"),
						generateEvent(t, "204"),
					},
					{
						generateEvent(t, "204"),
						generateEvent(t, "204"),
					},
					{
						generateEvent(t, "205"),
						generateEvent(t, "205"),
						generateEvent(t, "205"),
						generateEvent(t, "206"),
						generateEvent(t, "206"),
						generateEvent(t, "206"),
					},
					{
						generateEvent(t, "204"),
						generateEvent(t, "205"),
						generateEvent(t, "205"),
						generateEvent(t, "206"),
					},
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
			name: "clients push different aggregates",
			args: args{
				events: [][]*repository.Event{
					{
						generateEvent(t, "207"),
						generateEvent(t, "207"),
						generateEvent(t, "207"),
						generateEvent(t, "207"),
						generateEvent(t, "207"),
						generateEvent(t, "207"),
					},
					{
						generateEvent(t, "208"),
						generateEvent(t, "208"),
						generateEvent(t, "208"),
						generateEvent(t, "208"),
						generateEvent(t, "208"),
					},
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
					err := db.Push(context.Background(), events)
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

func TestCRDB_Filter(t *testing.T) {
	type args struct {
		searchQuery *repository.SearchQuery
	}
	type fields struct {
		existingEvents []*repository.Event
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
				searchQuery: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, "not found", repository.OperationEquals),
					},
				},
			},
			fields: fields{
				existingEvents: []*repository.Event{
					generateEvent(t, "300"),
					generateEvent(t, "300"),
					generateEvent(t, "300"),
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
				searchQuery: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, t.Name(), repository.OperationEquals),
						repository.NewFilter(repository.FieldAggregateID, "303", repository.OperationEquals),
					},
				},
			},
			fields: fields{
				existingEvents: []*repository.Event{
					generateEvent(t, "303"),
					generateEvent(t, "303"),
					generateEvent(t, "303"),
					generateEvent(t, "305"),
				},
			},
			res: res{
				eventCount: 3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{
				client: testCRDBClient,
			}

			// setup initial data for query
			if err := db.Push(context.Background(), tt.fields.existingEvents); err != nil {
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

func TestCRDB_LatestSequence(t *testing.T) {
	type args struct {
		searchQuery *repository.SearchQuery
	}
	type fields struct {
		existingEvents []*repository.Event
	}
	type res struct {
		sequence uint64
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
				searchQuery: &repository.SearchQuery{
					Columns: repository.ColumnsMaxSequence,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, "not found", repository.OperationEquals),
					},
				},
			},
			fields: fields{
				existingEvents: []*repository.Event{
					generateEvent(t, "400"),
					generateEvent(t, "400"),
					generateEvent(t, "400"),
				},
			},
			res: res{
				sequence: 0,
			},
			wantErr: false,
		},
		{
			name: "aggregate type filter sequence",
			args: args{
				searchQuery: &repository.SearchQuery{
					Columns: repository.ColumnsMaxSequence,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, t.Name(), repository.OperationEquals),
					},
				},
			},
			fields: fields{
				existingEvents: []*repository.Event{
					generateEvent(t, "401"),
					generateEvent(t, "401"),
					generateEvent(t, "401"),
				},
			},
			res: res{
				sequence: 3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{
				client: testCRDBClient,
			}

			// setup initial data for query
			if err := db.Push(context.Background(), tt.fields.existingEvents); err != nil {
				t.Errorf("error in setup = %v", err)
				return
			}

			sequence, err := db.LatestSequence(context.Background(), tt.args.searchQuery)
			if (err != nil) != tt.wantErr {
				t.Errorf("CRDB.query() error = %v, wantErr %v", err, tt.wantErr)
			}

			if sequence < tt.res.sequence {
				t.Errorf("CRDB.query() expected sequence: %d got %d", tt.res.sequence, sequence)
			}
		})
	}
}

func TestCRDB_Push_ResourceOwner(t *testing.T) {
	type args struct {
		events []*repository.Event
	}
	type res struct {
		resourceOwners []string
	}
	type fields struct {
		aggregateIDs  []string
		aggregateType string
	}
	tests := []struct {
		name   string
		args   args
		res    res
		fields fields
	}{
		{
			name: "two events of same aggregate same resource owner",
			args: args{
				events: []*repository.Event{
					generateEvent(t, "500", func(e *repository.Event) { e.ResourceOwner = "caos" }),
					generateEvent(t, "500", func(e *repository.Event) { e.ResourceOwner = "caos" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"500"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "caos"},
			},
		},
		{
			name: "two events of different aggregate same resource owner",
			args: args{
				events: []*repository.Event{
					generateEvent(t, "501", func(e *repository.Event) { e.ResourceOwner = "caos" }),
					generateEvent(t, "502", func(e *repository.Event) { e.ResourceOwner = "caos" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"501", "502"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "caos"},
			},
		},
		{
			name: "two events of different aggregate different resource owner",
			args: args{
				events: []*repository.Event{
					generateEvent(t, "503", func(e *repository.Event) { e.ResourceOwner = "caos" }),
					generateEvent(t, "504", func(e *repository.Event) { e.ResourceOwner = "zitadel" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"503", "504"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "zitadel"},
			},
		},
		{
			name: "events of different aggregate different resource owner",
			args: args{
				events: []*repository.Event{
					generateEvent(t, "505", func(e *repository.Event) { e.ResourceOwner = "caos" }),
					generateEvent(t, "505", func(e *repository.Event) { e.ResourceOwner = "caos" }),
					generateEvent(t, "506", func(e *repository.Event) { e.ResourceOwner = "zitadel" }),
					generateEvent(t, "506", func(e *repository.Event) { e.ResourceOwner = "zitadel" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"505", "506"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "caos", "zitadel", "zitadel"},
			},
		},
		{
			name: "events of different aggregate different resource owner per event",
			args: args{
				events: []*repository.Event{
					generateEvent(t, "507", func(e *repository.Event) { e.ResourceOwner = "caos" }),
					generateEvent(t, "507", func(e *repository.Event) { e.ResourceOwner = "ignored" }),
					generateEvent(t, "508", func(e *repository.Event) { e.ResourceOwner = "zitadel" }),
					generateEvent(t, "508", func(e *repository.Event) { e.ResourceOwner = "ignored" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"507", "508"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "caos", "zitadel", "zitadel"},
			},
		},
		{
			name: "events of one aggregate different resource owner per event",
			args: args{
				events: []*repository.Event{
					generateEvent(t, "509", func(e *repository.Event) { e.ResourceOwner = "caos" }),
					generateEvent(t, "509", func(e *repository.Event) { e.ResourceOwner = "ignored" }),
					generateEvent(t, "509", func(e *repository.Event) { e.ResourceOwner = "ignored" }),
					generateEvent(t, "509", func(e *repository.Event) { e.ResourceOwner = "ignored" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"509"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "caos", "caos", "caos"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{
				client: testCRDBClient,
			}
			if err := db.Push(context.Background(), tt.args.events); err != nil {
				t.Errorf("CRDB.Push() error = %v", err)
			}

			if len(tt.args.events) != len(tt.res.resourceOwners) {
				t.Errorf("length of events (%d) and resource owners (%d) must be equal", len(tt.args.events), len(tt.res.resourceOwners))
				return
			}

			for i, event := range tt.args.events {
				if event.ResourceOwner != tt.res.resourceOwners[i] {
					t.Errorf("resource owner not expected want: %q got: %q", tt.res.resourceOwners[i], event.ResourceOwner)
				}
			}

			rows, err := testCRDBClient.Query("SELECT resource_owner FROM eventstore.events WHERE aggregate_type = $1 AND aggregate_id = ANY($2) ORDER BY event_sequence", tt.fields.aggregateType, pq.Array(tt.fields.aggregateIDs))
			if err != nil {
				t.Error("unable to query inserted rows: ", err)
				return
			}

			eventCount := 0
			for i := 0; rows.Next(); i++ {
				var resourceOwner string
				err = rows.Scan(&resourceOwner)
				if err != nil {
					t.Error("unable to scan row: ", err)
					return
				}
				if resourceOwner != tt.res.resourceOwners[i] {
					t.Errorf("unexpected resource owner in queried event. want %q, got: %q", tt.res.resourceOwners[i], resourceOwner)
				}
				eventCount++
			}

			if eventCount != len(tt.res.resourceOwners) {
				t.Errorf("wrong queried event count: want %d, got %d", len(tt.res.resourceOwners), eventCount)
			}
		})
	}
}

func canceledCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func generateEvent(t *testing.T, aggregateID string, opts ...func(*repository.Event)) *repository.Event {
	t.Helper()
	e := &repository.Event{
		AggregateID:   aggregateID,
		AggregateType: repository.AggregateType(t.Name()),
		EditorService: "svc",
		EditorUser:    "user",
		ResourceOwner: "ro",
		Type:          "test.created",
		Version:       "v1",
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func generateEventWithData(t *testing.T, aggregateID string, data []byte) *repository.Event {
	t.Helper()
	return &repository.Event{
		AggregateID:   aggregateID,
		AggregateType: repository.AggregateType(t.Name()),
		EditorService: "svc",
		EditorUser:    "user",
		ResourceOwner: "ro",
		Type:          "test.created",
		Version:       "v1",
		Data:          data,
	}
}

func generateAddUniqueConstraint(t *testing.T, table, uniqueField string) *repository.UniqueConstraint {
	t.Helper()
	e := &repository.UniqueConstraint{
		UniqueType:  table,
		UniqueField: uniqueField,
		Action:      repository.UniqueConstraintAdd,
	}

	return e
}

func generateRemoveUniqueConstraint(t *testing.T, table, uniqueField string) *repository.UniqueConstraint {
	t.Helper()
	e := &repository.UniqueConstraint{
		UniqueType:  table,
		UniqueField: uniqueField,
		Action:      repository.UniqueConstraintRemoved,
	}

	return e
}
