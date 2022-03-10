package crdb

import (
	"fmt"
	"strings"

	"github.com/caos/zitadel/internal/eventstore/handler"
)

type Column struct {
	Name          string
	Type          ColumnType
	nullable      bool
	defaultValue  interface{}
	suffix        string
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

func TableSuffix(suffix string) ColumnOption {
	return func(c *Column) {
		c.suffix = suffix
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

func NewTableCheck(table *Table, opts ...execOption) *handler.Check {
	config := execConfig{}
	create := func(config execConfig) string {
		return createTable(table, config.tableName, "")
	}

	return &handler.Check{
		Execute: exec(config, create, opts),
	}
}

type Table struct {
	columns    []*Column
	primaryKey PrimaryKey
	suffix     string
}

func NewTable(columns []*Column, key PrimaryKey) *Table {
	return &Table{
		columns:    columns,
		primaryKey: key,
	}
}

func NewSecondaryTable(columns []*Column, key PrimaryKey, suffix string) *Table {
	return &Table{
		columns:    columns,
		primaryKey: key,
		suffix:     suffix,
	}
}

func NewMultiTableCheck(primaryTable *Table, secondaryTables ...*Table) *handler.Check {
	config := execConfig{}
	create := func(config execConfig) string {
		stmt := createTable(primaryTable, config.tableName, "")
		for _, table := range secondaryTables {
			stmt += createTable(table, config.tableName, "_"+table.suffix)
		}
		return stmt
	}

	return &handler.Check{
		Execute: exec(config, create, nil),
	}
}

func NewViewCheck(selectStmt string, secondaryTables ...*Table) *handler.Check {
	config := execConfig{}
	create := func(config execConfig) string {
		var stmt string
		for _, table := range secondaryTables {
			stmt += createTable(table, config.tableName, "_"+table.suffix)
		}
		stmt += createView(config.tableName, selectStmt)
		return stmt
	}

	return &handler.Check{
		Execute: exec(config, create, nil),
	}
}

func createTable(table *Table, tableName string, suffix string) string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s, PRIMARY KEY (%s));",
		tableName+suffix,
		createColumnsToQuery(table.columns, tableName),
		strings.Join(table.primaryKey, ", "),
	)
}

func createView(viewName string, selectStmt string) string {
	return fmt.Sprintf("CREATE VIEW IF NOT EXISTS %s AS %s",
		viewName,
		selectStmt,
	)
}

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

func NewIndexCheck(index *Index, opts ...execOption) *handler.Check {
	config := execConfig{}
	create := func(config execConfig) string {
		stmt := fmt.Sprintf("CREATE INDEX %s IF NOT EXISTS ON %s (%s)",
			index.Name,
			config.tableName,
			strings.Join(index.Columns, ","),
		)
		if index.bucketCount == 0 {
			return stmt
		}
		return fmt.Sprintf("SET experimental_enable_hash_sharded_indexes=on; %s USING HASH WITH BUCKET_COUNT = %d",
			stmt, index.bucketCount)
	}

	return &handler.Check{
		Execute: exec(config, create, opts),
	}
}

func createColumnsToQuery(cols []*Column, tableName string) string {
	columns := make([]string, len(cols))
	for i, col := range cols {
		column := col.Name + " " + columnType(col.Type)
		if !col.nullable {
			column += " NOT NULL"
		}
		if col.defaultValue != nil {
			column += " DEFAULT " + defaultValue(col.defaultValue)
		}
		if col.deleteCascade != "" {
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
