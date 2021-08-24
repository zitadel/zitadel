package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/lib/pq"
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
			want: "event_sequence > ?",
		},
		{
			name: "less",
			args: args{filter: repository.NewFilter(repository.FieldSequence, 5000, repository.OperationLess)},
			want: "event_sequence < ?",
		},
		{
			name: "in list",
			args: args{filter: repository.NewFilter(repository.FieldAggregateType, []repository.AggregateType{"movies", "actors"}, repository.OperationIn)},
			want: "aggregate_type = ANY(?)",
		},
		{
			name: "invalid operation",
			args: args{filter: repository.NewFilter(repository.FieldAggregateType, []repository.AggregateType{"movies", "actors"}, repository.Operation(-1))},
			want: "",
		},
		{
			name: "invalid field",
			args: args{filter: repository.NewFilter(repository.Field(-1), []repository.AggregateType{"movies", "actors"}, repository.OperationEquals)},
			want: "",
		},
		{
			name: "invalid field and operation",
			args: args{filter: repository.NewFilter(repository.Field(-1), []repository.AggregateType{"movies", "actors"}, repository.Operation(-1))},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{}
			if got := getCondition(db, tt.args.filter); got != tt.want {
				t.Errorf("getCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prepareColumns(t *testing.T) {
	type fields struct {
		dbRow []interface{}
	}
	type args struct {
		columns repository.Columns
		dest    interface{}
		dbErr   error
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
			args: args{columns: repository.Columns(-1)},
			res: res{
				query: "",
				dbErr: func(err error) bool { return err == nil },
			},
		},
		{
			name: "max column",
			args: args{
				columns: repository.ColumnsMaxSequence,
				dest:    new(Sequence),
			},
			res: res{
				query:    "SELECT MAX(event_sequence) FROM eventstore.events",
				expected: Sequence(5),
			},
			fields: fields{
				dbRow: []interface{}{Sequence(5)},
			},
		},
		{
			name: "max sequence wrong dest type",
			args: args{
				columns: repository.ColumnsMaxSequence,
				dest:    new(uint64),
			},
			res: res{
				query: "SELECT MAX(event_sequence) FROM eventstore.events",
				dbErr: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "events",
			args: args{
				columns: repository.ColumnsEvent,
				dest:    &[]*repository.Event{},
			},
			res: res{
				query: "SELECT creation_date, event_type, event_sequence, previous_aggregate_sequence, previous_aggregate_type_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events",
				expected: []*repository.Event{
					{AggregateID: "hodor", AggregateType: "user", Sequence: 5, Data: make(Data, 0)},
				},
			},
			fields: fields{
				dbRow: []interface{}{time.Time{}, repository.EventType(""), uint64(5), Sequence(0), Sequence(0), Data(nil), "", "", "", repository.AggregateType("user"), "hodor", repository.Version("")},
			},
		},
		{
			name: "events wrong dest type",
			args: args{
				columns: repository.ColumnsEvent,
				dest:    []*repository.Event{},
			},
			res: res{
				query: "SELECT creation_date, event_type, event_sequence, previous_aggregate_sequence, previous_aggregate_type_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events",
				dbErr: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "event query error",
			args: args{
				columns: repository.ColumnsEvent,
				dest:    &[]*repository.Event{},
				dbErr:   sql.ErrConnDone,
			},
			res: res{
				query: "SELECT creation_date, event_type, event_sequence, previous_aggregate_sequence, previous_aggregate_type_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events",
				dbErr: errors.IsInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crdb := &CRDB{}
			query, rowScanner := prepareColumns(crdb, tt.args.columns)
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
			if !reflect.DeepEqual(reflect.Indirect(reflect.ValueOf(tt.args.dest)).Interface(), tt.res.expected) {
				t.Errorf("unexpected result from rowScanner \nwant: %+v \ngot: %+v", tt.fields.dbRow, reflect.Indirect(reflect.ValueOf(tt.args.dest)).Interface())
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
			return errors.ThrowInvalidArgumentf(nil, "SQL-NML1q", "expected len %d got %d", len(res), len(dests))
		}
		for i, r := range res {
			reflect.ValueOf(dests[i]).Elem().Set(reflect.ValueOf(r))
		}

		return nil
	}
}

func Test_prepareCondition(t *testing.T) {
	type args struct {
		filters [][]*repository.Filter
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
				filters: nil,
			},
			res: res{
				clause: "",
				values: nil,
			},
		},
		{
			name: "empty filters",
			args: args{
				filters: [][]*repository.Filter{},
			},
			res: res{
				clause: "",
				values: nil,
			},
		},
		{
			name: "invalid condition",
			args: args{
				filters: [][]*repository.Filter{
					{
						repository.NewFilter(repository.FieldAggregateID, "wrong", repository.Operation(-1)),
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
				filters: [][]*repository.Filter{
					{
						repository.NewFilter(repository.FieldAggregateType, []repository.AggregateType{"user", "org"}, repository.OperationIn),
					},
				},
			},
			res: res{
				clause: " WHERE ( aggregate_type = ANY(?) )",
				values: []interface{}{pq.Array([]repository.AggregateType{"user", "org"})},
			},
		},
		{
			name: "multiple filters",
			args: args{
				filters: [][]*repository.Filter{
					{
						repository.NewFilter(repository.FieldAggregateType, []repository.AggregateType{"user", "org"}, repository.OperationIn),
						repository.NewFilter(repository.FieldAggregateID, "1234", repository.OperationEquals),
						repository.NewFilter(repository.FieldEventType, []repository.EventType{"user.created", "org.created"}, repository.OperationIn),
					},
				},
			},
			res: res{
				clause: " WHERE ( aggregate_type = ANY(?) AND aggregate_id = ? AND event_type = ANY(?) )",
				values: []interface{}{pq.Array([]repository.AggregateType{"user", "org"}), "1234", pq.Array([]repository.EventType{"user.created", "org.created"})},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crdb := &CRDB{}
			gotClause, gotValues := prepareCondition(crdb, tt.args.filters)
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
		searchQuery *repository.SearchQuery
	}
	type fields struct {
		existingEvents []*repository.Event
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
				searchQuery: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Filters: [][]*repository.Filter{
						{
							repository.NewFilter(repository.FieldAggregateType, "not found", repository.OperationEquals),
						},
					},
				},
			},
			fields: fields{
				client: testCRDBClient,
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
			name: "aggregate type filter events found",
			args: args{
				searchQuery: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Filters: [][]*repository.Filter{
						{
							repository.NewFilter(repository.FieldAggregateType, t.Name(), repository.OperationEquals),
						},
					},
				},
			},
			fields: fields{
				client: testCRDBClient,
				existingEvents: []*repository.Event{
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
				searchQuery: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Filters: [][]*repository.Filter{
						{
							repository.NewFilter(repository.FieldAggregateType, t.Name(), repository.OperationEquals),
							repository.NewFilter(repository.FieldAggregateID, "303", repository.OperationEquals),
						},
					},
				},
			},
			fields: fields{
				client: testCRDBClient,
				existingEvents: []*repository.Event{
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
				searchQuery: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Filters: [][]*repository.Filter{
						{
							repository.NewFilter(repository.FieldResourceOwner, "caos", repository.OperationEquals),
						},
					},
				},
			},
			fields: fields{
				client: testCRDBClient,
				existingEvents: []*repository.Event{
					generateEvent(t, "306", func(e *repository.Event) { e.ResourceOwner = "caos" }),
					generateEvent(t, "307", func(e *repository.Event) { e.ResourceOwner = "caos" }),
					generateEvent(t, "308", func(e *repository.Event) { e.ResourceOwner = "caos" }),
					generateEvent(t, "309", func(e *repository.Event) { e.ResourceOwner = "orgID" }),
					generateEvent(t, "309", func(e *repository.Event) { e.ResourceOwner = "orgID" }),
				},
			},
			res: res{
				eventCount: 3,
			},
			wantErr: false,
		},
		{
			name: "editor service filter events found",
			args: args{
				searchQuery: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Filters: [][]*repository.Filter{
						{
							repository.NewFilter(repository.FieldEditorService, "MANAGEMENT-API", repository.OperationEquals),
							repository.NewFilter(repository.FieldEditorService, "ADMIN-API", repository.OperationEquals),
						},
					},
				},
			},
			fields: fields{
				client: testCRDBClient,
				existingEvents: []*repository.Event{
					generateEvent(t, "307", func(e *repository.Event) { e.EditorService = "MANAGEMENT-API" }),
					generateEvent(t, "307", func(e *repository.Event) { e.EditorService = "MANAGEMENT-API" }),
					generateEvent(t, "308", func(e *repository.Event) { e.EditorService = "ADMIN-API" }),
					generateEvent(t, "309", func(e *repository.Event) { e.EditorService = "AUTHAPI" }),
					generateEvent(t, "309", func(e *repository.Event) { e.EditorService = "AUTHAPI" }),
				},
			},
			res: res{
				eventCount: 3,
			},
			wantErr: false,
		},
		{
			name: "editor user filter events found",
			args: args{
				searchQuery: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Filters: [][]*repository.Filter{
						{
							repository.NewFilter(repository.FieldEditorUser, "adlerhurst", repository.OperationEquals),
							repository.NewFilter(repository.FieldEditorUser, "nobody", repository.OperationEquals),
							repository.NewFilter(repository.FieldEditorUser, "", repository.OperationEquals),
						},
					},
				},
			},
			fields: fields{
				client: testCRDBClient,
				existingEvents: []*repository.Event{
					generateEvent(t, "310", func(e *repository.Event) { e.EditorUser = "adlerhurst" }),
					generateEvent(t, "310", func(e *repository.Event) { e.EditorUser = "adlerhurst" }),
					generateEvent(t, "310", func(e *repository.Event) { e.EditorUser = "nobody" }),
					generateEvent(t, "311", func(e *repository.Event) { e.EditorUser = "" }),
					generateEvent(t, "311", func(e *repository.Event) { e.EditorUser = "" }),
					generateEvent(t, "312", func(e *repository.Event) { e.EditorUser = "fforootd" }),
					generateEvent(t, "312", func(e *repository.Event) { e.EditorUser = "fforootd" }),
				},
			},
			res: res{
				eventCount: 5,
			},
			wantErr: false,
		},
		{
			name: "event type filter events found",
			args: args{
				searchQuery: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Filters: [][]*repository.Filter{
						{
							repository.NewFilter(repository.FieldEventType, repository.EventType("user.created"), repository.OperationEquals),
							repository.NewFilter(repository.FieldEventType, repository.EventType("user.updated"), repository.OperationEquals),
						},
					},
				},
			},
			fields: fields{
				client: testCRDBClient,
				existingEvents: []*repository.Event{
					generateEvent(t, "311", func(e *repository.Event) { e.Type = "user.created" }),
					generateEvent(t, "311", func(e *repository.Event) { e.Type = "user.updated" }),
					generateEvent(t, "311", func(e *repository.Event) { e.Type = "user.deactivated" }),
					generateEvent(t, "311", func(e *repository.Event) { e.Type = "user.locked" }),
					generateEvent(t, "312", func(e *repository.Event) { e.Type = "user.created" }),
					generateEvent(t, "312", func(e *repository.Event) { e.Type = "user.updated" }),
					generateEvent(t, "312", func(e *repository.Event) { e.Type = "user.deactivated" }),
					generateEvent(t, "312", func(e *repository.Event) { e.Type = "user.reactivated" }),
					generateEvent(t, "313", func(e *repository.Event) { e.Type = "user.locked" }),
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
				searchQuery: &repository.SearchQuery{},
			},
			fields: fields{
				client:         testCRDBClient,
				existingEvents: []*repository.Event{},
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
				client: tt.fields.client,
			}

			// setup initial data for query
			if err := db.Push(context.Background(), tt.fields.existingEvents); err != nil {
				t.Errorf("error in setup = %v", err)
				return
			}

			events := []*repository.Event{}
			if err := query(context.Background(), db, tt.args.searchQuery, &events); (err != nil) != tt.wantErr {
				t.Errorf("CRDB.query() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_query_events_mocked(t *testing.T) {
	type args struct {
		query *repository.SearchQuery
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
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    true,
					Filters: [][]*repository.Filter{
						{
							{
								Field:     repository.FieldAggregateType,
								Value:     repository.AggregateType("user"),
								Operation: repository.OperationEquals,
							},
						},
					},
				},
			},
			fields: fields{
				mock: newMockClient(t).expectQuery(t,
					`SELECT creation_date, event_type, event_sequence, previous_aggregate_sequence, previous_aggregate_type_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE \( aggregate_type = \$1 \) ORDER BY event_sequence DESC`,
					[]driver.Value{repository.AggregateType("user")},
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
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    false,
					Limit:   5,
					Filters: [][]*repository.Filter{
						{
							{
								Field:     repository.FieldAggregateType,
								Value:     repository.AggregateType("user"),
								Operation: repository.OperationEquals,
							},
						},
					},
				},
			},
			fields: fields{
				mock: newMockClient(t).expectQuery(t,
					`SELECT creation_date, event_type, event_sequence, previous_aggregate_sequence, previous_aggregate_type_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE \( aggregate_type = \$1 \) ORDER BY event_sequence LIMIT \$2`,
					[]driver.Value{repository.AggregateType("user"), uint64(5)},
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
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    true,
					Limit:   5,
					Filters: [][]*repository.Filter{
						{
							{
								Field:     repository.FieldAggregateType,
								Value:     repository.AggregateType("user"),
								Operation: repository.OperationEquals,
							},
						},
					},
				},
			},
			fields: fields{
				mock: newMockClient(t).expectQuery(t,
					`SELECT creation_date, event_type, event_sequence, previous_aggregate_sequence, previous_aggregate_type_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE \( aggregate_type = \$1 \) ORDER BY event_sequence DESC LIMIT \$2`,
					[]driver.Value{repository.AggregateType("user"), uint64(5)},
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
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    true,
					Limit:   0,
					Filters: [][]*repository.Filter{
						{
							{
								Field:     repository.FieldAggregateType,
								Value:     repository.AggregateType("user"),
								Operation: repository.OperationEquals,
							},
						},
					},
				},
			},
			fields: fields{
				mock: newMockClient(t).expectQueryErr(t,
					`SELECT creation_date, event_type, event_sequence, previous_aggregate_sequence, previous_aggregate_type_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE \( aggregate_type = \$1 \) ORDER BY event_sequence DESC`,
					[]driver.Value{repository.AggregateType("user")},
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
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    true,
					Limit:   0,
					Filters: [][]*repository.Filter{
						{
							{
								Field:     repository.FieldAggregateType,
								Value:     repository.AggregateType("user"),
								Operation: repository.OperationEquals,
							},
						},
					},
				},
			},
			fields: fields{
				mock: newMockClient(t).expectQuery(t,
					`SELECT creation_date, event_type, event_sequence, previous_aggregate_sequence, previous_aggregate_type_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE \( aggregate_type = \$1 \) ORDER BY event_sequence DESC`,
					[]driver.Value{repository.AggregateType("user")},
					&repository.Event{Sequence: 100}),
			},
			res: res{
				wantErr: true,
			},
		},
		{
			name: "error no columns",
			args: args{
				query: &repository.SearchQuery{
					Columns: repository.Columns(-1),
				},
			},
			res: res{
				wantErr: true,
			},
		},
		{
			name: "invalid condition",
			args: args{
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Filters: [][]*repository.Filter{
						{
							{},
						},
					},
				},
			},
			res: res{
				wantErr: true,
			},
		},
		{
			name: "with subqueries",
			args: args{
				dest: &[]*repository.Event{},
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    true,
					Limit:   5,
					Filters: [][]*repository.Filter{
						{
							{
								Field:     repository.FieldAggregateType,
								Value:     repository.AggregateType("user"),
								Operation: repository.OperationEquals,
							},
						},
						{
							{
								Field:     repository.FieldAggregateType,
								Value:     repository.AggregateType("org"),
								Operation: repository.OperationEquals,
							},
							{
								Field:     repository.FieldAggregateID,
								Value:     "asdf42",
								Operation: repository.OperationEquals,
							},
						},
					},
				},
			},
			fields: fields{
				mock: newMockClient(t).expectQuery(t,
					`SELECT creation_date, event_type, event_sequence, previous_aggregate_sequence, previous_aggregate_type_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE \( aggregate_type = \$1 \) OR \( aggregate_type = \$2 AND aggregate_id = \$3 \) ORDER BY event_sequence DESC LIMIT \$4`,
					[]driver.Value{repository.AggregateType("user"), repository.AggregateType("org"), "asdf42", uint64(5)},
				),
			},
			res: res{
				wantErr: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crdb := &CRDB{}
			if tt.fields.mock != nil {
				crdb.client = tt.fields.mock.client
			}

			err := query(context.Background(), crdb, tt.args.query, tt.args.dest)
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
	rows := sqlmock.NewRows([]string{"event_sequence"})
	for _, event := range events {
		rows = rows.AddRow(event.Sequence)
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
	db, mock, err := sqlmock.New()
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
