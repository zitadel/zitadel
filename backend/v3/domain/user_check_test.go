package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func getUserIDCondition(repo *domainmock.MockUserRepository, userID string) database.Condition {
	idCondition := getTextCondition("zitadel.users", "id", userID)

	repo.EXPECT().
		IDCondition(userID).
		AnyTimes().
		Return(idCondition)
	return idCondition
}

func getSessionIDCondition(repo *domainmock.MockSessionRepository, sessionID string) database.Condition {
	idCondition := getTextCondition("zitadel.sessions", "id", sessionID)

	repo.EXPECT().
		IDCondition(sessionID).
		AnyTimes().
		Return(idCondition)
	return idCondition
}

func getTextCondition(tableName, column, value string) database.Condition {
	return database.NewTextCondition(
		database.NewColumn(tableName, column),
		database.TextOperationEqual,
		value,
	)
}

func TestUserCheckCommand_Validate(t *testing.T) {
	t.Parallel()
	getErr := errors.New("get error")
	notFoundErr := database.NewNoRowFoundError(nil)

	tt := []struct {
		testName      string
		userRepo      func(ctrl *gomock.Controller) domain.UserRepository
		checkUser     *session_grpc.CheckUser
		expectedError error
	}{
		{
			testName:      "when CheckUser is nil should return no error",
			checkUser:     nil,
			expectedError: nil,
		},
		{
			testName: "when search type is UserId should fetch user by ID",
			checkUser: &session_grpc.CheckUser{
				Search: &session_grpc.CheckUser_UserId{
					UserId: "user-123",
				},
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getUserIDCondition(repo, "user-123"),
						),
					)).
					Times(1).
					Return(&domain.User{
						ID:             "user-123",
						OrganizationID: "org-1",
						State:          domain.UserStateActive,
					}, nil)
				return repo
			},
			expectedError: nil,
		},
		{
			testName: "when user not active should return precondition failed error",
			checkUser: &session_grpc.CheckUser{
				Search: &session_grpc.CheckUser_UserId{
					UserId: "user-123",
				},
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getUserIDCondition(repo, "user-123"),
						),
					)).
					Times(1).
					Return(&domain.User{
						ID:             "user-123",
						OrganizationID: "org-1",
						State:          domain.UserStateInactive,
					}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-vgDIu9", "Errors.User.NotActive"),
		},
		{
			testName: "when user retrieval fails should return error",
			checkUser: &session_grpc.CheckUser{
				Search: &session_grpc.CheckUser_LoginName{
					LoginName: "user-123",
				},
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)

				loginNameCondition := database.NewTextCondition(
					database.NewColumn("zitadel.users", "login_name"),
					database.TextOperationEqual,
					"user-123",
				)

				repo.EXPECT().
					LoginNameCondition(database.TextOperationEqual, "user-123").
					AnyTimes().
					Return(loginNameCondition)

				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							loginNameCondition,
						),
					)).
					Times(1).
					Return(nil, getErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(getErr, "DOM-Y846I0", "failed fetching user"),
		},
		{
			testName: "when user not found should return not found error",
			checkUser: &session_grpc.CheckUser{
				Search: &session_grpc.CheckUser_UserId{
					UserId: "user-123",
				},
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)

				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getUserIDCondition(repo, "user-123"),
						),
					)).
					Times(1).
					Return(nil, notFoundErr)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(notFoundErr, "DOM-lcZeXI", "user not found"),
		},
		{
			testName: "when search type is unknown should return invalid argument error",
			checkUser: &session_grpc.CheckUser{
				Search: nil,
			},
			expectedError: zerrors.ThrowInvalidArgumentf(nil, "DOM-7B2m0b", "user search %T not implemented", nil),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := &domain.UserCheckCommand{
				CheckUser:  tc.checkUser,
				SessionID:  "session-1",
				InstanceID: "instance-1",
			}

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.userRepo != nil {
				domain.WithUserRepo(tc.userRepo(ctrl))(opts)
			}

			// Test
			err := cmd.Validate(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUserCheckCommand_Execute(t *testing.T) {
	t.Parallel()

	sessionNotFoundErr := database.NewNoRowFoundError(nil)
	getErr := errors.New("get error")
	updateErr := errors.New("update error")

	tt := []struct {
		testName    string
		sessionRepo func(ctrl *gomock.Controller) domain.SessionRepository
		checkUser   *session_grpc.CheckUser
		sessionID   string
		instanceID  string
		fetchedUser domain.User

		expectedError             error
		expectedPreferredLanguage *language.Tag
	}{
		{
			testName:      "when CheckUser is nil should return no error",
			checkUser:     nil,
			expectedError: nil,
		},
		{
			testName: "when session retrieval fails should return error",
			checkUser: &session_grpc.CheckUser{
				Search: &session_grpc.CheckUser_UserId{
					UserId: "user-123",
				},
			},
			fetchedUser: domain.User{
				ID:             "user-123",
				OrganizationID: "org-1",
			},
			sessionID: "session-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				// TODO(IAM-Marco): uncomment once Session repository is implemented and remove the above
				// repo := domainmock.NewSessionRepo(ctrl)

				// TODO(IAM-Marco): Remove idCondition and expectation once Session repository is implemented
				idCondition := getSessionIDCondition(repo, "session-1")

				repo.EXPECT().
					// TODO(IAM-Marco): Uncomment this down and remove the other Get call once
					// session repository is implemented
					// Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
					// 	database.WithCondition(
					// 		repo.IDCondition("session-1"),
					// 	),
					// )).
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(idCondition),
					)).
					Times(1).
					Return(nil, getErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(getErr, "DOM-To1rLz", "failed fetching session"),
		},
		{
			testName: "when session not found should return not found error",
			checkUser: &session_grpc.CheckUser{
				Search: &session_grpc.CheckUser_UserId{
					UserId: "user-123",
				},
			},
			fetchedUser: domain.User{
				ID:             "user-123",
				OrganizationID: "org-1",
			},
			sessionID: "session-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				// TODO(IAM-Marco): uncomment once Session repository is implemented and remove the above
				// repo := domainmock.NewSessionRepo(ctrl)

				// TODO(IAM-Marco): Remove idCondition and expectation once Session repository is implemented
				idCondition := getSessionIDCondition(repo, "session-1")

				repo.EXPECT().
					// TODO(IAM-Marco): Uncomment this down and remove the other Get call once
					// session repository is implemented
					// Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
					// 	database.WithCondition(
					// 		repo.IDCondition("session-1"),
					// 	),
					// )).
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(idCondition),
					)).
					Times(1).
					Return(nil, sessionNotFoundErr)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(sessionNotFoundErr, "DOM-rbdCv3", "session not found"),
		},
		{
			testName: "when user change is attempted should return invalid argument error",
			checkUser: &session_grpc.CheckUser{
				Search: &session_grpc.CheckUser_UserId{
					UserId: "user-123",
				},
			},
			sessionID:  "session-1",
			instanceID: "instance-1",
			fetchedUser: domain.User{
				ID:             "user-123",
				OrganizationID: "org-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				// TODO(IAM-Marco): uncomment once Session repository is implemented and remove the above
				// repo := domainmock.NewSessionRepo(ctrl)

				// TODO(IAM-Marco): Remove idCondition and expectation once Session repository is implemented
				idCondition := getSessionIDCondition(repo, "session-1")

				repo.EXPECT().
					// TODO(IAM-Marco): Uncomment this down and remove the other Get call once
					// session repository is implemented
					// Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
					// 	database.WithCondition(
					// 		repo.IDCondition("session-1"),
					// 	),
					// )).
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(idCondition),
					)).
					Times(1).
					Return(&domain.Session{
						ID:     "session-1",
						UserID: "user-456",
					}, nil)
				return repo
			},
			expectedError: zerrors.ThrowInvalidArgument(nil, "DOM-78g1TV", "user change not possible"),
		},
		{
			testName: "when session update fails should return error",
			checkUser: &session_grpc.CheckUser{
				Search: &session_grpc.CheckUser_UserId{
					UserId: "user-123",
				},
			},
			sessionID:  "session-1",
			instanceID: "instance-1",
			fetchedUser: domain.User{
				ID:             "user-123",
				OrganizationID: "org-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				// TODO(IAM-Marco): uncomment once Session repository is implemented and remove the above
				// repo := domainmock.NewSessionRepo(ctrl)

				// TODO(IAM-Marco): Remove idCondition and expectation once Session repository is implemented
				idCondition := database.NewTextCondition(
					database.NewColumn("sessions", "id"),
					database.TextOperationEqual,
					"session-1",
				)
				repo.EXPECT().IDCondition("session-1").Times(2).Return(idCondition)

				repo.EXPECT().
					// TODO(IAM-Marco): Uncomment this down and remove the other Get call once
					// session repository is implemented
					// Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
					// 	database.WithCondition(
					// 		repo.IDCondition("session-1"),
					// 	),
					// )).
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(idCondition),
					)).
					Times(1).
					Return(&domain.Session{
						ID:     "session-1",
						UserID: "user-123",
					}, nil)

				userFactor := &domain.SessionFactorUser{UserID: "user-123", LastVerifiedAt: time.Now()}

				// TODO(IAM-Marco): Fix this once Session repo is implemented
				factorChange := database.NewChange(
					database.NewColumn("sessions", "user_factor"),
					"user-123",
				)

				repo.EXPECT().
					SetFactor(userFactor).
					Times(1).
					Return(factorChange)

				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, factorChange).
					Times(1).
					Return(int64(0), updateErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(updateErr, "DOM-netNam", "failed updating session"),
		},
		{
			testName: "when session update returns no rows should return not found error",
			checkUser: &session_grpc.CheckUser{
				Search: &session_grpc.CheckUser_UserId{
					UserId: "user-123",
				},
			},
			sessionID:  "session-1",
			instanceID: "instance-1",
			fetchedUser: domain.User{
				ID:             "user-123",
				OrganizationID: "org-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				// TODO(IAM-Marco): uncomment once Session repository is implemented and remove the above
				// repo := domainmock.NewSessionRepo(ctrl)

				// TODO(IAM-Marco): Remove idCondition and expectation once Session repository is implemented
				idCondition := database.NewTextCondition(
					database.NewColumn("sessions", "id"),
					database.TextOperationEqual,
					"session-1",
				)
				repo.EXPECT().IDCondition("session-1").Times(2).Return(idCondition)

				repo.EXPECT().
					// TODO(IAM-Marco): Uncomment this down and remove the other Get call once
					// session repository is implemented
					// Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
					// 	database.WithCondition(
					// 		repo.IDCondition("session-1"),
					// 	),
					// )).
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(idCondition),
					)).
					Times(1).
					Return(&domain.Session{
						ID:     "session-1",
						UserID: "",
					}, nil)

				userFactor := &domain.SessionFactorUser{UserID: "user-123", LastVerifiedAt: time.Now()}

				// TODO(IAM-Marco): Fix this once Session repo is implemented
				factorChange := database.NewChange(
					database.NewColumn("sessions", "user_factor"),
					"user-123",
				)

				repo.EXPECT().
					SetFactor(userFactor).
					Times(1).
					Return(factorChange)

				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, factorChange).
					Times(1).
					Return(int64(0), nil)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(nil, "DOM-FszyWS", "session not found"),
		},
		{
			testName: "when session update returns multiple rows should return internal error",
			checkUser: &session_grpc.CheckUser{
				Search: &session_grpc.CheckUser_UserId{
					UserId: "user-123",
				},
			},
			sessionID:  "session-1",
			instanceID: "instance-1",
			fetchedUser: domain.User{
				ID:             "user-123",
				OrganizationID: "org-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				// TODO(IAM-Marco): uncomment once Session repository is implemented and remove the above
				// repo := domainmock.NewSessionRepo(ctrl)

				// TODO(IAM-Marco): Remove idCondition and expectation once Session repository is implemented
				idCondition := database.NewTextCondition(
					database.NewColumn("sessions", "id"),
					database.TextOperationEqual,
					"session-1",
				)
				repo.EXPECT().IDCondition("session-1").Times(2).Return(idCondition)

				repo.EXPECT().
					// TODO(IAM-Marco): Uncomment this down and remove the other Get call once
					// session repository is implemented
					// Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
					// 	database.WithCondition(
					// 		repo.IDCondition("session-1"),
					// 	),
					// )).
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(idCondition),
					)).
					Times(1).
					Return(&domain.Session{
						ID:     "session-1",
						UserID: "",
					}, nil)

				userFactor := &domain.SessionFactorUser{UserID: "user-123", LastVerifiedAt: time.Now()}

				// TODO(IAM-Marco): Fix this once Session repo is implemented
				factorChange := database.NewChange(
					database.NewColumn("sessions", "user_factor"),
					"user-123",
				)

				repo.EXPECT().
					SetFactor(userFactor).
					Times(1).
					Return(factorChange)

				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, factorChange).
					Times(1).
					Return(int64(2), nil)
				return repo
			},
			expectedError: zerrors.ThrowInternal(domain.NewMultipleObjectsUpdatedError(1, 2), "DOM-SsIwDt", "unexpected number of rows updated"),
		},
		{
			testName: "when session has no user and user is set should execute successfully",
			checkUser: &session_grpc.CheckUser{
				Search: &session_grpc.CheckUser_UserId{
					UserId: "user-123",
				},
			},
			sessionID:  "session-1",
			instanceID: "instance-1",
			fetchedUser: domain.User{
				ID:             "user-123",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					PreferredLanguage: language.Albanian,
				},
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				// TODO(IAM-Marco): uncomment once Session repository is implemented and remove the above
				// repo := domainmock.NewSessionRepo(ctrl)

				// TODO(IAM-Marco): Remove idCondition and expectation once Session repository is implemented
				idCondition := database.NewTextCondition(
					database.NewColumn("sessions", "id"),
					database.TextOperationEqual,
					"session-1",
				)
				repo.EXPECT().IDCondition("session-1").Times(2).Return(idCondition)

				repo.EXPECT().
					// TODO(IAM-Marco): Uncomment this down and remove the other Get call once
					// session repository is implemented
					// Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
					// 	database.WithCondition(
					// 		repo.IDCondition("session-1"),
					// 	),
					// )).
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(idCondition),
					)).
					Times(1).
					Return(&domain.Session{
						ID:     "session-1",
						UserID: "",
					}, nil)

				userFactor := &domain.SessionFactorUser{UserID: "user-123", LastVerifiedAt: time.Now()}

				// TODO(IAM-Marco): Fix this once Session repo is implemented
				factorChange := database.NewChange(
					database.NewColumn("sessions", "user_factor"),
					"user-123",
				)

				repo.EXPECT().
					SetFactor(userFactor).
					Times(1).
					Return(factorChange)

				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, factorChange).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			expectedPreferredLanguage: &language.Albanian,
		},
		{
			testName: "when session user matches fetched user should execute successfully",
			checkUser: &session_grpc.CheckUser{
				Search: &session_grpc.CheckUser_UserId{
					UserId: "user-123",
				},
			},
			sessionID:  "session-1",
			instanceID: "instance-1",
			fetchedUser: domain.User{
				ID:             "user-123",
				OrganizationID: "org-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				// TODO(IAM-Marco): uncomment once Session repository is implemented and remove the above
				// repo := domainmock.NewSessionRepo(ctrl)

				// TODO(IAM-Marco): Remove idCondition and expectation once Session repository is implemented
				idCondition := database.NewTextCondition(
					database.NewColumn("sessions", "id"),
					database.TextOperationEqual,
					"session-1",
				)
				repo.EXPECT().IDCondition("session-1").Times(2).Return(idCondition)

				repo.EXPECT().
					// TODO(IAM-Marco): Uncomment this down and remove the other Get call once
					// session repository is implemented
					// Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
					// 	database.WithCondition(
					// 		repo.IDCondition("session-1"),
					// 	),
					// )).
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(idCondition),
					)).
					Times(1).
					Return(&domain.Session{
						ID:     "session-1",
						UserID: "user-123",
					}, nil)

				userFactor := &domain.SessionFactorUser{UserID: "user-123", LastVerifiedAt: time.Now()}

				// TODO(IAM-Marco): Fix this once Session repo is implemented
				factorChange := database.NewChange(
					database.NewColumn("sessions", "user_factor"),
					"user-123",
				)

				repo.EXPECT().
					SetFactor(userFactor).
					Times(1).
					Return(factorChange)

				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, factorChange).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := &domain.UserCheckCommand{
				CheckUser:   tc.checkUser,
				SessionID:   tc.sessionID,
				InstanceID:  tc.instanceID,
				FetchedUser: tc.fetchedUser,
			}

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.sessionRepo != nil {
				domain.WithSessionRepo(tc.sessionRepo(ctrl))(opts)
			}

			// Test
			err := cmd.Execute(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError == nil && tc.checkUser != nil {
				assert.NotZero(t, cmd.UserCheckedAt)
				assert.Equal(t, tc.expectedPreferredLanguage, cmd.PreferredUserLanguage)
			}
		})
	}
}

