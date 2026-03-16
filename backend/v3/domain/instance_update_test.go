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
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestUpdateInstanceCommand_Validate(t *testing.T) {
	t.Parallel()
	getErr := errors.New("get error")
	permissionErr := errors.New("permission error")

	tt := []struct {
		testName           string
		instanceRepo       func(ctrl *gomock.Controller) domain.InstanceRepository
		permissionChecker  func(ctrl *gomock.Controller) domain.PermissionChecker
		inputInstanceID    string
		inputInstanceName  string
		expectedError      error
		expectedUpdateSkip bool
	}{
		{
			testName:          "when no ID should return invalid argument error",
			inputInstanceID:   "",
			inputInstanceName: "test-name",
			expectedError:     zerrors.ThrowInvalidArgument(nil, "DOM-wSs6kG", "Errors.Instance.ID"),
		},
		{
			testName:          "when no name shuld return invalid argument error",
			inputInstanceID:   "test-id",
			inputInstanceName: "",
			expectedError:     zerrors.ThrowInvalidArgument(nil, "DOM-FPJcLC", "Errors.Instance.Name"),
		},
		{
			testName: "when user is missing permission should return permission denied",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)

				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.InstanceWritePermission).
					Times(1).
					Return(permissionErr)

				return permChecker
			},
			inputInstanceID:   "instance-1",
			inputInstanceName: "test instance update",
			expectedError:     zerrors.ThrowPermissionDenied(permissionErr, "DOM-M5ObLP", "Errors.PermissionDenied"),
		},
		{
			testName: "when retrieving instance fails should return error",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)

				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.InstanceWritePermission).
					Times(1).
					Return(nil)

				return permChecker
			},
			instanceRepo: func(ctrl *gomock.Controller) domain.InstanceRepository {
				repo := domainmock.NewInstanceRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							repo.IDCondition("instance-1"),
						),
					)).
					Times(1).
					Return(nil, getErr)
				return repo
			},
			inputInstanceID:   "instance-1",
			inputInstanceName: "test instance update",
			expectedError:     zerrors.ThrowInternal(getErr, "DOM-j05Hdo", "Errors.Instance.Get"),
		},
		{
			testName: "when instance name is not changed should return name not changed error",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)

				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.InstanceWritePermission).
					Times(1).
					Return(nil)

				return permChecker
			},
			instanceRepo: func(ctrl *gomock.Controller) domain.InstanceRepository {
				repo := domainmock.NewInstanceRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(
						repo.IDCondition("instance-1"),
					))).
					Times(1).
					Return(&domain.Instance{
						ID:   "instance-1",
						Name: "test instance update",
					}, nil)
				return repo
			},
			inputInstanceID:    "instance-1",
			inputInstanceName:  "test instance update",
			expectedUpdateSkip: true,
		},
		{
			testName: "when instance name is changed should validate successfully and return no error",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)

				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.InstanceWritePermission).
					Times(1).
					Return(nil)

				return permChecker
			},
			instanceRepo: func(ctrl *gomock.Controller) domain.InstanceRepository {
				repo := domainmock.NewInstanceRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(
						repo.IDCondition("instance-1"),
					))).
					Times(1).
					Return(&domain.Instance{
						ID:   "instance-1",
						Name: "old instance name",
					}, nil)
				return repo
			},
			inputInstanceID:   "instance-1",
			inputInstanceName: "test instance update",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewUpdateInstanceCommand(tc.inputInstanceID, tc.inputInstanceName)

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.instanceRepo != nil {
				domain.WithInstanceRepo(tc.instanceRepo(ctrl))(opts)
			}
			if tc.permissionChecker != nil {
				opts.Permissions = tc.permissionChecker(ctrl)
			}

			err := cmd.Validate(ctx, opts)
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedUpdateSkip, cmd.ShouldSkipUpdate)
		})
	}
}

