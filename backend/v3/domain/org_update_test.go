package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestUpdateOrgCommand_Execute(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName string

		expectations func(ctx context.Context, mockDB *dbmock.MockPool, txMock *dbmock.MockTransaction, args ...any)

		dbArgs []any

		inputID   string
		inputName string

		expectedError error
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
			inputID:   "org-1",
			inputName: "test org update",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			mockCtrl := gomock.NewController(t)
			cmd := &domain.UpdateOrgCommand{
				ID:   tc.inputID,
				Name: tc.inputName,
			}
			mockDB := dbmock.NewMockPool(mockCtrl)
			mockTX := dbmock.NewMockTransaction(mockCtrl)
			tc.expectations(ctx, mockDB, mockTX, tc.inputName, tc.inputID, "instance-1")

			opts := &domain.CommandOpts{
				DB:            mockDB,
				OrgRepository: repository.OrganizationRepository(mockTX),
			}

			// Test
			err := cmd.Execute(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
func TestUpdateOrgCommand_Validate(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name          string
		cmd           *domain.UpdateOrgCommand
		expectedError error
	}{
		{
			name:          "when no ID should return invalid argument error",
			cmd:           &domain.UpdateOrgCommand{ID: "", Name: "test-name"},
			expectedError: zerrors.ThrowInvalidArgument(nil, "DOM-lEMhVC", "invalid organization ID"),
		},
		{
			name:          "when no name shuld return invalid argument error",
			cmd:           &domain.UpdateOrgCommand{ID: "test-id", Name: ""},
			expectedError: zerrors.ThrowInvalidArgument(nil, "DOM-wfUntW", "invalid organization name"),
		},
		{
			name: "when validation succeeds should return no error",
			cmd:  &domain.UpdateOrgCommand{ID: "test-id", Name: "test-name"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.cmd.Validate()
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
