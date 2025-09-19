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

// TextOperation are operations that can be performed on text values.
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
	// TextOperationContains checks if the first string contains the second
	TextOperationContains
	// TextOperationContains checks if the first string contains the second, ignoring case.
	TextOperationContainsWithIgnoreCase
	// TextOperationEndsWith checks if the first string ends with the second.
	TextOperationEndsWith
	// TextOperationEndsWithIgnoreCase checks if the first string ends with the second, ignoring case.
	TextOperationEndsWithIgnoreCase
)

var textOperations = map[TextOperation]string{
	TextOperationEqual:                  " = ",
	TextOperationEqualIgnoreCase:        " LIKE ",
	TextOperationNotEqual:               " <> ",
	TextOperationNotEqualIgnoreCase:     " NOT LIKE ",
	TextOperationStartsWith:             " LIKE ",
	TextOperationStartsWithIgnoreCase:   " LIKE ",
	TextOperationContains:               " LIKE ",
	TextOperationContainsWithIgnoreCase: " ILIKE ",
	TextOperationEndsWith:               " LIKE ",
	TextOperationEndsWithIgnoreCase:     " LIKE ",
}

func writeTextOperation[T Text](builder *StatementBuilder, col Column, op TextOperation, value T) {
	switch op {
	case TextOperationEqual, TextOperationNotEqual:
		col.WriteQualified(builder)
		builder.WriteString(textOperations[op])
		builder.WriteArg(value)
	case TextOperationEqualIgnoreCase, TextOperationNotEqualIgnoreCase:
		builder.WriteString("LOWER(")
		col.WriteQualified(builder)
		builder.WriteString(")")

		builder.WriteString(textOperations[op])
		builder.WriteString("LOWER(")
		builder.WriteArg(value)
		builder.WriteString(")")
	case TextOperationStartsWith:
		col.WriteQualified(builder)
		builder.WriteString(textOperations[op])
		builder.WriteArg(value)
		builder.WriteString(" || '%'")
	case TextOperationStartsWithIgnoreCase:
		builder.WriteString("LOWER(")
		col.WriteQualified(builder)
		builder.WriteString(")")

		builder.WriteString(textOperations[op])
		builder.WriteString("LOWER(")
		builder.WriteArg(value)
		builder.WriteString(")")
		builder.WriteString(" || '%'")
	case TextOperationContains, TextOperationContainsWithIgnoreCase:
		col.WriteQualified(builder)
		builder.WriteString(textOperations[op])
		builder.WriteString("'%' || ")
		builder.WriteArg(value)
		builder.WriteString(" || '%'")
	case TextOperationEndsWith:
		col.WriteQualified(builder)
		builder.WriteString(textOperations[op])
		builder.WriteString("'%' || ")
		builder.WriteArg(value)
	case TextOperationEndsWithIgnoreCase:
		builder.WriteString("LOWER(")
		col.WriteQualified(builder)
		builder.WriteString(")")

		builder.WriteString(textOperations[op])
		builder.WriteString("'%' || ")
		builder.WriteString("LOWER(")
		builder.WriteArg(value)
		builder.WriteString(")")

	default:
		panic("unsupported text operation")
	}
}

type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex | time.Time | time.Duration
}

// NumberOperation are operations that can be performed on number values.
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
	col.WriteQualified(builder)
	builder.WriteString(numberOperations[op])
	builder.WriteArg(value)
}

type Boolean interface {
	~bool
}

// BooleanOperation are operations that can be performed on boolean values.
type BooleanOperation uint8

const (
	BooleanOperationIsTrue BooleanOperation = iota + 1
	BooleanOperationIsFalse
)

func writeBooleanOperation[T Boolean](builder *StatementBuilder, col Column, value T) {
	col.WriteQualified(builder)
	builder.WriteString(" = ")
	builder.WriteArg(value)
}
