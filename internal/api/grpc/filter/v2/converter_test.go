package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
)

func TestTextMethodPbToQuery(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name   string
		input  filter.TextFilterMethod
		output query.TextComparison
	}{
		{"Equals", filter.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS, query.TextEquals},
		{"EqualsIgnoreCase", filter.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS_IGNORE_CASE, query.TextEqualsIgnoreCase},
		{"StartsWith", filter.TextFilterMethod_TEXT_FILTER_METHOD_STARTS_WITH, query.TextStartsWith},
		{"StartsWithIgnoreCase", filter.TextFilterMethod_TEXT_FILTER_METHOD_STARTS_WITH_IGNORE_CASE, query.TextStartsWithIgnoreCase},
		{"Contains", filter.TextFilterMethod_TEXT_FILTER_METHOD_CONTAINS, query.TextContains},
		{"ContainsIgnoreCase", filter.TextFilterMethod_TEXT_FILTER_METHOD_CONTAINS_IGNORE_CASE, query.TextContainsIgnoreCase},
		{"EndsWith", filter.TextFilterMethod_TEXT_FILTER_METHOD_ENDS_WITH, query.TextEndsWith},
		{"EndsWithIgnoreCase", filter.TextFilterMethod_TEXT_FILTER_METHOD_ENDS_WITH_IGNORE_CASE, query.TextEndsWithIgnoreCase},
		{"Unknown", filter.TextFilterMethod(999), -1},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := TextMethodPbToQuery(tc.input)

			assert.Equal(t, tc.output, got)
		})
	}
}

func TestTimestampMethodPbToQuery(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name   string
		input  filter.TimestampFilterMethod
		output query.TimestampComparison
	}{
		{"Equals", filter.TimestampFilterMethod_TIMESTAMP_FILTER_METHOD_EQUALS, query.TimestampEquals},
		{"Before", filter.TimestampFilterMethod_TIMESTAMP_FILTER_METHOD_BEFORE, query.TimestampLess},
		{"After", filter.TimestampFilterMethod_TIMESTAMP_FILTER_METHOD_AFTER, query.TimestampGreater},
		{"BeforeOrEquals", filter.TimestampFilterMethod_TIMESTAMP_FILTER_METHOD_BEFORE_OR_EQUALS, query.TimestampLessOrEquals},
		{"AfterOrEquals", filter.TimestampFilterMethod_TIMESTAMP_FILTER_METHOD_AFTER_OR_EQUALS, query.TimestampGreaterOrEquals},
		{"Unknown", filter.TimestampFilterMethod(999), -1},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := TimestampMethodPbToQuery(tc.input)

			assert.Equal(t, tc.output, got)
		})
	}
}

func TestByteMethodPbToQuery(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name   string
		input  filter.ByteFilterMethod
		output query.BytesComparison
	}{
		{"Equals", filter.ByteFilterMethod_BYTE_FILTER_METHOD_EQUALS, query.BytesEquals},
		{"NotEquals", filter.ByteFilterMethod_BYTE_FILTER_METHOD_NOT_EQUALS, query.BytesNotEquals},
		{"Unknown", filter.ByteFilterMethod(999), -1},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := ByteMethodPbToQuery(tc.input)

			assert.Equal(t, tc.output, got)
		})
	}
}

func TestPaginationPbToQuery(t *testing.T) {
	t.Parallel()

	defaults := systemdefaults.SystemDefaults{
		DefaultQueryLimit: 10,
		MaxQueryLimit:     100,
	}
	tt := []struct {
		name    string
		query   *filter.PaginationRequest
		wantOff uint64
		wantLim uint64
		wantAsc bool
		wantErr bool
	}{
		{
			name:    "nil query",
			query:   nil,
			wantOff: 0,
			wantLim: 10,
			wantAsc: false,
			wantErr: false,
		},
		{
			name:    "limit not set",
			query:   &filter.PaginationRequest{Offset: 5, Limit: 0, Asc: true},
			wantOff: 5,
			wantLim: 10,
			wantAsc: true,
			wantErr: false,
		},
		{
			name:    "limit set below max",
			query:   &filter.PaginationRequest{Offset: 2, Limit: 50, Asc: false},
			wantOff: 2,
			wantLim: 50,
			wantAsc: false,
			wantErr: false,
		},
		{
			name:    "limit exceeds max",
			query:   &filter.PaginationRequest{Offset: 1, Limit: 101, Asc: true},
			wantOff: 0,
			wantLim: 0,
			wantAsc: false,
			wantErr: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			off, lim, asc, err := PaginationPbToQuery(defaults, tc.query)

			require.Equal(t, tc.wantErr, err != nil)

			assert.Equal(t, tc.wantOff, off)
			assert.Equal(t, tc.wantLim, lim)
			assert.Equal(t, tc.wantAsc, asc)
		})
	}
}

func TestQueryToPaginationPb(t *testing.T) {
	t.Parallel()

	req := query.SearchRequest{Limit: 20}
	resp := query.SearchResponse{Count: 123}
	got := QueryToPaginationPb(req, resp)

	assert.Equal(t, req.Limit, got.AppliedLimit)
	assert.Equal(t, resp.Count, got.TotalResult)
}
