package object

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/zerrors"
	v2_object "github.com/zitadel/zitadel/pkg/grpc/object/v2"
)

func TestTextQueryMethodToTextOperation(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name           string
		queryOperation v2_object.TextQueryMethod

		expectedOperation database.TextOperation
		expectedError     error
	}{
		{
			name:              "contains",
			queryOperation:    v2_object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS,
			expectedOperation: database.TextOperationContains,
		},
		{
			name:              "contains ignore case",
			queryOperation:    v2_object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE,
			expectedOperation: database.TextOperationContainsIgnoreCase,
		},
		{
			name:              "ends with",
			queryOperation:    v2_object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH,
			expectedOperation: database.TextOperationEndsWith,
		},
		{
			name:              "ends with ignore case",
			queryOperation:    v2_object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE,
			expectedOperation: database.TextOperationEndsWithIgnoreCase,
		},
		{
			name:              "equals",
			queryOperation:    v2_object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
			expectedOperation: database.TextOperationEqual,
		},
		{
			name:              "equals ignore case",
			queryOperation:    v2_object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE,
			expectedOperation: database.TextOperationEqualIgnoreCase,
		},
		{
			name:              "starts with",
			queryOperation:    v2_object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH,
			expectedOperation: database.TextOperationStartsWith,
		},
		{
			name:              "starts with ignore case",
			queryOperation:    v2_object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE,
			expectedOperation: database.TextOperationStartsWithIgnoreCase,
		},
		{
			name:           "unknown operation",
			queryOperation: v2_object.TextQueryMethod(99),
			expectedError:  zerrors.ThrowInvalidArgument(nil, "OBJ-iBRBVe", "invalid text query method"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := TextQueryMethodToTextOperation(tc.queryOperation)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOperation, got)
		})
	}
}
