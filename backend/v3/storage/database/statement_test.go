package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatementBuilder_AppendArg(t *testing.T) {
	t.Run("same arg returns same placeholder", func(t *testing.T) {
		var b StatementBuilder
		placeholder1 := b.AppendArg("same")
		placeholder2 := b.AppendArg("same")
		assert.Equal(t, placeholder1, placeholder2)
		assert.Len(t, b.args, 1)
		assert.Len(t, b.existingArgs, 1)
	})

	t.Run("same arg different types", func(t *testing.T) {
		var b StatementBuilder
		placeholder1 := b.AppendArg("same")
		placeholder2 := b.AppendArg([]byte("same"))
		placeholder3 := b.AppendArg("same")
		assert.NotEqual(t, placeholder1, placeholder2)
		assert.Equal(t, placeholder1, placeholder3)
		assert.Len(t, b.args, 2)
		assert.Len(t, b.existingArgs, 2)
	})
}
