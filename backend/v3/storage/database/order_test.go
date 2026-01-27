package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_orderBy_Write(t *testing.T) {
	tests := []struct {
		name  string
		want  string
		order Order
	}{
		{
			name:  "order by column",
			want:  " ORDER BY table.column",
			order: OrderBy(NewColumn("table", "column")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder StatementBuilder
			tt.order.Write(&builder)
			assert.Equal(t, tt.want, builder.String())
		})
	}
}
