package domain_test

import (
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
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestRemoveInstanceDomainCommand_Validate(t *testing.T) {
	t.Parallel()

	getErr := errors.New("get error")
	noRowFoundErr := &database.NoRowFoundError{}
	permissionErr := errors.New("permission error")

	tt := []struct {
		testName          string
		domainRepo        func(ctrl *gomock.Controller) domain.InstanceDomainRepository
		permissionChecker func(ctrl *gomock.Controller) domain.PermissionChecker
		inputInstanceID   string
		inputDomainName   string
		expectedError     error
	}{
		{
			testName:        "when no instance ID should return invalid argument error",
			inputInstanceID: "",
			inputDomainName: "test-domain.com",
			expectedError:   zerrors.ThrowInvalidArgument(nil, "DOM-VSsTTf", "Errors.Invalid.Argument"),
		},
		{
			testName:        "when no domain name should return invalid argument error",
			inputInstanceID: "instance-1",
			inputDomainName: "",
			expectedError:   zerrors.ThrowInvalidArgument(nil, "DOM-PLpYix", "Errors.Invalid.Argument"),
		},
		{
			testName:        "when instance ID does not match context should return invalid argument error",
			inputInstanceID: "different-instance",
			inputDomainName: "test-domain.com",
			expectedError:   zerrors.ThrowInvalidArgument(nil, "DOM-83FUdY", "Errors.Invalid.Argument"),
		},
		{
			testName: "when user is missing permission should return permission denied",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)
				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.DomainWritePermission).
					Times(1).
					Return(permissionErr)
				return permChecker
			},
			inputInstanceID: "instance-1",
			inputDomainName: "test-domain.com",
			expectedError:   zerrors.ThrowPermissionDenied(permissionErr, "DOM-eroxID", "permission denied"),
		},
		{
			testName: "when domain not found should return not found error",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)
				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.DomainWritePermission).
					Times(1).
					Return(nil)
				return permChecker
			},
			domainRepo: func(ctrl *gomock.Controller) domain.InstanceDomainRepository {
				repo := domainmock.NewMockInstanceDomainRepository(ctrl)
				domainCond := database.NewTextCondition(
					database.NewColumn("instance_domains", "domain"),
					database.TextOperationEqual,
					"test-domain.com")
				repo.EXPECT().
					DomainCondition(database.TextOperationEqual, "test-domain.com").
					Times(1).
					Return(domainCond)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(domainCond))).
					Times(1).
					Return(nil, noRowFoundErr)
				return repo
			},
			inputInstanceID: "instance-1",
			inputDomainName: "test-domain.com",
			expectedError:   zerrors.ThrowNotFound(noRowFoundErr, "DOM-nryNFt", "Errors.Instance.Domain.NotFound"),
		},
		{
			testName: "when get domain fails should return error",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)
				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.DomainWritePermission).
					Times(1).
					Return(nil)
				return permChecker
			},
			domainRepo: func(ctrl *gomock.Controller) domain.InstanceDomainRepository {
				repo := domainmock.NewMockInstanceDomainRepository(ctrl)
				domainCond := database.NewTextCondition(
					database.NewColumn("instance_domains", "domain"),
					database.TextOperationEqual,
					"test-domain.com")
				repo.EXPECT().
					DomainCondition(database.TextOperationEqual, "test-domain.com").
					Times(1).
					Return(domainCond)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(domainCond))).
					Times(1).
					Return(nil, getErr)
				return repo
			},
			inputInstanceID: "instance-1",
			inputDomainName: "test-domain.com",
			expectedError:   getErr,
		},
		{
			testName: "when domain is generated should return precondition failed error",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)
				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.DomainWritePermission).
					Times(1).
					Return(nil)
				return permChecker
			},
			domainRepo: func(ctrl *gomock.Controller) domain.InstanceDomainRepository {
				repo := domainmock.NewMockInstanceDomainRepository(ctrl)
				isGenerated := true
				domainCond := database.NewTextCondition(
					database.NewColumn("instance_domains", "domain"),
					database.TextOperationEqual,
					"test-domain.com")
				repo.EXPECT().
					DomainCondition(database.TextOperationEqual, "test-domain.com").
					Times(1).
					Return(domainCond)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(domainCond))).
					Times(1).
					Return(&domain.InstanceDomain{IsGenerated: &isGenerated}, nil)
				return repo
			},
			inputInstanceID: "instance-1",
			inputDomainName: "test-domain.com",
			expectedError:   zerrors.ThrowPreconditionFailed(nil, "DOM-cSfCVG", "Errors.Instance.Domain.GeneratedNotRemovable"),
		},
		{
			testName: "when domain exists and is not generated should validate successfully",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)
				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.DomainWritePermission).
					Times(1).
					Return(nil)
				return permChecker
			},
			domainRepo: func(ctrl *gomock.Controller) domain.InstanceDomainRepository {
				repo := domainmock.NewMockInstanceDomainRepository(ctrl)
				isGenerated := false
				domainCond := database.NewTextCondition(
					database.NewColumn("instance_domains", "domain"),
					database.TextOperationEqual,
					"test-domain.com")
				repo.EXPECT().
					DomainCondition(database.TextOperationEqual, "test-domain.com").
					Times(1).
					Return(domainCond)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(domainCond))).
					Times(1).
					Return(&domain.InstanceDomain{IsGenerated: &isGenerated}, nil)
				return repo
			},
			inputInstanceID: "instance-1",
			inputDomainName: "test-domain.com",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewRemoveInstanceDomainCommand(tc.inputInstanceID, tc.inputDomainName)

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.domainRepo != nil {
				domain.WithInstanceDomainRepo(tc.domainRepo(ctrl))(opts)
			}
			if tc.permissionChecker != nil {
				opts.Permissions = tc.permissionChecker(ctrl)
			}

			err := cmd.Validate(ctx, opts)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
