package database

import (
	"time"

	"golang.org/x/exp/constraints"
)

type Value interface {
	Boolean | Number | Text | Instruction
}

type Operation interface {
	BooleanOperation | NumberOperation | TextOperation
}

type Text interface {
	~string | ~[]byte
}

type TextOperation uint8

const (
	// TextOperationEqual compares two strings for equality.
	TextOperationEqual TextOperation = iota + 1
	// TextOperationEqualIgnoreCase compares two strings for equality, ignoring case.
	TextOperationEqualIgnoreCase
	// TextOperationNotEqual compares two strings for inequality.
	TextOperationNotEqual
	// TextOperationNotEqualIgnoreCase compares two strings for inequality, ignoring case.
	TextOperationNotEqualIgnoreCase
	// TextOperationStartsWith checks if the first string starts with the second.
	TextOperationStartsWith
	// TextOperationStartsWithIgnoreCase checks if the first string starts with the second, ignoring case.
	TextOperationStartsWithIgnoreCase
)

var textOperations = map[TextOperation]string{
	TextOperationEqual:                " = ",
	TextOperationEqualIgnoreCase:      " LIKE ",
	TextOperationNotEqual:             " <> ",
	TextOperationNotEqualIgnoreCase:   " NOT LIKE ",
	TextOperationStartsWith:           " LIKE ",
	TextOperationStartsWithIgnoreCase: " LIKE ",
}

func writeTextOperation[T Text](builder *StatementBuilder, col Column, op TextOperation, value T) {
	switch op {
	case TextOperationEqual, TextOperationNotEqual:
		col.Write(builder)
		builder.WriteString(textOperations[op])
		builder.WriteString(builder.AppendArg(value))
	case TextOperationEqualIgnoreCase, TextOperationNotEqualIgnoreCase:
		if ignoreCaseCol, ok := col.(ignoreCaseColumn); ok {
			ignoreCaseCol.WriteIgnoreCase(builder)
		} else {
			builder.WriteString("LOWER(")
			col.Write(builder)
			builder.WriteString(")")
		}
		builder.WriteString(textOperations[op])
		builder.WriteString("LOWER(")
		builder.WriteString(builder.AppendArg(value))
		builder.WriteString(")")
	case TextOperationStartsWith:
		col.Write(builder)
		builder.WriteString(textOperations[op])
		builder.WriteString(builder.AppendArg(value))
		builder.WriteString(" || '%'")
	case TextOperationStartsWithIgnoreCase:
		if ignoreCaseCol, ok := col.(ignoreCaseColumn); ok {
			ignoreCaseCol.WriteIgnoreCase(builder)
		} else {
			builder.WriteString("LOWER(")
			col.Write(builder)
			builder.WriteString(")")
		}
		builder.WriteString(textOperations[op])
		builder.WriteString("LOWER(")
		builder.WriteString(builder.AppendArg(value))
		builder.WriteString(")")
		builder.WriteString(" || '%'")
	default:
		panic("unsupported text operation")
	}
}

type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex | time.Time | time.Duration
}

type NumberOperation uint8

const (
	// NumberOperationEqual compares two numbers for equality.
	NumberOperationEqual NumberOperation = iota + 1
	// NumberOperationNotEqual compares two numbers for inequality.
	NumberOperationNotEqual
	// NumberOperationLessThan compares two numbers to check if the first is less than the second.
	NumberOperationLessThan
	// NumberOperationLessThanOrEqual compares two numbers to check if the first is less than or equal to the second.
	NumberOperationAtLeast
	// NumberOperationGreaterThan compares two numbers to check if the first is greater than the second.
	NumberOperationGreaterThan
	// NumberOperationGreaterThanOrEqual compares two numbers to check if the first is greater than or equal to the second.
	NumberOperationAtMost
)

var numberOperations = map[NumberOperation]string{
	NumberOperationEqual:       " = ",
	NumberOperationNotEqual:    " <> ",
	NumberOperationLessThan:    " < ",
	NumberOperationAtLeast:     " <= ",
	NumberOperationGreaterThan: " > ",
	NumberOperationAtMost:      " >= ",
}

func writeNumberOperation[T Number](builder *StatementBuilder, col Column, op NumberOperation, value T) {
	col.Write(builder)
	builder.WriteString(numberOperations[op])
	builder.WriteString(builder.AppendArg(value))
}

type Boolean interface {
	~bool
}

type BooleanOperation uint8

const (
	BooleanOperationIsTrue BooleanOperation = iota + 1
	BooleanOperationIsFalse
)

func writeBooleanOperation[T Boolean](builder *StatementBuilder, col Column, value T) {
	col.Write(builder)
	builder.WriteString(" IS ")
	builder.WriteString(builder.AppendArg(value))
}