func TestUserCheckCommand_Events(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName         string
		checkUser        *session_grpc.CheckUser
		sessionID        string
		instanceID       string
		fetchedUser      domain.User
		userCheckedAt    time.Time
		expectedEventLen int
		expectedError    error
	}{
		{
			testName:         "when CheckUser is nil should return nil",
			checkUser:        nil,
			expectedEventLen: 0,
			expectedError:    nil,
		},
		{
			testName: "when CheckUser is set should return user checked event",
			checkUser: &session_grpc.CheckUser{
				Search: &session_grpc.CheckUser_UserId{
					UserId: "user-123",
				},
			},
			sessionID:        "session-1",
			instanceID:       "instance-1",
			fetchedUser:      domain.User{ID: "user-123", OrganizationID: "org-1"},
			userCheckedAt:    time.Now(),
			expectedEventLen: 1,
			expectedError:    nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext(tc.instanceID, "", "")
			cmd := &domain.UserCheckCommand{
				CheckUser:     tc.checkUser,
				SessionID:     tc.sessionID,
				InstanceID:    tc.instanceID,
				FetchedUser:   tc.fetchedUser,
				UserCheckedAt: tc.userCheckedAt,
			}

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}

			// Test
			events, err := cmd.Events(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			require.Len(t, events, tc.expectedEventLen)
			if tc.expectedEventLen != 0 {
				usrCheckedEvent, ok := events[0].(*session.UserCheckedEvent)
				require.True(t, ok)

				assert.Equal(t, tc.fetchedUser.ID, usrCheckedEvent.UserID)
				assert.Equal(t, tc.fetchedUser.OrganizationID, usrCheckedEvent.UserResourceOwner)
				assert.NotZero(t, usrCheckedEvent.CheckedAt)
			}
		})
	}
}

