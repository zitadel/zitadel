package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatementBuilder_AppendArg(t *testing.T) {
	t.Run("same arg returns same placeholder", func(t *testing.T) {
		var b StatementBuilder
		placeholder1 := b.AppendArg("same")
		placeholder2 := b.AppendArg("same")
		assert.Equal(t, placeholder1, placeholder2)
		assert.Len(t, b.Args(), 1)
		assert.Len(t, b.existingArgs, 1)
	})

	t.Run("same arg different types", func(t *testing.T) {
		var b StatementBuilder
		placeholder1 := b.AppendArg("same")
		placeholder2 := b.AppendArg([]byte("same"))
		placeholder3 := b.AppendArg("same")
		assert.NotEqual(t, placeholder1, placeholder2)
		assert.Equal(t, placeholder1, placeholder3)
		assert.Len(t, b.Args(), 2)
		assert.Len(t, b.existingArgs, 2)
	})

	t.Run("Instruction args are always different", func(t *testing.T) {
		var b StatementBuilder
		placeholder1 := b.AppendArg(DefaultInstruction)
		placeholder2 := b.AppendArg(DefaultInstruction)
		assert.Equal(t, placeholder1, placeholder2)
		assert.Len(t, b.Args(), 0)
		assert.Len(t, b.existingArgs, 0)
	})
}

func TestStatementBuilder_AppendArgs(t *testing.T) {
	t.Run("same arg returns same placeholder", func(t *testing.T) {
		var b StatementBuilder
		b.AppendArgs("same", "same")
		assert.Len(t, b.Args(), 1)
		assert.Len(t, b.existingArgs, 1)
	})

	t.Run("same arg different types", func(t *testing.T) {
		var b StatementBuilder
		b.AppendArgs("same", []byte("same"), "same")
		assert.Len(t, b.Args(), 2)
		assert.Len(t, b.existingArgs, 2)
	})

	t.Run("Instruction args are always different", func(t *testing.T) {
		var b StatementBuilder
		b.AppendArgs(DefaultInstruction, DefaultInstruction)
		assert.Len(t, b.Args(), 0)
		assert.Len(t, b.existingArgs, 0)
	})
}

func TestStatementBuilder_WriteArg(t *testing.T) {
	for _, test := range []struct {
		name    string
		arg     any
		wantSQL string
		wantArg []any
	}{
		{
			name:    "primitive arg",
			arg:     "test",
			wantSQL: "$1",
			wantArg: []any{"test"},
		},
		{
			name:    "argWriter arg",
			arg:     SHA256Value("wrapped"),
			wantSQL: "SHA256($1)",
			wantArg: []any{"wrapped"},
		},
		{
			name:    "Instruction arg",
			arg:     DefaultInstruction,
			wantSQL: "DEFAULT",
			wantArg: []any{},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var b StatementBuilder
			b.WriteArg(test.arg)
			assert.Equal(t, test.wantSQL, b.String())
			require.Len(t, b.Args(), len(test.wantArg))
			for i := range test.wantArg {
				assert.Equal(t, test.wantArg[i], b.Args()[i])
			}
		})
	}
}

func TestStatementBuilder_WriteArgs(t *testing.T) {
	for _, test := range []struct {
		name    string
		args    []any
		wantSQL string
		wantArg []any
	}{
		{
			name:    "primitive args",
			args:    []any{"test", 123, true, uint32(123)},
			wantSQL: "$1, $2, $3, $4",
			wantArg: []any{"test", 123, true, uint32(123)},
		},
		{
			name:    "argWriter args",
			args:    []any{SHA256Value("wrapped"), LowerValue("ASDF")},
			wantSQL: "SHA256($1), LOWER($2)",
			wantArg: []any{"wrapped", "ASDF"},
		},
		{
			name:    "Instruction args",
			args:    []any{DefaultInstruction, NowInstruction, NullInstruction},
			wantSQL: "DEFAULT, NOW(), NULL",
			wantArg: []any{},
		},
		{
			name:    "mixed args",
			args:    []any{123, uint32(123), NowInstruction, NullInstruction, SHA256Value("wrapped"), LowerValue("ASDF")},
			wantSQL: "$1, $2, NOW(), NULL, SHA256($3), LOWER($4)",
			wantArg: []any{123, uint32(123), "wrapped", "ASDF"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var b StatementBuilder
			b.WriteArgs(test.args...)
			assert.Equal(t, test.wantSQL, b.String())
			require.Len(t, b.Args(), len(test.wantArg))
			for i := range test.wantArg {
				assert.Equal(t, test.wantArg[i], b.Args()[i])
			}
		})
	}
}
