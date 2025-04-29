package v4

import (
	"time"

	"golang.org/x/exp/constraints"
)

type Value interface {
	Boolean | Number | Text | databaseInstruction
}

type Operator interface {
	BooleanOperator | NumberOperator | TextOperator
}

type Text interface {
	~string | ~[]byte
}

type TextOperator uint8

const (
	// TextOperatorEqual compares two strings for equality.
	TextOperatorEqual TextOperator = iota + 1
	// TextOperatorEqualIgnoreCase compares two strings for equality, ignoring case.
	TextOperatorEqualIgnoreCase
	// TextOperatorNotEqual compares two strings for inequality.
	TextOperatorNotEqual
	// TextOperatorNotEqualIgnoreCase compares two strings for inequality, ignoring case.
	TextOperatorNotEqualIgnoreCase
	// TextOperatorStartsWith checks if the first string starts with the second.
	TextOperatorStartsWith
	// TextOperatorStartsWithIgnoreCase checks if the first string starts with the second, ignoring case.
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

func writeTextOperation[T Text](builder *statementBuilder, col Column, op TextOperator, value T) {
	switch op {
	case TextOperatorEqual, TextOperatorNotEqual:
		col.writeTo(builder)
		builder.WriteString(textOperators[op])
		builder.WriteString(builder.appendArg(value))
	case TextOperatorEqualIgnoreCase, TextOperatorNotEqualIgnoreCase:
		if ignoreCaseCol, ok := col.(ignoreCaseColumn); ok {
			ignoreCaseCol.writeIgnoreCaseTo(builder)
		} else {
			builder.WriteString("LOWER(")
			col.writeTo(builder)
			builder.WriteString(")")
		}
		builder.WriteString(textOperators[op])
		builder.WriteString("LOWER(")
		builder.WriteString(builder.appendArg(value))
		builder.WriteString(")")
	case TextOperatorStartsWith:
		col.writeTo(builder)
		builder.WriteString(textOperators[op])
		builder.WriteString(builder.appendArg(value))
		builder.WriteString(" || '%'")
	case TextOperatorStartsWithIgnoreCase:
		if ignoreCaseCol, ok := col.(ignoreCaseColumn); ok {
			ignoreCaseCol.writeIgnoreCaseTo(builder)
		} else {
			builder.WriteString("LOWER(")
			col.writeTo(builder)
			builder.WriteString(")")
		}
		builder.WriteString(textOperators[op])
		builder.WriteString("LOWER(")
		builder.WriteString(builder.appendArg(value))
		builder.WriteString(")")
		builder.WriteString(" || '%'")
	default:
		panic("unsupported text operation")
	}
}

type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex | time.Time | time.Duration
}

type NumberOperator uint8

const (
	// NumberOperatorEqual compares two numbers for equality.
	NumberOperatorEqual NumberOperator = iota + 1
	// NumberOperatorNotEqual compares two numbers for inequality.
	NumberOperatorNotEqual
	// NumberOperatorLessThan compares two numbers to check if the first is less than the second.
	NumberOperatorLessThan
	// NumberOperatorLessThanOrEqual compares two numbers to check if the first is less than or equal to the second.
	NumberOperatorAtLeast
	// NumberOperatorGreaterThan compares two numbers to check if the first is greater than the second.
	NumberOperatorGreaterThan
	// NumberOperatorGreaterThanOrEqual compares two numbers to check if the first is greater than or equal to the second.
	NumberOperatorAtMost
)

var numberOperators = map[NumberOperator]string{
	NumberOperatorEqual:       " = ",
	NumberOperatorNotEqual:    " <> ",
	NumberOperatorLessThan:    " < ",
	NumberOperatorAtLeast:     " <= ",
	NumberOperatorGreaterThan: " > ",
	NumberOperatorAtMost:      " >= ",
}

func writeNumberOperation[T Number](builder *statementBuilder, col Column, op NumberOperator, value T) {
	col.writeTo(builder)
	builder.WriteString(numberOperators[op])
	builder.WriteString(builder.appendArg(value))
}

type Boolean interface {
	~bool
}

type BooleanOperator uint8

const (
	BooleanOperatorIsTrue BooleanOperator = iota + 1
	BooleanOperatorIsFalse
)

func writeBooleanOperation[T Boolean](builder *statementBuilder, col Column, value T) {
	col.writeTo(builder)
	builder.WriteString(" IS ")
	builder.WriteString(builder.appendArg(value))
}
