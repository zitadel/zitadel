package domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
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

func TestDeleteSessionCommand_Validate(t *testing.T) {
	t.Parallel()
	ctx := authz.NewMockContext("inst-1", "org-1", gofakeit.UUID())
	getErr := errors.New("get error")
	validationErr := errors.New("validation error")

	type args struct {
		sessionID    string
		sessionToken *string
	}
	tt := []struct {
		testName             string
		sessionRepo          func(ctrl *gomock.Controller) domain.SessionRepository
		args                 args
		inputSessionVerifier func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error)
		expectedError        error
	}{
		{
			testName: "validate delete session, without token verifier, success",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				sessionRepositoryPrimaryKey(repo, "inst-1", "session-1")

				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(
						repo.PrimaryKeyCondition("inst-1", "session-1"),
					))).
					Times(1).
					Return(&domain.Session{
						InstanceID: "inst-1",
						ID:         "session-1",
						TokenID:    "token-1",
						Lifetime:   0,
						Expiration: time.Time{},
						UserID:     "",
						CreatedAt:  time.Time{},
						UpdatedAt:  time.Time{},
						Factors:    nil,
						Challenges: nil,
						Metadata:   nil,
						UserAgent:  nil,
					}, nil)
				return repo
			},
			args: args{
				sessionID:    "session-1",
				sessionToken: nil,
			},
			expectedError: nil,
		},
		{
			testName: "validate delete session, without session, error",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				sessionRepositoryPrimaryKey(repo, "inst-1", "session-1")

				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(
						repo.PrimaryKeyCondition("inst-1", "session-1"),
					))).
					Times(1).
					Return(nil, database.NewNoRowFoundError(getErr))
				return repo
			},
			args: args{
				sessionID:    "session-1",
				sessionToken: nil,
			},
			inputSessionVerifier: mockSessionVerification(nil),
			expectedError:        zerrors.ThrowNotFound(database.NewNoRowFoundError(getErr), "DOM-8KYOH3", "Errors.Session.NotFound"),
		},
		{
			testName: "validate delete session, without token, success",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				sessionRepositoryPrimaryKey(repo, "inst-1", "session-1")

				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(
						repo.PrimaryKeyCondition("inst-1", "session-1"),
					))).
					Times(1).
					Return(&domain.Session{
						InstanceID: "inst-1",
						ID:         "session-1",
						TokenID:    "token-1",
						Lifetime:   0,
						Expiration: time.Time{},
						UserID:     "",
						CreatedAt:  time.Time{},
						UpdatedAt:  time.Time{},
						Factors:    nil,
						Challenges: nil,
						Metadata:   nil,
						UserAgent:  nil,
					}, nil)
				return repo
			},
			args: args{
				sessionID:    "session-1",
				sessionToken: nil,
			},
			inputSessionVerifier: mockSessionVerification(nil),
			expectedError:        nil,
		},
		{
			testName: "validate delete session, with token, success",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				sessionRepositoryPrimaryKey(repo, "inst-1", "session-1")

				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(
						repo.PrimaryKeyCondition("inst-1", "session-1"),
					))).
					Times(1).
					Return(&domain.Session{
						InstanceID: "inst-1",
						ID:         "session-1",
						TokenID:    "token-1",
						Lifetime:   0,
						Expiration: time.Time{},
						UserID:     "",
						CreatedAt:  time.Time{},
						UpdatedAt:  time.Time{},
						Factors:    nil,
						Challenges: nil,
						Metadata:   nil,
						UserAgent:  nil,
					}, nil)
				return repo
			},
			args: args{
				sessionID:    "session-1",
				sessionToken: gu.Ptr("token"),
			},
			inputSessionVerifier: mockSessionVerification(nil),
			expectedError:        nil,
		},
		{
			testName: "validate delete session, with token, failure",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				sessionRepositoryPrimaryKey(repo, "inst-1", "session-1")

				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(
						repo.PrimaryKeyCondition("inst-1", "session-1"),
					))).
					Times(1).
					Return(&domain.Session{
						InstanceID: "inst-1",
						ID:         "session-1",
						TokenID:    "token-1",
						Lifetime:   0,
						Expiration: time.Time{},
						UserID:     "",
						CreatedAt:  time.Time{},
						UpdatedAt:  time.Time{},
						Factors:    nil,
						Challenges: nil,
						Metadata:   nil,
						UserAgent:  nil,
					}, nil)
				return repo
			},
			args: args{
				sessionID:    "session-1",
				sessionToken: gu.Ptr("token-1"),
			},
			inputSessionVerifier: mockSessionVerification(validationErr),
			expectedError:        validationErr,
		},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			d := domain.NewDeleteSessionCommand(tc.args.sessionID, tc.args.sessionToken)
			ctrl := gomock.NewController(t)
			opts := &domain.InvokeOpts{}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)

			if tc.sessionRepo != nil {
				domain.WithSessionRepo(tc.sessionRepo(ctrl))(opts)
			}
			if tc.inputSessionVerifier != nil {
				domain.WithSessionTokenVerifier(tc.inputSessionVerifier)(opts)
			}

			// Test
			err := d.Validate(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func mockSessionVerification(err error) func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
	return func(ctx context.Context, sessionToken, sessionID, tokenID string) error {
		return err
	}
}

