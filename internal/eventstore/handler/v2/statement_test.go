package handler

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type wantExecuter struct {
	params        []params
	i             int
	t             *testing.T
	wasExecuted   bool
	shouldExecute bool
}

type params struct {
	query string
	args  []interface{}
}

var errTest = errors.New("some error")

var _ eventstore.Event = &testEvent{}

type testEvent struct {
	eventstore.BaseEvent
	sequence         uint64
	previousSequence uint64
	aggregateType    eventstore.AggregateType
	instanceID       string
}

func (e *testEvent) Sequence() uint64 {
	return e.sequence
}

func (e *testEvent) Aggregate() *eventstore.Aggregate {
	return &eventstore.Aggregate{
		Type:       e.aggregateType,
		InstanceID: e.instanceID,
	}
}

func (e *testEvent) PreviousAggregateTypeSequence() uint64 {
	return e.previousSequence
}

func (ex *wantExecuter) check(t *testing.T) {
	t.Helper()
	switch {
	case ex.wasExecuted && !ex.shouldExecute:
		t.Error("executer should not be executed")
	case !ex.wasExecuted && ex.shouldExecute:
		t.Error("executer should be executed")
	case ex.wasExecuted != ex.shouldExecute:
		t.Errorf("executed missmatched should be %t, but was %t", ex.shouldExecute, ex.wasExecuted)
	}
}

func (ex *wantExecuter) Exec(query string, args ...interface{}) (sql.Result, error) {
	ex.t.Helper()
	if strings.Contains(query, "SAVEPOINT") {
		return nil, nil
	}
	ex.wasExecuted = true
	if ex.i >= len(ex.params) {
		ex.t.Errorf("did not expect more exec, but got:\n    %q with %q", query, args)
		return nil, nil
	}
	p := ex.params[ex.i]
	if query != p.query {
		ex.t.Errorf("wrong query:\n  expected:\n    %q\n  got:\n    %q", p.query, query)
	}
	if !reflect.DeepEqual(p.args, args) {
		ex.t.Errorf("wrong args:\n  expected:\n    %v\n  got:\n    %v", p.args, args)
	}
	ex.i++
	return nil, nil
}

