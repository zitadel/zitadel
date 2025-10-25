package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/muhlemmer/gu"
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

func TestUpdateOrgCommand_Execute(t *testing.T) {
	t.Parallel()

	txInitErr := errors.New("tx init error")
	getErr := errors.New("get error")
	updateErr := errors.New("update error")

	tt := []struct {
		testName string

		queryExecutor func(ctrl *gomock.Controller) database.QueryExecutor
		orgRepo       func(ctrl *gomock.Controller) domain.OrganizationRepository

		inputID   string
		inputName string

		expectedError          error
		expectedOldDomainName  *string
		expectedDomainVerified *bool
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
			testName: "when retrieving org fails should return error",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					LoadDomains().
					Times(1).
					Return(repo)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							database.And(
								repo.IDCondition("org-1"),
								repo.InstanceIDCondition("instance-1"),
							),
						))).
					Times(1).
					Return(nil, getErr)
				return repo
			},
			inputID:       "org-1",
			inputName:     "test org update",
			expectedError: getErr,
		},
		{
			testName: "when setting domain info fails should return error",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					LoadDomains().
					Times(1).
					Return(repo)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(
						database.And(
							repo.IDCondition("org-1"),
							repo.InstanceIDCondition("instance-1"),
						),
					))).
					Times(1).
					Return(&domain.Organization{
						ID:    "org-1",
						Name:  "",
						State: domain.OrgStateActive,
						Domains: []*domain.OrganizationDomain{
							{IsPrimary: true, IsVerified: true, Domain: "old org name"},
							{IsPrimary: false, IsVerified: true, Domain: "old org name"},
							{IsPrimary: false, IsVerified: true, Domain: "old primary org name"},
						},
					}, nil)
				return repo
			},
			inputID:       "org-1",
			inputName:     "test org update",
			expectedError: zerrors.ThrowInvalidArgument(nil, "ORG-RrfXY", "Errors.Org.Domain.EmptyString"),
		},
		{
			testName: "when org update fails should return error",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					LoadDomains().
					Times(1).
					Return(repo)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(
						database.And(
							repo.IDCondition("org-1"),
							repo.InstanceIDCondition("instance-1"),
						),
					))).
					Times(1).
					Return(&domain.Organization{
						ID:         "org-1",
						Name:       "old org name",
						InstanceID: "instance-1",
						State:      domain.OrgStateActive,
						Domains: []*domain.OrganizationDomain{
							{IsPrimary: true, IsVerified: true, Domain: "old-org-name."},
							{IsPrimary: false, IsVerified: true, Domain: "old-org-name-2."},
						},
					}, nil)

				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						database.And(
							repo.IDCondition("org-1"),
							repo.InstanceIDCondition("instance-1"),
						),
						repo.SetName("test org update"),
					).
					Times(1).
					Return(int64(0), updateErr)
				return repo
			},
			inputID:       "org-1",
			inputName:     "test org update",
			expectedError: updateErr,
		},
		{
			testName:  "when org update returns 0 rows updated should return not found error",
			inputID:   "org-1",
			inputName: "test org update",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					LoadDomains().
					Times(1).
					Return(repo)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(
						database.And(
							repo.IDCondition("org-1"),
							repo.InstanceIDCondition("instance-1"),
						),
					))).
					Times(1).
					Return(&domain.Organization{
						ID:         "org-1",
						Name:       "old org name",
						InstanceID: "instance-1",
						State:      domain.OrgStateActive,
						Domains: []*domain.OrganizationDomain{
							{IsPrimary: true, IsVerified: true, Domain: "old-org-name."},
							{IsPrimary: false, IsVerified: true, Domain: "old-org-name."},
							{IsPrimary: false, IsVerified: true, Domain: "old-org-name-2."},
						},
					}, nil)

				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						database.And(
							repo.IDCondition("org-1"),
							repo.InstanceIDCondition("instance-1"),
						),
						repo.SetName("test org update")).
					Times(1).
					Return(int64(0), nil)
				return repo
			},
			expectedError:          zerrors.ThrowNotFound(nil, "DOM-7PfSUn", "Errors.Org.NotFound"),
			expectedOldDomainName:  gu.Ptr("old-org-name."),
			expectedDomainVerified: gu.Ptr(true),
		},
		{
			testName: "when org update returns more than 1 row updated should return internal error",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					LoadDomains().
					Times(1).
					Return(repo)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(
						database.And(
							repo.IDCondition("org-1"),
							repo.InstanceIDCondition("instance-1"),
						),
					))).
					Times(1).
					Return(&domain.Organization{
						ID:         "org-1",
						Name:       "old org name",
						InstanceID: "instance-1",
						State:      domain.OrgStateActive,
						Domains: []*domain.OrganizationDomain{
							{IsPrimary: true, IsVerified: true, Domain: "old-org-name."},
							{IsPrimary: false, IsVerified: true, Domain: "old-org-name."},
							{IsPrimary: false, IsVerified: true, Domain: "old-org-name-2."},
						},
					}, nil)

				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						database.And(
							repo.IDCondition("org-1"),
							repo.InstanceIDCondition("instance-1"),
						),
						repo.SetName("test org update")).
					Times(1).
					Return(int64(2), nil)
				return repo
			},
			inputID:                "org-1",
			inputName:              "test org update",
			expectedError:          zerrors.ThrowInternal(domain.NewMultipleObjectsUpdatedError(1, 2), "DOM-QzITrx", "unexpected number of rows updated"),
			expectedOldDomainName:  gu.Ptr("old-org-name."),
			expectedDomainVerified: gu.Ptr(true),
		},
		{
			testName: "when org update returns 1 row updated should return no error and set non-primary verified domain",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					LoadDomains().
					Times(1).
					Return(repo)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(
						database.And(
							repo.IDCondition("org-1"),
							repo.InstanceIDCondition("instance-1"),
						),
					))).
					Times(1).
					Return(&domain.Organization{
						ID:         "org-1",
						Name:       "old org name",
						InstanceID: "instance-1",
						State:      domain.OrgStateActive,
						Domains: []*domain.OrganizationDomain{
							{IsPrimary: true, IsVerified: true, Domain: "old-org-name."},
							{IsPrimary: false, IsVerified: true, Domain: "old-org-name."},
							{IsPrimary: false, IsVerified: true, Domain: "old-org-name-2."},
						},
					}, nil)

				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						database.And(
							repo.IDCondition("org-1"),
							repo.InstanceIDCondition("instance-1"),
						),
						repo.SetName("test org update")).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			inputID:   "org-1",
			inputName: "test org update",

			expectedOldDomainName:  gu.Ptr("old-org-name."),
			expectedDomainVerified: gu.Ptr(true),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewUpdateOrgCommand(tc.inputID, tc.inputName)

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)

			if tc.orgRepo != nil {
				domain.WithOrganizationRepo(tc.orgRepo(ctrl))(opts)
			}
			if tc.queryExecutor != nil {
				domain.WithQueryExecutor(tc.queryExecutor(ctrl))(opts)
			}

			// Test
			err := opts.Invoke(ctx, cmd)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOldDomainName, cmd.OldDomainName)
			assert.Equal(t, tc.expectedDomainVerified, cmd.IsOldDomainVerified)
		})
	}
}

