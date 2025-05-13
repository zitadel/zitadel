package v3

import (
	"fmt"
	"strings"
)

type statement[T object] struct {
	tbl       Table
	columns   []Column
	condition Condition

	builder      strings.Builder
	args         []any
	existingArgs map[any]string
}

func newStatement[O object](t Table) *statement[O] {
	var o O
	return &statement[O]{
		tbl:     t,
		columns: o.Columns(t),
	}
}

// Where implements [Query].
func (stmt *statement[T]) Where(condition Condition) {
	stmt.condition = condition
}

func (stmt *statement[T]) writeFrom() {
	stmt.builder.WriteString(" FROM ")
	stmt.builder.WriteString(stmt.tbl.Schema())
	stmt.builder.WriteRune('.')
	stmt.builder.WriteString(stmt.tbl.Name())
	if stmt.tbl.Alias() != "" {
		stmt.builder.WriteString(" AS ")
		stmt.builder.WriteString(stmt.tbl.Alias())
	}
}

func (stmt *statement[T]) writeCondition() {
	if stmt.condition == nil {
		return
	}
	stmt.builder.WriteString(" WHERE ")
	stmt.condition.writeOn(stmt)
}

// appendArg implements [statementBuilder].
func (stmt *statement[T]) appendArg(arg any) (placeholder string) {
	if stmt.existingArgs == nil {
		stmt.existingArgs = make(map[any]string)
	}
	if placeholder, ok := stmt.existingArgs[arg]; ok {
		return placeholder
	}

	stmt.args = append(stmt.args, arg)
	placeholder = fmt.Sprintf("$%d", len(stmt.args))
	stmt.existingArgs[arg] = placeholder
	return placeholder
}

// table implements [statementBuilder].
func (stmt *statement[T]) table() Table {
	return stmt.tbl
}

// write implements [statementBuilder].
func (stmt *statement[T]) write(data []byte) {
	stmt.builder.Write(data)
}

// writeRune implements [statementBuilder].
func (stmt *statement[T]) writeRune(r rune) {
	stmt.builder.WriteRune(r)
}

// writeString implements [statementBuilder].
func (stmt *statement[T]) writeString(s string) {
	stmt.builder.WriteString(s)
}

var _ statementBuilder = (*statement[Instance])(nil)
