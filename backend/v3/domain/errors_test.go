package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestUnexpectedQueryTypeError_Error(t *testing.T) {
	tests := []struct {
		name     string
		typeVal  any
		expected string
	}{
		{
			name:     "string type",
			typeVal:  "test",
			expected: "Message=unexpected query type 'string'",
		},
		{
			name:     "int type",
			typeVal:  42,
			expected: "Message=unexpected query type 'int'",
		},
		{
			name:     "struct type",
			typeVal:  struct{}{},
			expected: "Message=unexpected query type 'struct {}'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewUnexpectedQueryTypeError(tt.typeVal)
			assert.Equal(t, tt.expected, err.Error())
		})
	}
}

func TestHandleUpdateError(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name             string
		inputErr         error
		expectedRowCount int64
		actualRowCount   int64
		errorID          string
		objectType       string
		expectedErr      error
	}{
		{
			name:             "no error and counts match",
			inputErr:         nil,
			expectedRowCount: 1,
			actualRowCount:   1,
			errorID:          "test-001",
			objectType:       "user",
			expectedErr:      nil,
		},
		{
			name:             "input error provided",
			inputErr:         errors.New("db error"),
			expectedRowCount: 1,
			actualRowCount:   1,
			errorID:          "test-002",
			objectType:       "session",
			expectedErr:      zerrors.ThrowInternalf(errors.New("db error"), "test-002", "failed updating %s", "session"),
		},
		{
			name:             "no rows affected",
			inputErr:         nil,
			expectedRowCount: 1,
			actualRowCount:   0,
			errorID:          "test-003",
			objectType:       "idp",
			expectedErr:      zerrors.ThrowNotFoundf(nil, "test-003", "%s not found", "idp"),
		},
		{
			name:             "unexpected number of rows updated",
			inputErr:         nil,
			expectedRowCount: 1,
			actualRowCount:   5,
			errorID:          "test-004",
			objectType:       "org",
			expectedErr:      zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(1, 5), "test-004", "unexpected number of rows updated"),
		},
		{
			name:             "counts mismatch",
			inputErr:         nil,
			expectedRowCount: 2,
			actualRowCount:   1,
			errorID:          "test-005",
			objectType:       "project",
			expectedErr:      zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(2, 1), "test-005", "unexpected number of rows updated"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := handleUpdateError(tc.inputErr, tc.expectedRowCount, tc.actualRowCount, tc.errorID, tc.objectType)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
