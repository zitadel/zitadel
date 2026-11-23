package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_editorsCapacity(t *testing.T) {
	tests := []struct {
		name  string
		limit uint64
		want  uint64
	}{
		{
			name:  "no limit",
			limit: 0,
			want:  0,
		},
		{
			name:  "limit below max",
			limit: 20,
			want:  20,
		},
		{
			name:  "limit above max",
			limit: 1 << 32,
			want:  maxEditorsCapacity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, editorsCapacity(tt.limit))
		})
	}
}
