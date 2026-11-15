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

type testPagination struct {
	offset uint64
	limit  uint32
	asc    bool
}

func (p *testPagination) GetOffset() uint64 {
	if p == nil {
		return 0
	}
	return p.offset
}

func (p *testPagination) GetLimit() uint32 {
	if p == nil {
		return 0
	}
	return p.limit
}

func (p *testPagination) GetAsc() bool {
	if p == nil {
		return false
	}
	return p.asc
}

func TestSystemDefaults_PaginationToQuery(t *testing.T) {
	t.Parallel()

	defaults := SystemDefaults{DefaultQueryLimit: 100, MaxQueryLimit: 1000}
	tests := []struct {
		name       string
		pagination *testPagination
		wantOffset uint64
		wantLimit  uint64
		wantAsc    bool
		wantErr    error
	}{
		{
			name:       "nil pagination, default limit",
			pagination: nil,
			wantLimit:  100,
		},
		{
			name:       "all fields set",
			pagination: &testPagination{offset: 10, limit: 50, asc: true},
			wantOffset: 10,
			wantLimit:  50,
			wantAsc:    true,
		},
		{
			name:       "limit exceeds max",
			pagination: &testPagination{offset: 10, limit: 1001, asc: true},
			wantErr:    zerrors.ThrowInvalidArgument(nil, "QUERY-4M0fs", "Errors.Query.LimitExceeded"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			offset, limit, asc, err := defaults.PaginationToQuery(tt.pagination)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantOffset, offset)
			assert.Equal(t, tt.wantLimit, limit)
			assert.Equal(t, tt.wantAsc, asc)
		})
	}
}

func TestSystemDefaults_V1QueryLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		defaults  SystemDefaults
		requested uint64
		want      uint64
		wantErr   error
	}{
		{
			name:      "no limit requested, max applied",
			defaults:  SystemDefaults{DefaultQueryLimit: 100, MaxQueryLimit: 1000},
			requested: 0,
			want:      1000,
		},
		{
			name:      "limit below max",
			defaults:  SystemDefaults{DefaultQueryLimit: 100, MaxQueryLimit: 1000},
			requested: 500,
			want:      500,
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
			requested: 0,
			want:      0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.defaults.V1QueryLimit(tt.requested)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
