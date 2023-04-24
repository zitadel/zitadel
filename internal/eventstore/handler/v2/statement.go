package handler

import (
	"database/sql"
	"encoding/json"
	errs "errors"
	"strconv"
	"strings"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

func (h *Handler) eventsToStatements(tx *sql.Tx, events []eventstore.Event, currentState *state) (statements []*Statement, err error) {
	statements = make([]*Statement, len(events))
	for i, event := range events {
		statements[i], err = h.reduce(event)
		if err != nil {
			h.logEvent(event).WithError(err).Debug("reduce failed")
			if h.handleFailedStmt(tx, currentState, failureFromEvent(event, err)) {
				statements[i] = NewNoOpStatement(event)
				continue
			}
			return nil, err
		}
	}
	return statements, nil
}

func (h *Handler) reduce(event eventstore.Event) (*Statement, error) {
	for _, reducer := range h.projection.Reducers() {
		if reducer.Aggregate != event.Aggregate().Type {
			continue
		}
		for _, reduce := range reducer.EventRedusers {
			if reduce.Event != event.Type() {
				continue
			}
			return reduce.Reduce(event)
		}
	}
	return NewNoOpStatement(event), nil
}

type Statement struct {
	AggregateType eventstore.AggregateType
	AggregateID   string
	Sequence      uint64
	CreationDate  time.Time
	InstanceID    string

	Execute func(ex Executer, projectionName string) error
}

func WithTableSuffix(name string) func(*execConfig) {
	return func(o *execConfig) {
		o.tableName += "_" + name
	}
}

var (
	ErrNoProjection = errs.New("no projection")
	ErrNoValues     = errs.New("no values")
	ErrNoCondition  = errs.New("no condition")
)

func NewCreateStatement(event eventstore.Event, values []Column, opts ...execOption) *Statement {
	cols, params, args := columnsToQuery(values)
	columnNames := strings.Join(cols, ", ")
	valuesPlaceholder := strings.Join(params, ", ")

	config := execConfig{
		args: args,
	}

	if len(values) == 0 {
		config.err = ErrNoValues
	}

	q := func(config execConfig) string {
		return "INSERT INTO " + config.tableName + " (" + columnNames + ") VALUES (" + valuesPlaceholder + ")"
	}

	return &Statement{
		AggregateType: event.Aggregate().Type,
		Sequence:      event.Sequence(),
		AggregateID:   event.Aggregate().ID,
		CreationDate:  event.CreationDate(),
		InstanceID:    event.Aggregate().InstanceID,
		Execute:       exec(config, q, opts),
	}
}

func NewUpsertStatement(event eventstore.Event, conflictCols []Column, values []Column, opts ...execOption) *Statement {
	cols, params, args := columnsToQuery(values)

	conflictTarget := make([]string, len(conflictCols))
	for i, col := range conflictCols {
		conflictTarget[i] = col.Name
	}

	config := execConfig{
		args: args,
	}

	if len(values) == 0 {
		config.err = ErrNoValues
	}

	updateCols, updateVals := getUpdateCols(cols, conflictTarget)
	if len(updateCols) == 0 || len(updateVals) == 0 {
		config.err = ErrNoValues
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

	return &Statement{
		AggregateType: event.Aggregate().Type,
		Sequence:      event.Sequence(),
		AggregateID:   event.Aggregate().ID,
		CreationDate:  event.CreationDate(),
		InstanceID:    event.Aggregate().InstanceID,
		Execute:       exec(config, q, opts),
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

func NewUpdateStatement(event eventstore.Event, values []Column, conditions []Condition, opts ...execOption) *Statement {
	cols, params, args := columnsToQuery(values)
	wheres, whereArgs := conditionsToWhere(conditions, len(args))
	args = append(args, whereArgs...)

	config := execConfig{
		args: args,
	}

	if len(values) == 0 {
		config.err = ErrNoValues
	}

	if len(conditions) == 0 {
		config.err = ErrNoCondition
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

	return &Statement{
		AggregateType: event.Aggregate().Type,
		Sequence:      event.Sequence(),
		AggregateID:   event.Aggregate().ID,
		CreationDate:  event.CreationDate(),
		InstanceID:    event.Aggregate().InstanceID,
		Execute:       exec(config, q, opts),
	}
}

func NewDeleteStatement(event eventstore.Event, conditions []Condition, opts ...execOption) *Statement {
	wheres, args := conditionsToWhere(conditions, 0)

	wheresPlaceholders := strings.Join(wheres, " AND ")

	config := execConfig{
		args: args,
	}

	if len(conditions) == 0 {
		config.err = ErrNoCondition
	}

	q := func(config execConfig) string {
		return "DELETE FROM " + config.tableName + " WHERE " + wheresPlaceholders
	}

	return &Statement{
		AggregateType: event.Aggregate().Type,
		Sequence:      event.Sequence(),
		AggregateID:   event.Aggregate().ID,
		CreationDate:  event.CreationDate(),
		InstanceID:    event.Aggregate().InstanceID,
		Execute:       exec(config, q, opts),
	}
}

func NewNoOpStatement(event eventstore.Event) *Statement {
	return &Statement{
		AggregateType: event.Aggregate().Type,
		Sequence:      event.Sequence(),
		AggregateID:   event.Aggregate().ID,
		CreationDate:  event.CreationDate(),
		InstanceID:    event.Aggregate().InstanceID,
		Execute:       func(Executer, string) error { return nil },
	}
}

func NewMultiStatement(event eventstore.Event, opts ...func(eventstore.Event) Exec) *Statement {
	if len(opts) == 0 {
		return NewNoOpStatement(event)
	}
	execs := make([]Exec, len(opts))
	for i, opt := range opts {
		execs[i] = opt(event)
	}
	return &Statement{
		AggregateType: event.Aggregate().Type,
		Sequence:      event.Sequence(),
		AggregateID:   event.Aggregate().ID,
		CreationDate:  event.CreationDate(),
		InstanceID:    event.Aggregate().InstanceID,
		Execute:       multiExec(execs),
	}
}

type Exec func(ex Executer, projectionName string) error

func AddCreateStatement(columns []Column, opts ...execOption) func(eventstore.Event) Exec {
	return func(event eventstore.Event) Exec {
		return NewCreateStatement(event, columns, opts...).Execute
	}
}

func AddUpsertStatement(indexCols []Column, values []Column, opts ...execOption) func(eventstore.Event) Exec {
	return func(event eventstore.Event) Exec {
		return NewUpsertStatement(event, indexCols, values, opts...).Execute
	}
}

func AddUpdateStatement(values []Column, conditions []Condition, opts ...execOption) func(eventstore.Event) Exec {
	return func(event eventstore.Event) Exec {
		return NewUpdateStatement(event, values, conditions, opts...).Execute
	}
}

func AddDeleteStatement(conditions []Condition, opts ...execOption) func(eventstore.Event) Exec {
	return func(event eventstore.Event) Exec {
		return NewDeleteStatement(event, conditions, opts...).Execute
	}
}

func AddCopyStatement(conflict, from, to []Column, conditions []Condition, opts ...execOption) func(eventstore.Event) Exec {
	return func(event eventstore.Event) Exec {
		return NewCopyStatement(event, conflict, from, to, conditions, opts...).Execute
	}
}

func NewArrayAppendCol(column string, value interface{}) Column {
	return Column{
		Name:  column,
		Value: value,
		ParameterOpt: func(placeholder string) string {
			return "array_append(" + column + ", " + placeholder + ")"
		},
	}
}

func NewArrayRemoveCol(column string, value interface{}) Column {
	return Column{
		Name:  column,
		Value: value,
		ParameterOpt: func(placeholder string) string {
			return "array_remove(" + column + ", " + placeholder + ")"
		},
	}
}

func NewArrayIntersectCol(column string, value interface{}) Column {
	var arrayType string
	switch value.(type) {

	case []string, database.StringArray:
		arrayType = "TEXT"
		//TODO: handle more types if necessary
	}
	return Column{
		Name:  column,
		Value: value,
		ParameterOpt: func(placeholder string) string {
			return "SELECT ARRAY( SELECT UNNEST(" + column + ") INTERSECT SELECT UNNEST (" + placeholder + "::" + arrayType + "[]))"
		},
	}
}

func NewCopyCol(column, from string) Column {
	return Column{
		Name:  column,
		Value: NewCol(from, nil),
	}
}

// NewCopyStatement creates a new upsert statement which updates a column from an existing row
// cols represent the columns which are objective to change.
// if the value of a col is empty the data will be copied from the selected row
// if the value of a col is not empty the data will be set by the static value
// conds represent the conditions for the selection subquery
func NewCopyStatement(event eventstore.Event, conflictCols, from, to []Column, conds []Condition, opts ...execOption) *Statement {
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

	wheres := make([]string, len(conds))
	for i, cond := range conds {
		argCounter++
		wheres[i] = "copy_table." + cond.Name + " = $" + strconv.Itoa(argCounter)
		args = append(args, cond.Value)
	}

	conflictTargets := make([]string, len(conflictCols))
	for i, conflictCol := range conflictCols {
		conflictTargets[i] = conflictCol.Name
	}

	config := execConfig{
		args: args,
	}

	if len(from) == 0 || len(to) == 0 || len(from) != len(to) {
		config.err = ErrNoValues
	}

	if len(conds) == 0 {
		config.err = ErrNoCondition
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

	return &Statement{
		AggregateType: event.Aggregate().Type,
		Sequence:      event.Sequence(),
		AggregateID:   event.Aggregate().ID,
		CreationDate:  event.CreationDate(),
		InstanceID:    event.Aggregate().InstanceID,
		Execute:       exec(config, q, opts),
	}
}

func columnsToQuery(cols []Column) (names []string, parameters []string, values []interface{}) {
	names = make([]string, len(cols))
	values = make([]interface{}, len(cols))
	parameters = make([]string, len(cols))
	var parameterIndex int
	for i, col := range cols {
		names[i] = col.Name
		if c, ok := col.Value.(Column); ok {
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

func conditionsToWhere(cols []Condition, paramOffset int) (wheres []string, values []interface{}) {
	wheres = make([]string, len(cols))
	values = make([]interface{}, len(cols))

	for i, col := range cols {
		wheres[i] = "(" + col.Name + " = $" + strconv.Itoa(i+1+paramOffset) + ")"
		values[i] = col.Value
	}

	return wheres, values
}

type Column struct {
	Name         string
	Value        interface{}
	ParameterOpt func(string) string
}

func NewCol(name string, value interface{}) Column {
	return Column{
		Name:  name,
		Value: value,
	}
}

func NewJSONCol(name string, value interface{}) Column {
	marshalled, err := json.Marshal(value)
	if err != nil {
		logging.WithFields("column", name).WithError(err).Panic("unable to marshal column")
	}

	return NewCol(name, marshalled)
}

type Condition Column

func NewCond(name string, value interface{}) Condition {
	return Condition{
		Name:  name,
		Value: value,
	}
}

type Executer interface {
	Exec(string, ...interface{}) (sql.Result, error)
}

type execOption func(*execConfig)
type execConfig struct {
	tableName string

	args []interface{}
	err  error
}

type query func(config execConfig) string

func exec(config execConfig, q query, opts []execOption) Exec {
	return func(ex Executer, projectionName string) error {
		if projectionName == "" {
			return ErrNoProjection
		}

		if config.err != nil {
			return config.err
		}

		config.tableName = projectionName
		for _, opt := range opts {
			opt(&config)
		}

		if _, err := ex.Exec(q(config), config.args...); err != nil {
			return errors.ThrowInternal(err, "CRDB-pKtsr", "exec failed")
		}

		return nil
	}
}

func multiExec(execList []Exec) Exec {
	return func(ex Executer, projectionName string) error {
		for _, exec := range execList {
			if err := exec(ex, projectionName); err != nil {
				return err
			}
		}
		return nil
	}
}
