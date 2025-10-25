package domain_test

import (
	"context"
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
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/milestone"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestDeleteInstanceCommand_Validate(t *testing.T) {
	t.Parallel()
	permissionErr := errors.New("permission error")
	tt := []struct {
		name              string
		permissionChecker func(ctrl *gomock.Controller) domain.PermissionChecker
		inputID           string
		expectedError     error
	}{
		{
			name:          "empty id",
			inputID:       "",
			expectedError: zerrors.ThrowInvalidArgument(nil, "DOM-VpQ9lF", "Errors.Invalid.Argument"),
		},
		{
			name:          "whitespace id",
			inputID:       "   ",
			expectedError: zerrors.ThrowInvalidArgument(nil, "DOM-VpQ9lF", "Errors.Invalid.Argument"),
		},
		{
			name: "when user is missing permission should return permission denied",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)

				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.InstanceWritePermission).
					Times(1).
					Return(permissionErr)

				return permChecker
			},
			expectedError: zerrors.ThrowPermissionDenied(permissionErr, "DOM-Yz8f1X", "permission denied"),
			inputID:       "instance-1",
		},
		{
			name: "valid id",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)

				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.InstanceWritePermission).
					Times(1).
					Return(nil)

				return permChecker
			},
			inputID: "instance-1",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Given
			d := &domain.DeleteInstanceCommand{ID: tc.inputID}

			cmdOpts := &domain.InvokeOpts{}
			if tc.permissionChecker != nil {
				ctrl := gomock.NewController(t)
				cmdOpts.Permissions = tc.permissionChecker(ctrl)
			}

			// Test
			err := d.Validate(context.Background(), cmdOpts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteInstanceCommand_Events(t *testing.T) {
	t.Parallel()
	// Given
	cmd := &domain.DeleteInstanceCommand{
		ID:              "instance-1",
		InstanceName:    "instance-name",
		InstanceDomains: []string{"domain1.com", "domain2.com"},
	}
	expectedEvents := []eventstore.Command{
		instance.NewInstanceRemovedEvent(context.Background(), &instance.NewAggregate(cmd.ID).Aggregate, cmd.InstanceName, cmd.InstanceDomains),
		milestone.NewReachedEvent(context.Background(), milestone.NewInstanceAggregate(cmd.ID), milestone.InstanceDeleted),
	}

	// Test
	events, err := cmd.Events(context.Background(), nil)

	// Verify
	assert.NoError(t, err)
	assert.Len(t, events, len(expectedEvents))
	assert.IsType(t, &instance.InstanceRemovedEvent{}, events[0])
	assert.IsType(t, &milestone.ReachedEvent{}, events[1])
}

func TestDeleteInstanceCommand_Execute(t *testing.T) {
	t.Parallel()

	ctx := authz.NewMockContext("inst-1", "org-1", gofakeit.UUID())
	txInitErr := errors.New("tx init error")
	deleteErr := errors.New("delete error")
	getErr := errors.New("get error")

	tests := []struct {
		testName     string
		mockTx       func(ctrl *gomock.Controller) database.QueryExecutor
		instanceRepo func(ctrl *gomock.Controller) domain.InstanceRepository

		inputInstanceID string

		expectedError           error
		expectedInstanceName    string
		expectedInstanceDomains []string
	}{
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
		{
			testName: "when retrieving instance fails should return error",
			instanceRepo: func(ctrl *gomock.Controller) domain.InstanceRepository {
				instanceRepo := domainmock.NewInstanceRepo(ctrl)

				instanceRepo.EXPECT().
					LoadDomains().
					Times(1).
					Return(instanceRepo)

				instanceRepo.EXPECT().
					Get(
						gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(
							instanceRepo.IDCondition("instance-1"),
						)),
					).
					Times(1).
					Return(nil, getErr)
				return instanceRepo
			},
			inputInstanceID: "instance-1",
			expectedError:   getErr,
		},
		{
			testName: "when delete instance fails should return error",
			instanceRepo: func(ctrl *gomock.Controller) domain.InstanceRepository {
				instanceRepo := domainmock.NewInstanceRepo(ctrl)

				instanceRepo.EXPECT().
					LoadDomains().
					Times(1).
					Return(instanceRepo)

				instanceRepo.EXPECT().
					Get(
						gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(
							instanceRepo.IDCondition("instance-1"),
						)),
					).
					Times(1).
					Return(&domain.Instance{
						ID:   "instance-1",
						Name: "My instance 1",
						Domains: []*domain.InstanceDomain{
							{Domain: "d1.example.com"},
							{Domain: "d2.example.com"},
							{Domain: "d3.example.com"},
						},
					}, nil)

				instanceRepo.EXPECT().
					Delete(gomock.Any(), gomock.Any(), "instance-1").
					Times(1).
					Return(int64(0), deleteErr)
				return instanceRepo
			},
			inputInstanceID:         "instance-1",
			expectedError:           deleteErr,
			expectedInstanceName:    "My instance 1",
			expectedInstanceDomains: []string{"d1.example.com", "d2.example.com", "d3.example.com"},
		},
		{
			testName: "when more than one row deleted should return internal error",
			instanceRepo: func(ctrl *gomock.Controller) domain.InstanceRepository {
				instanceRepo := domainmock.NewInstanceRepo(ctrl)

				instanceRepo.EXPECT().
					LoadDomains().
					Times(1).
					Return(instanceRepo)

				instanceRepo.EXPECT().
					Get(
						gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(
							instanceRepo.IDCondition("instance-1"),
						)),
					).
					Times(1).
					Return(&domain.Instance{
						ID:   "instance-1",
						Name: "My instance 1",
						Domains: []*domain.InstanceDomain{
							{Domain: "d1.example.com"},
							{Domain: "d2.example.com"},
							{Domain: "d3.example.com"},
						},
					}, nil)

				instanceRepo.EXPECT().
					Delete(gomock.Any(), gomock.Any(), "instance-1").
					Times(1).
					Return(int64(2), nil)
				return instanceRepo
			},
			inputInstanceID:         "instance-1",
			expectedError:           zerrors.ThrowInternalf(nil, "DOM-Od04Jx", "expecting 1 row deleted, got %d", 2),
			expectedInstanceName:    "My instance 1",
			expectedInstanceDomains: []string{"d1.example.com", "d2.example.com", "d3.example.com"},
		},
		{
			testName: "when no rows deleted should return not found error",
			instanceRepo: func(ctrl *gomock.Controller) domain.InstanceRepository {
				instanceRepo := domainmock.NewInstanceRepo(ctrl)

				instanceRepo.EXPECT().
					LoadDomains().
					Times(1).
					Return(instanceRepo)

				instanceRepo.EXPECT().
					Get(
						gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(
							instanceRepo.IDCondition("instance-1"),
						)),
					).
					Times(1).
					Return(&domain.Instance{
						ID:   "instance-1",
						Name: "My instance 1",
						Domains: []*domain.InstanceDomain{
							{Domain: "d1.example.com"},
							{Domain: "d2.example.com"},
							{Domain: "d3.example.com"},
						},
					}, nil)

				instanceRepo.EXPECT().
					Delete(gomock.Any(), gomock.Any(), "instance-1").
					Times(1).
					Return(int64(0), nil)
				return instanceRepo
			},
			inputInstanceID:         "instance-1",
			expectedError:           zerrors.ThrowNotFound(nil, "DOM-daglwD", "instance not found"),
			expectedInstanceName:    "My instance 1",
			expectedInstanceDomains: []string{"d1.example.com", "d2.example.com", "d3.example.com"},
		},
		{
			testName: "when one row deleted should execute successfully",
			instanceRepo: func(ctrl *gomock.Controller) domain.InstanceRepository {
				instanceRepo := domainmock.NewInstanceRepo(ctrl)

				instanceRepo.EXPECT().
					LoadDomains().
					Times(1).
					Return(instanceRepo)

				instanceRepo.EXPECT().
					Get(
						gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(
							instanceRepo.IDCondition("instance-1"),
						)),
					).
					Times(1).
					Return(&domain.Instance{
						ID:   "instance-1",
						Name: "My instance 1",
						Domains: []*domain.InstanceDomain{
							{Domain: "d1.example.com"},
							{Domain: "d2.example.com"},
							{Domain: "d3.example.com"},
						},
					}, nil)

				instanceRepo.EXPECT().
					Delete(gomock.Any(), gomock.Any(), "instance-1").
					Times(1).
					Return(int64(1), nil)
				return instanceRepo
			},
			inputInstanceID:         "instance-1",
			expectedInstanceName:    "My instance 1",
			expectedInstanceDomains: []string{"d1.example.com", "d2.example.com", "d3.example.com"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			d := domain.NewDeleteInstanceCommand(tc.inputInstanceID)
			ctrl := gomock.NewController(t)
			opts := &domain.InvokeOpts{DB: new(noopdb.Pool)}

			if tc.mockTx != nil {
				opts.DB = tc.mockTx(ctrl)
			}
			if tc.instanceRepo != nil {
				domain.WithInstanceRepo(tc.instanceRepo(ctrl))(opts)
			}

			// Test
			err := d.Execute(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedInstanceName, d.InstanceName)
			assert.ElementsMatch(t, tc.expectedInstanceDomains, d.InstanceDomains)
		})
	}
}
