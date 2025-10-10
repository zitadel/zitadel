package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
)

func TestBaseCommand_OperationMapper(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name           string
		queryOperation object.TextQueryMethod

		expectedOperation database.TextOperation
		expectedError     error
	}{
		{
			name:              "contains",
			queryOperation:    object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS,
			expectedOperation: database.TextOperationContains,
		},
		{
			name:              "contains ignore case",
			queryOperation:    object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE,
			expectedOperation: database.TextOperationContainsIgnoreCase,
		},
		{
			name:              "ends with",
			queryOperation:    object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH,
			expectedOperation: database.TextOperationEndsWith,
		},
		{
			name:              "ends with ignore case",
			queryOperation:    object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE,
			expectedOperation: database.TextOperationEndsWithIgnoreCase,
		},
		{
			name:              "equals",
			queryOperation:    object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
			expectedOperation: database.TextOperationEqual,
		},
		{
			name:              "equals ignore case",
			queryOperation:    object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE,
			expectedOperation: database.TextOperationEqual,
		},
		{
			name:              "starts with",
			queryOperation:    object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH,
			expectedOperation: database.TextOperationStartsWith,
		},
		{
			name:              "starts with ignore case",
			queryOperation:    object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE,
			expectedOperation: database.TextOperationStartsWithIgnoreCase,
		},
		{
			name:           "unknown operation",
			queryOperation: object.TextQueryMethod(99),
			expectedError:  zerrors.ThrowInvalidArgument(domain.NewUnexpectedTextQueryOperationError(object.TextQueryMethod(99)), "DOM-iBRBVe", "List.Query.Invalid"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			b := &domain.BaseCommand{}
			got, err := b.TextOperationMapper(tc.queryOperation)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOperation, got)
		})
	}
}
func TestBaseCommand_Pagination(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name           string
		inputLimit     uint32
		inputOffset    uint64
		expectedLimit  uint32
		expectedOffset uint32
	}{
		{
			name: "when zero values should return empty options",
		},
		{
			name:           "normal values",
			inputLimit:     10,
			inputOffset:    20,
			expectedLimit:  10,
			expectedOffset: 20,
		},
		{
			name:           "large offset",
			inputLimit:     5,
			inputOffset:    1000000,
			expectedLimit:  5,
			expectedOffset: 1000000,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			t.Parallel()
			b := &domain.BaseCommand{}
			opts := &database.QueryOpts{}

			// Test
			limitOpt, offsetOpt := b.Pagination(tc.inputLimit, tc.inputOffset)
			limitOpt(opts)
			offsetOpt(opts)

			// Verify
			assert.Equal(t, tc.expectedLimit, opts.Limit)
			assert.Equal(t, tc.expectedOffset, opts.Offset)
		})
	}
}
