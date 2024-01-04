package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (h *Handler) eventsToStatements(tx *sql.Tx, events []eventstore.Event, currentState *state) (statements []*Statement, err error) {
	statements = make([]*Statement, 0, len(events))

	previousPosition := currentState.position
	offset := currentState.offset
	for _, event := range events {
		statement, err := h.reduce(event)
		if err != nil {
			h.logEvent(event).WithError(err).Error("reduce failed")
			if shouldContinue := h.handleFailedStmt(tx, failureFromEvent(event, err)); shouldContinue {
				continue
			}
			return statements, err
		}
		offset++
		if previousPosition != event.Position() {
			// offset is 1 because we want to skip this event
			offset = 1
		}
		statement.offset = offset
		statement.Position = event.Position()
		previousPosition = event.Position()
		statements = append(statements, statement)
	}
	return statements, nil
}

func (h *Handler) reduce(event eventstore.Event) (*Statement, error) {
	for _, reducer := range h.projection.Reducers() {
		if reducer.Aggregate != event.Aggregate().Type {
			continue
		}
		for _, reduce := range reducer.EventReducers {
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
	Position      float64
	CreationDate  time.Time
	InstanceID    string

	offset uint32

	Execute Exec
}

type Exec func(ex Executer, projectionName string) error

func WithTableSuffix(name string) func(*execConfig) {
	return func(o *execConfig) {
		o.tableName += "_" + name
	}
}

var (
	ErrNoProjection = errors.New("no projection")
	ErrNoValues     = errors.New("no values")
	ErrNoCondition  = errors.New("no condition")
)

func NewStatement(event eventstore.Event, e Exec) *Statement {
	return &Statement{
		AggregateType: event.Aggregate().Type,
		Sequence:      event.Sequence(),
		Position:      event.Position(),
		AggregateID:   event.Aggregate().ID,
		CreationDate:  event.CreatedAt(),
		InstanceID:    event.Aggregate().InstanceID,
		Execute:       e,
	}
}

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

	return NewStatement(event, exec(config, q, opts))
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

	updateCols, updateVals := getUpdateCols(values, conflictTarget)
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

	return NewStatement(event, exec(config, q, opts))
}

var _ ValueContainer = (*onlySetValueOnInsert)(nil)

type onlySetValueOnInsert struct {
	Table string
	Value interface{}
}

func (c *onlySetValueOnInsert) GetValue() interface{} {
	return c.Value
}

func OnlySetValueOnInsert(table string, value interface{}) *onlySetValueOnInsert {
	return &onlySetValueOnInsert{
		Table: table,
		Value: value,
	}
}

func getUpdateCols(cols []Column, conflictTarget []string) (updateCols, updateVals []string) {
	updateCols = make([]string, len(cols))
	updateVals = make([]string, len(cols))

	for i := len(cols) - 1; i >= 0; i-- {
		col := cols[i]
		table := "EXCLUDED"
		if onlyOnInsert, ok := col.Value.(*onlySetValueOnInsert); ok {
			table = onlyOnInsert.Table
		}
		updateCols[i] = col.Name
		updateVals[i] = table + "." + col.Name
		for _, conflict := range conflictTarget {
			if conflict == col.Name {
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
	wheres, whereArgs := conditionsToWhere(conditions, len(args)+1)
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

	return NewStatement(event, exec(config, q, opts))
}

func NewDeleteStatement(event eventstore.Event, conditions []Condition, opts ...execOption) *Statement {
	wheres, args := conditionsToWhere(conditions, 1)

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

	return NewStatement(event, exec(config, q, opts))
}

func NewNoOpStatement(event eventstore.Event) *Statement {
	return NewStatement(event, nil)
}

func NewMultiStatement(event eventstore.Event, opts ...func(eventstore.Event) Exec) *Statement {
	if len(opts) == 0 {
		return NewNoOpStatement(event)
	}
	execs := make([]Exec, len(opts))
	for i, opt := range opts {
		execs[i] = opt(event)
	}
	return NewStatement(event, multiExec(execs))
}

func AddNoOpStatement() func(eventstore.Event) Exec {
	return func(event eventstore.Event) Exec {
		return NewNoOpStatement(event).Execute
	}
}

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

func AddCopyStatement(conflict, from, to []Column, conditions []NamespacedCondition, opts ...execOption) func(eventstore.Event) Exec {
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

	case []string, database.TextArray[string]:
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
func NewCopyStatement(event eventstore.Event, conflictCols, from, to []Column, nsCond []NamespacedCondition, opts ...execOption) *Statement {
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
	cond := make([]Condition, len(nsCond))
	for i := range nsCond {
		cond[i] = nsCond[i]("copy_table")
	}
	wheres, values := conditionsToWhere(cond, len(args)+1)
	args = append(args, values...)

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

	if len(cond) == 0 {
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

	return NewStatement(event, exec(config, q, opts))
}

type ValueContainer interface {
	GetValue() interface{}
}

func columnsToQuery(cols []Column) (names []string, parameters []string, values []interface{}) {
	names = make([]string, len(cols))
	values = make([]interface{}, len(cols))
	parameters = make([]string, len(cols))
	var parameterIndex int
	for i, col := range cols {
		names[i] = col.Name
		switch c := col.Value.(type) {
		case Column:
			parameters[i] = c.Name
			continue
		case ValueContainer:
			values[parameterIndex] = c.GetValue()
		default:
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

func conditionsToWhere(conds []Condition, paramOffset int) (wheres []string, values []interface{}) {
	wheres = make([]string, len(conds))
	values = make([]any, 0, len(conds))

	for i, cond := range conds {
		var args []any
		wheres[i], args = cond("$" + strconv.Itoa(paramOffset))
		paramOffset += len(args)
		values = append(values, args...)
		wheres[i] = "(" + wheres[i] + ")"
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

type Condition func(param string) (string, []any)

type NamespacedCondition func(namespace string) Condition

func NewCond(name string, value interface{}) Condition {
	return func(param string) (string, []any) {
		return name + " = " + param, []any{value}
	}
}

func NewNamespacedCondition(name string, value interface{}) NamespacedCondition {
	return func(namespace string) Condition {
		return NewCond(namespace+"."+name, value)
	}
}

func NewLessThanCond(column string, value interface{}) Condition {
	return func(param string) (string, []any) {
		return column + " < " + param, []any{value}
	}
}

func NewIsNullCond(column string) Condition {
	return func(string) (string, []any) {
		return column + " IS NULL", nil
	}
}

// NewTextArrayContainsCond returns a Condition that checks if the column that stores an array of text contains the given value
func NewTextArrayContainsCond(column string, value string) Condition {
	return func(param string) (string, []any) {
		return column + " @> " + param, []any{database.TextArray[string]{value}}
	}
}

// Not is a function and not a method, so that calling it is well readable
// For example conditions := []Condition{ Not(NewTextArrayContainsCond())}
func Not(condition Condition) Condition {
	return func(param string) (string, []any) {
		cond, value := condition(param)
		return "NOT (" + cond + ")", value
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
	return func(ex Executer, projectionName string) (err error) {
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

		_, err = ex.Exec("SAVEPOINT stmt_exec")
		if err != nil {
			return zerrors.ThrowInternal(err, "CRDB-YdOXD", "create savepoint failed")
		}
		defer func() {
			if err != nil {
				_, rollbackErr := ex.Exec("ROLLBACK TO SAVEPOINT stmt_exec")
				logging.OnError(rollbackErr).Debug("rollback failed")
				return
			}
			_, err = ex.Exec("RELEASE SAVEPOINT stmt_exec")
		}()
		_, err = ex.Exec(q(config), config.args...)
		if err != nil {
			return zerrors.ThrowInternal(err, "CRDB-pKtsr", "exec failed")
		}

		return nil
	}
}

func multiExec(execList []Exec) Exec {
	return func(ex Executer, projectionName string) error {
		for _, exec := range execList {
			if exec == nil {
				continue
			}
			if err := exec(ex, projectionName); err != nil {
				return err
			}
		}
		return nil
	}
}
