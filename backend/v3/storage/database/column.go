package database

type Columns []Column

// WriteQualified implements [Column].
func (m Columns) WriteQualified(builder *StatementBuilder) {
	for i, col := range m {
		if i > 0 {
			builder.WriteString(", ")
		}
		col.WriteQualified(builder)
	}
}

// WriteUnqualified implements [Column].
func (m Columns) WriteUnqualified(builder *StatementBuilder) {
	for i, col := range m {
		if i > 0 {
			builder.WriteString(", ")
		}
		col.WriteUnqualified(builder)
	}
}

// Column represents a column in a database table.
type Column interface {
	// Write(builder *StatementBuilder)
	WriteQualified(builder *StatementBuilder)
	WriteUnqualified(builder *StatementBuilder)
}

type column struct {
	table string
	name  string
}

func NewColumn(table, name string) Column {
	return column{table: table, name: name}
}

// WriteQualified implements [Column].
func (c column) WriteQualified(builder *StatementBuilder) {
	builder.Grow(len(c.table) + len(c.name) + 1)
	builder.WriteString(c.table)
	builder.WriteRune('.')
	builder.WriteString(c.name)
}

// WriteUnqualified implements [Column].
func (c column) WriteUnqualified(builder *StatementBuilder) {
	builder.WriteString(c.name)
}

var _ Column = (*column)(nil)

// // ignoreCaseColumn represents two database columns, one for the
// // original value and one for the lower case value.
// type ignoreCaseColumn interface {
// 	Column
// 	WriteIgnoreCase(builder *StatementBuilder)
// }

// func NewIgnoreCaseColumn(col Column, suffix string) ignoreCaseColumn {
// 	return ignoreCaseCol{
// 		column: col,
// 		suffix: suffix,
// 	}
// }

// type ignoreCaseCol struct {
// 	column Column
// 	suffix string
// }

// // WriteIgnoreCase implements [ignoreCaseColumn].
// func (c ignoreCaseCol) WriteIgnoreCase(builder *StatementBuilder) {
// 	c.column.WriteQualified(builder)
// 	builder.WriteString(c.suffix)
// }

// // WriteQualified implements [ignoreCaseColumn].
// func (c ignoreCaseCol) WriteQualified(builder *StatementBuilder) {
// 	c.column.WriteQualified(builder)
// 	builder.WriteString(c.suffix)
// }
