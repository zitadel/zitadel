package change

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/zerrors"
	change_pb "github.com/zitadel/zitadel/pkg/grpc/change"
)

func TestChangeQueryToModel(t *testing.T) {
	defaults := systemdefaults.SystemDefaults{
		DefaultQueryLimit: 100,
		MaxQueryLimit:     1000,
	}
	tests := []struct {
		name         string
		query        *change_pb.ChangeQuery
		wantLimit    uint64
		wantSequence uint64
		wantAsc      bool
		wantErr      error
	}{
		{
			name:      "nil query, max limit as default",
			query:     nil,
			wantLimit: 1000,
		},
		{
			name:         "all fields set",
			query:        &change_pb.ChangeQuery{Sequence: 5, Limit: 30, Asc: true},
			wantLimit:    30,
			wantSequence: 5,
			wantAsc:      true,
		},
		{
			name:    "limit exceeds max",
			query:   &change_pb.ChangeQuery{Sequence: 5, Limit: 1001, Asc: true},
			wantErr: zerrors.ThrowInvalidArgument(nil, "QUERY-4M0fs", "Errors.Query.LimitExceeded"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit, sequence, asc, err := ChangeQueryToModel(defaults, tt.query)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantLimit, limit)
			assert.Equal(t, tt.wantSequence, sequence)
			assert.Equal(t, tt.wantAsc, asc)
		})
	}
}
