package crdb

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
)

type wantExecuter struct {
	query         string
	args          []interface{}
	t             *testing.T
	wasExecuted   bool
	shouldExecute bool
}

var errTestErr = errors.New("some error")

func (ex *wantExecuter) check(t *testing.T) {
	t.Helper()
	if ex.wasExecuted && !ex.shouldExecute {
		t.Error("executer should not be executed")
	} else if !ex.wasExecuted && ex.shouldExecute {
		t.Error("executer should be executed")
	} else if ex.wasExecuted != ex.shouldExecute {
		t.Errorf("executed missmatched should be %t, but was %t", ex.shouldExecute, ex.wasExecuted)
	}
}

func (ex *wantExecuter) Exec(query string, args ...interface{}) (sql.Result, error) {
	ex.t.Helper()
	ex.wasExecuted = true
	if query != ex.query {
		ex.t.Errorf("wrong query:\n  expected:\n    %q\n  got:\n    %q", ex.query, query)
	}
	if !reflect.DeepEqual(ex.args, args) {
		ex.t.Errorf("wrong args:\n  expected:\n    %v\n  got:\n    %v", ex.args, args)
	}
	return nil, nil
}

func TestNewCreateStatement(t *testing.T) {
	type args struct {
		table  string
		event  *testEvent
		values []handler.Column
	}
	type want struct {
		aggregateType    eventstore.AggregateType
		sequence         uint64
		previousSequence uint64
		table            string
		executer         *wantExecuter
		isErr            func(error) bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no table",
			args: args{
				table: "",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				values: []handler.Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
			},
			want: want{
				table:            "",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 0,
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, handler.ErrNoProjection)
				},
			},
		},
		{
			name: "no values",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				values: []handler.Column{},
			},
			want: want{
				table:            "my_table",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 1,
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, handler.ErrNoValues)
				},
			},
		},
		{
			name: "correct",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				values: []handler.Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
			},
			want: want{
				table:            "my_table",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 1,
				executer: &wantExecuter{
					query:         "INSERT INTO my_table (col1) VALUES ($1)",
					shouldExecute: true,
					args:          []interface{}{"val"},
				},
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.executer.t = t
			stmt := NewCreateStatement(tt.args.event, tt.args.values)

			err := stmt.Execute(tt.want.executer, tt.args.table)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}
			tt.want.executer.check(t)
		})
	}
}

func TestNewUpsertStatement(t *testing.T) {
	type args struct {
		table  string
		event  *testEvent
		values []handler.Column
	}
	type want struct {
		aggregateType    eventstore.AggregateType
		sequence         uint64
		previousSequence uint64
		table            string
		executer         *wantExecuter
		isErr            func(error) bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no table",
			args: args{
				table: "",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				values: []handler.Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
			},
			want: want{
				table:            "",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 0,
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, handler.ErrNoProjection)
				},
			},
		},
		{
			name: "no values",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				values: []handler.Column{},
			},
			want: want{
				table:            "my_table",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 1,
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, handler.ErrNoValues)
				},
			},
		},
		{
			name: "correct",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				values: []handler.Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
			},
			want: want{
				table:            "my_table",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 1,
				executer: &wantExecuter{
					query:         "UPSERT INTO my_table (col1) VALUES ($1)",
					shouldExecute: true,
					args:          []interface{}{"val"},
				},
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.executer.t = t
			stmt := NewUpsertStatement(tt.args.event, tt.args.values)

			err := stmt.Execute(tt.want.executer, tt.args.table)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}
			tt.want.executer.check(t)
		})
	}
}

func TestNewUpdateStatement(t *testing.T) {
	type args struct {
		table      string
		event      *testEvent
		conditions []handler.Column
		values     []handler.Column
	}
	type want struct {
		table            string
		aggregateType    eventstore.AggregateType
		sequence         uint64
		previousSequence uint64
		executer         *wantExecuter
		isErr            func(error) bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no table",
			args: args{
				table: "",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				values: []handler.Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
				conditions: []handler.Column{
					{
						Name:  "col2",
						Value: 1,
					},
				},
			},
			want: want{
				table:            "",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 0,
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, handler.ErrNoProjection)
				},
			},
		},
		{
			name: "no values",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				values: []handler.Column{},
				conditions: []handler.Column{
					{
						Name:  "col2",
						Value: 1,
					},
				},
			},
			want: want{
				table:            "my_table",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 1,
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, handler.ErrNoValues)
				},
			},
		},
		{
			name: "no conditions",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				values: []handler.Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
				conditions: []handler.Column{},
			},
			want: want{
				table:            "my_table",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 1,
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, handler.ErrNoCondition)
				},
			},
		},
		{
			name: "correct",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				values: []handler.Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
				conditions: []handler.Column{
					{
						Name:  "col2",
						Value: 1,
					},
				},
			},
			want: want{
				table:            "my_table",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 1,
				executer: &wantExecuter{
					query:         "UPDATE my_table SET (col1) = ($1) WHERE (col2 = $2)",
					shouldExecute: true,
					args:          []interface{}{"val", 1},
				},
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.executer.t = t
			stmt := NewUpdateStatement(tt.args.event, tt.args.values, tt.args.conditions)

			err := stmt.Execute(tt.want.executer, tt.args.table)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}
			tt.want.executer.check(t)
		})
	}
}

