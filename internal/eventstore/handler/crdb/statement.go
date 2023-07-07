package crdb

import (
	"strconv"
	"strings"

	"github.com/zitadel/zitadel/internal/database"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
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
		InstanceID:       event.Aggregate().InstanceID,
		Execute:          exec(config, q, opts),
	}
}

func NewUpsertStatement(event eventstore.Event, conflictCols []handler.Column, values []handler.Column, opts ...execOption) *handler.Statement {
	cols, params, args := columnsToQuery(values)

	conflictTarget := make([]string, len(conflictCols))
	for i, col := range conflictCols {
		conflictTarget[i] = col.Name
	}

	config := execConfig{
		args: args,
	}

	if len(values) == 0 {
		config.err = handler.ErrNoValues
	}

	updateCols, updateVals := getUpdateCols(cols, conflictTarget)
	if len(updateCols) == 0 || len(updateVals) == 0 {
		config.err = handler.ErrNoValues
	}

	q := func(config execConfig) string {
		var updateStmt string
		// the postgres standard does not allow to update a single column using a multi-column update
		// discussion: https://www.postgresql.org/message-id/17451.1509381766%40sss.pgh.pa.us
		// see Compatibility in https://www.postgresql.org/docs/current/sql-update.html
		if len(updateCols) == 1 && !strings.HasPrefix(updateVals[0], "SELECT") {
			updateStmt = "UPDATE SET " + updateCols[0] + " = " + updateVals[0]
		} else {
			updateStmt = "UPDATE SET (" + strings.Join(updateCols, ", ") + ") = (" + strings.Join(updateVals, ", ") + ")"
		}
		return "INSERT INTO " + config.tableName + " (" + strings.Join(cols, ", ") + ") VALUES (" + strings.Join(params, ", ") + ")" +
			" ON CONFLICT (" + strings.Join(conflictTarget, ", ") + ") DO " + updateStmt
	}

	return &handler.Statement{
		AggregateType:    event.Aggregate().Type,
		Sequence:         event.Sequence(),
		PreviousSequence: event.PreviousAggregateTypeSequence(),
		InstanceID:       event.Aggregate().InstanceID,
		Execute:          exec(config, q, opts),
	}
}

func getUpdateCols(cols, conflictTarget []string) (updateCols, updateVals []string) {
	updateCols = make([]string, len(cols))
	updateVals = make([]string, len(cols))

	copy(updateCols, cols)

	for i := len(updateCols) - 1; i >= 0; i-- {
		updateVals[i] = "EXCLUDED." + updateCols[i]

		for _, conflict := range conflictTarget {
			if conflict == updateCols[i] {
				copy(updateCols[i:], updateCols[i+1:])
				updateCols[len(updateCols)-1] = ""
				updateCols = updateCols[:len(updateCols)-1]

				copy(updateVals[i:], updateVals[i+1:])
				updateVals[len(updateVals)-1] = ""
				updateVals = updateVals[:len(updateVals)-1]

				break
			}
		}
	}

	return updateCols, updateVals
}

func NewUpdateStatement(event eventstore.Event, values []handler.Column, conditions []handler.Condition, opts ...execOption) *handler.Statement {
	cols, params, args := columnsToQuery(values)
	wheres, whereArgs := conditionsToWhere(conditions, len(args))
	args = append(args, whereArgs...)

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
		// the postgres standard does not allow to update a single column using a multi-column update
		// discussion: https://www.postgresql.org/message-id/17451.1509381766%40sss.pgh.pa.us
		// see Compatibility in https://www.postgresql.org/docs/current/sql-update.html
		if len(cols) == 1 && !strings.HasPrefix(params[0], "SELECT") {
			return "UPDATE " + config.tableName + " SET " + cols[0] + " = " + params[0] + " WHERE " + strings.Join(wheres, " AND ")
		}
		return "UPDATE " + config.tableName + " SET (" + strings.Join(cols, ", ") + ") = (" + strings.Join(params, ", ") + ") WHERE " + strings.Join(wheres, " AND ")
	}

	return &handler.Statement{
		AggregateType:    event.Aggregate().Type,
		Sequence:         event.Sequence(),
		PreviousSequence: event.PreviousAggregateTypeSequence(),
		InstanceID:       event.Aggregate().InstanceID,
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
		InstanceID:       event.Aggregate().InstanceID,
		Execute:          exec(config, q, opts),
	}
}

func NewNoOpStatement(event eventstore.Event) *handler.Statement {
	return &handler.Statement{
		AggregateType:    event.Aggregate().Type,
		Sequence:         event.Sequence(),
		PreviousSequence: event.PreviousAggregateTypeSequence(),
		InstanceID:       event.Aggregate().InstanceID,
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
		InstanceID:       event.Aggregate().InstanceID,
		Execute:          multiExec(execs),
	}
}

type Exec func(ex handler.Executer, projectionName string) error

func AddCreateStatement(columns []handler.Column, opts ...execOption) func(eventstore.Event) Exec {
	return func(event eventstore.Event) Exec {
		return NewCreateStatement(event, columns, opts...).Execute
	}
}

