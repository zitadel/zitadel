package database

// Condition represents a SQL condition.
// Its written after the WHERE keyword in a SQL statement.
type Condition interface {
	Write(builder *StatementBuilder)
}

type and struct {
	conditions []Condition
}

// Write implements [Condition].
func (a *and) Write(builder *StatementBuilder) {
	if len(a.conditions) > 1 {
		builder.WriteString("(")
		defer builder.WriteString(")")
	}
	firstCondition := true
	for _, condition := range a.conditions {
		// if condition == nil {
		// 	continue
		// }
		if !firstCondition {
			builder.WriteString(" AND ")
		}
		condition.Write(builder)
		firstCondition = false
	}
}

// And combines multiple conditions with AND.
func And(conditions ...Condition) *and {
	return &and{conditions: conditions}
}

var _ Condition = (*and)(nil)

type or struct {
	conditions []Condition
}

// Write implements [Condition].
func (o *or) Write(builder *StatementBuilder) {
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

var _ Condition = (*or)(nil)

type isNull struct {
	column Column
}

// Write implements [Condition].
func (i *isNull) Write(builder *StatementBuilder) {
	i.column.Write(builder)
	builder.WriteString(" IS NULL")
}

// IsNull creates a condition that checks if a column is NULL.
func IsNull(column Column) *isNull {
	return &isNull{column: column}
}

var _ Condition = (*isNull)(nil)

type isNotNull struct {
	column Column
}

// Write implements [Condition].
func (i *isNotNull) Write(builder *StatementBuilder) {
	i.column.Write(builder)
	builder.WriteString(" IS NOT NULL")
}

// IsNotNull creates a condition that checks if a column is NOT NULL.
func IsNotNull(column Column) *isNotNull {
	return &isNotNull{column: column}
}

var _ Condition = (*isNotNull)(nil)

type valueCondition func(builder *StatementBuilder)

// NewTextCondition creates a condition that compares a text column with a value.
func NewTextCondition[V Text](col Column, op TextOperation, value V) Condition {
	return valueCondition(func(builder *StatementBuilder) {
		writeTextOperation(builder, col, op, value)
	})
}

// NewDateCondition creates a condition that compares a numeric column with a value.
func NewNumberCondition[V Number](col Column, op NumberOperation, value V) Condition {
	return valueCondition(func(builder *StatementBuilder) {
		writeNumberOperation(builder, col, op, value)
	})
}

// NewDateCondition creates a condition that compares a boolean column with a value.
func NewBooleanCondition[V Boolean](col Column, value V) Condition {
	return valueCondition(func(builder *StatementBuilder) {
		writeBooleanOperation(builder, col, value)
	})
}

// Write implements [Condition].
func (c valueCondition) Write(builder *StatementBuilder) {
	c(builder)
}

var _ Condition = (*valueCondition)(nil)
