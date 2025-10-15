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

func TestUpdateInstanceCommand_Validate(t *testing.T) {
	t.Parallel()
	txInitErr := errors.New("tx init error")
	getErr := errors.New("get error")

	tt := []struct {
		testName          string
		queryExecutor     func(ctrl *gomock.Controller) database.QueryExecutor
		instanceRepo      func(ctrl *gomock.Controller) domain.InstanceRepository
		inputInstanceID   string
		inputInstanceName string
		expectedError     error
	}{
		{
			testName:          "when no ID should return invalid argument error",
			inputInstanceID:   "",
			inputInstanceName: "test-name",
			expectedError:     zerrors.ThrowInvalidArgument(nil, "DOM-wSs6kG", "invalid instance ID"),
		},
		{
			testName:          "when no name shuld return invalid argument error",
			inputInstanceID:   "test-id",
			inputInstanceName: "",
			expectedError:     zerrors.ThrowInvalidArgument(nil, "DOM-FPJcLC", "invalid instance name"),
		},
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
			inputInstanceID:   "test-id",
			inputInstanceName: "test-name",
			expectedError:     txInitErr,
		},
		{
			testName: "when retrieving instance fails should return error",
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
			expectedError:     getErr,
		},
		{
			testName: "when instance name is not changed should return name not changed error",
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
			inputInstanceID:   "instance-1",
			inputInstanceName: "test instance update",
			expectedError:     zerrors.ThrowPreconditionFailed(nil, "DOM-5MrT21", "Errors.Instance.NotChanged"),
		},
		{
			testName: "when instance name is changed should validate successfully and return no error",
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

			opts := &domain.CommandOpts{
				DB: new(noopdb.Pool),
			}
			if tc.instanceRepo != nil {
				opts.SetInstanceRepo(tc.instanceRepo(ctrl))
			}
			if tc.queryExecutor != nil {
				opts.DB = tc.queryExecutor(ctrl)
			}

			err := cmd.Validate(ctx, opts)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateInstanceCommand_Execute(t *testing.T) {
	t.Parallel()

	txInitErr := errors.New("tx init error")
	updateErr := errors.New("update error")

	tt := []struct {
		testName string

		queryExecutor func(ctrl *gomock.Controller) database.QueryExecutor
		instanceRepo  func(ctrl *gomock.Controller) domain.InstanceRepository

		inputID   string
		inputName string

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
			expectedError: updateErr,
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
			expectedError: zerrors.ThrowInternal(domain.NewMultipleObjectsUpdatedError(1, 2), "DOM-HlrNmD", "unexpected number of rows updated"),
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

			opts := &domain.CommandOpts{
				DB: new(noopdb.Pool),
			}
			if tc.instanceRepo != nil {
				opts.SetInstanceRepo(tc.instanceRepo(ctrl))
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
func TestUpdateInstanceCommand_Events(t *testing.T) {
	t.Parallel()

	// Given
	ctx := authz.NewMockContext("instance-1", "", "")
	cmd := domain.NewUpdateInstanceCommand("instance-1", "test-name")

	// Test
	events, err := cmd.Events(ctx, &domain.CommandOpts{})

	// Verify
	assert.NoError(t, err)
	require.Len(t, events, 1)

	event := events[0].(*instance.InstanceChangedEvent)
	assert.Equal(t, "instance-1", event.Aggregate().ID)
	assert.Equal(t, "test-name", event.Name)
}
