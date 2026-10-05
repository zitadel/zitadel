package systemdefaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestSystemDefaults_QueryLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		defaults  SystemDefaults
		requested uint64
		want      uint64
		wantErr   error
	}{
		{
			name:      "no limit requested, default applied",
			defaults:  SystemDefaults{DefaultQueryLimit: 100, MaxQueryLimit: 1000},
			requested: 0,
			want:      100,
		},
		{
			name:      "limit below max",
			defaults:  SystemDefaults{DefaultQueryLimit: 100, MaxQueryLimit: 1000},
			requested: 500,
			want:      500,
		},
		{
			name:      "limit equals max",
			defaults:  SystemDefaults{DefaultQueryLimit: 100, MaxQueryLimit: 1000},
			requested: 1000,
			want:      1000,
		},
		{
			name:      "limit exceeds max",
			defaults:  SystemDefaults{DefaultQueryLimit: 100, MaxQueryLimit: 1000},
			requested: 1001,
			wantErr:   zerrors.ThrowInvalidArgument(nil, "QUERY-4M0fs", "Errors.Query.LimitExceeded"),
		},
		{
			name:      "no max configured",
			defaults:  SystemDefaults{DefaultQueryLimit: 100},
			requested: 5000,
			want:      5000,
		},
		{
			name:      "no default configured",
			defaults:  SystemDefaults{MaxQueryLimit: 1000},
			requested: 0,
			want:      0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.defaults.QueryLimit(tt.requested)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
