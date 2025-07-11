package database

// Order represents a SQL condition.
// Its written after the ORDER BY keyword in a SQL statement.
type Order interface {
	Write(builder *StatementBuilder)
}

type orderBY struct {
	column Column
}

func OrderBY(column Column) Order {
	return &orderBY{column: column}
}

// Write implements [Order].
func (o *orderBY) Write(builder *StatementBuilder) {
	builder.WriteString(" ORDER BY ")
	o.column.Write(builder)
}
