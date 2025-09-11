package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestUpdateOrgCommand_Execute(t *testing.T) {
	tt := []struct {
		testName string

		expectations func(ctx context.Context, mockDB *dbmock.MockPool, txMock *dbmock.MockTransaction, args ...any)

		dbArgs []any

		inputID   string
		inputName string

		expectedError    error
		expectedCacheObj *Organization
	}{
		{
			testName: "when EnsureTx fails should return error",
			expectations: func(ctx context.Context, mockDB *dbmock.MockPool, _ *dbmock.MockTransaction, args ...any) {
				mockDB.EXPECT().
					Begin(ctx, gomock.Any()).
					Times(1).
					Return(nil, errors.New("mock tx init error"))
			},
			expectedError: errors.New("mock tx init error"),
		},
		{
			testName: "when org update fails should return error",
			expectations: func(ctx context.Context, mockDB *dbmock.MockPool, txMock *dbmock.MockTransaction, args ...any) {
				mockUpdateErr := errors.New("mock update failed")
				mockDB.EXPECT().
					Begin(ctx, gomock.Any()).
					Times(1).
					Return(txMock, nil)
				txMock.EXPECT().
					Exec(ctx, "UPDATE zitadel.organizations SET name = $1 WHERE (organizations.id = $2 AND organizations.instance_id = $3)", args...).
					Times(1).
					Return(int64(-1), mockUpdateErr)
				txMock.EXPECT().End(ctx, mockUpdateErr).Return(mockUpdateErr)
			},
			inputID:       "org-1",
			inputName:     "test org update",
			expectedError: errors.New("mock update failed"),
		},
		{
			testName: "when org update returns 0 rows updated should return not found error",
			expectations: func(ctx context.Context, mockDB *dbmock.MockPool, txMock *dbmock.MockTransaction, args ...any) {
				mockUpdateErr := zerrors.ThrowNotFound(nil, "DOM-7PfSUn", "organization not found")
				mockDB.EXPECT().
					Begin(ctx, gomock.Any()).
					Times(1).
					Return(txMock, nil)
				txMock.EXPECT().
					Exec(ctx, "UPDATE zitadel.organizations SET name = $1 WHERE (organizations.id = $2 AND organizations.instance_id = $3)", args...).
					Times(1).
					Return(int64(0), nil)
				txMock.EXPECT().End(ctx, mockUpdateErr).Return(mockUpdateErr)
			},
			inputID:       "org-1",
			inputName:     "test org update",
			expectedError: zerrors.ThrowNotFound(nil, "DOM-7PfSUn", "organization not found"),
		},
		{
			testName: "when org update returns more than 1 row updated should return internal error",
			expectations: func(ctx context.Context, mockDB *dbmock.MockPool, txMock *dbmock.MockTransaction, args ...any) {
				mockUpdateErr := zerrors.ThrowInternalf(nil, "DOM-QzITrx", "expecting 1 row updated, got %d", 2)
				mockDB.EXPECT().
					Begin(ctx, gomock.Any()).
					Times(1).
					Return(txMock, nil)
				txMock.EXPECT().
					Exec(ctx, "UPDATE zitadel.organizations SET name = $1 WHERE (organizations.id = $2 AND organizations.instance_id = $3)", args...).
					Times(1).
					Return(int64(2), nil)
				txMock.EXPECT().End(ctx, mockUpdateErr).Return(mockUpdateErr)
			},
			inputID:       "org-1",
			inputName:     "test org update",
			expectedError: zerrors.ThrowInternalf(nil, "DOM-QzITrx", "expecting 1 row updated, got %d", 2),
		},
		{
			testName: "when org update returns 1 row updated should return no error and set cache",
			expectations: func(ctx context.Context, mockDB *dbmock.MockPool, txMock *dbmock.MockTransaction, args ...any) {
				mockDB.EXPECT().
					Begin(ctx, gomock.Any()).
					Times(1).
					Return(txMock, nil)
				txMock.EXPECT().
					Exec(ctx, "UPDATE zitadel.organizations SET name = $1 WHERE (organizations.id = $2 AND organizations.instance_id = $3)", args...).
					Times(1).
					Return(int64(1), nil)
				txMock.EXPECT().End(ctx, nil).Return(nil)
			},
			inputID:          "org-1",
			inputName:        "test org update",
			expectedCacheObj: &Organization{ID: "org-1", Name: "test org update"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			mockCtrl := gomock.NewController(t)
			cmd := &UpdateOrgCommand{
				ID:   tc.inputID,
				Name: tc.inputName,
			}
			mockDB := dbmock.NewMockPool(mockCtrl)
			mockTX := dbmock.NewMockTransaction(mockCtrl)
			tc.expectations(ctx, mockDB, mockTX, tc.inputName, tc.inputID, "instance-1")

			opts := &CommandOpts{
				DB: mockDB,
			}

			// Test
			err := cmd.Execute(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			if tc.expectedCacheObj != nil {
				cachedOrg, found := orgCache.Get(ctx, orgCacheIndexID, tc.inputID)
				require.True(t, found)
				assert.Equal(t, tc.expectedCacheObj, cachedOrg)
			}
		})
	}
}