func TestUpdateInstanceCommand_Execute(t *testing.T) {
	t.Parallel()

	updateErr := errors.New("update error")

	tt := []struct {
		testName string

		instanceRepo func(ctrl *gomock.Controller) domain.InstanceRepository

		inputID          string
		inputName        string
		shouldSkipUpdate bool

		expectedError error
	}{
		{
			testName:         "when ShouldSkipUpdate is true should return nil",
			shouldSkipUpdate: true,
		},
		{
			testName: "when instance update fails should return error",
			instanceRepo: func(ctrl *gomock.Controller) domain.InstanceRepository {
				repo := domainmock.NewInstanceRepo(ctrl)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						"instance-1",
						repo.SetName("test instance update"),
					).
					Times(1).
					Return(int64(0), updateErr)
				return repo
			},
			inputID:       "instance-1",
			inputName:     "test instance update",
			expectedError: zerrors.ThrowInternal(updateErr, "DOM-PkVMNR", "Errors.Instance.Update"),
		},
		{
			testName:  "when instance update returns 0 rows updated should return not found error",
			inputID:   "instance-1",
			inputName: "test instance update",
			instanceRepo: func(ctrl *gomock.Controller) domain.InstanceRepository {
				repo := domainmock.NewInstanceRepo(ctrl)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						"instance-1",
						repo.SetName("test instance update"),
					).
					Times(1).
					Return(int64(0), nil)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(nil, "DOM-ghfov1", "Errors.Instance.NotFound"),
		},
		{
			testName: "when instance update returns more than 1 row updated should return internal error",
			instanceRepo: func(ctrl *gomock.Controller) domain.InstanceRepository {
				repo := domainmock.NewInstanceRepo(ctrl)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						"instance-1",
						repo.SetName("test instance update"),
					).
					Times(1).
					Return(int64(2), nil)
				return repo
			},
			inputID:       "instance-1",
			inputName:     "test instance update",
			expectedError: zerrors.ThrowInternal(domain.NewMultipleObjectsUpdatedError(1, 2), "DOM-HlrNmD", "Errors.Instance.UpdateMismatch"),
		},
		{
			testName: "when instance update returns 1 row updated should return no error and set non-primary verified domain",
			instanceRepo: func(ctrl *gomock.Controller) domain.InstanceRepository {
				repo := domainmock.NewInstanceRepo(ctrl)

				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						"instance-1",
						repo.SetName("test instance update"),
					).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			inputID:   "instance-1",
			inputName: "test instance update",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewUpdateInstanceCommand(tc.inputID, tc.inputName)
			cmd.ShouldSkipUpdate = tc.shouldSkipUpdate

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.instanceRepo != nil {
				domain.WithInstanceRepo(tc.instanceRepo(ctrl))(opts)
			}

			// Test
			err := cmd.Execute(ctx, opts)

			// Verify
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
func TestUpdateInstanceCommand_Events(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName         string
		inputName        string
		inputInstanceID  string
		shouldSkipUpdate bool
		expectedEvent    eventstore.Command
	}{
		{
			testName: "when ShouldSkipUpdate is true should return no event and no error",
		},
		{
			testName:        "when ShouldSkipUpdate is true should return expected event",
			inputName:       "test-name",
			inputInstanceID: "instance-1",
			expectedEvent: instance.NewInstanceChangedEvent(
				t.Context(),
				&instance.NewAggregate("instance-1").Aggregate,
				"test-name",
			),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Given
			ctx := authz.NewMockContext(tc.inputInstanceID, "", "")
			cmd := domain.NewUpdateInstanceCommand(tc.inputInstanceID, tc.inputName)

			// Test
			events, err := cmd.Events(ctx, &domain.InvokeOpts{})

			// Verify
			assert.NoError(t, err)
			if tc.expectedEvent != nil {
				require.Len(t, events, 1)

				event := events[0].(*instance.InstanceChangedEvent)
				assert.Equal(t, tc.inputInstanceID, event.Aggregate().ID)
				assert.Equal(t, tc.inputName, event.Name)
			}
		})
	}
}
