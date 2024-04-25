package database

import (
	"fmt"
	"strings"

	"github.com/zitadel/logging"
)

type TextFilter[T text] struct {
	Filter[textCompare, T]
}

func NewTextEqual[T text](t T) *TextFilter[T] {
	return newTextFilter(textEqual, t)
}

func NewTextUnequal[T text](t T) *TextFilter[T] {
	return newTextFilter(textUnequal, t)
}

func NewTextEqualInsensitive[T text](t T) *TextFilter[string] {
	return newTextFilter(textEqualInsensitive, strings.ToLower(string(t)))
}

func NewTextUnequalInsensitive[T text](t T) *TextFilter[string] {
	return newTextFilter(textUnequalInsensitive, strings.ToLower(string(t)))
}

func NewTextStartsWith[T text](t T) *TextFilter[T] {
	return newTextFilter(textStartsWith, t)
}

func NewTextStartsWithInsensitive[T text](t T) *TextFilter[string] {
	return newTextFilter(textStartsWithInsensitive, strings.ToLower(string(t)))
}

func NewTextEndsWith[T text](t T) *TextFilter[T] {
	return newTextFilter(textEndsWith, t)
}

func NewTextEndsWithInsensitive[T text](t T) *TextFilter[string] {
	return newTextFilter(textEndsWithInsensitive, strings.ToLower(string(t)))
}

func NewTextContains[T text](t T) *TextFilter[T] {
	return newTextFilter(textContains, t)
}

func NewTextContainsInsensitive[T text](t T) *TextFilter[string] {
	return newTextFilter(textContainsInsensitive, strings.ToLower(string(t)))
}

func newTextFilter[T text](comp textCompare, t T) *TextFilter[T] {
	return &TextFilter[T]{
		Filter: Filter[textCompare, T]{
			comp:  comp,
			value: t,
		},
	}
}

func (f *TextFilter[T]) Write(stmt *Statement, columnName string) {
	if f.comp.isInsensitive() {
		f.writeCaseInsensitive(stmt, columnName)
		return
	}
	f.Filter.Write(stmt, columnName)
}

func (f *TextFilter[T]) writeCaseInsensitive(stmt *Statement, columnName string) {
	stmt.WriteString("LOWER(")
	stmt.WriteString(columnName)
	stmt.WriteString(") ")
	stmt.WriteString(f.comp.String())
	stmt.WriteRune(' ')
	f.writeArg(stmt)
}

func (f *TextFilter[T]) writeArg(stmt *Statement) {
	// TODO: condition must know if it's args are named parameters or not
	// var v any = f.value
	// workaround for placeholder
	// if placeholder, ok := v.(placeholder); ok {
	// 	stmt.Builder.WriteString(" LOWER(")
	// 	stmt.WriteArg(placeholder)
	// 	stmt.Builder.WriteString(")")
	// }
	stmt.WriteArg(strings.ToLower(fmt.Sprint(f.value)))
}

type textCompare uint8

const (
	textEqual textCompare = iota
	textUnequal
	textEqualInsensitive
	textUnequalInsensitive
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
	case textUnequal, textUnequalInsensitive:
		return "<>"
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
	~string
	// TODO: condition must know if it's args are named parameters or not
	// ~string | placeholder
}
