package signal

import (
	"testing"

	object "github.com/zitadel/zitadel/pkg/grpc/object/v2"
)

func TestListQueryToOffsetLimit(t *testing.T) {
	tests := []struct {
		name       string
		query      *object.ListQuery
		wantOffset int
		wantLimit  int
	}{
		{
			name:       "nil query uses defaults",
			query:      nil,
			wantOffset: 0,
			wantLimit:  20,
		},
		{
			name:       "zero limit defaults to 20",
			query:      &object.ListQuery{Offset: 0, Limit: 0},
			wantOffset: 0,
			wantLimit:  20,
		},
		{
			name:       "explicit values preserved",
			query:      &object.ListQuery{Offset: 50, Limit: 100},
			wantOffset: 50,
			wantLimit:  100,
		},
		{
			name:       "limit capped at 1000",
			query:      &object.ListQuery{Offset: 0, Limit: 5000},
			wantOffset: 0,
			wantLimit:  1000,
		},
		{
			name:       "negative limit defaults to 20",
			query:      &object.ListQuery{Offset: 0, Limit: 0},
			wantOffset: 0,
			wantLimit:  20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset, limit := listQueryToOffsetLimit(tt.query)
			if offset != tt.wantOffset {
				t.Errorf("offset = %d, want %d", offset, tt.wantOffset)
			}
			if limit != tt.wantLimit {
				t.Errorf("limit = %d, want %d", limit, tt.wantLimit)
			}
		})
	}
}
