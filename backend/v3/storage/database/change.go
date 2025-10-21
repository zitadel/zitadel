package database

import "slices"

// Change represents a change to a column in a database table.
// Its written in the SET clause of an UPDATE statement.
type Change interface {
	// Write writes the change to the given statement builder.
	Write(builder *StatementBuilder)
	// IsOnColumn checks if the change is on the given column.
	IsOnColumn(col Column) bool
}

type change[V Value] struct {
	column Column
	value  V
}

var _ Change = (*change[string])(nil)

// NewChange creates a new Change for the given column and value.
// If you want to set a column to NULL, use [NewChangePtr].
func NewChange[V Value](col Column, value V) Change {
	return &change[V]{
		column: col,
		value:  value,
	}
}

// NewChangePtr creates a new Change for the given column and value pointer.
// If the value pointer is nil, the column will be set to NULL.
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

// IsOnColumn implements [Change].
func (c change[V]) IsOnColumn(col Column) bool {
	return c.column.Equals(col)
}

type Changes []Change

func NewChanges(cols ...Change) Change {
	return Changes(cols)
}

// IsOnColumn implements [Change].
func (c Changes) IsOnColumn(col Column) bool {
	return slices.ContainsFunc(c, func(change Change) bool {
		return change.IsOnColumn(col)
	})
}

// Write implements [Change].
func (m Changes) Write(builder *StatementBuilder) {
	for i, change := range m {
		if i > 0 {
			builder.WriteString(", ")
		}
		change.Write(builder)
	}
}

var _ Change = Changes(nil)
