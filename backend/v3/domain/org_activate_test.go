package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestActivateOrgCommand_Events(t *testing.T) {
	t.Parallel()
	// Given
	expected := []eventstore.Command{org.NewOrgReactivatedEvent(context.Background(), &org.NewAggregate("some-id").Aggregate)}
	activateCmd := domain.NewActivateOrgCommand("some-id")

	// Test
	actual, err := activateCmd.Events(context.Background(), &domain.InvokeOpts{})

	// Verify

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestActivateOrgCommand_Validate(t *testing.T) {
	t.Parallel()
	getErr := errors.New("get error")

	tt := []struct {
		testName   string
		inputOrgID string

		orgRepo func(ctrl *gomock.Controller) domain.OrganizationRepository

		expectedError error
	}{
		{
			testName:      "empty org id",
			expectedError: zerrors.ThrowInvalidArgument(nil, "DOM-hJuuAv", "invalid organization ID"),
		},
		{
			testName:      "whitespace org id",
			inputOrgID:    "   ",
			expectedError: zerrors.ThrowInvalidArgument(nil, "DOM-hJuuAv", "invalid organization ID"),
		},
		{
			testName: "when retrieving org fails with generic error should return error",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
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
			inputOrgID:    "org-1",
			expectedError: getErr,
		},
		{
			testName: "when retrieving org fails with not found error should return not found error",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							database.And(
								repo.IDCondition("org-1"),
								repo.InstanceIDCondition("instance-1"),
							),
						))).
					Times(1).
					Return(nil, database.NewNoRowFoundError(getErr))
				return repo
			},
			inputOrgID:    "org-1",
			expectedError: zerrors.ThrowNotFound(database.NewNoRowFoundError(getErr), "DOM-86HVfs", "Errors.Org.NotFound"),
		},
		{
			testName: "when org state is active should return already active error",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
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
						State: domain.OrgStateActive,
					}, nil)
				return repo
			},
			inputOrgID:    "org-1",
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-Ixfbxh", "Errors.Org.AlreadyActive"),
		},
		{
			testName: "when org state inactive should successfully validate organization",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
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
						State: domain.OrgStateInactive,
					}, nil)
				return repo
			},
			inputOrgID: "org-1",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			cmd := domain.NewActivateOrgCommand(tc.inputOrgID)
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			opts := &domain.InvokeOpts{
				DB: new(noopdb.Pool),
			}
			if tc.orgRepo != nil {
				domain.WithOrganizationRepo(tc.orgRepo(ctrl))(opts)
			}

			// Test
			err := cmd.Validate(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestActivateOrgCommand_Execute(t *testing.T) {
	t.Parallel()
	txInitErr := errors.New("tx init error")
	updateErr := errors.New("update error")

	tt := []struct {
		testName string

		queryExecutor func(ctrl *gomock.Controller) database.QueryExecutor
		orgRepo       func(ctrl *gomock.Controller) domain.OrganizationRepository

		inputID string

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
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						database.And(
							repo.IDCondition("org-1"),
							repo.InstanceIDCondition("instance-1"),
						),
						repo.SetState(domain.OrgStateActive),
					).
					Times(1).
					Return(int64(0), updateErr)
				return repo
			},
			inputID:       "org-1",
			expectedError: updateErr,
		},
		{
			testName: "when org update returns 0 rows updated should return not found error",
			inputID:  "org-1",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						database.And(
							repo.IDCondition("org-1"),
							repo.InstanceIDCondition("instance-1"),
						),
						repo.SetState(domain.OrgStateActive)).
					Times(1).
					Return(int64(0), nil)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(nil, "DOM-CGumXG", "Errors.Org.NotFound"),
		},
		{
			testName: "when org update returns more than 1 row updated should return internal error",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						database.And(
							repo.IDCondition("org-1"),
							repo.InstanceIDCondition("instance-1"),
						),
						repo.SetState(domain.OrgStateActive)).
					Times(1).
					Return(int64(2), nil)
				return repo
			},
			inputID:       "org-1",
			expectedError: zerrors.ThrowInternal(domain.NewMultipleObjectsUpdatedError(1, 2), "DOM-SEWCLp", "unexpected number of rows updated"),
		},
		{
			testName: "when org update returns 1 row updated should return no error",
			orgRepo: func(ctrl *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						database.And(
							repo.IDCondition("org-1"),
							repo.InstanceIDCondition("instance-1"),
						),
						repo.SetState(domain.OrgStateActive)).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			inputID: "org-1",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := &domain.ActivateOrgCommand{
				ID: tc.inputID,
			}

			opts := &domain.InvokeOpts{
				DB: new(noopdb.Pool),
			}
			if tc.orgRepo != nil {
				domain.WithOrganizationRepo(tc.orgRepo(ctrl))(opts)
			}
			if tc.queryExecutor != nil {
				opts.DB = tc.queryExecutor(ctrl)
			}

			// Test
			err := cmd.Execute(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
