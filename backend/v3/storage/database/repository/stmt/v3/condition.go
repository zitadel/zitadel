package v3

type statementBuilder interface {
	write([]byte)
	writeString(string)
	writeRune(rune)

	appendArg(any) (placeholder string)
	table() Table
}

type Condition interface {
	writeOn(builder statementBuilder)
}

type and struct {
	conditions []Condition
}

func And(conditions ...Condition) *and {
	return &and{conditions: conditions}
}

// writeOn implements [Condition].
func (a *and) writeOn(builder statementBuilder) {
	if len(a.conditions) > 1 {
		builder.writeString("(")
		defer builder.writeString(")")
	}

	for i, condition := range a.conditions {
		if i > 0 {
			builder.writeString(" AND ")
		}
		condition.writeOn(builder)
	}
}

var _ Condition = (*and)(nil)

type or struct {
	conditions []Condition
}

func Or(conditions ...Condition) *or {
	return &or{conditions: conditions}
}

// writeOn implements [Condition].
func (o *or) writeOn(builder statementBuilder) {
	if len(o.conditions) > 1 {
		builder.writeString("(")
		defer builder.writeString(")")
	}

	for i, condition := range o.conditions {
		if i > 0 {
			builder.writeString(" OR ")
		}
		condition.writeOn(builder)
	}
}

var _ Condition = (*or)(nil)

type isNull struct {
	column Column
}

func IsNull(column Column) *isNull {
	return &isNull{column: column}
}

// writeOn implements [Condition].
func (cond *isNull) writeOn(builder statementBuilder) {
	cond.column.Write(builder)
	builder.writeString(" IS NULL")
}

var _ Condition = (*isNull)(nil)

type isNotNull struct {
	column Column
}

func IsNotNull(column Column) *isNotNull {
	return &isNotNull{column: column}
}

// writeOn implements [Condition].
func (cond *isNotNull) writeOn(builder statementBuilder) {
	cond.column.Write(builder)
	builder.writeString(" IS NOT NULL")
}

var _ Condition = (*isNotNull)(nil)

type condition[Op Operator, V Value] struct {
	column   Column
	operator Op
	value    V
}

// writeOn implements [Condition].
func (cond condition[Op, V]) writeOn(builder statementBuilder) {
	cond.column.Write(builder)
	builder.writeString(cond.operator.String())
	builder.writeString(builder.appendArg(cond.value))
}

var _ Condition = (*condition[TextOperator, string])(nil)

type textCondition[V Text] struct {
	condition[TextOperator, V]
}

func NewTextCondition[V Text](column Column, operator TextOperator, value V) *textCondition[V] {
	return &textCondition[V]{
		condition: condition[TextOperator, V]{
			column:   column,
			operator: operator,
			value:    value,
		},
	}
}

// writeOn implements [Condition].
func (cond *textCondition[V]) writeOn(builder statementBuilder) {
	switch cond.operator {
	case TextOperatorEqual, TextOperatorNotEqual:
		cond.column.Write(builder)
		builder.writeString(cond.operator.String())
		builder.writeString(builder.appendArg(cond.value))
	case TextOperatorEqualIgnoreCase, TextOperatorNotEqualIgnoreCase:
		if col, ok := cond.column.(ignoreCaseColumn); ok {
			col.WriteIgnoreCase(builder)
		} else {
			builder.writeString("LOWER(")
			cond.column.Write(builder)
			builder.writeString(")")
		}
		builder.writeString(cond.operator.String())
		builder.writeString("LOWER(")
		builder.writeString(builder.appendArg(cond.value))
		builder.writeString(")")
	case TextOperatorStartsWith:
		cond.column.Write(builder)
		builder.writeString(cond.operator.String())
		builder.writeString(builder.appendArg(cond.value))
		builder.writeString(" || '%'")
	case TextOperatorStartsWithIgnoreCase:
		if col, ok := cond.column.(ignoreCaseColumn); ok {
			col.WriteIgnoreCase(builder)
		} else {
			builder.writeString("LOWER(")
			cond.column.Write(builder)
			builder.writeString(")")
		}
		builder.writeString(cond.operator.String())
		builder.writeString("LOWER(")
		builder.writeString(builder.appendArg(cond.value))
		builder.writeString(") || '%'")
	}
}

var _ Condition = (*textCondition[string])(nil)

type numberCondition[V Number] struct {
	condition[NumberOperator, V]
}

func NewNumberCondition[V Number](column Column, operator NumberOperator, value V) *numberCondition[V] {
	return &numberCondition[V]{
		condition: condition[NumberOperator, V]{
			column:   column,
			operator: operator,
			value:    value,
		},
	}
}

var _ Condition = (*numberCondition[int])(nil)