func AddUpsertStatement(indexCols []handler.Column, values []handler.Column, opts ...execOption) func(eventstore.Event) Exec {
	return func(event eventstore.Event) Exec {
		return NewUpsertStatement(event, indexCols, values, opts...).Execute
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

	case []string, database.StringArray:
		arrayType = "TEXT"
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

func NewCopyCol(column, from string) handler.Column {
	return handler.Column{
		Name:  column,
		Value: handler.NewCol(from, nil),
	}
}

func NewLessThanCond(column string, value interface{}) handler.Condition {
	return func(param string) (string, interface{}) {
		return column + " < " + param, value
	}
}

func NewIsNullCond(column string) handler.Condition {
	return func(param string) (string, interface{}) {
		return column + " IS NULL", nil
	}
}

// NewTextArrayContainsCond returns a handler.Condition that checks if the column that stores an array of text contains the given value
func NewTextArrayContainsCond(column string, value string) handler.Condition {
	return func(param string) (string, interface{}) {
		return column + " @> " + param, database.StringArray{value}
	}
}

// Not is a function and not a method, so that calling it is well readable
// For example conditions := []handler.Condition{ Not(NewTextArrayContainsCond())}
func Not(condition handler.Condition) handler.Condition {
	return func(param string) (string, interface{}) {
		cond, value := condition(param)
		return "NOT (" + cond + ")", value
	}
}

// NewCopyStatement creates a new upsert statement which updates a column from an existing row
// cols represent the columns which are objective to change.
// if the value of a col is empty the data will be copied from the selected row
// if the value of a col is not empty the data will be set by the static value
// conds represent the conditions for the selection subquery
func NewCopyStatement(event eventstore.Event, conflictCols, from, to []handler.Column, nsCond []handler.NamespacedCondition, opts ...execOption) *handler.Statement {
	columnNames := make([]string, len(to))
	selectColumns := make([]string, len(from))
	updateColumns := make([]string, len(columnNames))
	argCounter := 0
	args := []interface{}{}

	for i, col := range from {
		columnNames[i] = to[i].Name
		selectColumns[i] = from[i].Name
		updateColumns[i] = "EXCLUDED." + col.Name
		if col.Value != nil {
			argCounter++
			selectColumns[i] = "$" + strconv.Itoa(argCounter)
			updateColumns[i] = selectColumns[i]
			args = append(args, col.Value)
		}

	}
	cond := make([]handler.Condition, len(nsCond))
	for i := range nsCond {
		cond[i] = nsCond[i]("copy_table")
	}
	wheres, values := conditionsToWhere(cond, len(args))
	args = append(args, values...)

	conflictTargets := make([]string, len(conflictCols))
	for i, conflictCol := range conflictCols {
		conflictTargets[i] = conflictCol.Name
	}

	config := execConfig{
		args: args,
	}

	if len(from) == 0 || len(to) == 0 || len(from) != len(to) {
		config.err = handler.ErrNoValues
	}

	if len(cond) == 0 {
		config.err = handler.ErrNoCondition
	}

	q := func(config execConfig) string {
		return "INSERT INTO " +
			config.tableName +
			" (" +
			strings.Join(columnNames, ", ") +
			") SELECT " +
			strings.Join(selectColumns, ", ") +
			" FROM " +
			config.tableName + " AS copy_table WHERE " +
			strings.Join(wheres, " AND ") +
			" ON CONFLICT (" +
			strings.Join(conflictTargets, ", ") +
			") DO UPDATE SET (" +
			strings.Join(columnNames, ", ") +
			") = (" +
			strings.Join(updateColumns, ", ") +
			")"
	}

	return &handler.Statement{
		AggregateType:    event.Aggregate().Type,
		Sequence:         event.Sequence(),
		PreviousSequence: event.PreviousAggregateTypeSequence(),
		InstanceID:       event.Aggregate().InstanceID,
		Execute:          exec(config, q, opts),
	}
}

func columnsToQuery(cols []handler.Column) (names []string, parameters []string, values []interface{}) {
	names = make([]string, len(cols))
	values = make([]interface{}, len(cols))
	parameters = make([]string, len(cols))
	var parameterIndex int
	for i, col := range cols {
		names[i] = col.Name
		if c, ok := col.Value.(handler.Column); ok {
			parameters[i] = c.Name
			continue
		} else {
			values[parameterIndex] = col.Value
		}
		parameters[i] = "$" + strconv.Itoa(parameterIndex+1)
		if col.ParameterOpt != nil {
			parameters[i] = col.ParameterOpt(parameters[i])
		}
		parameterIndex++
	}
	return names, parameters, values[:parameterIndex]
}

func conditionsToWhere(conditions []handler.Condition, paramOffset int) (wheres []string, values []interface{}) {
	wheres = make([]string, len(conditions))
	values = make([]interface{}, 0, len(conditions))
	for i, conditionFunc := range conditions {
		condition, value := conditionFunc("$" + strconv.Itoa(i+1+paramOffset))
		wheres[i] = "(" + condition + ")"
		if value != nil {
			values = append(values, value)
		}
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