func TestDeleteSessionCommand_Execute(t *testing.T) {
	t.Parallel()

	ctx := authz.NewMockContext("inst-1", "org-1", gofakeit.UUID())
	deleteErr := errors.New("delete error")

	tt := []struct {
		testName    string
		mockTx      func(ctrl *gomock.Controller) database.QueryExecutor
		sessionRepo func(ctrl *gomock.Controller) domain.SessionRepository

		inputSessionID string

		expectedError error
	}{
		{
			testName: "when delete session fails should return error",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				sessionRepositoryPrimaryKey(repo, "inst-1", "session-1")

				repo.EXPECT().
					Delete(gomock.Any(), gomock.Any(),
						repo.PrimaryKeyCondition("inst-1", "session-1"),
					).
					Times(1).
					Return(int64(0), deleteErr)
				return repo
			},
			inputSessionID: "session-1",
			expectedError:  deleteErr,
		},
		{
			testName: "when more than one row deleted should return internal error",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				sessionRepositoryPrimaryKey(repo, "inst-1", "session-1")

				repo.EXPECT().
					Delete(gomock.Any(), gomock.Any(),
						repo.PrimaryKeyCondition("inst-1", "session-1"),
					).
					Times(1).
					Return(int64(2), nil)
				return repo
			},
			inputSessionID: "session-1",
			expectedError:  zerrors.ThrowInternalf(nil, "DOM-wv33rsKpRw", "expecting 1 row deleted, got %d", 2),
		},
		{
			testName: "when no rows deleted should return not found error",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				sessionRepositoryPrimaryKey(repo, "inst-1", "session-1")

				repo.EXPECT().
					Delete(gomock.Any(), gomock.Any(),
						repo.PrimaryKeyCondition("inst-1", "session-1"),
					).
					Times(1).
					Return(int64(0), nil)
				return repo
			},
			inputSessionID: "session-1",
			expectedError:  zerrors.ThrowNotFound(nil, "DOM-g1lDb1qs1f", "session not found"),
		},
		{
			testName: "when one row deleted should execute successfully",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				// for the expected call
				sessionRepositoryPrimaryKey(repo, "inst-1", "session-1")

				repo.EXPECT().
					Delete(gomock.Any(), gomock.Any(),
						repo.PrimaryKeyCondition("inst-1", "session-1"),
					).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			inputSessionID: "session-1",
		},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			cmd := domain.NewDeleteSessionCommand(tc.inputSessionID, nil)
			ctrl := gomock.NewController(t)
			opts := &domain.InvokeOpts{}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.mockTx != nil {
				domain.WithQueryExecutor(tc.mockTx(ctrl))(opts)
			}
			if tc.sessionRepo != nil {
				domain.WithSessionRepo(tc.sessionRepo(ctrl))(opts)
			}

			// Test
			err := opts.Invoke(ctx, cmd)

			// Verify
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func sessionRepositoryPrimaryKey(repo *domainmock.MockSessionRepository, instanceID, sessionID string) {
	repo.EXPECT().PrimaryKeyCondition(instanceID, sessionID).Times(2).Return(
		database.And(
			database.NewTextCondition(
				database.NewColumn("zitadel.sessions", "instance_id"),
				database.TextOperationEqual, instanceID,
			),
			database.NewTextCondition(
				database.NewColumn("zitadel.sessions", "id"),
				database.TextOperationEqual, sessionID,
			),
		),
	)
}

func TestDeleteSessionCommand_Events(t *testing.T) {
	t.Parallel()
	ctx := authz.NewMockContext("inst-1", "org-1", gofakeit.UUID())

	tt := []struct {
		testName      string
		mockTx        func(ctrl *gomock.Controller) database.QueryExecutor
		command       *domain.DeleteSessionCommand
		expectedError error
		expectedCount int
	}{
		{
			testName: "should create session removed event",
			command: &domain.DeleteSessionCommand{
				ID: "session-1",
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
