package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
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
			expectedOperation: database.TextOperationEqualIgnoreCase,
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
			expectedError:  domain.NewUnexpectedTextQueryOperationError("DOM-iBRBVe", object.TextQueryMethod(99)),
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
