package sql

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/lib/pq"
)

func Test_numberPlaceholder(t *testing.T) {
	type args struct {
		query string
		old   string
		new   string
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
			name: "no replaces",
			args: args{
				new:   "$",
				old:   "?",
				query: "SELECT * FROM eventstore.events",
			},
			res: res{
				query: "SELECT * FROM eventstore.events",
			},
		},
		{
			name: "two replaces",
			args: args{
				new:   "$",
				old:   "?",
				query: "SELECT * FROM eventstore.events WHERE aggregate_type = ? AND LIMIT = ?",
			},
			res: res{
				query: "SELECT * FROM eventstore.events WHERE aggregate_type = $1 AND LIMIT = $2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := numberPlaceholder(tt.args.query, tt.args.old, tt.args.new); got != tt.res.query {
				t.Errorf("numberPlaceholder() = %v, want %v", got, tt.res.query)
			}
		})
	}
}

func Test_getOperation(t *testing.T) {
	t.Run("all ops", func(t *testing.T) {
		for op, expected := range map[es_models.Operation]string{
			es_models.Operation_Equals:  "=",
			es_models.Operation_In:      "=",
			es_models.Operation_Greater: ">",
			es_models.Operation_Less:    "<",
			es_models.Operation(-1):     "",
		} {
			if got := getOperation(op); got != expected {
				t.Errorf("getOperation() = %v, want %v", got, expected)
			}
		}
	})
}

func Test_getField(t *testing.T) {
	t.Run("all fields", func(t *testing.T) {
		for field, expected := range map[es_models.Field]string{
			es_models.Field_AggregateType:  "aggregate_type",
			es_models.Field_AggregateID:    "aggregate_id",
			es_models.Field_LatestSequence: "event_sequence",
			es_models.Field_ResourceOwner:  "resource_owner",
			es_models.Field_EditorService:  "editor_service",
			es_models.Field_EditorUser:     "editor_user",
			es_models.Field_EventType:      "event_type",
			es_models.Field(-1):            "",
		} {
			if got := getField(field); got != expected {
				t.Errorf("getField() = %v, want %v", got, expected)
			}
		}
	})
}

