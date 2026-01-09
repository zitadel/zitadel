package domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
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
	permissionErr := errors.New("permission error")

	type args struct {
		sessionID      string
		sessionToken   string
		mustCheckToken bool
	}
	type fields struct {
		sessionRepo          func(ctrl *gomock.Controller) domain.SessionRepository
		sessionTokenVerifier domain.SessionTokenVerifier
		permissionCheck      func(ctrl *gomock.Controller) domain.PermissionChecker
	}
	tt := []struct {
		name          string
		args          args
		fields        fields
		expectedError error
	}{
		{
			name: "when sessionID is empty should return invalid argument error",
			args: args{
				sessionID: "",
			},
			fields:        fields{},
			expectedError: zerrors.ThrowInvalidArgument(nil, "SESS-3n9fs", "Errors.IDMissing"),
		},
		{
			name: "when user without permission deletes own session without token should validate successfully",
			args: args{
				sessionID:      "session-1",
				sessionToken:   "",
				mustCheckToken: true,
			},
			fields: fields{
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
				sessionTokenVerifier: mockSessionVerification(nil),
			},
			expectedError: nil,
		},
		{
			name: "when user without permission deletes own session, which does not exist anymore, should validate successfully",
			args: args{
				sessionID:      "session-1",
				sessionToken:   "",
				mustCheckToken: true,
			},
			fields: fields{
				sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
					repo := domainmock.NewMockSessionRepository(ctrl)
					sessionRepositoryPrimaryKey(repo, "inst-1", "session-1")

					repo.EXPECT().
						Get(gomock.Any(), gomock.Any(),
							dbmock.QueryOptions(database.WithCondition(
								repo.PrimaryKeyCondition("inst-1", "session-1"),
							)),
						).
						Times(1).
						Return(nil, database.NewNoRowFoundError(nil))
					return repo
				},
				sessionTokenVerifier: mockSessionVerification(nil),
				permissionCheck: func(ctrl *gomock.Controller) domain.PermissionChecker {
					permChecker := domainmock.NewMockPermissionChecker(ctrl)

					permChecker.EXPECT().
						CheckOrganizationPermission(gomock.Any(), domain.SessionDeletePermission, "").
						Times(1).
						Return(permissionErr)

					return permChecker
				},
			},
			expectedError: nil,
		},
		{
			name: "when user without permission deletes other / non existing session, should result in permission error",
			args: args{
				sessionID:      "session-1",
				mustCheckToken: true,
			},
			fields: fields{
				sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
					repo := domainmock.NewMockSessionRepository(ctrl)
					sessionRepositoryPrimaryKey(repo, "inst-1", "session-1")

					repo.EXPECT().
						Get(gomock.Any(), gomock.Any(),
							dbmock.QueryOptions(database.WithCondition(
								repo.PrimaryKeyCondition("inst-1", "session-1"),
							)),
						).
						Times(1).
						Return(nil, database.NewNoRowFoundError(nil))
					return repo
				},
				sessionTokenVerifier: mockSessionVerification(nil),
				permissionCheck: func(ctrl *gomock.Controller) domain.PermissionChecker {
					permChecker := domainmock.NewMockPermissionChecker(ctrl)

					permChecker.EXPECT().
						CheckOrganizationPermission(gomock.Any(), domain.SessionDeletePermission, "").
						Times(1).
						Return(permissionErr)

					return permChecker
				},
			},
			expectedError: permissionErr,
		},
		{
			name: "when user with permission deletes session, should validate successfully",
			args: args{
				sessionID:      "session-1",
				mustCheckToken: true,
			},
			fields: fields{
				sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
					repo := domainmock.NewMockSessionRepository(ctrl)
					sessionRepositoryPrimaryKey(repo, "inst-1", "session-1")

					repo.EXPECT().
						Get(gomock.Any(), gomock.Any(),
							dbmock.QueryOptions(database.WithCondition(
								repo.PrimaryKeyCondition("inst-1", "session-1"),
							)),
						).
						Times(1).
						Return(&domain.Session{
							InstanceID: "inst-1",
							ID:         "session-1",
							TokenID:    "token-1",
							Lifetime:   0,
							Expiration: time.Time{},
							UserID:     "user-2",
							CreatedAt:  time.Time{},
							UpdatedAt:  time.Time{},
							Factors:    nil,
							Challenges: nil,
							Metadata:   nil,
							UserAgent:  nil,
						}, nil)
					return repo
				},
				sessionTokenVerifier: mockSessionVerification(nil),
				permissionCheck: func(ctrl *gomock.Controller) domain.PermissionChecker {
					permChecker := domainmock.NewMockPermissionChecker(ctrl)

					permChecker.EXPECT().
						CheckOrganizationPermission(gomock.Any(), domain.SessionDeletePermission, "").
						Times(1).
						Return(nil)

					return permChecker
				},
			},
			expectedError: nil,
		},
		{
			name: "when user with permission deletes session, which does not exist anymore, should validate successfully",
			args: args{
				sessionID:      "session-1",
				mustCheckToken: true,
			},
			fields: fields{
				sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
					repo := domainmock.NewMockSessionRepository(ctrl)
					sessionRepositoryPrimaryKey(repo, "inst-1", "session-1")

					repo.EXPECT().
						Get(gomock.Any(), gomock.Any(),
							dbmock.QueryOptions(database.WithCondition(
								repo.PrimaryKeyCondition("inst-1", "session-1"),
							)),
						).
						Times(1).
						Return(nil, database.NewNoRowFoundError(nil))
					return repo
				},
				sessionTokenVerifier: mockSessionVerification(nil),
				permissionCheck: func(ctrl *gomock.Controller) domain.PermissionChecker {
					permChecker := domainmock.NewMockPermissionChecker(ctrl)

					permChecker.EXPECT().
						CheckOrganizationPermission(gomock.Any(), domain.SessionDeletePermission, "").
						Times(1).
						Return(nil)

					return permChecker
				},
			},
			expectedError: nil,
		},
		{
			name: "when user with permission deletes other session (different org), should result in permission error",
			args: args{
				sessionID:      "session-1",
				mustCheckToken: true,
			},
			fields: fields{
				sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
					repo := domainmock.NewMockSessionRepository(ctrl)
					sessionRepositoryPrimaryKey(repo, "inst-1", "session-1")

					repo.EXPECT().
						Get(gomock.Any(), gomock.Any(),
							dbmock.QueryOptions(database.WithCondition(
								repo.PrimaryKeyCondition("inst-1", "session-1"),
							)),
						).
						Times(1).
						Return(&domain.Session{
							InstanceID: "inst-1",
							ID:         "session-1",
							TokenID:    "token-1",
							Lifetime:   0,
							Expiration: time.Time{},
							UserID:     "user-2",
							CreatedAt:  time.Time{},
							UpdatedAt:  time.Time{},
							Factors:    nil,
							Challenges: nil,
							Metadata:   nil,
							UserAgent:  nil,
						}, nil)
					return repo
				},
				sessionTokenVerifier: mockSessionVerification(nil),
				permissionCheck: func(ctrl *gomock.Controller) domain.PermissionChecker {
					permChecker := domainmock.NewMockPermissionChecker(ctrl)

					permChecker.EXPECT().
						CheckOrganizationPermission(gomock.Any(), domain.SessionDeletePermission, "").
						Times(1).
						Return(permissionErr)

					return permChecker
				},
			},
			expectedError: permissionErr,
		},
		{
			name: "when user with permission deletes non existing session, should result in permission error",
			args: args{
				sessionID:      "session-1",
				mustCheckToken: true,
			},
			fields: fields{
				sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
					repo := domainmock.NewMockSessionRepository(ctrl)
					sessionRepositoryPrimaryKey(repo, "inst-1", "session-1")

					repo.EXPECT().
						Get(gomock.Any(), gomock.Any(),
							dbmock.QueryOptions(database.WithCondition(
								repo.PrimaryKeyCondition("inst-1", "session-1"),
							)),
						).
						Times(1).
						Return(nil, database.NewNoRowFoundError(nil))
					return repo
				},
				sessionTokenVerifier: mockSessionVerification(nil),
				permissionCheck: func(ctrl *gomock.Controller) domain.PermissionChecker {
					permChecker := domainmock.NewMockPermissionChecker(ctrl)

					permChecker.EXPECT().
						CheckOrganizationPermission(gomock.Any(), domain.SessionDeletePermission, "").
						Times(1).
						Return(permissionErr)

					return permChecker
				},
			},
			expectedError: permissionErr,
		},
		{
			name: "when deleting session with valid token, should validate successfully",
			args: args{
				sessionID:      "session-1",
				sessionToken:   "valid-token",
				mustCheckToken: true,
			},
			fields: fields{
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
				sessionTokenVerifier: mockSessionVerification(nil),
			},
			expectedError: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Given
			d := domain.NewDeleteSessionCommand(tc.args.sessionID, tc.args.sessionToken, tc.args.mustCheckToken)
			ctrl := gomock.NewController(t)
			opts := domain.DefaultOpts(nil)
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)

			if tc.fields.sessionRepo != nil {
				domain.WithSessionRepo(tc.fields.sessionRepo(ctrl))(opts)
			}
			if tc.fields.sessionTokenVerifier != nil {
				domain.WithSessionTokenVerifier(tc.fields.sessionTokenVerifier)(opts)
			}
			if tc.fields.permissionCheck != nil {
				domain.WithPermissionCheck(tc.fields.permissionCheck(ctrl))(opts)
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

		expectedError   error
		expectDeletedAt bool
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
			testName: "when no rows deleted should execute successfully",
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
			expectedError:  nil,
		},
		{
			testName: "when one row deleted should execute successfully and set DeletedAt",
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
			inputSessionID:  "session-1",
			expectDeletedAt: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			cmd := domain.NewDeleteSessionCommand(tc.inputSessionID, "", false)
			ctrl := gomock.NewController(t)
			opts := domain.DefaultOpts(nil)
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.mockTx != nil {
				domain.WithQueryExecutor(tc.mockTx(ctrl))(opts)
			}
			if tc.sessionRepo != nil {
				domain.WithSessionRepo(tc.sessionRepo(ctrl))(opts)
			}

			// Test
			err := cmd.Execute(ctx, opts)

			// Verify
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectDeletedAt, !cmd.DeletedAt.IsZero())
		})
	}
}

func sessionRepositoryPrimaryKey(repo *domainmock.MockSessionRepository, instanceID, sessionID string) {
	repo.EXPECT().PrimaryKeyCondition(instanceID, sessionID).MaxTimes(4).Return(
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
