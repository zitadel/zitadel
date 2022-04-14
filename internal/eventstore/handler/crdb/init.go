package crdb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/caos/logging"
	"github.com/lib/pq"

	caos_errs "github.com/caos/zitadel/internal/errors"

	"github.com/caos/zitadel/internal/eventstore/handler"
)

type Table struct {
	columns     []*Column
	primaryKey  PrimaryKey
	indices     []*Index
	constraints []*Constraint
}

func NewTable(columns []*Column, key PrimaryKey, opts ...TableOption) *Table {
	t := &Table{
		columns:    columns,
		primaryKey: key,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

type SuffixedTable struct {
	Table
	suffix string
}

func NewSuffixedTable(columns []*Column, key PrimaryKey, suffix string, opts ...TableOption) *SuffixedTable {
	return &SuffixedTable{
		Table:  *NewTable(columns, key, opts...),
		suffix: suffix,
	}
}

type TableOption func(*Table)

func WithIndex(index *Index) TableOption {
	return func(table *Table) {
		table.indices = append(table.indices, index)
	}
}

func WithConstraint(constraint *Constraint) TableOption {
	return func(table *Table) {
		table.constraints = append(table.constraints, constraint)
	}
}

type Column struct {
	Name          string
	Type          ColumnType
	nullable      bool
	defaultValue  interface{}
	deleteCascade string
}

type ColumnOption func(*Column)

func NewColumn(name string, columnType ColumnType, opts ...ColumnOption) *Column {
	column := &Column{
		Name:         name,
		Type:         columnType,
		nullable:     false,
		defaultValue: nil,
	}
	for _, opt := range opts {
		opt(column)
	}
	return column
}

func Nullable() ColumnOption {
	return func(c *Column) {
		c.nullable = true
	}
}

func Default(value interface{}) ColumnOption {
	return func(c *Column) {
		c.defaultValue = value
	}
}

func DeleteCascade(column string) ColumnOption {
	return func(c *Column) {
		c.deleteCascade = column
	}
}

type PrimaryKey []string

func NewPrimaryKey(columnNames ...string) PrimaryKey {
	return columnNames
}

type ColumnType int32

const (
	ColumnTypeText ColumnType = iota
	ColumnTypeTextArray
	ColumnTypeJSONB
	ColumnTypeBytes
	ColumnTypeTimestamp
	ColumnTypeEnum
	ColumnTypeEnumArray
	ColumnTypeInt64
	ColumnTypeBool
)

func NewIndex(name string, columns []string, opts ...indexOpts) *Index {
	i := &Index{
		Name:        name,
		Columns:     columns,
		bucketCount: 0,
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

type Index struct {
	Name        string
	Columns     []string
	bucketCount uint16
}

type indexOpts func(*Index)

func Hash(bucketsCount uint16) indexOpts {
	return func(i *Index) {
		i.bucketCount = bucketsCount
	}
}

func NewConstraint(name string, columns []string) *Constraint {
	i := &Constraint{
		Name:    name,
		Columns: columns,
	}
	return i
}

type Constraint struct {
	Name    string
	Columns []string
}

//Init implements handler.Init
func (h *StatementHandler) Init(ctx context.Context, checks ...*handler.Check) error {
	for _, check := range checks {
		if check == nil || check.IsNoop() {
			return nil
		}
		tx, err := h.client.BeginTx(ctx, nil)
		if err != nil {
			return caos_errs.ThrowInternal(err, "CRDB-SAdf2", "begin failed")
		}
		for i, execute := range check.Executes {
			logging.WithFields("projection", h.ProjectionName, "execute", i).Debug("executing check")
			next, err := execute(h.client, h.ProjectionName)
			if err != nil {
				tx.Rollback()
				return err
			}
			if !next {
				logging.WithFields("projection", h.ProjectionName, "execute", i).Debug("skipping next check")
				break
			}
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func NewTableCheck(table *Table, opts ...execOption) *handler.Check {
	config := execConfig{}
	create := func(config execConfig) string {
		return createTableStatement(table, config.tableName, "")
	}
	executes := make([]func(handler.Executer, string) (bool, error), len(table.indices)+1)
	executes[0] = execNextIfExists(config, create, opts, true)
	for i, index := range table.indices {
		executes[i+1] = execNextIfExists(config, createIndexStatement(index), opts, true)
	}
	return &handler.Check{
		Executes: executes,
	}
}

func NewMultiTableCheck(primaryTable *Table, secondaryTables ...*SuffixedTable) *handler.Check {
	config := execConfig{}
	create := func(config execConfig) string {
		stmt := createTableStatement(primaryTable, config.tableName, "")
		for _, table := range secondaryTables {
			stmt += createTableStatement(&table.Table, config.tableName, "_"+table.suffix)
		}
		return stmt
	}

	return &handler.Check{
		Executes: []func(handler.Executer, string) (bool, error){
			execNextIfExists(config, create, nil, true),
		},
	}
}

func NewViewCheck(selectStmt string, secondaryTables ...*SuffixedTable) *handler.Check {
	config := execConfig{}
	create := func(config execConfig) string {
		var stmt string
		for _, table := range secondaryTables {
			stmt += createTableStatement(&table.Table, config.tableName, "_"+table.suffix)
		}
		stmt += createViewStatement(config.tableName, selectStmt)
		return stmt
	}

	return &handler.Check{
		Executes: []func(handler.Executer, string) (bool, error){
			execNextIfExists(config, create, nil, false),
		},
	}
}

func execNextIfExists(config execConfig, q query, opts []execOption, executeNext bool) func(handler.Executer, string) (bool, error) {
	return func(handler handler.Executer, name string) (bool, error) {
		err := exec(config, q, opts)(handler, name)
		if isErrAlreadyExists(err) {
			return executeNext, nil
		}
		return false, err
	}
}

func isErrAlreadyExists(err error) bool {
	caosErr := &caos_errs.CaosError{}
	if !errors.As(err, &caosErr) {
		return false
	}
	sqlErr, ok := caosErr.GetParent().(*pq.Error)
	if !ok {
		return false
	}
	return sqlErr.Routine == "NewRelationAlreadyExistsError"
}

func createTableStatement(table *Table, tableName string, suffix string) string {
	stmt := fmt.Sprintf("CREATE TABLE %s (%s, PRIMARY KEY (%s)",
		tableName+suffix,
		createColumnsStatement(table.columns, tableName),
		strings.Join(table.primaryKey, ", "),
	)
	for _, index := range table.indices {
		stmt += fmt.Sprintf(", INDEX %s (%s)", index.Name, strings.Join(index.Columns, ","))
	}
	for _, constraint := range table.constraints {
		stmt += fmt.Sprintf(", CONSTRAINT %s UNIQUE (%s)", constraint.Name, strings.Join(constraint.Columns, ","))
	}
	return stmt + ");"
}

func createViewStatement(viewName string, selectStmt string) string {
	return fmt.Sprintf("CREATE VIEW %s AS %s",
		viewName,
		selectStmt,
	)
}

func createIndexStatement(index *Index) func(config execConfig) string {
	return func(config execConfig) string {
		stmt := fmt.Sprintf("CREATE INDEX %s ON %s (%s)",
			index.Name,
			config.tableName,
			strings.Join(index.Columns, ","),
		)
		if index.bucketCount == 0 {
			return stmt + ";"
		}
		return fmt.Sprintf("SET experimental_enable_hash_sharded_indexes=on; %s USING HASH WITH BUCKET_COUNT = %d;",
			stmt, index.bucketCount)
	}
}

func createColumnsStatement(cols []*Column, tableName string) string {
	columns := make([]string, len(cols))
	for i, col := range cols {
		column := col.Name + " " + columnType(col.Type)
		if !col.nullable {
			column += " NOT NULL"
		}
		if col.defaultValue != nil {
			column += " DEFAULT " + defaultValue(col.defaultValue)
		}
		if len(col.deleteCascade) != 0 {
			column += fmt.Sprintf(" REFERENCES %s (%s) ON DELETE CASCADE", tableName, col.deleteCascade)
		}
		columns[i] = column
	}
	return strings.Join(columns, ",")
}

func defaultValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return "'" + v + "'"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func columnType(columnType ColumnType) string {
	switch columnType {
	case ColumnTypeText:
		return "TEXT"
	case ColumnTypeTextArray:
		return "TEXT[]"
	case ColumnTypeTimestamp:
		return "TIMESTAMPTZ"
	case ColumnTypeEnum:
		return "SMALLINT"
	case ColumnTypeEnumArray:
		return "SMALLINT[]"
	case ColumnTypeInt64:
		return "BIGINT"
	case ColumnTypeBool:
		return "BOOLEAN"
	case ColumnTypeJSONB:
		return "JSONB"
	case ColumnTypeBytes:
		return "BYTES"
	default:
		panic("") //TODO: remove?
		return ""
	}
}
