package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Statement struct {
	Sequence         uint64
	PreviousSequence uint64

	execute func(ex executer, projectionName string) error
}

type executer interface {
	Exec(string, ...interface{}) (sql.Result, error)
}

var (
	ErrNoTable      = errors.New("no table")
	ErrPrevSeqGtSeq = errors.New("prev seq >= seq")
	ErrNoValues     = errors.New("no values")
	ErrNoCondition  = errors.New("no condition")
)

func NewCreateStatement(values []Column, sequence, previousSequence uint64) Statement {
	cols, params, args := columnsToQuery(values)
	columnNames := strings.Join(cols, ", ")
	valuesPlaceholder := strings.Join(params, ", ")

	return Statement{
		Sequence:         sequence,
		PreviousSequence: previousSequence,
		execute: func(ex executer, projectionName string) error {
			if projectionName == "" {
				return ErrNoTable
			}
			if previousSequence >= sequence {
				return ErrPrevSeqGtSeq
			}
			if len(values) == 0 {
				return ErrNoValues
			}
			query := "INSERT INTO " + projectionName + " (" + columnNames + ") VALUES (" + valuesPlaceholder + ")"
			_, err := ex.Exec(query, args...)
			return err
		},
	}
}

func NewUpdateStatement(conditions []Column, values []Column, sequence, previousSequence uint64) Statement {
	cols, params, args := columnsToQuery(values)
	wheres, whereArgs := columnsToWhere(conditions, len(params))
	args = append(args, whereArgs...)

	columnNames := strings.Join(cols, ", ")
	valuesPlaceholder := strings.Join(params, ", ")
	wheresPlaceholders := strings.Join(wheres, " AND ")

	return Statement{
		Sequence:         sequence,
		PreviousSequence: previousSequence,
		execute: func(ex executer, projectionName string) error {
			if projectionName == "" {
				return ErrNoTable
			}
			if previousSequence >= sequence {
				return ErrPrevSeqGtSeq
			}
			if len(values) == 0 {
				return ErrNoValues
			}
			if len(conditions) == 0 {
				return ErrNoCondition
			}
			query := "UPDATE " + projectionName + " SET (" + columnNames + ") = (" + valuesPlaceholder + ") WHERE " + wheresPlaceholders
			_, err := ex.Exec(query, args...)
			return err
		},
	}
}

func NewDeleteStatement(conditions []Column, sequence, previousSequence uint64) Statement {
	wheres, args := columnsToWhere(conditions, 0)

	wheresPlaceholders := strings.Join(wheres, " AND ")

	return Statement{
		Sequence:         sequence,
		PreviousSequence: previousSequence,
		execute: func(ex executer, projectionName string) error {
			if projectionName == "" {
				return ErrNoTable
			}
			if previousSequence >= sequence {
				return ErrPrevSeqGtSeq
			}
			if len(conditions) == 0 {
				return ErrNoCondition
			}
			query := fmt.Sprintf("DELETE FROM " + projectionName + " WHERE " + wheresPlaceholders)
			_, err := ex.Exec(query, args...)
			return err
		},
	}
}

func NewNoOpStatement(sequence, previousSequence uint64) Statement {
	return Statement{
		Sequence:         sequence,
		PreviousSequence: previousSequence,
	}
}

func (stmt *Statement) Execute(ex executer, projectionName string) error {
	if stmt.execute == nil {
		return nil
	}
	return stmt.execute(ex, projectionName)
}

type Column struct {
	Name  string
	Value interface{}
}

func columnsToQuery(cols []Column) (names []string, parameters []string, values []interface{}) {
	names = make([]string, len(cols))
	values = make([]interface{}, len(cols))
	parameters = make([]string, len(cols))
	for i, col := range cols {
		names[i] = col.Name
		values[i] = col.Value
		parameters[i] = "$" + strconv.Itoa(i+1)

	}
	return names, parameters, values
}

func columnsToWhere(cols []Column, paramOffset int) (wheres []string, values []interface{}) {
	wheres = make([]string, len(cols))
	values = make([]interface{}, len(cols))

	for i, col := range cols {
		wheres[i] = "(" + col.Name + " = $" + strconv.Itoa(i+1+paramOffset) + ")"
		values[i] = col.Value
	}

	return wheres, values
}
