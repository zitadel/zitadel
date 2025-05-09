package stmt

import (
	"time"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex | time.Time | time.Duration
}

type between[N Number] struct {
	min, max N
}

type NumberBetween[V Number, T any] struct {
	condition[between[V], T, NumberOperation]
}

func (nb *NumberBetween[V, T]) Apply(stmt *statement[T]) {
	nb.field.Apply(stmt)
	stmt.builder.WriteString(" BETWEEN ")
	stmt.builder.WriteString(stmt.appendArg(nb.value.min))
	stmt.builder.WriteString(" AND ")
	stmt.builder.WriteString(stmt.appendArg(nb.value.max))
}

type NumberCondition[V Number, T any] struct {
	condition[V, T, NumberOperation]
}

func (nc *NumberCondition[V, T]) Apply(stmt *statement[T]) {
	nc.field.Apply(stmt)
	operation[T, NumberOperation]{nc.op}.Apply(stmt)
	stmt.builder.WriteString(stmt.appendArg(nc.value))
}

type NumberOperation uint8

const (
	NumberOperationEqual NumberOperation = iota + 1
	NumberOperationNotEqual
	NumberOperationLessThan
	NumberOperationLessThanOrEqual
	NumberOperationGreaterThan
	NumberOperationGreaterThanOrEqual
)

var numberOperations = map[NumberOperation]string{
	NumberOperationEqual:              " = ",
	NumberOperationNotEqual:           " <> ",
	NumberOperationLessThan:           " < ",
	NumberOperationLessThanOrEqual:    " <= ",
	NumberOperationGreaterThan:        " > ",
	NumberOperationGreaterThanOrEqual: " >= ",
}

func (no NumberOperation) String() string {
	return numberOperations[no]
}
