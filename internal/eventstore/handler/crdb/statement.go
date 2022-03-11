package crdb

import (
	"strconv"
	"strings"

	"github.com/lib/pq"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
)

type execOption func(*execConfig)
type execConfig struct {
	tableName string

	args []interface{}
	err  error
}

func WithTableSuffix(name string) func(*execConfig) {
	return func(o *execConfig) {
		o.tableName += "_" + name
	}
}

func NewCreateStatement(event eventstore.Event, values []handler.Column, opts ...execOption) *handler.Statement {
	cols, params, args := columnsToQuery(values)
	columnNames := strings.Join(cols, ", ")
	valuesPlaceholder := strings.Join(params, ", ")

	config := execConfig{
		args: args,
	}

	if len(values) == 0 {
		config.err = handler.ErrNoValues
	}

	q := func(config execConfig) string {
		return "INSERT INTO " + config.tableName + " (" + columnNames + ") VALUES (" + valuesPlaceholder + ")"
	}

	return &handler.Statement{
		AggregateType:    event.Aggregate().Type,
		Sequence:         event.Sequence(),
		PreviousSequence: event.PreviousAggregateTypeSequence(),
		Execute:          exec(config, q, opts),
	}
}

func NewUpsertStatement(event eventstore.Event, values []handler.Column, opts ...execOption) *handler.Statement {
	cols, params, args := columnsToQuery(values)
	columnNames := strings.Join(cols, ", ")
	valuesPlaceholder := strings.Join(params, ", ")

	config := execConfig{
		args: args,
	}

	if len(values) == 0 {
		config.err = handler.ErrNoValues
	}

	q := func(config execConfig) string {
		return "UPSERT INTO " + config.tableName + " (" + columnNames + ") VALUES (" + valuesPlaceholder + ")"
	}

	return &handler.Statement{
		AggregateType:    event.Aggregate().Type,
		Sequence:         event.Sequence(),
		PreviousSequence: event.PreviousAggregateTypeSequence(),
		Execute:          exec(config, q, opts),
	}
}

func NewUpdateStatement(event eventstore.Event, values []handler.Column, conditions []handler.Condition, opts ...execOption) *handler.Statement {
	cols, params, args := columnsToQuery(values)
	wheres, whereArgs := conditionsToWhere(conditions, len(params))
	args = append(args, whereArgs...)

	columnNames := strings.Join(cols, ", ")
	valuesPlaceholder := strings.Join(params, ", ")
	wheresPlaceholders := strings.Join(wheres, " AND ")

	config := execConfig{
		args: args,
	}

	if len(values) == 0 {
		config.err = handler.ErrNoValues
	}

	if len(conditions) == 0 {
		config.err = handler.ErrNoCondition
	}

	q := func(config execConfig) string {
		return "UPDATE " + config.tableName + " SET (" + columnNames + ") = (" + valuesPlaceholder + ") WHERE " + wheresPlaceholders
	}

	return &handler.Statement{
		AggregateType:    event.Aggregate().Type,
		Sequence:         event.Sequence(),
		PreviousSequence: event.PreviousAggregateTypeSequence(),
		Execute:          exec(config, q, opts),
	}
}

func NewDeleteStatement(event eventstore.Event, conditions []handler.Condition, opts ...execOption) *handler.Statement {
	wheres, args := conditionsToWhere(conditions, 0)

	wheresPlaceholders := strings.Join(wheres, " AND ")

	config := execConfig{
		args: args,
	}

	if len(conditions) == 0 {
		config.err = handler.ErrNoCondition
	}

	q := func(config execConfig) string {
		return "DELETE FROM " + config.tableName + " WHERE " + wheresPlaceholders
	}

	return &handler.Statement{
		AggregateType:    event.Aggregate().Type,
		Sequence:         event.Sequence(),
		PreviousSequence: event.PreviousAggregateTypeSequence(),
		Execute:          exec(config, q, opts),
	}
}

func NewNoOpStatement(event eventstore.Event) *handler.Statement {
	return &handler.Statement{
		AggregateType:    event.Aggregate().Type,
		Sequence:         event.Sequence(),
		PreviousSequence: event.PreviousAggregateTypeSequence(),
	}
}

func NewMultiStatement(event eventstore.Event, opts ...func(eventstore.Event) Exec) *handler.Statement {
	if len(opts) == 0 {
		return NewNoOpStatement(event)
	}
	execs := make([]Exec, len(opts))
	for i, opt := range opts {
		execs[i] = opt(event)
	}
	return &handler.Statement{
		AggregateType:    event.Aggregate().Type,
		Sequence:         event.Sequence(),
		PreviousSequence: event.PreviousAggregateTypeSequence(),
		Execute:          multiExec(execs),
	}
}

type Exec func(ex handler.Executer, projectionName string) error

func AddCreateStatement(columns []handler.Column, opts ...execOption) func(eventstore.Event) Exec {
	return func(event eventstore.Event) Exec {
		return NewCreateStatement(event, columns, opts...).Execute
	}
}

func AddUpsertStatement(values []handler.Column, opts ...execOption) func(eventstore.Event) Exec {
	return func(event eventstore.Event) Exec {
		return NewUpsertStatement(event, values, opts...).Execute
	}
}

