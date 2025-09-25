package domain_test

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
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

func TestDeleteOrgCommand_Validate(t *testing.T) {
	ctx := authz.NewMockContext("inst-1", "org-default", gofakeit.UUID())
	txInitErr := errors.New("tx init error")
	getErr := errors.New("get error")

	tt := []struct {
		testName string
		mockTx   func(ctrl *gomock.Controller) database.QueryExecutor
		orgRepo  func(ctrl *gomock.Controller) func(client database.QueryExecutor) domain.OrganizationRepository
		// projectRepo       func(ctrl *gomock.Controller) func(client database.QueryExecutor) domain.ProjectRepository
		inputOrganizationID string
		expectedError       error
	}{
		{
			testName: "validate delete default org, precondition failed",

			inputOrganizationID: "org-default",
			expectedError:       zerrors.ThrowPreconditionFailed(nil, "DOM-LCkE69", "Errors.Org.DefaultOrgNotDeletable"),
		},
		{
			testName: "when EnsureTx fails should return error",
			mockTx: func(ctrl *gomock.Controller) database.QueryExecutor {
				mockDB := dbmock.NewMockPool(ctrl)
				mockDB.EXPECT().
					Begin(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, txInitErr)
				return mockDB
			},
			expectedError: txInitErr,
		},
		// TODO(IAM-Marco): Fix when relational table for projects is available
		// {
		// 	testName: "when fetching project fails with NON precondition error should return error",
		// },
		// TODO(IAM-Marco): Fix when relational table for projects is available
		// {
		// 	testName: "when fetching project succeeds should precondition error",
		// },
		{
			testName: "when fetching organization fails should return error",
			orgRepo: func(ctrl *gomock.Controller) func(client database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.IDCondition("org-1")))).
					Times(1).
					Return(nil, getErr)
				return func(_ database.QueryExecutor) domain.OrganizationRepository { return repo }
			},
			inputOrganizationID: "org-1",
			expectedError:       getErr,
		},
		{
			testName: "when organization is neither active nor inactive should return not found error",
			orgRepo: func(ctrl *gomock.Controller) func(client database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.IDCondition("org-1")))).
					Times(1).
					Return(&domain.Organization{
						ID:    "org-1",
						State: domain.OrgStateRemoved,
					}, nil)
				return func(_ database.QueryExecutor) domain.OrganizationRepository { return repo }
			},
			inputOrganizationID: "org-1",
			expectedError:       zerrors.ThrowNotFound(nil, "DOM-8KYOH3", "Errors.Org.NotFound"),
		},
		{
			testName: "when organization is active should validate successfully",
			orgRepo: func(ctrl *gomock.Controller) func(client database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.IDCondition("org-1")))).
					Times(1).
					Return(&domain.Organization{
						ID:    "org-1",
						State: domain.OrgStateActive,
					}, nil)
				return func(_ database.QueryExecutor) domain.OrganizationRepository { return repo }
			},
			inputOrganizationID: "org-1",
		},
		{
			testName: "when organization is inactive should validate successfully",
			orgRepo: func(ctrl *gomock.Controller) func(client database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.IDCondition("org-1")))).
					Times(1).
					Return(&domain.Organization{
						ID:    "org-1",
						State: domain.OrgStateInactive,
					}, nil)
				return func(_ database.QueryExecutor) domain.OrganizationRepository { return repo }
			},
			inputOrganizationID: "org-1",
		},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Given
			d := domain.NewDeleteOrgCommand(tc.inputOrganizationID)
			ctrl := gomock.NewController(t)
			opts := &domain.CommandOpts{DB: new(noopdb.Pool)}

			if tc.mockTx != nil {
				opts.DB = tc.mockTx(ctrl)
			}
			if tc.orgRepo != nil {
				opts.SetOrgRepo(tc.orgRepo(ctrl))
			}

			// Test
			err := d.Validate(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
