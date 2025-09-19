package database

type Columns []*Column

func (m Columns) WriteQualified(builder *StatementBuilder) {
	for i, col := range m {
		if i > 0 {
			builder.WriteString(", ")
		}
		col.WriteQualified(builder)
	}
}

func (m Columns) WriteUnqualified(builder *StatementBuilder) {
	for i, col := range m {
		if i > 0 {
			builder.WriteString(", ")
		}
		col.WriteUnqualified(builder)
	}
}

type Column struct {
	table string
	name  string
}

func NewColumn(table, name string) *Column {
	return &Column{table: table, name: name}
}

func (c Column) WriteQualified(builder *StatementBuilder) {
	builder.Grow(len(c.table) + len(c.name) + 1)
	builder.WriteString(c.table)
	builder.WriteRune('.')
	builder.WriteString(c.name)
}

func (c Column) WriteUnqualified(builder *StatementBuilder) {
	builder.WriteString(c.name)
}

func (c *Column) Equals(col *Column) bool {
	if col == nil {
		return c == nil
	}
	return c.table == col.table && c.name == col.name
}

func Lower(col *Column) *lowerColumn {
	return &lowerColumn{Column: col}
}

type lowerColumn struct {
	*Column
}

func (c lowerColumn) WriteQualified(builder *StatementBuilder) {
	builder.Grow(len("lower()") + len(c.table) + len(c.name) + 1)
	builder.WriteString("LOWER(")
	c.Column.WriteQualified(builder)
	builder.WriteRune(')')
}

func (c lowerColumn) WriteUnqualified(builder *StatementBuilder) {
	builder.WriteString("LOWER(")
	c.Column.WriteUnqualified(builder)
	builder.WriteRune(')')
}

func (c *lowerColumn) Equals(col *Column) bool {
	if col == nil {
		return c == nil
	}
	return c.table == col.table && c.name == col.name
}
