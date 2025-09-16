package database

import (
	"strconv"
	"strings"
)

type Instruction string

const (
	DefaultInstruction Instruction = "DEFAULT"
	NowInstruction     Instruction = "NOW()"
	NullInstruction    Instruction = "NULL"
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

// WriteArgs adds the arguments to the statement and writes the placeholders to the query.
// The placeholders are comma separated.
func (b *StatementBuilder) WriteArgs(args ...any) {
	for i, arg := range args {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteArg(arg)
	}
}

// AppendArg adds the argument to the statement and returns the placeholder.
func (b *StatementBuilder) AppendArg(arg any) (placeholder string) {
	if b.existingArgs == nil {
		b.existingArgs = make(map[any]string)
	}
	// the key is used to work around the following panic:
	// panic: runtime error: hash of unhashable type []uint8
	key := arg
	if argBytes, ok := arg.([]uint8); ok {
		key = string(argBytes)
	}
	if placeholder, ok := b.existingArgs[key]; ok {
		return placeholder
	}
	if instruction, ok := arg.(Instruction); ok {
		return string(instruction)
	}

	b.args = append(b.args, arg)
	placeholder = "$" + strconv.Itoa(len(b.args))
	b.existingArgs[key] = placeholder
	return placeholder
}

// AppendArgs adds the arguments to the statement and doesn't return the placeholders.
// If an argument is already added, it will not be added again.
func (b *StatementBuilder) AppendArgs(args ...any) {
	for _, arg := range args {
		b.AppendArg(arg)
	}
}

// Args returns the arguments added to the statement.
func (b *StatementBuilder) Args() []any {
	return b.args
}
