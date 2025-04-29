package v4

type Change interface {
	Column
}

type change[V Value] struct {
	column Column
	value  V
}

func newChange[V Value](col Column, value V) Change {
	return &change[V]{
		column: col,
		value:  value,
	}
}

func newUpdatePtrColumn[V Value](col Column, value *V) Change {
	if value == nil {
		return newChange(col, nullDBInstruction)
	}
	return newChange(col, *value)
}

// writeTo implements [Change].
func (c change[V]) writeTo(builder *statementBuilder) {
	c.column.writeTo(builder)
	builder.WriteString(" = ")
	builder.writeArg(c.value)
}

type Changes []Change

func newChanges(cols ...Change) Change {
	return Changes(cols)
}

// writeTo implements [Change].
func (m Changes) writeTo(builder *statementBuilder) {
	for i, col := range m {
		if i > 0 {
			builder.WriteString(", ")
		}
		col.writeTo(builder)
	}
}

var _ Change = Changes(nil)

var _ Change = (*change[string])(nil)

type Column interface {
	writeTo(builder *statementBuilder)
}

type column struct {
	name string
}

func (c column) writeTo(builder *statementBuilder) {
	builder.WriteString(c.name)
}

type ignoreCaseColumn interface {
	Column
	writeIgnoreCaseTo(builder *statementBuilder)
}

type ignoreCaseCol struct {
	column
	suffix string
}

func (c ignoreCaseCol) writeIgnoreCaseTo(builder *statementBuilder) {
	c.column.writeTo(builder)
	builder.WriteString(c.suffix)
}
