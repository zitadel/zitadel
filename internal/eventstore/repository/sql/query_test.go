package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/cockroach"
	db_mock "github.com/zitadel/zitadel/internal/database/mock"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func Test_getCondition(t *testing.T) {
	type args struct {
		filter *repository.Filter
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "equals",
			args: args{filter: repository.NewFilter(repository.FieldAggregateID, "", repository.OperationEquals)},
			want: "aggregate_id = ?",
		},
		{
			name: "greater",
			args: args{filter: repository.NewFilter(repository.FieldSequence, 0, repository.OperationGreater)},
			want: `"sequence" > ?`,
		},
		{
			name: "less",
			args: args{filter: repository.NewFilter(repository.FieldSequence, 5000, repository.OperationLess)},
			want: `"sequence" < ?`,
		},
		{
			name: "in list",
			args: args{filter: repository.NewFilter(repository.FieldAggregateType, []eventstore.AggregateType{"movies", "actors"}, repository.OperationIn)},
			want: "aggregate_type = ANY(?)",
		},
		{
			name: "invalid operation",
			args: args{filter: repository.NewFilter(repository.FieldAggregateType, []eventstore.AggregateType{"movies", "actors"}, repository.Operation(-1))},
			want: "",
		},
		{
			name: "invalid field",
			args: args{filter: repository.NewFilter(repository.Field(-1), []eventstore.AggregateType{"movies", "actors"}, repository.OperationEquals)},
			want: "",
		},
		{
			name: "invalid field and operation",
			args: args{filter: repository.NewFilter(repository.Field(-1), []eventstore.AggregateType{"movies", "actors"}, repository.Operation(-1))},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{}
			if got := getCondition(db, tt.args.filter, false); got != tt.want {
				t.Errorf("getCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prepareColumns(t *testing.T) {
	var reducedEvents []eventstore.Event

	type fields struct {
		dbRow []interface{}
	}
	type args struct {
		columns eventstore.Columns
		dest    interface{}
		dbErr   error
		useV1   bool
	}
	type res struct {
		query    string
		expected interface{}
		dbErr    func(error) bool
	}
	tests := []struct {
		name   string
		args   args
		res    res
		fields fields
	}{
		{
			name: "invalid columns",
			args: args{columns: eventstore.Columns(-1)},
			res: res{
				query: "",
				dbErr: func(err error) bool { return err == nil },
			},
		},
		{
			name: "max column",
			args: args{
				columns: eventstore.ColumnsMaxSequence,
				dest:    new(sql.NullFloat64),
				useV1:   true,
			},
			res: res{
				query:    `SELECT event_sequence FROM eventstore.events`,
				expected: sql.NullFloat64{Float64: 43, Valid: true},
			},
			fields: fields{
				dbRow: []interface{}{sql.NullFloat64{Float64: 43, Valid: true}},
			},
		},
		{
			name: "max column v2",
			args: args{
				columns: eventstore.ColumnsMaxSequence,
				dest:    new(sql.NullFloat64),
			},
			res: res{
				query:    `SELECT "position" FROM eventstore.events2`,
				expected: sql.NullFloat64{Float64: 43, Valid: true},
			},
			fields: fields{
				dbRow: []interface{}{sql.NullFloat64{Float64: 43, Valid: true}},
			},
		},
		{
			name: "max sequence wrong dest type",
			args: args{
				columns: eventstore.ColumnsMaxSequence,
				dest:    new(uint64),
			},
			res: res{
				query: `SELECT "position" FROM eventstore.events2`,
				dbErr: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "events",
			args: args{
				columns: eventstore.ColumnsEvent,
				dest: eventstore.Reducer(func(event eventstore.Event) error {
					reducedEvents = append(reducedEvents, event)
					return nil
				}),
				useV1: true,
			},
			res: res{
				query: `SELECT creation_date, event_type, event_sequence, event_data, editor_user, resource_owner, instance_id, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events`,
				expected: []eventstore.Event{
					&repository.Event{AggregateID: "hodor", AggregateType: "user", Seq: 5, Data: nil},
				},
			},
			fields: fields{
				dbRow: []interface{}{time.Time{}, eventstore.EventType(""), uint64(5), sql.RawBytes(nil), "", sql.NullString{}, "", eventstore.AggregateType("user"), "hodor", eventstore.Version("")},
			},
		},
		{
			name: "events v2",
			args: args{
				columns: eventstore.ColumnsEvent,
				dest: eventstore.Reducer(func(event eventstore.Event) error {
					reducedEvents = append(reducedEvents, event)
					return nil
				}),
			},
			res: res{
				query: `SELECT created_at, event_type, "sequence", "position", payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2`,
				expected: []eventstore.Event{
					&repository.Event{AggregateID: "hodor", AggregateType: "user", Seq: 5, Pos: 42, Data: nil, Version: "v1"},
				},
			},
			fields: fields{
				dbRow: []interface{}{time.Time{}, eventstore.EventType(""), uint64(5), sql.NullFloat64{Float64: 42, Valid: true}, sql.RawBytes(nil), "", sql.NullString{}, "", eventstore.AggregateType("user"), "hodor", uint8(1)},
			},
		},
		{
			name: "event null position",
			args: args{
				columns: eventstore.ColumnsEvent,
				dest: eventstore.Reducer(func(event eventstore.Event) error {
					reducedEvents = append(reducedEvents, event)
					return nil
				}),
			},
			res: res{
				query: `SELECT created_at, event_type, "sequence", "position", payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2`,
				expected: []eventstore.Event{
					&repository.Event{AggregateID: "hodor", AggregateType: "user", Seq: 5, Pos: 0, Data: nil, Version: "v1"},
				},
			},
			fields: fields{
				dbRow: []interface{}{time.Time{}, eventstore.EventType(""), uint64(5), sql.NullFloat64{Float64: 0, Valid: false}, sql.RawBytes(nil), "", sql.NullString{}, "", eventstore.AggregateType("user"), "hodor", uint8(1)},
			},
		},
		{
			name: "events wrong dest type",
			args: args{
				columns: eventstore.ColumnsEvent,
				dest:    []*repository.Event{},
				useV1:   true,
			},
			res: res{
				query: `SELECT creation_date, event_type, event_sequence, event_data, editor_user, resource_owner, instance_id, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events`,
				dbErr: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "event query error",
			args: args{
				columns: eventstore.ColumnsEvent,
				dest: eventstore.Reducer(func(event eventstore.Event) error {
					reducedEvents = append(reducedEvents, event)
					return nil
				}),
				dbErr: sql.ErrConnDone,
				useV1: true,
			},
			res: res{
				query: `SELECT creation_date, event_type, event_sequence, event_data, editor_user, resource_owner, instance_id, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events`,
				dbErr: zerrors.IsInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crdb := &CRDB{}
			query, rowScanner := prepareColumns(crdb, tt.args.columns, tt.args.useV1)
			if query != tt.res.query {
				t.Errorf("prepareColumns() got = %s, want %s", query, tt.res.query)
			}
			if tt.res.query == "" && rowScanner != nil {
				t.Errorf("row scanner should be nil")
			}
			if rowScanner == nil {
				return
			}
			err := rowScanner(prepareTestScan(tt.args.dbErr, tt.fields.dbRow), tt.args.dest)
			if err != nil && tt.res.dbErr == nil || err != nil && !tt.res.dbErr(err) || err == nil && tt.res.dbErr != nil {
				t.Errorf("wrong error type in rowScanner got: %v", err)
				return
			}
			if tt.res.dbErr != nil && tt.res.dbErr(err) {
				return
			}
			if equalizer, ok := tt.res.expected.(interface{ Equal(time.Time) bool }); ok {
				equalizer.Equal(tt.args.dest.(*sql.NullTime).Time)
				return
			}
			if _, ok := tt.args.dest.(eventstore.Reducer); ok {
				assert.Equal(t, tt.res.expected, reducedEvents)
				reducedEvents = nil
				return
			}

			got := reflect.Indirect(reflect.ValueOf(tt.args.dest)).Interface()
			if !reflect.DeepEqual(got, tt.res.expected) {
				t.Errorf("unexpected result from rowScanner \nwant: %+v \ngot: %+v", tt.res.expected, got)
			}
		})
	}
}

func prepareTestScan(err error, res []interface{}) scan {
	return func(dests ...interface{}) error {
		if err != nil {
			return err
		}
		if len(dests) != len(res) {
			return zerrors.ThrowInvalidArgumentf(nil, "SQL-NML1q", "expected len %d got %d", len(res), len(dests))
		}
		for i, r := range res {
			_, ok := dests[i].(*eventstore.Version)
			if ok {
				val, ok := r.(uint8)
				if ok {
					r = eventstore.Version("" + strconv.Itoa(int(val)))
				}
			}
			reflect.ValueOf(dests[i]).Elem().Set(reflect.ValueOf(r))
		}

		return nil
	}
}

func Test_prepareCondition(t *testing.T) {
	type args struct {
		query *repository.SearchQuery
		useV1 bool
	}
	type res struct {
		clause string
		values []interface{}
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "nil filters",
			args: args{
				query: &repository.SearchQuery{},
				useV1: true,
			},
			res: res{
				clause: "",
				values: nil,
			},
		},
		{
			name: "nil filters v2",
			args: args{
				query: &repository.SearchQuery{},
			},
			res: res{
				clause: "",
				values: nil,
			},
		},
		{
			name: "empty filters",
			args: args{
				query: &repository.SearchQuery{
					SubQueries: [][]*repository.Filter{},
				},
				useV1: true,
			},
			res: res{
				clause: "",
				values: nil,
			},
		},
		{
			name: "empty filters v2",
			args: args{
				query: &repository.SearchQuery{
					SubQueries: [][]*repository.Filter{},
				},
			},
			res: res{
				clause: "",
				values: nil,
			},
		},
		{
			name: "invalid condition",
			args: args{
				query: &repository.SearchQuery{
					SubQueries: [][]*repository.Filter{
						{
							repository.NewFilter(repository.FieldAggregateID, "wrong", repository.Operation(-1)),
						},
					},
				},
				useV1: true,
			},
			res: res{
				clause: "",
				values: nil,
			},
		},
		{
			name: "invalid condition v2",
			args: args{
				query: &repository.SearchQuery{
					SubQueries: [][]*repository.Filter{
						{
							repository.NewFilter(repository.FieldAggregateID, "wrong", repository.Operation(-1)),
						},
					},
				},
			},
			res: res{
				clause: "",
				values: nil,
			},
		},
		{
			name: "array as condition value",
			args: args{
				query: &repository.SearchQuery{
					AwaitOpenTransactions: true,
					SubQueries: [][]*repository.Filter{
						{
							repository.NewFilter(repository.FieldAggregateType, []eventstore.AggregateType{"user", "org"}, repository.OperationIn),
						},
					},
				},
				useV1: true,
			},
			res: res{
				clause: " WHERE aggregate_type = ANY(?) AND creation_date::TIMESTAMP < (SELECT COALESCE(MIN(start), NOW())::TIMESTAMP FROM crdb_internal.cluster_transactions where application_name = 'zitadel_es_pusher')",
				values: []interface{}{[]eventstore.AggregateType{"user", "org"}},
			},
		},
		{
			name: "array as condition value v2",
			args: args{
				query: &repository.SearchQuery{
					AwaitOpenTransactions: true,
					SubQueries: [][]*repository.Filter{
						{
							repository.NewFilter(repository.FieldAggregateType, []eventstore.AggregateType{"user", "org"}, repository.OperationIn),
						},
					},
				},
			},
			res: res{
				clause: ` WHERE aggregate_type = ANY(?) AND hlc_to_timestamp("position") < (SELECT COALESCE(MIN(start), NOW())::TIMESTAMP FROM crdb_internal.cluster_transactions where application_name = 'zitadel_es_pusher')`,
				values: []interface{}{[]eventstore.AggregateType{"user", "org"}},
			},
		},
		{
			name: "multiple filters",
			args: args{
				query: &repository.SearchQuery{
					AwaitOpenTransactions: true,
					SubQueries: [][]*repository.Filter{
						{
							repository.NewFilter(repository.FieldAggregateType, []eventstore.AggregateType{"user", "org"}, repository.OperationIn),
							repository.NewFilter(repository.FieldAggregateID, "1234", repository.OperationEquals),
							repository.NewFilter(repository.FieldEventType, []eventstore.EventType{"user.created", "org.created"}, repository.OperationIn),
						},
					},
				},
				useV1: true,
			},
			res: res{
				clause: " WHERE aggregate_type = ANY(?) AND aggregate_id = ? AND event_type = ANY(?) AND creation_date::TIMESTAMP < (SELECT COALESCE(MIN(start), NOW())::TIMESTAMP FROM crdb_internal.cluster_transactions where application_name = 'zitadel_es_pusher')",
				values: []interface{}{[]eventstore.AggregateType{"user", "org"}, "1234", []eventstore.EventType{"user.created", "org.created"}},
			},
		},
		{
			name: "multiple filters v2",
			args: args{
				query: &repository.SearchQuery{
					AwaitOpenTransactions: true,
					SubQueries: [][]*repository.Filter{
						{
							repository.NewFilter(repository.FieldAggregateType, []eventstore.AggregateType{"user", "org"}, repository.OperationIn),
							repository.NewFilter(repository.FieldAggregateID, "1234", repository.OperationEquals),
							repository.NewFilter(repository.FieldEventType, []eventstore.EventType{"user.created", "org.created"}, repository.OperationIn),
						},
					},
				},
			},
			res: res{
				clause: ` WHERE aggregate_type = ANY(?) AND aggregate_id = ? AND event_type = ANY(?) AND hlc_to_timestamp("position") < (SELECT COALESCE(MIN(start), NOW())::TIMESTAMP FROM crdb_internal.cluster_transactions where application_name = 'zitadel_es_pusher')`,
				values: []interface{}{[]eventstore.AggregateType{"user", "org"}, "1234", []eventstore.EventType{"user.created", "org.created"}},
			},
		},
	}
	crdb := NewCRDB(&database.DB{Database: new(cockroach.Config)})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClause, gotValues := prepareConditions(crdb, tt.args.query, tt.args.useV1)
			if gotClause != tt.res.clause {
				t.Errorf("prepareCondition() gotClause = %v, want %v", gotClause, tt.res.clause)
			}
			if len(gotValues) != len(tt.res.values) {
				t.Errorf("wrong length of gotten values got = %d, want %d", len(gotValues), len(tt.res.values))
				return
			}
			for i, value := range gotValues {
				if !reflect.DeepEqual(value, tt.res.values[i]) {
					t.Errorf("prepareCondition() gotValues = %v, want %v", gotValues, tt.res.values)
				}
			}
		})
	}
}

func Test_query_events_with_crdb(t *testing.T) {
	type args struct {
		searchQuery *eventstore.SearchQueryBuilder
	}
	type fields struct {
		existingEvents []eventstore.Command
		client         *sql.DB
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
				client: testCRDBClient,
				existingEvents: []eventstore.Command{
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
			name: "aggregate type filter events found",
			args: args{
				searchQuery: eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
					AddQuery().
					AggregateTypes(eventstore.AggregateType(t.Name())).
					Builder(),
			},
			fields: fields{
				client: testCRDBClient,
				existingEvents: []eventstore.Command{
					generateEvent(t, "301"),
					generateEvent(t, "302"),
					generateEvent(t, "302"),
					generateEvent(t, "303", func(e *repository.Event) { e.AggregateType = "not in list" }),
				},
			},
			res: res{
				eventCount: 3,
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
				client: testCRDBClient,
				existingEvents: []eventstore.Command{
					generateEvent(t, "303"),
					generateEvent(t, "303"),
					generateEvent(t, "303"),
					generateEvent(t, "304", func(e *repository.Event) { e.AggregateType = "not in list" }),
					generateEvent(t, "305"),
				},
			},
			res: res{
				eventCount: 3,
			},
			wantErr: false,
		},
		{
			name: "resource owner filter events found",
			args: args{
				searchQuery: eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
					ResourceOwner("caos"),
			},
			fields: fields{
				client: testCRDBClient,
				existingEvents: []eventstore.Command{
					generateEvent(t, "306", func(e *repository.Event) { e.ResourceOwner = sql.NullString{String: "caos", Valid: true} }),
					generateEvent(t, "307", func(e *repository.Event) { e.ResourceOwner = sql.NullString{String: "caos", Valid: true} }),
					generateEvent(t, "308", func(e *repository.Event) { e.ResourceOwner = sql.NullString{String: "caos", Valid: true} }),
					generateEvent(t, "309", func(e *repository.Event) { e.ResourceOwner = sql.NullString{String: "orgID", Valid: true} }),
					generateEvent(t, "309", func(e *repository.Event) { e.ResourceOwner = sql.NullString{String: "orgID", Valid: true} }),
				},
			},
			res: res{
				eventCount: 3,
			},
			wantErr: false,
		},
		{
			name: "event type filter events found",
			args: args{
				searchQuery: eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
					AddQuery().
					EventTypes("user.created", "user.updated").
					Builder(),
			},
			fields: fields{
				client: testCRDBClient,
				existingEvents: []eventstore.Command{
					generateEvent(t, "311", func(e *repository.Event) { e.Typ = "user.created" }),
					generateEvent(t, "311", func(e *repository.Event) { e.Typ = "user.updated" }),
					generateEvent(t, "311", func(e *repository.Event) { e.Typ = "user.deactivated" }),
					generateEvent(t, "311", func(e *repository.Event) { e.Typ = "user.locked" }),
					generateEvent(t, "312", func(e *repository.Event) { e.Typ = "user.created" }),
					generateEvent(t, "312", func(e *repository.Event) { e.Typ = "user.updated" }),
					generateEvent(t, "312", func(e *repository.Event) { e.Typ = "user.deactivated" }),
					generateEvent(t, "312", func(e *repository.Event) { e.Typ = "user.reactivated" }),
					generateEvent(t, "313", func(e *repository.Event) { e.Typ = "user.locked" }),
				},
			},
			res: res{
				eventCount: 7,
			},
			wantErr: false,
		},
		{
			name: "fail because no filter",
			args: args{
				searchQuery: eventstore.NewSearchQueryBuilder(eventstore.Columns(-1)),
			},
			fields: fields{
				client:         testCRDBClient,
				existingEvents: []eventstore.Command{},
			},
			res: res{
				eventCount: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{
				DB: &database.DB{
					DB:       tt.fields.client,
					Database: new(testDB),
				},
			}

			// setup initial data for query
			if _, err := db.Push(context.Background(), tt.fields.existingEvents...); err != nil {
				t.Errorf("error in setup = %v", err)
				return
			}

			events := []eventstore.Event{}
			if err := query(context.Background(), db, tt.args.searchQuery, eventstore.Reducer(func(event eventstore.Event) error {
				events = append(events, event)
				return nil
			}), true); (err != nil) != tt.wantErr {
				t.Errorf("CRDB.query() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_query_events_mocked(t *testing.T) {
	type args struct {
		query *eventstore.SearchQueryBuilder
		dest  interface{}
	}
	type res struct {
		wantErr bool
	}
	type fields struct {
		mock *dbMock
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		res    res
	}{
		{
			name: "with order by desc",
			args: args{
				dest: &[]*repository.Event{},
				query: eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
					OrderDesc().
					AwaitOpenTransactions().
					AddQuery().
					AggregateTypes("user").
					Builder(),
			},
			fields: fields{
				mock: newMockClient(t).expectQuery(t,
					`SELECT creation_date, event_type, event_sequence, event_data, editor_user, resource_owner, instance_id, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE aggregate_type = \$1 AND creation_date::TIMESTAMP < \(SELECT COALESCE\(MIN\(start\), NOW\(\)\)::TIMESTAMP FROM crdb_internal\.cluster_transactions where application_name = 'zitadel_es_pusher'\) ORDER BY event_sequence DESC`,
					[]driver.Value{eventstore.AggregateType("user")},
				),
			},
			res: res{
				wantErr: false,
			},
		},
		{
			name: "with limit",
			args: args{
				dest: &[]*repository.Event{},
				query: eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
					OrderAsc().
					AwaitOpenTransactions().
					Limit(5).
					AddQuery().
					AggregateTypes("user").
					Builder(),
			},
			fields: fields{
				mock: newMockClient(t).expectQuery(t,
					`SELECT creation_date, event_type, event_sequence, event_data, editor_user, resource_owner, instance_id, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE aggregate_type = \$1 AND creation_date::TIMESTAMP < \(SELECT COALESCE\(MIN\(start\), NOW\(\)\)::TIMESTAMP FROM crdb_internal\.cluster_transactions where application_name = 'zitadel_es_pusher'\) ORDER BY event_sequence LIMIT \$2`,
					[]driver.Value{eventstore.AggregateType("user"), uint64(5)},
				),
			},
			res: res{
				wantErr: false,
			},
		},
		{
			name: "with limit and order by desc",
			args: args{
				dest: &[]*repository.Event{},
				query: eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
					OrderDesc().
					AwaitOpenTransactions().
					Limit(5).
					AddQuery().
					AggregateTypes("user").
					Builder(),
			},
			fields: fields{
				mock: newMockClient(t).expectQuery(t,
					`SELECT creation_date, event_type, event_sequence, event_data, editor_user, resource_owner, instance_id, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE aggregate_type = \$1 AND creation_date::TIMESTAMP < \(SELECT COALESCE\(MIN\(start\), NOW\(\)\)::TIMESTAMP FROM crdb_internal\.cluster_transactions where application_name = 'zitadel_es_pusher'\) ORDER BY event_sequence DESC LIMIT \$2`,
					[]driver.Value{eventstore.AggregateType("user"), uint64(5)},
				),
			},
			res: res{
				wantErr: false,
			},
		},
		{
			name: "with limit and order by desc as of system time",
			args: args{
				dest: &[]*repository.Event{},
				query: eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
					OrderDesc().
					AwaitOpenTransactions().
					Limit(5).
					AllowTimeTravel().
					AddQuery().
					AggregateTypes("user").
					Builder(),
			},
			fields: fields{
				mock: newMockClient(t).expectQuery(t,
					`SELECT creation_date, event_type, event_sequence, event_data, editor_user, resource_owner, instance_id, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events AS OF SYSTEM TIME '-1 ms' WHERE aggregate_type = \$1 AND creation_date::TIMESTAMP < \(SELECT COALESCE\(MIN\(start\), NOW\(\)\)::TIMESTAMP FROM crdb_internal\.cluster_transactions where application_name = 'zitadel_es_pusher'\) ORDER BY event_sequence DESC LIMIT \$2`,
					[]driver.Value{eventstore.AggregateType("user"), uint64(5)},
				),
			},
			res: res{
				wantErr: false,
			},
		},
		{
			name: "error sql conn closed",
			args: args{
				dest: &[]*repository.Event{},
				query: eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
					OrderDesc().
					AwaitOpenTransactions().
					Limit(0).
					AddQuery().
					AggregateTypes("user").
					Builder(),
			},
			fields: fields{
				mock: newMockClient(t).expectQueryErr(t,
					`SELECT creation_date, event_type, event_sequence, event_data, editor_user, resource_owner, instance_id, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE aggregate_type = \$1 AND creation_date::TIMESTAMP < \(SELECT COALESCE\(MIN\(start\), NOW\(\)\)::TIMESTAMP FROM crdb_internal\.cluster_transactions where application_name = 'zitadel_es_pusher'\) ORDER BY event_sequence DESC`,
					[]driver.Value{eventstore.AggregateType("user")},
					sql.ErrConnDone),
			},
			res: res{
				wantErr: true,
			},
		},
		{
			name: "error unexpected dest",
			args: args{
				dest: nil,
				query: eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
					OrderDesc().
					AwaitOpenTransactions().
					Limit(0).
					AddQuery().
					AggregateTypes("user").
					Builder(),
			},
			fields: fields{
				mock: newMockClient(t).expectQueryScanErr(t,
					`SELECT creation_date, event_type, event_sequence, event_data, editor_user, resource_owner, instance_id, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE aggregate_type = \$1 AND creation_date::TIMESTAMP < \(SELECT COALESCE\(MIN\(start\), NOW\(\)\)::TIMESTAMP FROM crdb_internal\.cluster_transactions where application_name = 'zitadel_es_pusher'\) ORDER BY event_sequence DESC`,
					[]driver.Value{eventstore.AggregateType("user")},
					&repository.Event{Seq: 100}),
			},
			res: res{
				wantErr: true,
			},
		},
		{
			name: "error no columns",
			args: args{
				query: eventstore.NewSearchQueryBuilder(eventstore.Columns(-1)),
			},
			res: res{
				wantErr: true,
			},
		},
		{
			name: "with subqueries",
			args: args{
				dest: &[]*repository.Event{},
				query: eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
					OrderDesc().
					AwaitOpenTransactions().
					Limit(5).
					AddQuery().
					AggregateTypes("user").
					Or().
					AggregateTypes("org").
					AggregateIDs("asdf42").
					Builder(),
			},
			fields: fields{
				mock: newMockClient(t).expectQuery(t,
					`SELECT creation_date, event_type, event_sequence, event_data, editor_user, resource_owner, instance_id, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE \(aggregate_type = \$1 OR \(aggregate_type = \$2 AND aggregate_id = \$3\)\) AND creation_date::TIMESTAMP < \(SELECT COALESCE\(MIN\(start\), NOW\(\)\)::TIMESTAMP FROM crdb_internal\.cluster_transactions where application_name = 'zitadel_es_pusher'\) ORDER BY event_sequence DESC LIMIT \$4`,
					[]driver.Value{eventstore.AggregateType("user"), eventstore.AggregateType("org"), "asdf42", uint64(5)},
				),
			},
			res: res{
				wantErr: false,
			},
		},
	}
	crdb := NewCRDB(&database.DB{Database: new(testDB)})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fields.mock != nil {
				crdb.DB.DB = tt.fields.mock.client
			}

			err := query(context.Background(), crdb, tt.args.query, tt.args.dest, true)
			if (err != nil) != tt.res.wantErr {
				t.Errorf("query() error = %v, wantErr %v", err, tt.res.wantErr)
			}

			if tt.fields.mock == nil {
				return
			}

			if err := tt.fields.mock.mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all expectaions met: %v", err)
			}
		})
	}
}

type dbMock struct {
	mock   sqlmock.Sqlmock
	client *sql.DB
}

func (m *dbMock) expectQuery(t *testing.T, expectedQuery string, args []driver.Value, events ...*repository.Event) *dbMock {
	query := m.mock.ExpectQuery(expectedQuery).WithArgs(args...)
	rows := m.mock.NewRows([]string{"sequence"})
	for _, event := range events {
		rows = rows.AddRow(event.Seq)
	}
	query.WillReturnRows(rows).RowsWillBeClosed()
	return m
}

func (m *dbMock) expectQueryScanErr(t *testing.T, expectedQuery string, args []driver.Value, events ...*repository.Event) *dbMock {
	query := m.mock.ExpectQuery(expectedQuery).WithArgs(args...)
	rows := m.mock.NewRows([]string{"sequence"})
	for _, event := range events {
		rows = rows.AddRow(event.Seq)
	}
	query.WillReturnRows(rows).RowsWillBeClosed()
	return m
}

func (m *dbMock) expectQueryErr(t *testing.T, expectedQuery string, args []driver.Value, err error) *dbMock {
	m.mock.ExpectQuery(expectedQuery).WithArgs(args...).WillReturnError(err)
	return m
}

func newMockClient(t *testing.T) *dbMock {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.ValueConverterOption(new(db_mock.TypeConverter)))
	if err != nil {
		t.Errorf("unable to create mock client: %v", err)
		t.FailNow()
		return nil
	}

	return &dbMock{
		mock:   mock,
		client: db,
	}
}
