package database

// Condition represents a SQL condition.
// Its written after the WHERE keyword in a SQL statement.
type Condition interface {
	Write(builder *StatementBuilder)
	// IsRestrictingColumn is used to check if the condition filters for a specific column.
	// It acts as a save guard database operations that should be specific on the given column.
	IsRestrictingColumn(col Column) bool
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

// IsRestrictingColumn implements [Condition].
func (a and) IsRestrictingColumn(col Column) bool {
	for _, condition := range a.conditions {
		if condition.IsRestrictingColumn(col) {
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

// IsRestrictingColumn implements [Condition].
// It returns true only if all conditions are restricting the given column.
func (o or) IsRestrictingColumn(col Column) bool {
	for _, condition := range o.conditions {
		if !condition.IsRestrictingColumn(col) {
			return false
		}
	}
	return true
}

var _ Condition = (*or)(nil)

type isNull struct {
	column Column
}

// Write implements [Condition].
func (i isNull) Write(builder *StatementBuilder) {
	i.column.WriteQualified(builder)
	builder.WriteString(" IS NULL")
}

// IsNull creates a condition that checks if a column is NULL.
func IsNull(column Column) *isNull {
	return &isNull{column: column}
}

// IsRestrictingColumn implements [Condition].
// It returns false because it cannot be used for restricting a column.
func (i isNull) IsRestrictingColumn(col Column) bool {
	return false
}

var _ Condition = (*isNull)(nil)

type isNotNull struct {
	column Column
}

// Write implements [Condition].
func (i isNotNull) Write(builder *StatementBuilder) {
	i.column.WriteQualified(builder)
	builder.WriteString(" IS NOT NULL")
}

// IsNotNull creates a condition that checks if a column is NOT NULL.
func IsNotNull(column Column) *isNotNull {
	return &isNotNull{column: column}
}

// IsRestrictingColumn implements [Condition].
// It returns false because it cannot be used for restricting a column.
func (i isNotNull) IsRestrictingColumn(col Column) bool {
	return false
}

var _ Condition = (*isNotNull)(nil)

type valueCondition struct {
	write func(builder *StatementBuilder)
	col   Column
}

// NewTextCondition creates a condition that compares a text column with a value.
// If you want to use ignore case operations, consider using [NewTextIgnoreCaseCondition].
func NewTextCondition[T Text](col Column, op TextOperation, value T) Condition {
	return valueCondition{
		col: col,
		write: func(builder *StatementBuilder) {
			writeTextOperation[T](builder, col, op, value)
		},
	}
}

// NewTextIgnoreCaseCondition creates a condition that compares a text column with a value, ignoring case by lowercasing both.
func NewTextIgnoreCaseCondition[T Text](col Column, op TextOperation, value T) Condition {
	return valueCondition{
		col: col,
		write: func(builder *StatementBuilder) {
			writeTextOperation[T](builder, LowerColumn(col), op, LowerValue(value))
		},
	}
}

// NewDateCondition creates a condition that compares a numeric column with a value.
func NewNumberCondition[V Number](col Column, op NumberOperation, value V) Condition {
	return valueCondition{
		col: col,
		write: func(builder *StatementBuilder) {
			writeNumberOperation[V](builder, col, op, value)
		},
	}
}

// NewDateCondition creates a condition that compares a boolean column with a value.
func NewBooleanCondition[V Boolean](col Column, value V) Condition {
	return valueCondition{
		col: col,
		write: func(builder *StatementBuilder) {
			writeBooleanOperation[V](builder, col, value)
		},
	}
}

// NewBytesCondition creates a condition that compares a BYTEA column with a value.
func NewBytesCondition[V Bytes](col Column, op BytesOperation, value any) Condition {
	return valueCondition{
		col: col,
		write: func(builder *StatementBuilder) {
			writeBytesOperation[V](builder, col, op, value)
		},
	}
}

// NewColumnCondition creates a condition that compares two columns on equality.
func NewColumnCondition(col1, col2 Column) Condition {
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

// IsRestrictingColumn implements [Condition].
func (i valueCondition) IsRestrictingColumn(col Column) bool {
	return i.col.Equals(col)
}

var _ Condition = (*valueCondition)(nil)

// existsCondition is a helper to write an EXISTS (SELECT 1 FROM <table> WHERE <condition>) clause.
// It implements Condition so it can be composed with other conditions using And/Or.
type existsCondition struct {
	table     string
	condition Condition
}

// Exists creates a condition that checks for the existence of rows in a subquery.
func Exists(table string, condition Condition) Condition {
	return &existsCondition{
		table:     table,
		condition: condition,
	}
}

// Write implements [Condition].
func (e existsCondition) Write(builder *StatementBuilder) {
	builder.WriteString(" EXISTS (SELECT 1 FROM ")
	builder.WriteString(e.table)
	builder.WriteString(" WHERE ")
	e.condition.Write(builder)
	builder.WriteString(")")
}

// IsRestrictingColumn implements [Condition].
func (e existsCondition) IsRestrictingColumn(col Column) bool {
	// Forward to the inner condition so safety checks (like instance_id presence) can still work.
	return e.condition.IsRestrictingColumn(col)
}

var _ Condition = (*existsCondition)(nil)
