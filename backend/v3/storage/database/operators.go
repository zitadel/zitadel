package database

import (
	"time"

	"golang.org/x/exp/constraints"
)

type wrappedValue[V Value] struct {
	value V
	fn    function
}

func LowerValue[T Value](v T) wrappedValue[T] {
	return wrappedValue[T]{value: v, fn: functionLower}
}

func SHA256Value[T Value](v T) wrappedValue[T] {
	return wrappedValue[T]{value: v, fn: functionSHA256}
}

func (b wrappedValue[V]) WriteArg(builder *StatementBuilder) {
	builder.Grow(len(b.fn) + 5)
	builder.WriteString(string(b.fn))
	builder.WriteRune('(')
	builder.WriteArg(b.value)
	builder.WriteRune(')')
}

var _ argWriter = (*wrappedValue[string])(nil)

type Value interface {
	Boolean | Number | Text | Instruction | Bytes
}

type Operation interface {
	BooleanOperation | NumberOperation | TextOperation | BytesOperation
}

type Text interface {
	~string | Bytes
}

// TextOperation are operations that can be performed on text values.
type TextOperation uint8

const (
	// TextOperationEqual compares two strings for equality.
	TextOperationEqual TextOperation = iota + 1
	// TextOperationNotEqual compares two strings for inequality.
	TextOperationNotEqual
	// TextOperationStartsWith checks if the first string starts with the second.
	TextOperationStartsWith
)

var textOperations = map[TextOperation]string{
	TextOperationEqual:      " = ",
	TextOperationNotEqual:   " <> ",
	TextOperationStartsWith: " LIKE ",
}

func writeTextOperation[T Text](builder *StatementBuilder, col Column, op TextOperation, value any) {
	switch value.(type) {
	case T, argWriter:
	default:
		panic("unsupported text value type")
	}

	switch op {
	case TextOperationEqual, TextOperationNotEqual:
		col.WriteQualified(builder)
		builder.WriteString(textOperations[op])
		builder.WriteArg(value)
	case TextOperationStartsWith:
		col.WriteQualified(builder)
		builder.WriteString(textOperations[op])
		builder.WriteArg(value)
		builder.WriteString(" || '%'")
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

type Bytes interface {
	~[]byte
}

// BytesOperation are operations that can be performed on bytea values.
type BytesOperation uint8

const (
	BytesOperationEqual BytesOperation = iota + 1
	BytesOperationNotEqual
)

var bytesOperations = map[BytesOperation]string{
	BytesOperationEqual:    " = ",
	BytesOperationNotEqual: " <> ",
}

func writeBytesOperation[B Bytes](builder *StatementBuilder, col Column, op BytesOperation, value any) {
	col.WriteQualified(builder)
	builder.WriteString(bytesOperations[op])
	switch value.(type) {
	case B, argWriter:
		builder.WriteArg(value)
		return
	default:
		panic("unsupported bytes value type")
	}
}
