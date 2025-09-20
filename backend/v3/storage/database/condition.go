package database

import "go.uber.org/mock/gomock"

// Condition represents a SQL condition.
// Its written after the WHERE keyword in a SQL statement.
type Condition interface {
	gomock.Matcher
	Write(builder *StatementBuilder)
}

type and struct {
	conditions []Condition
}

// Matches implements Condition.
func (a *and) Matches(x any) bool {
	toMatch, ok := x.(*and)
	if !ok {
		return false
	}
	if len(a.conditions) != len(toMatch.conditions) {
		return false
	}
	for i, condition := range a.conditions {
		if !condition.Matches(toMatch.conditions[i]) {
			return false
		}
	}
	return true
}

// String implements Condition.
func (a *and) String() string {
	return "database.and"
}

// Write implements [Condition].
func (a *and) Write(builder *StatementBuilder) {
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

var _ Condition = (*and)(nil)

type or struct {
	conditions []Condition
}

// Matches implements Condition.
func (o *or) Matches(x any) bool {
	toMatch, ok := x.(*or)
	if !ok {
		return false
	}
	if len(o.conditions) != len(toMatch.conditions) {
		return false
	}
	for i, condition := range o.conditions {
		if !condition.Matches(toMatch.conditions[i]) {
			return false
		}
	}
	return true
}

// String implements Condition.
func (o *or) String() string {
	return "database.or"
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

// Matches implements Condition.
func (i *isNull) Matches(x any) bool {
	toMatch, ok := x.(*isNull)
	if !ok {
		return false
	}
	return i.column.Matches(toMatch.column)
}

// String implements Condition.
func (i *isNull) String() string {
	return "database.isNull"
}

// Write implements [Condition].
func (i *isNull) Write(builder *StatementBuilder) {
	i.column.WriteQualified(builder)
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

// Matches implements Condition.
func (i *isNotNull) Matches(x any) bool {
	toMatch, ok := x.(*isNotNull)
	if !ok {
		return false
	}
	return i.column.Matches(toMatch.column)
}

// String implements Condition.
func (i *isNotNull) String() string {
	return "database.isNotNull"
}

// Write implements [Condition].
func (i *isNotNull) Write(builder *StatementBuilder) {
	i.column.WriteQualified(builder)
	builder.WriteString(" IS NOT NULL")
}

// IsNotNull creates a condition that checks if a column is NOT NULL.
func IsNotNull(column Column) *isNotNull {
	return &isNotNull{column: column}
}

var _ Condition = (*isNotNull)(nil)

type valueCondition func(builder *StatementBuilder)

// Matches implements Condition.
func (c valueCondition) Matches(x any) bool {
	toMatch, ok := x.(valueCondition)
	if !ok {
		return false
	}
	return c.String() == toMatch.String()
}

// String implements Condition.
func (c valueCondition) String() string {
	return "database.valueCondition"
}

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

// NewColumnCondition creates a condition that compares two columns on equality.
func NewColumnCondition(col1, col2 Column) Condition {
	return valueCondition(func(builder *StatementBuilder) {
		col1.WriteQualified(builder)
		builder.WriteString(" = ")
		col2.WriteQualified(builder)
	})
}

// Write implements [Condition].
func (c valueCondition) Write(builder *StatementBuilder) {
	c(builder)
}

var _ Condition = (*valueCondition)(nil)

// existsCondition is a helper to write an EXISTS (SELECT 1 FROM <table> WHERE <condition>) clause.
// It implements Condition so it can be composed with other conditions using And/Or.
type existsCondition struct {
	table     string
	condition Condition
}

// Matches implements Condition.
func (e *existsCondition) Matches(x any) bool {
	toMatch, ok := x.(existsCondition)
	if !ok {
		return false
	}
	return e.String() == toMatch.String()
}

// String implements Condition.
func (e *existsCondition) String() string {
	return "database.existsCondition"
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
