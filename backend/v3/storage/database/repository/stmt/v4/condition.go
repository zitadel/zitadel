package v4

type Condition interface {
	writeTo(builder *statementBuilder)
}

type and struct {
	conditions []Condition
}

// writeTo implements [Condition].
func (a *and) writeTo(builder *statementBuilder) {
	if len(a.conditions) > 1 {
		builder.WriteString("(")
		defer builder.WriteString(")")
	}
	for i, condition := range a.conditions {
		if i > 0 {
			builder.WriteString(" AND ")
		}
		condition.writeTo(builder)
	}
}

func And(conditions ...Condition) *and {
	return &and{conditions: conditions}
}

var _ Condition = (*and)(nil)

type or struct {
	conditions []Condition
}

// writeTo implements [Condition].
func (o *or) writeTo(builder *statementBuilder) {
	if len(o.conditions) > 1 {
		builder.WriteString("(")
		defer builder.WriteString(")")
	}
	for i, condition := range o.conditions {
		if i > 0 {
			builder.WriteString(" OR ")
		}
		condition.writeTo(builder)
	}
}

func Or(conditions ...Condition) *or {
	return &or{conditions: conditions}
}

var _ Condition = (*or)(nil)

type isNull struct {
	column Column
}

// writeTo implements [Condition].
func (i *isNull) writeTo(builder *statementBuilder) {
	i.column.writeTo(builder)
	builder.WriteString(" IS NULL")
}

func IsNull(column Column) *isNull {
	return &isNull{column: column}
}

var _ Condition = (*isNull)(nil)

type isNotNull struct {
	column Column
}

// writeTo implements [Condition].
func (i *isNotNull) writeTo(builder *statementBuilder) {
	i.column.writeTo(builder)
	builder.WriteString(" IS NOT NULL")
}

func IsNotNull(column Column) *isNotNull {
	return &isNotNull{column: column}
}

var _ Condition = (*isNotNull)(nil)

type valueCondition func(builder *statementBuilder)

func newTextCondition[V Text](col Column, op TextOperator, value V) Condition {
	return valueCondition(func(builder *statementBuilder) {
		writeTextOperation(builder, col, op, value)
	})
}

func newNumberCondition[V Number](col Column, op NumberOperator, value V) Condition {
	return valueCondition(func(builder *statementBuilder) {
		writeNumberOperation(builder, col, op, value)
	})
}

func newBooleanCondition[V Boolean](col Column, value V) Condition {
	return valueCondition(func(builder *statementBuilder) {
		writeBooleanOperation(builder, col, value)
	})
}

// writeTo implements [Condition].
func (c valueCondition) writeTo(builder *statementBuilder) {
	c(builder)
}

var _ Condition = (*valueCondition)(nil)
