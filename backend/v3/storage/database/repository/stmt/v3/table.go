package v3

import "github.com/zitadel/zitadel/backend/v3/storage/database"

type object interface {
	User | Org | Instance
	Columns(t Table) []Column
	Scan(s database.Scanner) error
}

type Table interface {
	Schema() string
	Name() string
	Alias() string
	Columns() []Column

	writeOn(builder statementBuilder)
}

type table struct {
	schema string
	name   string
	alias  string

	possibleJoins func(table Table) map[string]Column

	columns map[string]Column
	colList []Column
}

func newTable[O object](schema, name string) *table {
	t := &table{
		schema: schema,
		name:   name,
	}

	var o O
	t.colList = o.Columns(t)
	t.columns = make(map[string]Column, len(t.colList))
	for _, col := range t.colList {
		t.columns[col.Name()] = col
	}

	return t
}

// Columns implements [Table].
func (t *table) Columns() []Column {
	if len(t.colList) > 0 {
		return t.colList
	}

	t.colList = make([]Column, 0, len(t.columns))
	for _, column := range t.columns {
		t.colList = append(t.colList, column)
	}

	return t.colList
}

// Name implements [Table].
func (t *table) Name() string {
	return t.name
}

// Schema implements [Table].
func (t *table) Schema() string {
	return t.schema
}

// Alias implements [Table].
func (t *table) Alias() string {
	if t.alias != "" {
		return t.alias
	}
	return t.schema + "." + t.name
}

// writeOn implements [Table].
func (t *table) writeOn(builder statementBuilder) {
	builder.writeString(t.Alias())
}

var _ Table = (*table)(nil)
