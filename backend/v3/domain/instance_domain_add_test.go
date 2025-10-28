package domain_test

import (
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
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestAddInstanceDomainCommand_Validate(t *testing.T) {
	t.Parallel()
	getErr := errors.New("get error")
	permissionErr := errors.New("permission error")

	tt := []struct {
		testName           string
		instanceDomainRepo func(ctrl *gomock.Controller) domain.InstanceDomainRepository
		permissionChecker  func(ctrl *gomock.Controller) domain.PermissionChecker
		inputInstanceID    string
		inputDomainName    string
		expectedError      error
	}{
		{
			testName:        "when no name shuld return invalid argument error",
			inputInstanceID: " ",
			inputDomainName: " ",
			expectedError:   zerrors.ThrowInvalidArgument(nil, "DOM-jieuM8", "Errors.Invalid.Argument"),
		},
		{
			testName:        "when no ID should return invalid argument error",
			inputInstanceID: " ",
			inputDomainName: " domain.name",
			expectedError:   zerrors.ThrowInvalidArgument(nil, "DOM-YaUBp5", "Errors.Invalid.Argument"),
		},
		{
			testName:        "when domain name contains invalid characters should return invalid argument error",
			inputInstanceID: "instance-1 ",
			inputDomainName: "?",
			expectedError:   zerrors.ThrowInvalidArgument(nil, "DOM-98VcSQ", "Errors.Instance.Domain.InvalidCharacter"),
		},
		{
			testName:        "when input instance ID doesn't match instance in context should return invalid argument error",
			inputInstanceID: " instance-2",
			inputDomainName: "valid.domain",
			expectedError:   zerrors.ThrowInvalidArgument(nil, "DOM-x01cai", "Errors.Invalid.Argument"),
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
			inputDomainName: "valid.domain",
			expectedError:   zerrors.ThrowPermissionDenied(permissionErr, "DOM-c83vPX", "permission denied"),
		},
		{
			testName: "when retrieving instance domain fails with generic error should return error",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)

				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.DomainWritePermission).
					Times(1).
					Return(nil)

				return permChecker
			},
			instanceDomainRepo: func(ctrl *gomock.Controller) domain.InstanceDomainRepository {
				repo := domainmock.NewInstancesDomainRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							repo.DomainCondition(database.TextOperationEqual, "valid.domain"),
						),
					)).
					Times(1).
					Return(nil, getErr)
				return repo
			},
			inputInstanceID: "instance-1",
			inputDomainName: "valid.domain",
			expectedError:   getErr,
		},
		{
			testName: "when retrieving instance domain succeeds should return already exists error",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)

				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.DomainWritePermission).
					Times(1).
					Return(nil)

				return permChecker
			},
			instanceDomainRepo: func(ctrl *gomock.Controller) domain.InstanceDomainRepository {
				repo := domainmock.NewInstancesDomainRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							repo.DomainCondition(database.TextOperationEqual, "valid.domain"),
						),
					)).
					Times(1).
					Return(&domain.InstanceDomain{}, nil)
				return repo
			},
			inputInstanceID: "instance-1",
			inputDomainName: "valid.domain",
			expectedError:   zerrors.ThrowAlreadyExists(nil, "DOM-CvQ8tf", "Errors.Instance.Domain.AlreadyExists"),
		},
		{
			testName: "when retrieving instance domain fails with not found error validation should succeed",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)

				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.DomainWritePermission).
					Times(1).
					Return(nil)

				return permChecker
			},
			instanceDomainRepo: func(ctrl *gomock.Controller) domain.InstanceDomainRepository {
				repo := domainmock.NewInstancesDomainRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							repo.DomainCondition(database.TextOperationEqual, "valid.domain"),
						),
					)).
					Times(1).
					Return(nil, &database.NoRowFoundError{})
				return repo
			},
			inputInstanceID: "instance-1",
			inputDomainName: "valid.domain",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "org-1", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewAddInstanceDomainCommand(tc.inputInstanceID, tc.inputDomainName)

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.instanceDomainRepo != nil {
				domain.WithInstanceDomainRepo(tc.instanceDomainRepo(ctrl))(opts)
			}
			if tc.permissionChecker != nil {
				opts.Permissions = tc.permissionChecker(ctrl)
			}

			err := cmd.Validate(ctx, opts)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestAddInstanceDomainCommand_Events(t *testing.T) {
	t.Parallel()
	// Given
	instanceID := "instance-1"
	domainName := "test.domain"
	ctx := authz.NewMockContext(instanceID, "org-1", "")
	cmd := domain.NewAddInstanceDomainCommand(instanceID, domainName)

	// When
	events, err := cmd.Events(ctx, &domain.InvokeOpts{})

	// Then
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.NotNil(t, events[0].Aggregate())
	assert.Equal(t, instanceID, events[0].Aggregate().ID)
	assert.Equal(t, instanceID, events[0].Aggregate().InstanceID)
	assert.Equal(t, instanceID, events[0].Aggregate().ResourceOwner)
	assert.EqualValues(t, instance.AggregateType, events[0].Aggregate().Type)
	assert.EqualValues(t, instance.AggregateVersion, events[0].Aggregate().Version)
}

func TestAddInstanceDomainCommand_Execute(t *testing.T) {
	t.Parallel()
	addErr := errors.New("add error")

	tt := []struct {
		testName           string
		instanceDomainRepo func(ctrl *gomock.Controller) domain.InstanceDomainRepository
		expectedError      error
	}{
		{
			testName: "when adding domain succeeds should return nil",
			instanceDomainRepo: func(ctrl *gomock.Controller) domain.InstanceDomainRepository {
				repo := domainmock.NewInstancesDomainRepo(ctrl)
				repo.EXPECT().
					Add(gomock.Any(), gomock.Any(), &domain.AddInstanceDomain{
						InstanceID:  "instance-1",
						Domain:      "valid.domain",
						IsPrimary:   gu.Ptr(false),
						IsGenerated: gu.Ptr(false),
						Type:        domain.DomainTypeCustom,
					}).
					Times(1).
					Return(nil)
				return repo
			},
		},
		{
			testName: "when adding domain fails should return error",
			instanceDomainRepo: func(ctrl *gomock.Controller) domain.InstanceDomainRepository {
				repo := domainmock.NewInstancesDomainRepo(ctrl)
				repo.EXPECT().
					Add(gomock.Any(), gomock.Any(), &domain.AddInstanceDomain{
						InstanceID:  "instance-1",
						Domain:      "valid.domain",
						IsPrimary:   gu.Ptr(false),
						IsGenerated: gu.Ptr(false),
						Type:        domain.DomainTypeCustom,
					}).
					Times(1).
					Return(addErr)
				return repo
			},
			expectedError: addErr,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "org-1", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewAddInstanceDomainCommand("instance-1", "valid.domain")

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.instanceDomainRepo != nil {
				domain.WithInstanceDomainRepo(tc.instanceDomainRepo(ctrl))(opts)
			}

			// When
			err := cmd.Execute(ctx, opts)

			// Then
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
