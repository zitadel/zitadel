package database

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/zitadel/logging"
	"golang.org/x/exp/constraints"
)

type NumberFilter[N number] struct {
	Filter[numberCompare, N]
}

func NewNumberEquals[N number](n N) *NumberFilter[N] {
	return newNumberFilter(numberEqual, n)
}

func NewNumberAtLeast[N number](n N) *NumberFilter[N] {
	return newNumberFilter(numberAtLeast, n)
}

func NewNumberAtMost[N number](n N) *NumberFilter[N] {
	return newNumberFilter(numberAtMost, n)
}

func NewNumberGreater[N number](n N) *NumberFilter[N] {
	return newNumberFilter(numberGreater, n)
}

func NewNumberLess[N number](n N) *NumberFilter[N] {
	return newNumberFilter(numberLess, n)
}

func NewNumberUnequal[N number](n N) *NumberFilter[N] {
	return newNumberFilter(numberUnequal, n)
}

func newNumberFilter[N number](comp numberCompare, n N) *NumberFilter[N] {
	return &NumberFilter[N]{
		Filter: Filter[numberCompare, N]{
			comp:  comp,
			value: n,
		},
	}
}

// NumberBetweenFilter combines [AtLeast] and [AtMost] comparisons
type NumberBetweenFilter[N number] struct {
	min, max N
}

func NewNumberBetween[N number](min, max N) *NumberBetweenFilter[N] {
	return &NumberBetweenFilter[N]{
		min: min,
		max: max,
	}
}

func (f NumberBetweenFilter[N]) Write(stmt *Statement, columnName string) {
	NewNumberAtLeast[N](f.min).Write(stmt, columnName)
	stmt.WriteString(" AND ")
	NewNumberAtMost[N](f.max).Write(stmt, columnName)
}

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
	constraints.Integer | constraints.Float | time.Time | decimal.Decimal
	// TODO: condition must know if it's args are named parameters or not
	// constraints.Integer | constraints.Float | time.Time | placeholder
}