func AddUpdateStatement(values []handler.Column, conditions []handler.Condition, opts ...execOption) func(eventstore.Event) Exec {
	return func(event eventstore.Event) Exec {
		return NewUpdateStatement(event, values, conditions, opts...).Execute
	}
}

func AddDeleteStatement(conditions []handler.Condition, opts ...execOption) func(eventstore.Event) Exec {
	return func(event eventstore.Event) Exec {
		return NewDeleteStatement(event, conditions, opts...).Execute
	}
}

func NewArrayAppendCol(column string, value interface{}) handler.Column {
	return handler.Column{
		Name:  column,
		Value: value,
		ParameterOpt: func(placeholder string) string {
			return "array_append(" + column + ", " + placeholder + ")"
		},
	}
}

func NewArrayRemoveCol(column string, value interface{}) handler.Column {
	return handler.Column{
		Name:  column,
		Value: value,
		ParameterOpt: func(placeholder string) string {
			return "array_remove(" + column + ", " + placeholder + ")"
		},
	}
}

func NewArrayIntersectCol(column string, value interface{}) handler.Column {
	var arrayType string
	switch value.(type) {
	case pq.StringArray:
		arrayType = "STRING"
	case pq.Int32Array,
		pq.Int64Array:
		arrayType = "INT"
		//TODO: handle more types if necessary
	}
	return handler.Column{
		Name:  column,
		Value: value,
		ParameterOpt: func(placeholder string) string {
			return "SELECT ARRAY( SELECT UNNEST(" + column + ") INTERSECT SELECT UNNEST (" + placeholder + "::" + arrayType + "[]))"
		},
	}
}

//NewCopyStatement creates a new upsert statement which updates a column from an existing row
// cols represent the columns which are objective to change.
// if the value of a col is empty the data will be copied from the selected row
// if the value of a col is not empty the data will be set by the static value
// conds represent the conditions for the selection subquery
func NewCopyStatement(event eventstore.Event, cols []handler.Column, conds []handler.Condition, opts ...execOption) *handler.Statement {
	columnNames := make([]string, len(cols))
	selectColumns := make([]string, len(cols))
	argCounter := 0
	args := []interface{}{}

	for i, col := range cols {
		columnNames[i] = col.Name
		selectColumns[i] = col.Name
		if col.Value != nil {
			argCounter++
			selectColumns[i] = "$" + strconv.Itoa(argCounter)
			args = append(args, col.Value)
		}
	}

	wheres := make([]string, len(conds))
	for i, cond := range conds {
		argCounter++
		wheres[i] = "copy_table." + cond.Name + " = $" + strconv.Itoa(argCounter)
		args = append(args, cond.Value)
	}

	config := execConfig{
		args: args,
	}

	if len(cols) == 0 {
		config.err = handler.ErrNoValues
	}

	if len(conds) == 0 {
		config.err = handler.ErrNoCondition
	}

	q := func(config execConfig) string {
		return "UPSERT INTO " +
			config.tableName +
			" (" +
			strings.Join(columnNames, ", ") +
			") SELECT " +
			strings.Join(selectColumns, ", ") +
			" FROM " +
			config.tableName + " AS copy_table WHERE " +
			strings.Join(wheres, " AND ")
	}

	return &handler.Statement{
		AggregateType:    event.Aggregate().Type,
		Sequence:         event.Sequence(),
		PreviousSequence: event.PreviousAggregateTypeSequence(),
		Execute:          exec(config, q, opts),
	}
}

func columnsToQuery(cols []handler.Column) (names []string, parameters []string, values []interface{}) {
	names = make([]string, len(cols))
	values = make([]interface{}, len(cols))
	parameters = make([]string, len(cols))
	for i, col := range cols {
		names[i] = col.Name
		values[i] = col.Value
		parameters[i] = "$" + strconv.Itoa(i+1)
		if col.ParameterOpt != nil {
			parameters[i] = col.ParameterOpt(parameters[i])
		}
	}
	return names, parameters, values
}

func conditionsToWhere(cols []handler.Condition, paramOffset int) (wheres []string, values []interface{}) {
	wheres = make([]string, len(cols))
	values = make([]interface{}, len(cols))

	for i, col := range cols {
		wheres[i] = "(" + col.Name + " = $" + strconv.Itoa(i+1+paramOffset) + ")"
		values[i] = col.Value
	}

	return wheres, values
}

type query func(config execConfig) string

func exec(config execConfig, q query, opts []execOption) Exec {
	return func(ex handler.Executer, projectionName string) error {
		if projectionName == "" {
			return handler.ErrNoProjection
		}

		if config.err != nil {
			return config.err
		}

		config.tableName = projectionName
		for _, opt := range opts {
			opt(&config)
		}

		if _, err := ex.Exec(q(config), config.args...); err != nil {
			return caos_errs.ThrowInternal(err, "CRDB-pKtsr", "exec failed")
		}

		return nil
	}
}

func multiExec(execList []Exec) Exec {
	return func(ex handler.Executer, projectionName string) error {
		for _, exec := range execList {
			if err := exec(ex, projectionName); err != nil {
				return err
			}
		}
		return nil
	}
}
