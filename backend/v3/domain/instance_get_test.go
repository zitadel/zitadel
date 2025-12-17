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
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestGetInstanceCommand_Validate(t *testing.T) {
	t.Parallel()

	permissionErr := errors.New("permission error")

	tests := []struct {
		name              string
		permissionChecker func(ctrl *gomock.Controller) domain.PermissionChecker
		inputInstanceID   string
		ctx               context.Context
		expectedError     error
	}{
		{
			name:            "empty instance id, error",
			inputInstanceID: "",
			ctx:             context.Background(),
			expectedError:   zerrors.ThrowInvalidArgument(nil, "DOM-32a0o2", "Errors.Instance.ID"),
		},
		{
			name:            "whitespace instance id, error",
			inputInstanceID: "   ",
			ctx:             context.Background(),
			expectedError:   zerrors.ThrowInvalidArgument(nil, "DOM-32a0o2", "Errors.Instance.ID"),
		},
		{
			name:            "instance id mismatch, error",
			inputInstanceID: "instance1",
			ctx:             authz.NewMockContext("instance2", "", ""),
			expectedError:   zerrors.ThrowPermissionDenied(nil, "DOM-n0SvVB", "input instance ID doesn't match context instance"),
		},
		{
			name: "when user is missing permission should return permission denied",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)

				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.InstanceReadPermission).
					Times(1).
					Return(permissionErr)

				return permChecker
			},
			inputInstanceID: "instance1",
			ctx:             authz.NewMockContext("instance1", "", ""),
			expectedError:   zerrors.ThrowPermissionDenied(permissionErr, "DOM-Uq6b00", "Errors.PermissionDenied"),
		},
		{
			name: "valid instance id, success",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)

				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.InstanceReadPermission).
					Times(1).
					Return(nil)

				return permChecker
			},
			inputInstanceID: "instance1",
			ctx:             authz.NewMockContext("instance1", "", ""),
			expectedError:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Given
			cmd := domain.NewGetInstanceCommand(tt.inputInstanceID)

			cmdOpts := &domain.InvokeOpts{}
			if tt.permissionChecker != nil {
				ctrl := gomock.NewController(t)
				cmdOpts.Permissions = tt.permissionChecker(ctrl)
			}

			// Test
			err := cmd.Validate(tt.ctx, cmdOpts)

			// Verify
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestGetInstanceCommand_Execute(t *testing.T) {
	t.Parallel()

	ctx := authz.NewMockContext("inst-1", "org-1", gofakeit.UUID())
	getErr := errors.New("get error")
	notFoundErr := database.NewNoRowFoundError(nil)

	tt := []struct {
		testName string

		instanceRepo func(ctrl *gomock.Controller) domain.InstanceRepository

		inputInstanceID string

		expectedError    error
		expectedInstance *domain.Instance
	}{
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
			expectedError:   zerrors.ThrowInternal(getErr, "DOM-lvsRce", "Errors.Instance.Get"),
		},
		{
			testName: "when instance is not found should return not found error",
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
					Return(nil, notFoundErr)
				return instanceRepo
			},
			inputInstanceID: "instance-1",
			expectedError:   zerrors.ThrowNotFound(notFoundErr, "DOM-QVrUwc", "Errors.Instance.NotFound"),
		},
		{
			testName: "when retrieving instance succeeds should set instance on struct and return no error",
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

				return instanceRepo
			},
			inputInstanceID: "instance-1",
			expectedInstance: &domain.Instance{
				ID:   "instance-1",
				Name: "My instance 1",
				Domains: []*domain.InstanceDomain{
					{Domain: "d1.example.com"},
					{Domain: "d2.example.com"},
					{Domain: "d3.example.com"},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			g := domain.NewGetInstanceCommand(tc.inputInstanceID)
			ctrl := gomock.NewController(t)
			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)

			if tc.instanceRepo != nil {
				domain.WithInstanceRepo(tc.instanceRepo(ctrl))(opts)
			}

			// Test
			err := g.Execute(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedInstance, g.Result())
		})
	}
}
