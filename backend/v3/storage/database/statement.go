package database

import (
	"strconv"
	"strings"
)

type Instruction string

const (
	NowInstruction  Instruction = "NOW()"
	NullInstruction Instruction = "NULL"
)

// StatementBuilder is a helper to build SQL statement.
type StatementBuilder struct {
	strings.Builder
	args         []any
	existingArgs map[any]string
}

// WriteArgs adds the argument to the statement and writes the placeholder to the query.
func (b *StatementBuilder) WriteArg(arg any) {
	b.WriteString(b.AppendArg(arg))
}

// AppebdArg adds the argument to the statement and returns the placeholder.
func (b *StatementBuilder) AppendArg(arg any) (placeholder string) {
	if b.existingArgs == nil {
		b.existingArgs = make(map[any]string)
	}

	if instruction, ok := arg.(Instruction); ok {
		return string(instruction)
	}

	b.args = append(b.args, arg)
	placeholder = "$" + strconv.Itoa(len(b.args))
	b.existingArgs[arg] = placeholder
	return placeholder
}

// AppendArgs adds the arguments to the statement and doesn't return the placeholders.
func (b *StatementBuilder) AppendArgs(args ...any) {
	for _, arg := range args {
		b.AppendArg(arg)
	}
}

// Args returns the arguments added to the statement.
func (b *StatementBuilder) Args() []any {
	return b.args
}
