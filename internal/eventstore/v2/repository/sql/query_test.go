package sql

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
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
				query: "SELECT creation_date, event_type, event_sequence, previous_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events",
				expected: []*repository.Event{
					{AggregateID: "hodor", AggregateType: "user", Sequence: 5, Data: make(Data, 0)},
				},
			},
			fields: fields{
				dbRow: []interface{}{time.Time{}, repository.EventType(""), uint64(5), Sequence(0), Data(nil), "", "", "", repository.AggregateType("user"), "hodor", repository.Version("")},
			},
		},
		{
			name: "events wrong dest type",
			args: args{
				columns: repository.ColumnsEvent,
				dest:    []*repository.Event{},
			},
			res: res{
				query: "SELECT creation_date, event_type, event_sequence, previous_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events",
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
				query: "SELECT creation_date, event_type, event_sequence, previous_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events",
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
		filters []*repository.Filter
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
				filters: []*repository.Filter{},
			},
			res: res{
				clause: "",
				values: nil,
			},
		},
		{
			name: "invalid condition",
			args: args{
				filters: []*repository.Filter{
					repository.NewFilter(repository.FieldAggregateID, "wrong", repository.Operation(-1)),
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
				filters: []*repository.Filter{
					repository.NewFilter(repository.FieldAggregateType, []repository.AggregateType{"user", "org"}, repository.OperationIn),
				},
			},
			res: res{
				clause: " WHERE aggregate_type = ANY(?)",
				values: []interface{}{pq.Array([]repository.AggregateType{"user", "org"})},
			},
		},
		{
			name: "multiple filters",
			args: args{
				filters: []*repository.Filter{
					repository.NewFilter(repository.FieldAggregateType, []repository.AggregateType{"user", "org"}, repository.OperationIn),
					repository.NewFilter(repository.FieldAggregateID, "1234", repository.OperationEquals),
					repository.NewFilter(repository.FieldEventType, []repository.EventType{"user.created", "org.created"}, repository.OperationIn),
				},
			},
			res: res{
				clause: " WHERE aggregate_type = ANY(?) AND aggregate_id = ? AND event_type = ANY(?)",
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

func Test_buildQuery(t *testing.T) {
	type args struct {
		query *repository.SearchQuery
	}
	type res struct {
		query      string
		values     []interface{}
		rowScanner bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "with order by desc",
			args: args{
				//  NewSearchQueryFactory("user").OrderDesc()
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    true,
					Filters: []*repository.Filter{
						{
							Field:     repository.FieldAggregateType,
							Value:     repository.AggregateType("user"),
							Operation: repository.OperationEquals,
						},
					},
				},
			},
			res: res{
				query:      "SELECT creation_date, event_type, event_sequence, previous_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE aggregate_type = $1 ORDER BY event_sequence DESC",
				rowScanner: true,
				values:     []interface{}{repository.AggregateType("user")},
			},
		},
		{
			name: "with limit",
			args: args{
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    false,
					Limit:   5,
					Filters: []*repository.Filter{
						{
							Field:     repository.FieldAggregateType,
							Value:     repository.AggregateType("user"),
							Operation: repository.OperationEquals,
						},
					},
				},
			},
			res: res{
				query:      "SELECT creation_date, event_type, event_sequence, previous_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE aggregate_type = $1 ORDER BY event_sequence LIMIT $2",
				rowScanner: true,
				values:     []interface{}{repository.AggregateType("user"), uint64(5)},
			},
		},
		{
			name: "with limit and order by desc",
			args: args{
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    true,
					Limit:   5,
					Filters: []*repository.Filter{
						{
							Field:     repository.FieldAggregateType,
							Value:     repository.AggregateType("user"),
							Operation: repository.OperationEquals,
						},
					},
				},
			},
			res: res{
				query:      "SELECT creation_date, event_type, event_sequence, previous_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE aggregate_type = $1 ORDER BY event_sequence DESC LIMIT $2",
				rowScanner: true,
				values:     []interface{}{repository.AggregateType("user"), uint64(5)},
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
				query:      "",
				rowScanner: false,
				values:     []interface{}(nil),
			},
		},
		{
			name: "invalid condition",
			args: args{
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Filters: []*repository.Filter{
						{},
					},
				},
			},
			res: res{
				query:      "",
				rowScanner: false,
				values:     []interface{}(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crdb := &CRDB{}
			gotQuery, gotValues, gotRowScanner := buildQuery(crdb, tt.args.query)
			if gotQuery != tt.res.query {
				t.Errorf("buildQuery() gotQuery = %v, want %v", gotQuery, tt.res.query)
			}
			if len(gotValues) != len(tt.res.values) {
				t.Errorf("wrong length of gotten values got = %d, want %d", len(gotValues), len(tt.res.values))
				return
			}
			if !reflect.DeepEqual(gotValues, tt.res.values) {
				t.Errorf("prepareCondition() gotValues = %T: %v, want %T: %v", gotValues, gotValues, tt.res.values, tt.res.values)
			}
			if (tt.res.rowScanner && gotRowScanner == nil) || (!tt.res.rowScanner && gotRowScanner != nil) {
				t.Errorf("rowScanner should be nil==%v got nil==%v", tt.res.rowScanner, gotRowScanner == nil)
			}
		})
	}
}