func TestNewUserCheckCommand(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName   string
		sessionID  string
		instanceID string
		expected   *domain.UserCheckCommand
	}{
		{
			testName:   "should create command with provided sessionID and instanceID",
			sessionID:  "session-1",
			instanceID: "instance-1",
			expected: &domain.UserCheckCommand{
				SessionID:  "session-1",
				InstanceID: "instance-1",
			},
		},
		{
			testName:   "should create command with empty sessionID",
			sessionID:  "",
			instanceID: "instance-1",
			expected: &domain.UserCheckCommand{
				SessionID:  "",
				InstanceID: "instance-1",
			},
		},
		{
			testName:   "should create command with empty instanceID",
			sessionID:  "session-1",
			instanceID: "",
			expected: &domain.UserCheckCommand{
				SessionID:  "session-1",
				InstanceID: "",
			},
		},
		{
			testName:   "should create command with both empty",
			sessionID:  "",
			instanceID: "",
			expected: &domain.UserCheckCommand{
				SessionID:  "",
				InstanceID: "",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Test
			cmd := domain.NewUserCheckCommand(tc.sessionID, tc.instanceID)

			// Verify
			assert.Equal(t, tc.expected.SessionID, cmd.SessionID)
			assert.Equal(t, tc.expected.InstanceID, cmd.InstanceID)
			assert.Nil(t, cmd.CheckUser)
			assert.Zero(t, cmd.UserCheckedAt)
		})
	}
}