func TestUpdateOrgCommand_Validate(t *testing.T) {
	// t.Parallel()
	// txInitErr := errors.New("tx init error")
	getErr := errors.New("get error")

	tt := []struct {
		testName      string
		queryExecutor func(ctrl *gomock.Controller) database.QueryExecutor
		orgRepo       func(ctrl *gomock.Controller) domain.OrganizationRepository
		inputOrgID    string
		inputOrgName  string
		expectedError error
	}{
		{
			testName:      "when no ID should return invalid argument error",
			inputOrgID:    "",
			inputOrgName:  "test-name",
			expectedError: zerrors.ThrowInvalidArgument(nil, "DOM-lEMhVC", "invalid organization ID"),
		},
		{
			testName:      "when no name should return invalid argument error",
			inputOrgID:    "test-id",
			inputOrgName:  "",
			expectedError: zerrors.ThrowInvalidArgument(nil, "DOM-wfUntW", "invalid organization name"),
		},
		{
			testName: "when retrieving org fails should return error",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							repo.PrimaryKeyCondition("instance-1", "org-1"),
						))).
					Times(1).
					Return(nil, getErr)
				return repo
			},
			inputOrgID:    "org-1",
			inputOrgName:  "test org update",
			expectedError: getErr,
		},
		{
			testName: "when org name is not changed should return name not changed error",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(
						repo.PrimaryKeyCondition("instance-1", "org-1"),
					))).
					Times(1).
					Return(&domain.Organization{
						ID:   "org-1",
						Name: "test org update",
					}, nil)
				return repo
			},
			inputOrgID:    "org-1",
			inputOrgName:  "test org update",
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-nDzwIu", "Errors.Org.NotChanged"),
		},
		{
			testName: "when org name is changed should validate successfully and return no error",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(
						repo.PrimaryKeyCondition("instance-1", "org-1"),
					))).
					Times(1).
					Return(&domain.Organization{
						ID:   "org-1",
						Name: "old org name",
					}, nil)
				return repo
			},
			inputOrgID:   "org-1",
			inputOrgName: "test org update",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)

			cmd := domain.NewUpdateOrgCommand(tc.inputOrgID, tc.inputOrgName)

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(domain.NewValidatorInvoker(nil)),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.orgRepo != nil {
				domain.WithOrganizationRepo(tc.orgRepo(ctrl))(opts)
			}
			if tc.queryExecutor != nil {
				domain.WithQueryExecutor(tc.queryExecutor(ctrl))(opts)
			}
			err := cmd.Validate(ctx, opts)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateOrgCommand_Events(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		cmd           *domain.UpdateOrgCommand
		expectedCount int
	}{
		{
			name: "no old name, no events",
			cmd: &domain.UpdateOrgCommand{
				ID:            "org-1",
				Name:          "new-name",
				OldDomainName: nil,
			},
			expectedCount: 0,
		},
		{
			name: "old name equals new name, no events",
			cmd: &domain.UpdateOrgCommand{
				ID:            "org-1",
				Name:          "same-name",
				OldDomainName: gu.Ptr("same-name"),
			},
			expectedCount: 0,
		},
		{
			name: "old name different from new name, returns event",
			cmd: &domain.UpdateOrgCommand{
				ID:            "org-1",
				Name:          "new-name",
				OldDomainName: gu.Ptr("old-name"),
			},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			events, err := tt.cmd.Events(context.Background(), &domain.InvokeOpts{})
			require.Nil(t, err)
			assert.Len(t, events, tt.expectedCount)
		})
	}
}
