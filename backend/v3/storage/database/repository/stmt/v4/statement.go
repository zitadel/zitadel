package v4

import (
	"strconv"
	"strings"
)

type databaseInstruction string

const (
	nowDBInstruction  databaseInstruction = "NOW()"
	nullDBInstruction databaseInstruction = "NULL"
)

type statementBuilder struct {
	strings.Builder
	args         []any
	existingArgs map[any]string
}

func (b *statementBuilder) writeArg(arg any) {
	b.WriteString(b.appendArg(arg))
}

func (b *statementBuilder) appendArg(arg any) (placeholder string) {
	if b.existingArgs == nil {
		b.existingArgs = make(map[any]string)
	}
	if placeholder, ok := b.existingArgs[arg]; ok {
		return placeholder
	}
	if instruction, ok := arg.(databaseInstruction); ok {
		return string(instruction)
	}

	b.args = append(b.args, arg)
	placeholder = "$" + strconv.Itoa(len(b.args))
	b.existingArgs[arg] = placeholder
	return placeholder
}

func (b *statementBuilder) appendArgs(args ...any) {
	for _, arg := range args {
		b.appendArg(arg)
	}
}
