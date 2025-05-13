package database

type Columns []Column

// Write implements [Column].
func (m Columns) Write(builder *StatementBuilder) {
	for i, col := range m {
		if i > 0 {
			builder.WriteString(", ")
		}
		col.Write(builder)
	}
}

type Column interface {
	Write(builder *StatementBuilder)
}

type column struct {
	name string
}

func NewColumn(name string) Column {
	return column{name: name}
}

// Write implements [Column].
func (c column) Write(builder *StatementBuilder) {
	builder.WriteString(c.name)
}

var _ Column = (*column)(nil)

type ignoreCaseColumn interface {
	Column
	WriteIgnoreCase(builder *StatementBuilder)
}

func NewIgnoreCaseColumn(name, suffix string) ignoreCaseColumn {
	return ignoreCaseCol{
		column: column{name: name},
		suffix: suffix,
	}
}

type ignoreCaseCol struct {
	column
	suffix string
}

// WriteIgnoreCase implements [ignoreCaseColumn].
func (c ignoreCaseCol) WriteIgnoreCase(builder *StatementBuilder) {
	c.column.Write(builder)
	builder.WriteString(c.suffix)
}
