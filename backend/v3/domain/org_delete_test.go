package domain_test

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	t.Parallel()
	ctx := authz.NewMockContext("inst-1", "org-default", gofakeit.UUID(), authz.WithMockProjectID("prj-1"))
	getErr := errors.New("get error")

	tt := []struct {
		testName            string
		orgRepo             func(ctrl *gomock.Controller) domain.OrganizationRepository
		projectRepo         func(ctrl *gomock.Controller) domain.ProjectRepository
		inputOrganizationID string
		expectedError       error
	}{
		{
			testName: "validate delete default org, precondition failed",

			inputOrganizationID: "org-default",
			expectedError:       zerrors.ThrowPreconditionFailed(nil, "DOM-LCkE69", "Errors.Org.DefaultOrgNotDeletable"),
		},
		{
			testName:            "when fetching project fails with NON precondition error should return error",
			inputOrganizationID: "org-1",
			projectRepo: func(ctrl *gomock.Controller) domain.ProjectRepository {
				repo := domainmock.NewProjectRepo(ctrl)

				repo.EXPECT().Get(gomock.Any(), gomock.Any(),
					dbmock.QueryOptions(database.WithCondition(
						database.And(
							repo.IDCondition("prj-1"),
							repo.OrganizationIDCondition("org-1"),
							repo.InstanceIDCondition("inst-1"),
						),
					))).
					Times(1).
					Return(nil, getErr)

				return repo
			},
			expectedError: getErr,
		},
		{
			testName:            "when fetching project succeeds should return precondition failed error",
			inputOrganizationID: "org-1",
			projectRepo: func(ctrl *gomock.Controller) domain.ProjectRepository {
				repo := domainmock.NewProjectRepo(ctrl)

				repo.EXPECT().Get(gomock.Any(), gomock.Any(),
					dbmock.QueryOptions(database.WithCondition(
						database.And(
							repo.IDCondition("prj-1"),
							repo.OrganizationIDCondition("org-1"),
							repo.InstanceIDCondition("inst-1"),
						),
					))).
					Times(1).
					Return(&domain.Project{}, nil)

				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-X7YXxC", "Errors.Org.ZitadelOrgNotDeletable"),
		},
		{
			testName: "when fetching organization fails with not found error should return not found error",
			projectRepo: func(ctrl *gomock.Controller) domain.ProjectRepository {
				repo := domainmock.NewProjectRepo(ctrl)

				repo.EXPECT().Get(gomock.Any(), gomock.Any(),
					dbmock.QueryOptions(database.WithCondition(
						database.And(
							repo.IDCondition("prj-1"),
							repo.OrganizationIDCondition("org-1"),
							repo.InstanceIDCondition("inst-1"),
						),
					))).
					Times(1).
					Return(nil, database.NewNoRowFoundError(getErr))

				return repo
			},
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(
						repo.PrimaryKeyCondition("inst-1", "org-1"),
					))).
					Times(1).
					Return(nil, database.NewNoRowFoundError(getErr))
				return repo
			},
			inputOrganizationID: "org-1",
			expectedError:       zerrors.ThrowNotFound(database.NewNoRowFoundError(getErr), "DOM-8KYOH3", "Errors.Org.NotFound"),
		},
		{
			testName: "when organization is active should validate successfully",
			projectRepo: func(ctrl *gomock.Controller) domain.ProjectRepository {
				repo := domainmock.NewProjectRepo(ctrl)

				repo.EXPECT().Get(gomock.Any(), gomock.Any(),
					dbmock.QueryOptions(database.WithCondition(
						database.And(
							repo.IDCondition("prj-1"),
							repo.OrganizationIDCondition("org-1"),
							repo.InstanceIDCondition("inst-1"),
						),
					))).
					Times(1).
					Return(nil, database.NewNoRowFoundError(getErr))

				return repo
			},
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(
						repo.PrimaryKeyCondition("inst-1", "org-1"),
					))).
					Times(1).
					Return(&domain.Organization{
						ID:    "org-1",
						State: domain.OrgStateActive,
					}, nil)
				return repo
			},
			inputOrganizationID: "org-1",
		},
		{
			testName: "when organization is inactive should validate successfully",
			projectRepo: func(ctrl *gomock.Controller) domain.ProjectRepository {
				repo := domainmock.NewProjectRepo(ctrl)

				repo.EXPECT().Get(gomock.Any(), gomock.Any(),
					dbmock.QueryOptions(database.WithCondition(
						database.And(
							repo.IDCondition("prj-1"),
							repo.OrganizationIDCondition("org-1"),
							repo.InstanceIDCondition("inst-1"),
						),
					))).
					Times(1).
					Return(nil, database.NewNoRowFoundError(getErr))

				return repo
			},
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(
						repo.PrimaryKeyCondition("inst-1", "org-1"),
					))).
					Times(1).
					Return(&domain.Organization{
						ID:    "org-1",
						State: domain.OrgStateInactive,
					}, nil)
				return repo
			},
			inputOrganizationID: "org-1",
		},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			d := domain.NewDeleteOrgCommand(tc.inputOrganizationID)
			ctrl := gomock.NewController(t)
			opts := &domain.InvokeOpts{}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)

			if tc.orgRepo != nil {
				domain.WithOrganizationRepo(tc.orgRepo(ctrl))(opts)
			}
			if tc.projectRepo != nil {
				domain.WithProjectRepo(tc.projectRepo(ctrl))(opts)
			}

			// Test
			err := d.Validate(ctx, opts)

			// Verify
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestDeleteOrgCommand_Execute(t *testing.T) {
	t.Parallel()

	ctx := authz.NewMockContext("inst-1", "org-1", gofakeit.UUID())
	deleteErr := errors.New("delete error")
	getErr := errors.New("get error")

	tt := []struct {
		testName string
		mockTx   func(ctrl *gomock.Controller) database.QueryExecutor
		orgRepo  func(ctrl *gomock.Controller) domain.OrganizationRepository

		inputOrganizationID string

		expectedError   error
		expectedOrgName string
	}{
		{
			testName: "when retrieving organization fails should return error",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)

				repo.EXPECT().
					LoadDomains().
					Times(1).
					Return(repo)

				repo.EXPECT().
					Get(
						gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(
							repo.PrimaryKeyCondition("inst-1", "org-1"),
						)),
					).
					Times(1).
					Return(nil, getErr)
				return repo
			},
			inputOrganizationID: "org-1",
			expectedError:       getErr,
		},
		{
			testName: "when delete organization fails should return error",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)

				repo.EXPECT().
					LoadDomains().
					Times(1).
					Return(repo)

				repo.EXPECT().
					Get(
						gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(
							repo.PrimaryKeyCondition("inst-1", "org-1"),
						)),
					).
					Times(1).
					Return(&domain.Organization{
						ID:   "org-1",
						Name: "organization 1",
					}, nil)
				repo.EXPECT().
					Delete(gomock.Any(), gomock.Any(),
						repo.PrimaryKeyCondition("inst-1", "org-1"),
					).
					Times(1).
					Return(int64(0), deleteErr)
				return repo
			},
			inputOrganizationID: "org-1",
			expectedError:       deleteErr,
			expectedOrgName:     "organization 1",
		},
		{
			testName: "when more than one row deleted should return internal error",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)

				repo.EXPECT().
					LoadDomains().
					Times(1).
					Return(repo)

				repo.EXPECT().
					Get(
						gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(
							repo.PrimaryKeyCondition("inst-1", "org-1"),
						)),
					).
					Times(1).
					Return(&domain.Organization{
						ID:   "org-1",
						Name: "organization 1",
					}, nil)
				repo.EXPECT().
					Delete(gomock.Any(), gomock.Any(), repo.PrimaryKeyCondition("inst-1", "org-1")).
					Times(1).
					Return(int64(2), nil)
				return repo
			},
			inputOrganizationID: "org-1",
			expectedError:       zerrors.ThrowInternalf(nil, "DOM-5cE9u6", "expecting 1 row deleted, got %d", 2),
			expectedOrgName:     "organization 1",
		},
		{
			testName: "when no rows deleted should return not found error",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)

				repo.EXPECT().
					LoadDomains().
					Times(1).
					Return(repo)

				repo.EXPECT().
					Get(
						gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(
							repo.PrimaryKeyCondition("inst-1", "org-1"),
						)),
					).
					Times(1).
					Return(&domain.Organization{
						ID:   "org-1",
						Name: "organization 1",
					}, nil)
				repo.EXPECT().
					Delete(gomock.Any(), gomock.Any(), repo.PrimaryKeyCondition("inst-1", "org-1")).
					Times(1).
					Return(int64(0), nil)
				return repo
			},
			inputOrganizationID: "org-1",
			expectedError:       zerrors.ThrowNotFoundf(nil, "DOM-ur6Qyv", "organization not found"),
			expectedOrgName:     "organization 1",
		},
		{
			testName: "when one row deleted should execute successfully",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)

				repo.EXPECT().
					LoadDomains().
					Times(1).
					Return(repo)

				repo.EXPECT().
					Get(
						gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(
							repo.PrimaryKeyCondition("inst-1", "org-1"),
						)),
					).
					Times(1).
					Return(&domain.Organization{
						ID:   "org-1",
						Name: "organization 1",
					}, nil)
				repo.EXPECT().
					Delete(gomock.Any(), gomock.Any(), repo.PrimaryKeyCondition("inst-1", "org-1")).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			inputOrganizationID: "org-1",
			expectedOrgName:     "organization 1",
		},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			cmd := domain.NewDeleteOrgCommand(tc.inputOrganizationID)
			ctrl := gomock.NewController(t)
			opts := &domain.InvokeOpts{}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.mockTx != nil {
				domain.WithQueryExecutor(tc.mockTx(ctrl))(opts)
			}
			if tc.orgRepo != nil {
				domain.WithOrganizationRepo(tc.orgRepo(ctrl))(opts)
			}

			// Test
			err := opts.Invoke(ctx, cmd)

			// Verify
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedOrgName, cmd.OrganizationName)
		})
	}
}

// TODO(IAM-Marco): Expand these tests once the needed repositories (policies, org settings, idp links and entities) are available
func TestDeleteOrgCommand_Events(t *testing.T) {
	t.Parallel()
	ctx := authz.NewMockContext("inst-1", "org-1", gofakeit.UUID())

	tt := []struct {
		testName      string
		mockTx        func(ctrl *gomock.Controller) database.QueryExecutor
		command       *domain.DeleteOrgCommand
		expectedError error
		expectedCount int
	}{
		{
			testName: "should create org removed event",
			command: &domain.DeleteOrgCommand{
				ID:               "org-1",
				OrganizationName: "org name",
				Domains: []*domain.OrganizationDomain{
					{Domain: "domain1.com"},
					{Domain: "domain2.com"},
				},
			},
			expectedCount: 1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// Given
			ctrl := gomock.NewController(t)
			opts := &domain.InvokeOpts{}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)

			if tc.mockTx != nil {
				domain.WithQueryExecutor(tc.mockTx(ctrl))(opts)
			}

			// Test
			cmds, err := tc.command.Events(ctx, opts)

			// Verify
			require.Equal(t, tc.expectedError, err)
			assert.Len(t, cmds, tc.expectedCount)
		})
	}
}
