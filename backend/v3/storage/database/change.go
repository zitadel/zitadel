package database

// Change represents a change to a column in a database table.
// Its written in the SET clause of an UPDATE statement.
type Change interface {
	Write(builder *StatementBuilder)
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

type Changes []Change

func NewChanges(cols ...Change) Change {
	return Changes(cols)
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
