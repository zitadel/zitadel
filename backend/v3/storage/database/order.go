package database

// Order represents a SQL condition.
// Its written after the ORDER BY keyword in a SQL statement.
type Order interface {
	Write(builder *StatementBuilder)
}

type orderBy struct {
	column Column
}

func OrderBy(column Column) Order {
	return &orderBy{column: column}
}

// Write implements [Order].
func (o *orderBy) Write(builder *StatementBuilder) {
	builder.WriteString(" ORDER BY ")
	o.column.WriteQualified(builder)
}
