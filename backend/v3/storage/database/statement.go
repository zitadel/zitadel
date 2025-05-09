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

type StatementBuilder struct {
	strings.Builder
	args         []any
	existingArgs map[any]string
}

func (b *StatementBuilder) WriteArg(arg any) {
	b.WriteString(b.AppendArg(arg))
}

func (b *StatementBuilder) AppendArg(arg any) (placeholder string) {
	if b.existingArgs == nil {
		b.existingArgs = make(map[any]string)
	}
	if placeholder, ok := b.existingArgs[arg]; ok {
		return placeholder
	}
	if instruction, ok := arg.(Instruction); ok {
		return string(instruction)
	}

	b.args = append(b.args, arg)
	placeholder = "$" + strconv.Itoa(len(b.args))
	b.existingArgs[arg] = placeholder
	return placeholder
}

func (b *StatementBuilder) AppendArgs(args ...any) {
	for _, arg := range args {
		b.AppendArg(arg)
	}
}

func (b *StatementBuilder) Args() []any {
	return b.args
}
