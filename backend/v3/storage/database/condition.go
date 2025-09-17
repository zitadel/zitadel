package database

// Condition represents a SQL condition.
// Its written after the WHERE keyword in a SQL statement.
type Condition interface {
	Write(builder *StatementBuilder)
	// ContainsColumn is used to check if the condition filters for a specific column.
	// It acts as a save guard database operations that should be specific on the given column.
	ContainsColumn(col *Column) bool
}

type and struct {
	conditions []Condition
}

// Write implements [Condition].
func (a and) Write(builder *StatementBuilder) {
	if len(a.conditions) > 1 {
		builder.WriteString("(")
		defer builder.WriteString(")")
	}
	for i, condition := range a.conditions {
		if i > 0 {
			builder.WriteString(" AND ")
		}
		condition.Write(builder)
	}
}

// And combines multiple conditions with AND.
func And(conditions ...Condition) *and {
	return &and{conditions: conditions}
}

func (a and) ContainsColumn(col *Column) bool {
	for _, condition := range a.conditions {
		if condition.ContainsColumn(col) {
			return true
		}
	}
	return false
}

var _ Condition = (*and)(nil)

type or struct {
	conditions []Condition
}

// Write implements [Condition].
func (o or) Write(builder *StatementBuilder) {
	if len(o.conditions) > 1 {
		builder.WriteString("(")
		defer builder.WriteString(")")
	}
	for i, condition := range o.conditions {
		if i > 0 {
			builder.WriteString(" OR ")
		}
		condition.Write(builder)
	}
}

// Or combines multiple conditions with OR.
func Or(conditions ...Condition) *or {
	return &or{conditions: conditions}
}

// ContainsColumn implements [Condition].
// It always returns false because OR conditions
func (o or) ContainsColumn(col *Column) bool {
	for _, condition := range o.conditions {
		if !condition.ContainsColumn(col) {
			return false
		}
	}
	return true
}

var _ Condition = (*or)(nil)

type isNull struct {
	column *Column
}

// Write implements [Condition].
func (i isNull) Write(builder *StatementBuilder) {
	i.column.WriteQualified(builder)
	builder.WriteString(" IS NULL")
}

// IsNull creates a condition that checks if a column is NULL.
func IsNull(column *Column) *isNull {
	return &isNull{column: column}
}

func (i isNull) ContainsColumn(col *Column) bool {
	return i.column.Equals(col)
}

var _ Condition = (*isNull)(nil)

type isNotNull struct {
	column *Column
}

// Write implements [Condition].
func (i isNotNull) Write(builder *StatementBuilder) {
	i.column.WriteQualified(builder)
	builder.WriteString(" IS NOT NULL")
}

// IsNotNull creates a condition that checks if a column is NOT NULL.
func IsNotNull(column *Column) *isNotNull {
	return &isNotNull{column: column}
}

// ContainsColumn implements [Condition].
func (i isNotNull) ContainsColumn(col *Column) bool {
	return i.column.Equals(col)
}

var _ Condition = (*isNotNull)(nil)

type valueCondition struct {
	write func(builder *StatementBuilder)
	col   *Column
}

// NewTextCondition creates a condition that compares a text column with a value.
func NewTextCondition[V Text](col *Column, op TextOperation, value V) Condition {
	return valueCondition{
		col: col,
		write: func(builder *StatementBuilder) {
			writeTextOperation(builder, col, op, value)
		},
	}
}

// NewDateCondition creates a condition that compares a numeric column with a value.
func NewNumberCondition[V Number](col *Column, op NumberOperation, value V) Condition {
	return valueCondition{
		col: col,
		write: func(builder *StatementBuilder) {
			writeNumberOperation(builder, col, op, value)
		},
	}
}

// NewDateCondition creates a condition that compares a boolean column with a value.
func NewBooleanCondition[V Boolean](col *Column, value V) Condition {
	return valueCondition{
		col: col,
		write: func(builder *StatementBuilder) {
			writeBooleanOperation(builder, col, value)
		},
	}
}

// NewBytesCondition creates a condition that compares a BYTEA column with a value.
func NewBytesCondition[V Bytes](col *Column, op BytesOperation, value V) Condition {
	return valueCondition{
		col: col,
		write: func(builder *StatementBuilder) {
			writeBytesOperation(builder, col, op, value)
		},
	}
}

// NewColumnCondition creates a condition that compares two columns on equality.
func NewColumnCondition(col1, col2 *Column) Condition {
	return valueCondition{
		col: col1,
		write: func(builder *StatementBuilder) {
			col1.WriteQualified(builder)
			builder.WriteString(" = ")
			col2.WriteQualified(builder)
		},
	}
}

// Write implements [Condition].
func (c valueCondition) Write(builder *StatementBuilder) {
	c.write(builder)
}

// ContainsColumn implements [Condition].
func (i valueCondition) ContainsColumn(col *Column) bool {
	return i.col.Equals(col)
}

var _ Condition = (*valueCondition)(nil)
