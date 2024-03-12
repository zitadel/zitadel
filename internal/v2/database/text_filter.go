package database

import (
	"fmt"
	"strings"

	"github.com/zitadel/logging"
)

type TextFilter[T text] interface {
	Condition
}

type TextCondition[T text] struct {
	Filter[textCompare, T]
}

func NewTextEqual[T text](t T) *TextCondition[T] {
	return newTextFilter(textEqual, t)
}

func NewTextEqualInsensitive[T text](t T) *TextCondition[T] {
	return newTextFilter(textEqualInsensitive, t)
}

func NewTextStartsWith[T text](t T) *TextCondition[T] {
	return newTextFilter(textStartsWith, t)
}

func NewTextStartsWithInsensitive[T text](t T) *TextCondition[T] {
	return newTextFilter(textStartsWithInsensitive, t)
}

func NewTextEndsWith[T text](t T) *TextCondition[T] {
	return newTextFilter(textEndsWith, t)
}

func NewTextEndsWithInsensitive[T text](t T) *TextCondition[T] {
	return newTextFilter(textEndsWithInsensitive, t)
}

func NewTextContains[T text](t T) *TextCondition[T] {
	return newTextFilter(textContains, t)
}

func NewTextContainsInsensitive[T text](t T) *TextCondition[T] {
	return newTextFilter(textContainsInsensitive, t)
}

func newTextFilter[T text](comp textCompare, t T) *TextCondition[T] {
	return &TextCondition[T]{
		Filter: Filter[textCompare, T]{
			comp:  comp,
			value: t,
		},
	}
}

func (f TextCondition[T]) Write(stmt *Statement, columnName string) {
	if f.comp.isInsensitive() {
		f.writeCaseInsensitive(stmt, columnName)
		return
	}
	f.Filter.Write(stmt, columnName)
}

func (f *TextCondition[T]) writeCaseInsensitive(stmt *Statement, columnName string) {
	stmt.Builder.WriteString("LOWER(")
	stmt.Builder.WriteString(columnName)
	stmt.Builder.WriteString(") ")
	stmt.Builder.WriteString(f.comp.String())
	stmt.Builder.WriteString(" LOWER(")
	f.writeArg(stmt)
	stmt.Builder.WriteString(")")
}

func (f *TextCondition[T]) writeArg(stmt *Statement) {
	var v any = f.value
	// workaround for placeholder
	if placeholder, ok := v.(placeholder); ok {
		stmt.WriteArg(placeholder)
	}
	stmt.WriteArg(strings.ToLower(fmt.Sprint(f.value)))
}

type textCompare uint8

const (
	textEqual textCompare = iota
	textEqualInsensitive
	textStartsWith
	textStartsWithInsensitive
	textEndsWith
	textEndsWithInsensitive
	textContains
	textContainsInsensitive
)

func (c textCompare) String() string {
	switch c {
	case textEqual, textEqualInsensitive:
		return "="
	case textStartsWith, textStartsWithInsensitive, textEndsWith, textEndsWithInsensitive, textContains, textContainsInsensitive:
		return "LIKE"
	default:
		logging.WithFields("compare", c).Panic("comparison type not implemented")
		return ""
	}
}

func (c textCompare) isInsensitive() bool {
	return c == textEqualInsensitive ||
		c == textStartsWithInsensitive ||
		c == textEndsWithInsensitive ||
		c == textContainsInsensitive
}

type text interface {
	~string | placeholder
}
