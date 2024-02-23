package database

import (
	"time"

	"github.com/zitadel/logging"
	"golang.org/x/exp/constraints"
)

type NumberFilter[N number] interface {
	Condition
	implementsNumberFilter()
}

type NumberCondition[N number] struct {
	Filter[numberCompare, N]
}

func NewNumberEquals[N number](n N) *NumberCondition[N] {
	return newNumberFilter(numberEqual, n)
}

func NewNumberAtLeast[N number](n N) *NumberCondition[N] {
	return newNumberFilter(numberAtLeast, n)
}

func NewNumberAtMost[N number](n N) *NumberCondition[N] {
	return newNumberFilter(numberAtMost, n)
}

func NewNumberGreater[N number](n N) *NumberCondition[N] {
	return newNumberFilter(numberGreater, n)
}

func NewNumberLess[N number](n N) *NumberCondition[N] {
	return newNumberFilter(numberLess, n)
}

func NewNumberUnequal[N number](n N) *NumberCondition[N] {
	return newNumberFilter(numberUnequal, n)
}

func (NumberCondition[N]) implementsNumberFilter() {}

func newNumberFilter[N number](comp numberCompare, n N) *NumberCondition[N] {
	return &NumberCondition[N]{
		Filter: Filter[numberCompare, N]{
			comp:  comp,
			value: n,
		},
	}
}

// NumberBetweenCondition combines [AtLeast] and [AtMost] comparisons
type NumberBetweenCondition[N number] struct {
	min, max N
}

func NewNumberBetween[N number](min, max N) *NumberBetweenCondition[N] {
	return &NumberBetweenCondition[N]{
		min: min,
		max: max,
	}
}

func (f NumberBetweenCondition[N]) Write(stmt *Statement, columnName string) {
	NewNumberAtLeast[N](f.min).Write(stmt, columnName)
	stmt.Builder.WriteString(" AND ")
	NewNumberAtMost[N](f.max).Write(stmt, columnName)
}

func (NumberBetweenCondition[N]) implementsNumberFilter() {}

type numberCompare uint8

const (
	numberEqual numberCompare = iota
	numberAtLeast
	numberAtMost
	numberGreater
	numberLess
	numberUnequal
)

func (c numberCompare) String() string {
	switch c {
	case numberEqual:
		return "="
	case numberAtLeast:
		return ">="
	case numberAtMost:
		return "<="
	case numberGreater:
		return ">"
	case numberLess:
		return "<"
	case numberUnequal:
		return "<>"
	default:
		logging.WithFields("compare", c).Panic("comparison type not implemented")
		return ""
	}
}

type number interface {
	constraints.Integer | constraints.Float | time.Time
}
