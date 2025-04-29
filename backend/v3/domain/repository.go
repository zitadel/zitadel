package domain

import (
	"time"

	"golang.org/x/exp/constraints"
)

type Operation interface {
	// TextOperation |
	// 	NumberOperation |
	// 	BoolOperation

	op()
}

type clause[F ~uint8, Op Operation] struct {
	field F
	op    Op
}

func (c *clause[F, Op]) Field() F {
	return c.field
}

func (c *clause[F, Op]) Operation() Op {
	return c.op
}

type Text interface {
	~string | ~[]byte
}

type TextOperation uint8

const (
	TextOperationEqual TextOperation = iota
	TextOperationNotEqual
	TextOperationStartsWith
	TextOperationStartsWithIgnoreCase
)

func (TextOperation) op() {}

type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex | time.Time
}

type NumberOperation uint8

const (
	NumberOperationEqual NumberOperation = iota
	NumberOperationNotEqual
	NumberOperationLessThan
	NumberOperationLessThanOrEqual
	NumberOperationGreaterThan
	NumberOperationGreaterThanOrEqual
)

func (NumberOperation) op() {}

type Bool interface {
	~bool
}

type BoolOperation uint8

const (
	BoolOperationIs BoolOperation = iota
	BoolOperationNot
)

func (BoolOperation) op() {}

type ListOperation uint8

const (
	ListOperationContains ListOperation = iota
	ListOperationNotContains
)

func (ListOperation) op() {}
