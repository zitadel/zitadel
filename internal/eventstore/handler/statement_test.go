package handler

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"
)

type wantExecuter struct {
	query         string
	args          []interface{}
	t             *testing.T
	wasExecuted   bool
	shouldExecute bool
}

var testErr = errors.New("some error")

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
		table            string
		values           []Column
		sequence         uint64
		previousSequence uint64
	}
	type want struct {
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
				values: []Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
				sequence:         1,
				previousSequence: 0,
			},
			want: want{
				sequence:         1,
				previousSequence: 0,
				table:            "",
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, ErrNoTable)
				},
			},
		},
		{
			name: "sequence equal prev seq",
			args: args{
				table: "my_table",
				values: []Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
				sequence:         1,
				previousSequence: 1,
			},
			want: want{
				sequence:         1,
				previousSequence: 1,
				table:            "my_table",
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, ErrPrevSeqGtSeq)
				},
			},
		},
		{
			name: "sequence less prev seq",
			args: args{
				table: "my_table",
				values: []Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
				sequence:         1,
				previousSequence: 2,
			},
			want: want{
				sequence:         1,
				previousSequence: 2,
				table:            "my_table",
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, ErrPrevSeqGtSeq)
				},
			},
		},
		{
			name: "no values",
			args: args{
				table:            "my_table",
				values:           []Column{},
				sequence:         1,
				previousSequence: 0,
			},
			want: want{
				sequence:         1,
				previousSequence: 1,
				table:            "my_table",
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, ErrNoValues)
				},
			},
		},
		{
			name: "correct",
			args: args{
				table: "my_table",
				values: []Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
				sequence:         1,
				previousSequence: 0,
			},
			want: want{
				sequence:         1,
				previousSequence: 1,
				table:            "my_table",
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
			stmt := NewCreateStatement(tt.args.values, tt.args.sequence, tt.args.previousSequence)

			err := stmt.execute(tt.want.executer, tt.args.table)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}
			tt.want.executer.check(t)
		})
	}
}

func TestNewUpdateStatement(t *testing.T) {
	type args struct {
		table            string
		conditions       []Column
		values           []Column
		sequence         uint64
		previousSequence uint64
	}
	type want struct {
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
				values: []Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
				conditions: []Column{
					{
						Name:  "col2",
						Value: 1,
					},
				},
				sequence:         1,
				previousSequence: 0,
			},
			want: want{
				sequence:         1,
				previousSequence: 0,
				table:            "",
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, ErrNoTable)
				},
			},
		},
		{
			name: "sequence equal prev seq",
			args: args{
				table: "my_table",
				values: []Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
				conditions: []Column{
					{
						Name:  "col2",
						Value: 1,
					},
				},
				sequence:         1,
				previousSequence: 1,
			},
			want: want{
				sequence:         1,
				previousSequence: 1,
				table:            "my_table",
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, ErrPrevSeqGtSeq)
				},
			},
		},
		{
			name: "sequence less prev seq",
			args: args{
				table: "my_table",
				values: []Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
				conditions: []Column{
					{
						Name:  "col2",
						Value: 1,
					},
				},
				sequence:         1,
				previousSequence: 2,
			},
			want: want{
				sequence:         1,
				previousSequence: 2,
				table:            "my_table",
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, ErrPrevSeqGtSeq)
				},
			},
		},
		{
			name: "no values",
			args: args{
				table:  "my_table",
				values: []Column{},
				conditions: []Column{
					{
						Name:  "col2",
						Value: 1,
					},
				},
				sequence:         1,
				previousSequence: 0,
			},
			want: want{
				sequence:         1,
				previousSequence: 1,
				table:            "my_table",
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, ErrNoValues)
				},
			},
		},
		{
			name: "no conditions",
			args: args{
				table: "my_table",
				values: []Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
				conditions:       []Column{},
				sequence:         1,
				previousSequence: 0,
			},
			want: want{
				sequence:         1,
				previousSequence: 1,
				table:            "my_table",
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, ErrNoCondition)
				},
			},
		},
		{
			name: "correct",
			args: args{
				table: "my_table",
				values: []Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
				conditions: []Column{
					{
						Name:  "col2",
						Value: 1,
					},
				},
				sequence:         1,
				previousSequence: 0,
			},
			want: want{
				sequence:         1,
				previousSequence: 1,
				table:            "my_table",
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
			stmt := NewUpdateStatement(tt.args.conditions, tt.args.values, tt.args.sequence, tt.args.previousSequence)

			err := stmt.execute(tt.want.executer, tt.args.table)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}
			tt.want.executer.check(t)
		})
	}
}