func TestNewDeleteStatement(t *testing.T) {
	type args struct {
		table      string
		event      *testEvent
		conditions []handler.Column
	}

	type want struct {
		table            string
		aggregateType    eventstore.AggregateType
		sequence         uint64
		previousSequence uint64
		executer         *wantExecuter
		isErr            func(error) bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no table",
			args: args{
				table: "",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				conditions: []handler.Column{
					{
						Name:  "col2",
						Value: 1,
					},
				},
			},
			want: want{
				table:            "",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 0,
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, handler.ErrNoProjection)
				},
			},
		},
		{
			name: "no conditions",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				conditions: []handler.Column{},
			},
			want: want{
				table:            "my_table",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 1,
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, handler.ErrNoCondition)
				},
			},
		},
		{
			name: "correct",
			args: args{
				table: "my_table",
				event: &testEvent{
					sequence:         1,
					previousSequence: 0,
					aggregateType:    "agg",
				},
				conditions: []handler.Column{
					{
						Name:  "col1",
						Value: 1,
					},
				},
			},
			want: want{
				table:            "my_table",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 1,
				executer: &wantExecuter{
					query:         "DELETE FROM my_table WHERE (col1 = $1)",
					shouldExecute: true,
					args:          []interface{}{1},
				},
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.executer.t = t
			stmt := NewDeleteStatement(tt.args.event, tt.args.conditions)

			err := stmt.Execute(tt.want.executer, tt.args.table)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}
			tt.want.executer.check(t)
		})
	}
}

func TestNewNoOpStatement(t *testing.T) {
	type args struct {
		event *testEvent
	}
	tests := []struct {
		name string
		args args
		want handler.Statement
	}{
		{
			name: "generate correctly",
			args: args{
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         5,
					previousSequence: 3,
				},
			},
			want: handler.Statement{
				AggregateType:    "agg",
				Execute:          nil,
				Sequence:         5,
				PreviousSequence: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNoOpStatement(tt.args.event); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNoOpStatement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatement_Execute(t *testing.T) {
	type fields struct {
		execute func(ex handler.Executer, projectionName string) error
	}
	type want struct {
		isErr func(error) bool
	}
	type args struct {
		projectionName string
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   want
	}{
		{
			name: "execute returns no error",
			fields: fields{
				execute: func(ex handler.Executer, projectionName string) error { return nil },
			},
			args: args{
				projectionName: "my_projection",
			},
			want: want{
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
		{
			name: "execute returns error",
			args: args{
				projectionName: "my_projection",
			},
			fields: fields{
				execute: func(ex handler.Executer, projectionName string) error { return errTestErr },
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, errTestErr)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := &handler.Statement{
				Execute: tt.fields.execute,
			}
			if err := stmt.Execute(nil, tt.args.projectionName); !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func Test_columnsToQuery(t *testing.T) {
	type args struct {
		cols []handler.Column
	}
	type want struct {
		names  []string
		params []string
		values []interface{}
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no columns",
			args: args{},
			want: want{
				names:  []string{},
				params: []string{},
				values: []interface{}{},
			},
		},
		{
			name: "one column",
			args: args{
				cols: []handler.Column{
					{
						Name:  "col1",
						Value: 1,
					},
				},
			},
			want: want{
				names:  []string{"col1"},
				params: []string{"$1"},
				values: []interface{}{1},
			},
		},
		{
			name: "multiple columns",
			args: args{
				cols: []handler.Column{
					{
						Name:  "col1",
						Value: 1,
					},
					{
						Name:  "col2",
						Value: 3.14,
					},
				},
			},
			want: want{
				names:  []string{"col1", "col2"},
				params: []string{"$1", "$2"},
				values: []interface{}{1, 3.14},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNames, gotParameters, gotValues := columnsToQuery(tt.args.cols)
			if !reflect.DeepEqual(gotNames, tt.want.names) {
				t.Errorf("columnsToQuery() gotNames = %v, want %v", gotNames, tt.want.names)
			}
			if !reflect.DeepEqual(gotParameters, tt.want.params) {
				t.Errorf("columnsToQuery() gotParameters = %v, want %v", gotParameters, tt.want.params)
			}
			if !reflect.DeepEqual(gotValues, tt.want.values) {
				t.Errorf("columnsToQuery() gotValues = %v, want %v", gotValues, tt.want.values)
			}
		})
	}
}

func Test_columnsToWhere(t *testing.T) {
	type args struct {
		cols        []handler.Column
		paramOffset int
	}
	type want struct {
		wheres []string
		values []interface{}
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no wheres",
			args: args{},
			want: want{
				wheres: []string{},
				values: []interface{}{},
			},
		},
		{
			name: "no offset",
			args: args{
				cols: []handler.Column{
					{
						Name:  "col1",
						Value: "val1",
					},
				},
				paramOffset: 0,
			},
			want: want{
				wheres: []string{"(col1 = $1)"},
				values: []interface{}{"val1"},
			},
		},
		{
			name: "multiple cols",
			args: args{
				cols: []handler.Column{
					{
						Name:  "col1",
						Value: "val1",
					},
					{
						Name:  "col2",
						Value: "val2",
					},
				},
				paramOffset: 0,
			},
			want: want{
				wheres: []string{"(col1 = $1)", "(col2 = $2)"},
				values: []interface{}{"val1", "val2"},
			},
		},
		{
			name: "2 offset",
			args: args{
				cols: []handler.Column{
					{
						Name:  "col1",
						Value: "val1",
					},
				},
				paramOffset: 2,
			},
			want: want{
				wheres: []string{"(col1 = $3)"},
				values: []interface{}{"val1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWheres, gotValues := columnsToWhere(tt.args.cols, tt.args.paramOffset)
			if !reflect.DeepEqual(gotWheres, tt.want.wheres) {
				t.Errorf("columnsToWhere() gotWheres = %v, want %v", gotWheres, tt.want.wheres)
			}
			if !reflect.DeepEqual(gotValues, tt.want.values) {
				t.Errorf("columnsToWhere() gotValues = %v, want %v", gotValues, tt.want.values)
			}
		})
	}
}
