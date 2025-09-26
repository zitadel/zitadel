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

//go:generate  enumer -type NumberOperation,TextOperation,BytesOperation -linecomment -output ./operators_enumer.go
type Operation interface {
	NumberOperation | TextOperation | BytesOperation
}

type Text interface {
	~string | Bytes
}

// TextOperation are operations that can be performed on text values.
type TextOperation uint8

const (
	// TextOperationEqual compares two strings for equality.
	TextOperationEqual TextOperation = iota + 1 // =
	// TextOperationNotEqual compares two strings for inequality.
	TextOperationNotEqual // <>
	// TextOperationStartsWith checks if the first string starts with the second.
	TextOperationStartsWith // LIKE
)

func writeTextOperation[T Text](builder *StatementBuilder, col Column, op TextOperation, value any) {
	writeOperation[T](builder, col, op.String(), value)
	if op == TextOperationStartsWith {
		builder.WriteString(" || '%'")
	}
}

type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex | time.Time | time.Duration
}

// NumberOperation are operations that can be performed on number values.
type NumberOperation uint8

const (
	// NumberOperationEqual compares two numbers for equality.
	NumberOperationEqual NumberOperation = iota + 1 // =
	// NumberOperationNotEqual compares two numbers for inequality.
	NumberOperationNotEqual // <>
	// NumberOperationLessThan compares two numbers to check if the first is less than the second.
	NumberOperationLessThan // <
	// NumberOperationLessThanOrEqual compares two numbers to check if the first is less than or equal to the second.
	NumberOperationAtLeast // <=
	// NumberOperationGreaterThan compares two numbers to check if the first is greater than the second.
	NumberOperationGreaterThan // >
	// NumberOperationGreaterThanOrEqual compares two numbers to check if the first is greater than or equal to the second.
	NumberOperationAtMost // >=
)

func writeNumberOperation[T Number](builder *StatementBuilder, col Column, op NumberOperation, value any) {
	writeOperation[T](builder, col, op.String(), value)
}

type Boolean interface {
	~bool
}

func writeBooleanOperation[T Boolean](builder *StatementBuilder, col Column, value any) {
	writeOperation[T](builder, col, "=", value)
}

type Bytes interface {
	~[]byte
}

// BytesOperation are operations that can be performed on bytea values.
type BytesOperation uint8

const (
	BytesOperationEqual    BytesOperation = iota + 1 // =
	BytesOperationNotEqual                           // <>
)

func writeBytesOperation[T Bytes](builder *StatementBuilder, col Column, op BytesOperation, value any) {
	writeOperation[T](builder, col, op.String(), value)
}

func writeOperation[V Value](builder *StatementBuilder, col Column, op string, value any) {
	if op == "" {
		panic("unsupported operation")
	}

	switch value.(type) {
	case V, wrappedValue[V], *wrappedValue[V]:
	default:
		panic("unsupported value type")
	}
	col.WriteQualified(builder)
	builder.WriteRune(' ')
	builder.WriteString(op)
	builder.WriteRune(' ')
	builder.WriteArg(value)
}
