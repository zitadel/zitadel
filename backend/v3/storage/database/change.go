package database

import (
	"reflect"

	"go.uber.org/mock/gomock"
)

// Change represents a change to a column in a database table.
// Its written in the SET clause of an UPDATE statement.
type Change interface {
	gomock.Matcher
	Write(builder *StatementBuilder)
}

type change[V Value] struct {
	column Column
	value  V
}

// Matches implements [gomock.Matcher].
func (c *change[V]) Matches(x any) bool {
	toMatch, ok := x.(*change[V])
	if !ok {
		return false
	}
	return c.column == toMatch.column && reflect.DeepEqual(c.value, toMatch.value)
}

// String implements [gomock.Matcher].
func (c *change[V]) String() string {
	return "database.Change"
}

var (
	_ Change         = (*change[string])(nil)
	_ gomock.Matcher = (*change[string])(nil)
)

func NewChange[V Value](col Column, value V) Change {
	return &change[V]{
		column: col,
		value:  value,
	}
}

func NewChangePtr[V Value](col Column, value *V) Change {
	if value == nil {
		return NewChange(col, NullInstruction)
	}
	return NewChange(col, *value)
}

// Write implements [Change].
func (c change[V]) Write(builder *StatementBuilder) {
	c.column.WriteUnqualified(builder)
	builder.WriteString(" = ")
	builder.WriteArg(c.value)
}

type Changes []Change

func NewChanges(cols ...Change) Change {
	return Changes(cols)
}

// Write implements [Change].
func (m Changes) Write(builder *StatementBuilder) {
	for i, col := range m {
		if i > 0 {
			builder.WriteString(", ")
		}
		col.Write(builder)
	}
}

// Matches implements [gomock.Matcher].
func (c Changes) Matches(x any) bool {
	toMatch, ok := x.(*Changes)
	if !ok {
		return false
	}
	if len(c) != len(*toMatch) {
		return false
	}
	for i := range c {
		if !c[i].Matches((*toMatch)[i]) {
			return false
		}
	}
	return true
}

// String implements [gomock.Matcher].
func (c Changes) String() string {
	return "database.Change"
}

var _ Change = Changes(nil)