func TestRemoveInstanceDomainCommand_Execute(t *testing.T) {
	t.Parallel()

	removeErr := errors.New("remove error")

	tt := []struct {
		testName        string
		domainRepo      func(ctrl *gomock.Controller) domain.InstanceDomainRepository
		inputInstanceID string
		inputDomainName string
		expectedError   error
	}{
		{
			testName: "when domain remove fails should return error",
			domainRepo: func(ctrl *gomock.Controller) domain.InstanceDomainRepository {
				repo := domainmock.NewMockInstanceDomainRepository(ctrl)
				domainCond := database.NewTextCondition(
					database.NewColumn("instance_domains", "domain"),
					database.TextOperationEqual,
					"test-domain.com")
				repo.EXPECT().
					PrimaryKeyCondition("test-domain.com").
					Times(1).
					Return(domainCond)
				repo.EXPECT().
					Remove(gomock.Any(), gomock.Any(), domainCond).
					Times(1).
					Return(int64(0), removeErr)
				return repo
			},
			inputInstanceID: "instance-1",
			inputDomainName: "test-domain.com",
			expectedError:   removeErr,
		},
		{
			testName: "when domain remove returns 0 rows removed should return not found error",
			domainRepo: func(ctrl *gomock.Controller) domain.InstanceDomainRepository {
				repo := domainmock.NewMockInstanceDomainRepository(ctrl)
				domainCond := database.NewTextCondition(
					database.NewColumn("instance_domains", "domain"),
					database.TextOperationEqual,
					"test-domain.com")
				repo.EXPECT().
					PrimaryKeyCondition("test-domain.com").
					Times(1).
					Return(domainCond)
				repo.EXPECT().
					Remove(gomock.Any(), gomock.Any(), domainCond).
					Times(1).
					Return(int64(0), nil)
				return repo
			},
			inputInstanceID: "instance-1",
			inputDomainName: "test-domain.com",
			expectedError:   zerrors.ThrowNotFound(nil, "DOM-ZUteYg", "instance domain not found"),
		},
		{
			testName: "when domain remove returns more than 1 row removed should return internal error",
			domainRepo: func(ctrl *gomock.Controller) domain.InstanceDomainRepository {
				repo := domainmock.NewMockInstanceDomainRepository(ctrl)
				domainCond := database.NewTextCondition(
					database.NewColumn("instance_domains", "domain"),
					database.TextOperationEqual,
					"test-domain.com")
				repo.EXPECT().
					PrimaryKeyCondition("test-domain.com").
					Times(1).
					Return(domainCond)
				repo.EXPECT().
					Remove(gomock.Any(), gomock.Any(), domainCond).
					Times(1).
					Return(int64(2), nil)
				return repo
			},
			inputInstanceID: "instance-1",
			inputDomainName: "test-domain.com",
			expectedError:   zerrors.ThrowInternalf(nil, "DOM-XSCnJB", "expecting 1 row deleted, got %d", 2),
		},
		{
			testName: "when domain remove returns 1 row removed should return no error",
			domainRepo: func(ctrl *gomock.Controller) domain.InstanceDomainRepository {
				repo := domainmock.NewMockInstanceDomainRepository(ctrl)
				domainCond := database.NewTextCondition(
					database.NewColumn("instance_domains", "domain"),
					database.TextOperationEqual,
					"test-domain.com")
				repo.EXPECT().
					PrimaryKeyCondition("test-domain.com").
					Times(1).
					Return(domainCond)
				repo.EXPECT().
					Remove(gomock.Any(), gomock.Any(), domainCond).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			inputInstanceID: "instance-1",
			inputDomainName: "test-domain.com",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewRemoveInstanceDomainCommand(tc.inputInstanceID, tc.inputDomainName)

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.domainRepo != nil {
				domain.WithInstanceDomainRepo(tc.domainRepo(ctrl))(opts)
			}

			err := cmd.Execute(ctx, opts)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestRemoveInstanceDomainCommand_Events(t *testing.T) {
	t.Parallel()

	// Given
	ctx := authz.NewMockContext("instance-1", "", "")
	cmd := domain.NewRemoveInstanceDomainCommand("instance-1", "test-domain.com")

	// Test
	events, err := cmd.Events(ctx, &domain.InvokeOpts{})

	// Verify
	assert.NoError(t, err)
	require.Len(t, events, 1)

	event := events[0].(*instance.DomainRemovedEvent)
	assert.Equal(t, "instance-1", event.Aggregate().ID)
	assert.Equal(t, "test-domain.com", event.Domain)
}
