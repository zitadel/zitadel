package stmt

import "fmt"

type statementApplier[T any] interface {
	// Apply writes the statement to the builder.
	Apply(stmt *statement[T])
}

type Condition[T any] interface {
	statementApplier[T]
}

type op interface {
	TextOperation | NumberOperation | ListOperation
	fmt.Stringer
}

type operation[T any, O op] struct {
	o O
}

func (o operation[T, O]) String() string {
	return o.o.String()
}

func (o operation[T, O]) Apply(stmt *statement[T]) {
	stmt.builder.WriteString(o.o.String())
}

type condition[V, T any, OP op] struct {
	field Column[T]
	op    OP
	value V
}

func (c *condition[V, T, OP]) Apply(stmt *statement[T]) {
	// placeholder := stmt.appendArg(c.value)
	stmt.builder.WriteString(stmt.columnPrefix())
	stmt.builder.WriteString(c.field.String())
	// stmt.builder.WriteString(c.op)
	// stmt.builder.WriteString(placeholder)
}

type and[T any] struct {
	conditions []Condition[T]
}

func And[T any](conditions ...Condition[T]) *and[T] {
	return &and[T]{
		conditions: conditions,
	}
}

// Apply implements [Condition].
func (a *and[T]) Apply(stmt *statement[T]) {
	if len(a.conditions) > 1 {
		stmt.builder.WriteString("(")
		defer stmt.builder.WriteString(")")
	}

	for i, condition := range a.conditions {
		if i > 0 {
			stmt.builder.WriteString(" AND ")
		}
		condition.Apply(stmt)
	}
}

var _ Condition[any] = (*and[any])(nil)

type or[T any] struct {
	conditions []Condition[T]
}

func Or[T any](conditions ...Condition[T]) *or[T] {
	return &or[T]{
		conditions: conditions,
	}
}

// Apply implements [Condition].
func (o *or[T]) Apply(stmt *statement[T]) {
	if len(o.conditions) > 1 {
		stmt.builder.WriteString("(")
		defer stmt.builder.WriteString(")")
	}

	for i, condition := range o.conditions {
		if i > 0 {
			stmt.builder.WriteString(" OR ")
		}
		condition.Apply(stmt)
	}
}

var _ Condition[any] = (*or[any])(nil)