func Test_getConditionFormat(t *testing.T) {
	type args struct {
		operation es_models.Operation
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no in operation",
			args: args{
				operation: es_models.Operation_Equals,
			},
			want: "%s %s ?",
		},
		{
			name: "in operation",
			args: args{
				operation: es_models.Operation_In,
			},
			want: "%s %s ANY(?)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getConditionFormat(tt.args.operation); got != tt.want {
				t.Errorf("prepareConditionFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCondition(t *testing.T) {
	type args struct {
		filter *es_models.Filter
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "equals",
			args: args{filter: es_models.NewFilter(es_models.Field_AggregateID, "", es_models.Operation_Equals)},
			want: "aggregate_id = ?",
		},
		{
			name: "greater",
			args: args{filter: es_models.NewFilter(es_models.Field_LatestSequence, 0, es_models.Operation_Greater)},
			want: "event_sequence > ?",
		},
		{
			name: "less",
			args: args{filter: es_models.NewFilter(es_models.Field_LatestSequence, 5000, es_models.Operation_Less)},
			want: "event_sequence < ?",
		},
		{
			name: "in list",
			args: args{filter: es_models.NewFilter(es_models.Field_AggregateType, []es_models.AggregateType{"movies", "actors"}, es_models.Operation_In)},
			want: "aggregate_type = ANY(?)",
		},
		{
			name: "invalid operation",
			args: args{filter: es_models.NewFilter(es_models.Field_AggregateType, []es_models.AggregateType{"movies", "actors"}, es_models.Operation(-1))},
			want: "",
		},
		{
			name: "invalid field",
			args: args{filter: es_models.NewFilter(es_models.Field(-1), []es_models.AggregateType{"movies", "actors"}, es_models.Operation_Equals)},
			want: "",
		},
		{
			name: "invalid field and operation",
			args: args{filter: es_models.NewFilter(es_models.Field(-1), []es_models.AggregateType{"movies", "actors"}, es_models.Operation(-1))},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCondition(tt.args.filter); got != tt.want {
				t.Errorf("getCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prepareColumns(t *testing.T) {
	type args struct {
		columns models.Columns
		dest    interface{}
		dbErr   error
	}
	type res struct {
		query    string
		dbRow    []interface{}
		expected interface{}
		dbErr    func(error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "invalid columns",
			args: args{columns: es_models.Columns(-1)},
			res: res{
				query: "",
				dbErr: func(err error) bool { return err == nil },
			},
		},
		{
			name: "max column",
			args: args{
				columns: es_models.Columns_Max_Sequence,
				dest:    new(Sequence),
			},
			res: res{
				query:    "SELECT MAX(event_sequence) FROM eventstore.events",
				dbRow:    []interface{}{Sequence(5)},
				expected: Sequence(5),
			},
		},
		{
			name: "max sequence wrong dest type",
			args: args{
				columns: es_models.Columns_Max_Sequence,
				dest:    new(uint64),
			},
			res: res{
				query: "SELECT MAX(event_sequence) FROM eventstore.events",
				dbErr: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "event",
			args: args{
				columns: es_models.Columns_Event,
				dest:    new(models.Event),
			},
			res: res{
				query:    "SELECT creation_date, event_type, event_sequence, previous_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events",
				dbRow:    []interface{}{time.Time{}, models.EventType(""), uint64(5), Sequence(0), Data(nil), "", "", "", models.AggregateType("user"), "hodor", models.Version("")},
				expected: models.Event{AggregateID: "hodor", AggregateType: "user", Sequence: 5, Data: make(Data, 0)},
			},
		},
		{
			name: "event wrong dest type",
			args: args{
				columns: es_models.Columns_Event,
				dest:    new(uint64),
			},
			res: res{
				query: "SELECT creation_date, event_type, event_sequence, previous_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events",
				dbErr: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "event query error",
			args: args{
				columns: es_models.Columns_Event,
				dest:    new(models.Event),
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
			query, rowScanner := prepareColumns(tt.args.columns, true)
			if query != tt.res.query {
				t.Errorf("prepareColumns() got = %v, want %v", query, tt.res.query)
			}
			if tt.res.query == "" && rowScanner != nil {
				t.Errorf("row scanner should be nil")
			}
			if rowScanner == nil {
				return
			}
			err := rowScanner(prepareTestScan(tt.args.dbErr, tt.res.dbRow), tt.args.dest)
			if tt.res.dbErr != nil {
				if !tt.res.dbErr(err) {
					t.Errorf("wrong error type in rowScanner got: %v", err)
				}
			} else {
				if !reflect.DeepEqual(reflect.Indirect(reflect.ValueOf(tt.args.dest)).Interface(), tt.res.expected) {
					t.Errorf("unexpected result from rowScanner want: %v got: %v", tt.res.dbRow, tt.args.dest)
				}
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
		filters []*models.Filter
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
				filters: []*es_models.Filter{},
			},
			res: res{
				clause: "",
				values: nil,
			},
		},
		{
			name: "invalid condition",
			args: args{
				filters: []*es_models.Filter{
					es_models.NewFilter(es_models.Field_AggregateID, "wrong", es_models.Operation(-1)),
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
				filters: []*es_models.Filter{
					es_models.NewFilter(es_models.Field_AggregateType, []es_models.AggregateType{"user", "org"}, es_models.Operation_In),
				},
			},
			res: res{
				clause: " WHERE aggregate_type = ANY(?)",
				values: []interface{}{pq.Array([]es_models.AggregateType{"user", "org"})},
			},
		},
		{
			name: "multiple filters",
			args: args{
				filters: []*es_models.Filter{
					es_models.NewFilter(es_models.Field_AggregateType, []es_models.AggregateType{"user", "org"}, es_models.Operation_In),
					es_models.NewFilter(es_models.Field_AggregateID, "1234", es_models.Operation_Equals),
					es_models.NewFilter(es_models.Field_EventType, []es_models.EventType{"user.created", "org.created"}, es_models.Operation_In),
				},
			},
			res: res{
				clause: " WHERE aggregate_type = ANY(?) AND aggregate_id = ? AND event_type = ANY(?)",
				values: []interface{}{pq.Array([]es_models.AggregateType{"user", "org"}), "1234", pq.Array([]es_models.EventType{"user.created", "org.created"})},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClause, gotValues := prepareCondition(tt.args.filters)
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
		queryFactory *models.SearchQueryFactory
	}
	type res struct {
		query      string
		limit      uint64
		values     []interface{}
		rowScanner bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "invalid query factory",
			args: args{
				queryFactory: nil,
			},
			res: res{
				query:      "",
				limit:      0,
				rowScanner: false,
				values:     nil,
			},
		},
		{
			name: "with order by desc",
			args: args{
				queryFactory: es_models.NewSearchQueryFactory("user").OrderDesc(),
			},
			res: res{
				query:      "SELECT creation_date, event_type, event_sequence, previous_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE aggregate_type = $1 ORDER BY event_sequence DESC",
				rowScanner: true,
				values:     []interface{}{es_models.AggregateType("user")},
			},
		},
		{
			name: "with limit",
			args: args{
				queryFactory: es_models.NewSearchQueryFactory("user").Limit(5),
			},
			res: res{
				query:      "SELECT creation_date, event_type, event_sequence, previous_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE aggregate_type = $1 ORDER BY event_sequence LIMIT $2",
				rowScanner: true,
				values:     []interface{}{es_models.AggregateType("user"), uint64(5)},
				limit:      5,
			},
		},
		{
			name: "with limit and order by desc",
			args: args{
				queryFactory: es_models.NewSearchQueryFactory("user").Limit(5).OrderDesc(),
			},
			res: res{
				query:      "SELECT creation_date, event_type, event_sequence, previous_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore.events WHERE aggregate_type = $1 ORDER BY event_sequence DESC LIMIT $2",
				rowScanner: true,
				values:     []interface{}{es_models.AggregateType("user"), uint64(5)},
				limit:      5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotLimit, gotValues, gotRowScanner := buildQuery(tt.args.queryFactory)
			if gotQuery != tt.res.query {
				t.Errorf("buildQuery() gotQuery = %v, want %v", gotQuery, tt.res.query)
			}
			if gotLimit != tt.res.limit {
				t.Errorf("buildQuery() gotLimit = %v, want %v", gotLimit, tt.res.limit)
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
			if (tt.res.rowScanner && gotRowScanner == nil) || (!tt.res.rowScanner && gotRowScanner != nil) {
				t.Errorf("rowScanner should be nil==%v got nil==%v", tt.res.rowScanner, gotRowScanner == nil)
			}
		})
	}
}
