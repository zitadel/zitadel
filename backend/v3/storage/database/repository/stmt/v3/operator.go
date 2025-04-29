package v3

import (
	"fmt"
	"time"

	"golang.org/x/exp/constraints"
)

type Value interface {
	Bool | Number | Text
}

type Text interface {
	~string | ~[]byte
}

type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex | time.Time | time.Duration
}

type Bool interface {
	~bool
}

type Operator interface {
	fmt.Stringer
}

type TextOperator uint8

// String implements [Operator].
func (t TextOperator) String() string {
	return textOperators[t]
}

const (
	TextOperatorEqual TextOperator = iota + 1
	TextOperatorEqualIgnoreCase
	TextOperatorNotEqual
	TextOperatorNotEqualIgnoreCase
	TextOperatorStartsWith
	TextOperatorStartsWithIgnoreCase
)

var textOperators = map[TextOperator]string{
	TextOperatorEqual:                " = ",
	TextOperatorEqualIgnoreCase:      " LIKE ",
	TextOperatorNotEqual:             " <> ",
	TextOperatorNotEqualIgnoreCase:   " NOT LIKE ",
	TextOperatorStartsWith:           " LIKE ",
	TextOperatorStartsWithIgnoreCase: " LIKE ",
}

var _ Operator = TextOperator(0)

type NumberOperator uint8

// String implements Operator.
func (n NumberOperator) String() string {
	return numberOperators[n]
}

const (
	NumberOperatorEqual NumberOperator = iota + 1
	NumberOperatorNotEqual
	NumberOperatorLessThan
	NumberOperatorLessThanOrEqual
	NumberOperatorGreaterThan
	NumberOperatorGreaterThanOrEqual
)

var numberOperators = map[NumberOperator]string{
	NumberOperatorEqual:              " = ",
	NumberOperatorNotEqual:           " <> ",
	NumberOperatorLessThan:           " < ",
	NumberOperatorLessThanOrEqual:    " <= ",
	NumberOperatorGreaterThan:        " > ",
	NumberOperatorGreaterThanOrEqual: " >= ",
}

var _ Operator = NumberOperator(0)
