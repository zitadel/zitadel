package eventstore

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type Statement interface {
	Prepare() func(context.Context, *sql.Tx) error
}

type CreateStatement struct {
	TableName string
	Values    []Column
}

func (stmt *CreateStatement) Prepare() func(ctx context.Context, tx *sql.Tx) error {
	cols, params, args := columnsToQuery(stmt.Values)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", stmt.TableName, strings.Join(cols, ", "), strings.Join(params, ", "))

	return func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.QueryContext(ctx, query, args...)
		return err
	}
}

type UpdateStatement struct {
	TableName string
	PK        []Column
	Values    []Column
}

func (stmt *UpdateStatement) Prepare() func(context.Context, *sql.Tx) error {
	cols, params, args := columnsToQuery(stmt.Values)
	wheres, whereArgs := columnsToWhere(stmt.PK, len(params))
	args = append(args, whereArgs)
	query := fmt.Sprintf("UPDATE %s SET (%s) = (%s) WHERE %s", stmt.TableName, strings.Join(cols, ", "), strings.Join(params, ", "), strings.Join(wheres, " AND "))

	return func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, query, args...)
		return err
	}
}

type DeleteStatement struct {
	TableName string
	PK        []Column
}

func (stmt *DeleteStatement) Prepare() func(context.Context, *sql.Tx) error {
	wheres, args := columnsToWhere(stmt.PK, 0)
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", stmt.TableName, strings.Join(wheres, " AND "))

	return func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, query, args...)
		return err
	}
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
		wheres[i] = "(" + col.Name + " = " + strconv.Itoa(i+1+paramOffset) + ")"
		values[i] = col.Value
	}

	return wheres, values
}