func TestNewCreateStatement(t *testing.T) {
	type args struct {
		table  string
		event  *testEvent
		values []Column
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
				values: []Column{
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
					return errors.Is(err, ErrNoProjection)
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
				values: []Column{},
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
					return errors.Is(err, ErrNoValues)
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
				values: []Column{
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
					params: []params{
						{
							query: "INSERT INTO my_table (col1) VALUES ($1)",
							args:  []interface{}{"val"},
						},
					},
					shouldExecute: true,
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
		table        string
		event        *testEvent
		conflictCols []Column
		values       []Column
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
				values: []Column{
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
					return errors.Is(err, ErrNoProjection)
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
				values: []Column{},
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
					return errors.Is(err, ErrNoValues)
				},
			},
		},
		{
			name: "no update cols",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				conflictCols: []Column{
					NewCol("col1", nil),
				},
				values: []Column{
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
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return errors.Is(err, ErrNoValues)
				},
			},
		},
		{
			name: "correct UPDATE multi col",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				conflictCols: []Column{
					NewCol("col1", nil),
				},
				values: []Column{
					{
						Name:  "col1",
						Value: "val",
					},
					{
						Name:  "col2",
						Value: "val",
					},
					{
						Name:  "col3",
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
					params: []params{
						{
							query: "INSERT INTO my_table (col1, col2, col3) VALUES ($1, $2, $3) ON CONFLICT (col1) DO UPDATE SET (col2, col3) = (EXCLUDED.col2, EXCLUDED.col3)",
							args:  []interface{}{"val", "val", "val"},
						},
					},
					shouldExecute: true,
				},
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
		{
			name: "correct UPDATE single col",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				conflictCols: []Column{
					NewCol("col1", nil),
				},
				values: []Column{
					{
						Name:  "col1",
						Value: "val",
					},
					{
						Name:  "col2",
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
					params: []params{
						{
							query: "INSERT INTO my_table (col1, col2) VALUES ($1, $2) ON CONFLICT (col1) DO UPDATE SET col2 = EXCLUDED.col2",
							args:  []interface{}{"val", "val"},
						},
					},
					shouldExecute: true,
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
			stmt := NewUpsertStatement(tt.args.event, tt.args.conflictCols, tt.args.values)

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
		values     []Column
		conditions []Condition
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
				values: []Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
				conditions: []Condition{
					NewCond("col2", 1),
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
					return errors.Is(err, ErrNoProjection)
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
				values: []Column{},
				conditions: []Condition{
					NewCond("col2", 1),
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
					return errors.Is(err, ErrNoValues)
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
				values: []Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
				conditions: []Condition{},
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
					return errors.Is(err, ErrNoCondition)
				},
			},
		},
		{
			name: "correct single column",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				values: []Column{
					{
						Name:  "col1",
						Value: "val",
					},
				},
				conditions: []Condition{
					NewCond("col2", 1),
				},
			},
			want: want{
				table:            "my_table",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 1,
				executer: &wantExecuter{
					params: []params{
						{
							query: "UPDATE my_table SET col1 = $1 WHERE (col2 = $2)",
							args:  []interface{}{"val", 1},
						},
					},
					shouldExecute: true,
				},
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
		{
			name: "correct multi column",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				values: []Column{
					{
						Name:  "col1",
						Value: "val",
					},
					{
						Name:  "col3",
						Value: "val5",
					},
				},
				conditions: []Condition{
					NewCond("col2", 1),
				},
			},
			want: want{
				table:            "my_table",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 1,
				executer: &wantExecuter{
					params: []params{
						{
							query: "UPDATE my_table SET (col1, col3) = ($1, $2) WHERE (col2 = $3)",
							args:  []interface{}{"val", "val5", 1},
						},
					},
					shouldExecute: true,
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
		conditions []Condition
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
				conditions: []Condition{
					NewCond("col2", 1),
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
					return errors.Is(err, ErrNoProjection)
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
				conditions: []Condition{},
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
					return errors.Is(err, ErrNoCondition)
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
				conditions: []Condition{
					NewCond("col1", 1),
				},
			},
			want: want{
				table:            "my_table",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 1,
				executer: &wantExecuter{
					params: []params{
						{
							query: "DELETE FROM my_table WHERE (col1 = $1)",
							args:  []interface{}{1},
						},
					},
					shouldExecute: true,
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
		table string
	}
	type want struct {
		executer *wantExecuter
		isErr    func(error) bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "generate correctly",
			args: args{
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         5,
					previousSequence: 3,
					instanceID:       "instanceID",
				},
			},
			want: want{
				executer: nil,
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := NewNoOpStatement(tt.args.event)
			if tt.want.executer != nil && stmt.Execute == nil {
				t.Error("expected executer, but was nil")
			}
			if stmt.Execute == nil {
				return
			}
			tt.want.executer.t = t
			err := stmt.Execute(tt.want.executer, tt.args.table)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}
			tt.want.executer.check(t)
		})
	}
}

func TestNewMultiStatement(t *testing.T) {
	type args struct {
		table string
		event *testEvent
		execs []func(eventstore.Event) Exec
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
			name: "no op",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				execs: []func(eventstore.Event) Exec{
					AddNoOpStatement(),
				},
			},
			want: want{
				executer: &wantExecuter{
					shouldExecute: false,
				},
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
		{
			name: "no condition",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				execs: []func(eventstore.Event) Exec{
					AddDeleteStatement(
						[]Condition{},
					),
					AddCreateStatement(
						[]Column{
							{
								Name:  "col1",
								Value: 1,
							},
						}),
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
					return errors.Is(err, ErrNoCondition)
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
				execs: []func(eventstore.Event) Exec{
					AddDeleteStatement(
						[]Condition{
							NewCond("col1", 1),
						}),
					AddCreateStatement(
						[]Column{
							{
								Name:  "col1",
								Value: 1,
							},
						}),
					AddUpsertStatement(
						[]Column{
							NewCol("col1", nil),
						},
						[]Column{
							{
								Name:  "col1",
								Value: 1,
							},
							{
								Name:  "col2",
								Value: 2,
							},
						}),
					AddUpdateStatement(
						[]Column{
							{
								Name:  "col1",
								Value: 1,
							},
						},
						[]Condition{
							NewCond("col1", 1),
						}),
				},
			},
			want: want{
				table:            "my_table",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 1,
				executer: &wantExecuter{
					params: []params{
						{
							query: "DELETE FROM my_table WHERE (col1 = $1)",
							args:  []interface{}{1},
						},
						{
							query: "INSERT INTO my_table (col1) VALUES ($1)",
							args:  []interface{}{1},
						},
						{
							query: "INSERT INTO my_table (col1, col2) VALUES ($1, $2) ON CONFLICT (col1) DO UPDATE SET col2 = EXCLUDED.col2",
							args:  []interface{}{1, 2},
						},
						{
							query: "UPDATE my_table SET col1 = $1 WHERE (col1 = $2)",
							args:  []interface{}{1, 1},
						},
					},
					shouldExecute: true,
				},
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := NewMultiStatement(tt.args.event, tt.args.execs...)

			if tt.want.executer != nil && stmt.Execute == nil {
				t.Error("expected executer, but was nil")
			}
			if stmt.Execute == nil {
				return
			}
			tt.want.executer.t = t
			err := stmt.Execute(tt.want.executer, tt.args.table)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}
			tt.want.executer.check(t)
		})
	}
}

func TestNewCopyStatement(t *testing.T) {
	type args struct {
		table           string
		event           *testEvent
		conflictingCols []Column
		from            []Column
		to              []Column
		conds           []NamespacedCondition
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
				conds: []NamespacedCondition{
					NewNamespacedCondition("col2", 1),
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
					return errors.Is(err, ErrNoProjection)
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
				conds: []NamespacedCondition{},
				from: []Column{
					{
						Name: "col",
					},
				},
				to: []Column{
					{
						Name: "col",
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
					return errors.Is(err, ErrNoCondition)
				},
			},
		},
		{
			name: "more to than from cols",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				conds: []NamespacedCondition{},
				from: []Column{
					{
						Name: "col",
					},
				},
				to: []Column{
					{
						Name: "col",
					},
					{
						Name: "col2",
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
					return errors.Is(err, ErrNoCondition)
				},
			},
		},
		{
			name: "no columns",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				conds: []NamespacedCondition{
					NewNamespacedCondition("col2", nil),
				},
				from: []Column{},
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
					return errors.Is(err, ErrNoValues)
				},
			},
		},
		{
			name: "correct same column names",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				from: []Column{
					{
						Name:  "state",
						Value: 1,
					},
					{
						Name: "id",
					},
					{
						Name: "col_a",
					},
					{
						Name: "col_b",
					},
				},
				to: []Column{
					{
						Name: "state",
					},
					{
						Name: "id",
					},
					{
						Name: "col_a",
					},
					{
						Name: "col_b",
					},
				},
				conds: []NamespacedCondition{
					NewNamespacedCondition("id", 2),
					NewNamespacedCondition("state", 3),
				},
			},
			want: want{
				table:            "my_table",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 1,
				executer: &wantExecuter{
					params: []params{
						{
							query: "INSERT INTO my_table (state, id, col_a, col_b) SELECT $1, id, col_a, col_b FROM my_table AS copy_table WHERE (copy_table.id = $2) AND (copy_table.state = $3) ON CONFLICT () DO UPDATE SET (state, id, col_a, col_b) = ($1, EXCLUDED.id, EXCLUDED.col_a, EXCLUDED.col_b)",
							args:  []interface{}{1, 2, 3},
						},
					},
					shouldExecute: true,
				},
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
		{
			name: "correct different column names",
			args: args{
				table: "my_table",
				event: &testEvent{
					aggregateType:    "agg",
					sequence:         1,
					previousSequence: 0,
				},
				from: []Column{
					{
						Value: 1,
					},
					{
						Name: "id",
					},
					{
						Name: "col_a",
					},
					{
						Name: "col_b",
					},
				},
				to: []Column{
					{
						Name: "state",
					},
					{
						Name: "id",
					},
					{
						Name: "col_c",
					},
					{
						Name: "col_d",
					},
				},
				conds: []NamespacedCondition{
					NewNamespacedCondition("id", 2),
					NewNamespacedCondition("state", 3),
				},
			},
			want: want{
				table:            "my_table",
				aggregateType:    "agg",
				sequence:         1,
				previousSequence: 1,
				executer: &wantExecuter{
					params: []params{
						{
							query: "INSERT INTO my_table (state, id, col_c, col_d) SELECT $1, id, col_a, col_b FROM my_table AS copy_table WHERE (copy_table.id = $2) AND (copy_table.state = $3) ON CONFLICT () DO UPDATE SET (state, id, col_c, col_d) = ($1, EXCLUDED.id, EXCLUDED.col_a, EXCLUDED.col_b)",
							args:  []interface{}{1, 2, 3},
						},
					},
					shouldExecute: true,
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
			stmt := NewCopyStatement(tt.args.event, tt.args.conflictingCols, tt.args.from, tt.args.to, tt.args.conds)

			err := stmt.Execute(tt.want.executer, tt.args.table)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}
			tt.want.executer.check(t)
		})
	}
}

func TestStatement_Execute(t *testing.T) {
	type fields struct {
		execute func(ex Executer, projectionName string) error
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
				execute: func(ex Executer, projectionName string) error { return nil },
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
				execute: func(ex Executer, projectionName string) error { return errTest },
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, errTest)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := &Statement{
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
		{
			name: "with copy column",
			args: args{
				cols: []Column{
					{
						Name:  "col1",
						Value: 1,
					},
					{
						Name: "col2",
						Value: Column{
							Name: "col1",
						},
					},
					{
						Name:  "col3",
						Value: "something",
					},
				},
			},
			want: want{
				names:  []string{"col1", "col2", "col3"},
				params: []string{"$1", "col1", "$2"},
				values: []interface{}{1, "something"},
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
		conds       []Condition
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
				conds: []Condition{
					NewCond("col1", "val1"),
				},
				paramOffset: 1,
			},
			want: want{
				wheres: []string{"(col1 = $1)"},
				values: []interface{}{"val1"},
			},
		},
		{
			name: "multiple cols",
			args: args{
				conds: []Condition{
					NewCond("col1", "val1"),
					NewCond("col2", "val2"),
				},
				paramOffset: 1,
			},
			want: want{
				wheres: []string{"(col1 = $1)", "(col2 = $2)"},
				values: []interface{}{"val1", "val2"},
			},
		},
		{
			name: "2 offset",
			args: args{
				conds: []Condition{
					NewCond("col1", "val1"),
				},
				paramOffset: 3,
			},
			want: want{
				wheres: []string{"(col1 = $3)"},
				values: []interface{}{"val1"},
			},
		},
		{
			name: "less than",
			args: args{
				conds: []Condition{
					NewLessThanCond("col1", "val1"),
				},
				paramOffset: 1,
			},
			want: want{
				wheres: []string{"(col1 < $1)"},
				values: []interface{}{"val1"},
			},
		},
		{
			name: "is null",
			args: args{
				conds: []Condition{
					NewIsNullCond("col1"),
				},
			},
			want: want{
				wheres: []string{"(col1 IS NULL)"},
				values: []interface{}{},
			},
		},
		{
			name: "text array contains",
			args: args{
				conds: []Condition{
					NewTextArrayContainsCond("col1", "val1"),
				},
				paramOffset: 1,
			},
			want: want{
				wheres: []string{"(col1 @> $1)"},
				values: []interface{}{database.TextArray[string]{"val1"}},
			},
		},
		{
			name: "not",
			args: args{
				conds: []Condition{
					Not(NewCond("col1", "val1")),
				},
				paramOffset: 1,
			},
			want: want{
				wheres: []string{"(NOT (col1 = $1))"},
				values: []interface{}{"val1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWheres, gotValues := conditionsToWhere(tt.args.conds, tt.args.paramOffset)
			if !reflect.DeepEqual(gotWheres, tt.want.wheres) {
				t.Errorf("columnsToWhere() gotWheres = %v, want %v", gotWheres, tt.want.wheres)
			}
			if !reflect.DeepEqual(gotValues, tt.want.values) {
				t.Errorf("columnsToWhere() gotValues = %v, want %v", gotValues, tt.want.values)
			}
		})
	}
}

func TestParameterOpts(t *testing.T) {
	type args struct {
		column      string
		value       interface{}
		placeholder string
	}
	tests := []struct {
		name        string
		args        args
		constructor func(column string, value interface{}) Column
		want        string
	}{
		{
			name: "NewArrayAppendCol",
			args: args{
				column:      "testCol",
				value:       "val",
				placeholder: "$1",
			},
			constructor: NewArrayAppendCol,
			want:        "array_append(testCol, $1)",
		},
		{
			name: "NewArrayRemoveCol",
			args: args{
				column:      "testCol",
				value:       "val",
				placeholder: "$1",
			},
			constructor: NewArrayRemoveCol,
			want:        "array_remove(testCol, $1)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col := tt.constructor(tt.args.column, tt.args.value)
			if param := col.ParameterOpt(tt.args.placeholder); param != tt.want {
				t.Errorf("constructor() = %v, want %v", param, tt.want)
			}
		})
	}
}

// func TestHandler_reduce(t *testing.T) {
// 	type fields struct {
// 		projection Projection
// 	}
// 	type args struct {
// 		event eventstore.Event
// 	}
// 	tests := []struct {
// 		name           string
// 		fields         fields
// 		args           args
// 		isErr          func(t *testing.T, err error)
// 		shouldBeCalled bool
// 	}{
// 		{
// 			name: "",
// 			fields: fields{
// 				projection: &projection{
// 					reducers: []AggregateReducer{
// 						{
// 							Aggregate: "aggregate",
// 							EventRedusers: []EventReducer{
// 								{
// 									Event: "event",
// 									Reduce: (&mockEventReducer{
// 										statement: new(Statement),
// 									}).reduce,
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		if tt.isErr == nil {
// 			tt.isErr = func(t *testing.T, err error) {
// 				if err != nil {
// 					t.Error("expected no error got:", err)
// 				}
// 			}
// 		}
// 		t.Run(tt.name, func(t *testing.T) {
// 			h := &Handler{
// 				projection: tt.fields.projection,
// 			}
// 			got, err := h.reduce(tt.args.event)
// 			tt.isErr(t, err)
// 			if tt.shouldBeCalled != tt.
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Handler.reduce() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