func TestNewDeleteStatement(t *testing.T) {
	type args struct {
		table            string
		conditions       []Column
		sequence         uint64
		previousSequence uint64
	}

	type want struct {
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
				conditions: []Column{
					{
						Name:  "col2",
						Value: 1,
					},
				},
				sequence:         1,
				previousSequence: 0,
			},
			want: want{
				sequence:         1,
				previousSequence: 0,
				table:            "",
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, ErrNoTable)
				},
			},
		},
		{
			name: "sequence equal prev seq",
			args: args{
				table: "my_table",
				conditions: []Column{
					{
						Name:  "col2",
						Value: 1,
					},
				},
				sequence:         1,
				previousSequence: 1,
			},
			want: want{
				sequence:         1,
				previousSequence: 1,
				table:            "my_table",
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, ErrPrevSeqGtSeq)
				},
			},
		},
		{
			name: "sequence less prev seq",
			args: args{
				table: "my_table",
				conditions: []Column{
					{
						Name:  "col2",
						Value: 1,
					},
				},
				sequence:         1,
				previousSequence: 2,
			},
			want: want{
				sequence:         1,
				previousSequence: 2,
				table:            "my_table",
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, ErrPrevSeqGtSeq)
				},
			},
		},
		{
			name: "no conditions",
			args: args{
				table:            "my_table",
				conditions:       []Column{},
				sequence:         1,
				previousSequence: 0,
			},
			want: want{
				sequence:         1,
				previousSequence: 1,
				table:            "my_table",
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, ErrNoCondition)
				},
			},
		},
		{
			name: "correct",
			args: args{
				table: "my_table",
				conditions: []Column{
					{
						Name:  "col1",
						Value: 1,
					},
				},
				sequence:         1,
				previousSequence: 0,
			},
			want: want{
				sequence:         1,
				previousSequence: 1,
				table:            "my_table",
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
			stmt := NewDeleteStatement(tt.args.conditions, tt.args.sequence, tt.args.previousSequence)

			err := stmt.execute(tt.want.executer, tt.args.table)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}
			tt.want.executer.check(t)
		})
	}
}

func TestNewNoOpStatement(t *testing.T) {
	type args struct {
		sequence         uint64
		previousSequence uint64
	}
	tests := []struct {
		name string
		args args
		want Statement
	}{
		{
			name: "generate correctly",
			args: args{
				sequence:         5,
				previousSequence: 3,
			},
			want: Statement{
				execute:          nil,
				Sequence:         5,
				PreviousSequence: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNoOpStatement(tt.args.sequence, tt.args.previousSequence); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNoOpStatement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatement_Execute(t *testing.T) {
	type fields struct {
		execute func(ex executer, projectionName string) error
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
			name:   "no execute",
			fields: fields{},
			want: want{
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
		{
			name: "execute returns no error",
			fields: fields{
				execute: func(ex executer, projectionName string) error { return nil },
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
				execute: func(ex executer, projectionName string) error { return testErr },
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, testErr)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := &Statement{
				execute: tt.fields.execute,
			}
			if err := stmt.Execute(nil, tt.args.projectionName); !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func Test_columnsToQuery(t *testing.T) {
	type args struct {
		cols []Column
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
				cols: []Column{
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
				cols: []Column{
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
		cols        []Column
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
				cols: []Column{
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
				cols: []Column{
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
				cols: []Column{
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
