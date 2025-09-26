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
		orgRepo       func(ctrl *gomock.Controller) func(client database.QueryExecutor) domain.OrganizationRepository

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
			orgRepo: func(ctrl *gomock.Controller) func(client database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.IDCondition("org-1")))).
					Times(1).
					Return(nil, getErr)
				return func(_ database.QueryExecutor) domain.OrganizationRepository { return repo }
			},
			inputID:       "org-1",
			inputName:     "test org update",
			expectedError: getErr,
		},
		{
			testName: "when org name is not changed should return name not changed error",
			orgRepo: func(ctrl *gomock.Controller) func(client database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.IDCondition("org-1")))).
					Times(1).
					Return(&domain.Organization{
						ID:   "org-1",
						Name: "test org update",
					}, nil)
				return func(_ database.QueryExecutor) domain.OrganizationRepository { return repo }
			},
			inputID:       "org-1",
			inputName:     "test org update",
			expectedError: domain.NewOrgNameNotChangedError("DOM-nDzwIu"),
		},
		{
			testName: "when org is inactive should return not found error",
			orgRepo: func(ctrl *gomock.Controller) func(client database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.IDCondition("org-1")))).
					Times(1).
					Return(&domain.Organization{
						ID:    "org-1",
						Name:  "old org name",
						State: domain.OrgStateInactive,
					}, nil)
				return func(_ database.QueryExecutor) domain.OrganizationRepository { return repo }
			},
			inputID:       "org-1",
			inputName:     "test org update",
			expectedError: domain.NewOrgNotFoundError("DOM-OcA1jq"),
		},
		{
			testName: "when setting domain info fails should return error",
			orgRepo: func(ctrl *gomock.Controller) func(client database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.IDCondition("org-1")))).
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
				return func(_ database.QueryExecutor) domain.OrganizationRepository { return repo }
			},
			inputID:       "org-1",
			inputName:     "test org update",
			expectedError: zerrors.ThrowInvalidArgument(nil, "ORG-RrfXY", "Errors.Org.Domain.EmptyString"),
		},
		{
			testName: "when org update fails should return error",
			orgRepo: func(ctrl *gomock.Controller) func(client database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.IDCondition("org-1")))).
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
					Update(gomock.Any(), repo.IDCondition("org-1"), "instance-1", repo.SetName("test org update")).
					Times(1).
					Return(int64(0), updateErr)
				return func(_ database.QueryExecutor) domain.OrganizationRepository { return repo }
			},
			inputID:       "org-1",
			inputName:     "test org update",
			expectedError: updateErr,
		},
		{
			testName:  "when org update returns 0 rows updated should return not found error",
			inputID:   "org-1",
			inputName: "test org update",
			orgRepo: func(ctrl *gomock.Controller) func(client database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.IDCondition("org-1")))).
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
					Update(gomock.Any(), repo.IDCondition("org-1"), "instance-1", repo.SetName("test org update")).
					Times(1).
					Return(int64(0), nil)
				return func(_ database.QueryExecutor) domain.OrganizationRepository { return repo }
			},
			expectedError:          domain.NewOrgNotFoundError("DOM-7PfSUn"),
			expectedOldDomainName:  gu.Ptr("old-org-name."),
			expectedDomainVerified: gu.Ptr(true),
		},
		{
			testName: "when org update returns more than 1 row updated should return internal error",
			orgRepo: func(ctrl *gomock.Controller) func(client database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.IDCondition("org-1")))).
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
					Update(gomock.Any(), repo.IDCondition("org-1"), "instance-1", repo.SetName("test org update")).
					Times(1).
					Return(int64(2), nil)
				return func(_ database.QueryExecutor) domain.OrganizationRepository { return repo }
			},
			inputID:                "org-1",
			inputName:              "test org update",
			expectedError:          domain.NewMultipleOrgsUpdatedError("DOM-QzITrx", 1, 2),
			expectedOldDomainName:  gu.Ptr("old-org-name."),
			expectedDomainVerified: gu.Ptr(true),
		},
		{
			testName: "when org update returns 1 row updated should return no error and set non-primary verified domain",
			orgRepo: func(ctrl *gomock.Controller) func(client database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.IDCondition("org-1")))).
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
					Update(gomock.Any(), repo.IDCondition("org-1"), "instance-1", repo.SetName("test org update")).
					Times(1).
					Return(int64(1), nil)
				return func(_ database.QueryExecutor) domain.OrganizationRepository { return repo }
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

			// Test
			err := cmd.Execute(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOldDomainName, cmd.OldDomainName)
			assert.Equal(t, tc.expectedDomainVerified, cmd.IsOldDomainVerified)
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
			err := tc.cmd.Validate(context.Background(), &domain.CommandOpts{})
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
			events, err := tt.cmd.Events(context.Background(), &domain.CommandOpts{})
			require.Nil(t, err)
			assert.Len(t, events, tt.expectedCount)
		})
	}
}
