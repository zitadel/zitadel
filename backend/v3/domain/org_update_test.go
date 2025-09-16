package domain_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestUpdateOrgCommand_Execute(t *testing.T) {
	txInitErr := errors.New("tx init error")
	updateErr := errors.New("update error")

	tt := []struct {
		testName string

		queryExecutor func(ctrl *gomock.Controller) database.QueryExecutor
		orgRepo       func(ctrl *gomock.Controller) func(database.QueryExecutor) domain.OrganizationRepository

		inputID   string
		inputName string

		expectedError error
	}{
		{
			testName: "when EnsureTx fails should return error",
			queryExecutor: func(ctrl *gomock.Controller) database.QueryExecutor {
				mockDB := dbmock.NewMockPool(ctrl)
				mockDB.EXPECT().
					Begin(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, txInitErr)
				return mockDB
			},
			expectedError: txInitErr,
		},
		{
			testName: "when org update fails should return error",
			orgRepo: func(ctrl *gomock.Controller) func(database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Update(gomock.Any(), repo.IDCondition("org-1"), "instance-1", repo.SetName("test org update")).
					Return(int64(0), updateErr).
					AnyTimes()
				return func(_ database.QueryExecutor) domain.OrganizationRepository {
					return repo
				}
			},
			inputID:       "org-1",
			inputName:     "test org update",
			expectedError: updateErr,
		},
		{
			testName:  "when org update returns 0 rows updated should return not found error",
			inputID:   "org-1",
			inputName: "test org update",
			orgRepo: func(ctrl *gomock.Controller) func(database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Update(gomock.Any(), repo.IDCondition("org-1"), "instance-1", repo.SetName("test org update")).
					Return(int64(0), nil).
					AnyTimes()
				return func(_ database.QueryExecutor) domain.OrganizationRepository {
					return repo
				}
			},
			expectedError: domain.NewOrgNotFoundError("DOM-7PfSUn"),
		},
		{
			testName: "when org update returns more than 1 row updated should return internal error",
			orgRepo: func(ctrl *gomock.Controller) func(database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Update(gomock.Any(), repo.IDCondition("org-1"), "instance-1", repo.SetName("test org update")).
					Return(int64(2), nil).
					AnyTimes()
				return func(_ database.QueryExecutor) domain.OrganizationRepository {
					return repo
				}
			},
			inputID:       "org-1",
			inputName:     "test org update",
			expectedError: domain.NewMultipleOrgsUpdatedError("DOM-QzITrx", 1, 2),
		},
		{
			testName: "when org update returns 1 row updated should return no error and set cache",
			orgRepo: func(ctrl *gomock.Controller) func(database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Update(gomock.Any(), repo.IDCondition("org-1"), "instance-1", repo.SetName("test org update")).
					Return(int64(1), nil).
					AnyTimes()
				return func(_ database.QueryExecutor) domain.OrganizationRepository {
					return repo
				}
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
			ctrl := gomock.NewController(t)
			cmd := &domain.UpdateOrgCommand{
				ID:   tc.inputID,
				Name: tc.inputName,
			}

			opts := &domain.CommandOpts{
				DB: new(noopdb.Pool),
			}
			if tc.orgRepo != nil {
				opts.SetOrgRepo(tc.orgRepo(ctrl))
			}
			if tc.queryExecutor != nil {
				opts.DB = tc.queryExecutor(ctrl)
			}

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
