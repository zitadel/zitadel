package database

import (
	"database/sql"
	"strings"

	"github.com/zitadel/logging"
)

type TextFilter[T text] struct {
	Filter[textCompare, T]
}

func NewTextEqual[T text](t T) *TextFilter[T] {
	return newTextFilter(textEqual, t)
}
func NewTextEqualInsensitive[T text](t T) *TextFilter[T] {
	return newTextFilter(textEqualInsensitive, t)
}
func NewTextStartsWith[T text](t T) *TextFilter[T] {
	return newTextFilter(textStartsWith, t)
}
func NewTextStartsWithInsensitive[T text](t T) *TextFilter[T] {
	return newTextFilter(textStartsWithInsensitive, t)
}
func NewTextEndsWith[T text](t T) *TextFilter[T] {
	return newTextFilter(textEndsWith, t)
}
func NewTextEndsWithInsensitive[T text](t T) *TextFilter[T] {
	return newTextFilter(textEndsWithInsensitive, t)
}
func NewTextContains[T text](t T) *TextFilter[T] {
	return newTextFilter(textContains, t)
}
func NewTextContainsInsensitive[T text](t T) *TextFilter[T] {
	return newTextFilter(textContainsInsensitive, t)
}

func newTextFilter[T text](comp textCompare, t T) *TextFilter[T] {
	return &TextFilter[T]{
		Filter: Filter[textCompare, T]{
			comp:  comp,
			value: t,
		},
	}
}

func (f TextFilter[T]) Write(stmt *Statement, columnName string) {
	if f.comp.isInsensitive() {
		f.writeCaseInsensitive(stmt, columnName)
		return
	}
	f.Filter.Write(stmt, columnName)
}

func (f TextFilter[T]) WriteNamed(stmt *Statement, columnName string) {
	if f.comp.isInsensitive() {
		f.writeCaseInsensitiveNamed(stmt, columnName)
		return
	}
	f.Filter.WriteNamed(stmt, columnName)
}

func (f *TextFilter[T]) writeCaseInsensitive(stmt *Statement, columnName string) {
	prepareCaseInsensitiveWrite(stmt, columnName, f.comp)
	stmt.AppendArg(strings.ToLower(string(f.value)))
}

func (f *TextFilter[T]) writeCaseInsensitiveNamed(stmt *Statement, columnName string) {
	prepareCaseInsensitiveWrite(stmt, columnName, f.comp)
	stmt.AppendArg(sql.Named(columnName, strings.ToLower(string(f.value))))
}

func prepareCaseInsensitiveWrite(stmt *Statement, columnName string, comp textCompare) {
	stmt.Builder.WriteString("lower(")
	stmt.Builder.WriteString(columnName)
	stmt.Builder.WriteString(") ")
	stmt.Builder.WriteString(comp.String())
	stmt.Builder.WriteString(" ")
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
	~string
}
