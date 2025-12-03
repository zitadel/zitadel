package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/passwap"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestPasswordCheckCommand_Validate(t *testing.T) {
	t.Parallel()
	sessionGetErr := errors.New("session get error")
	userGetErr := errors.New("user get error")
	notFoundErr := database.NewNoRowFoundError(nil)

	tt := []struct {
		testName      string
		sessionRepo   func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo      func(ctrl *gomock.Controller) domain.UserRepository
		checkPassword *session_grpc.CheckPassword
		sessionID     string
		expectedError error
		expectedUser  domain.User
	}{
		{
			testName:      "when checkPassword is nil should return no error",
			checkPassword: nil,
			expectedError: nil,
		},
		{
			testName:      "when retrieving session fails should return error",
			checkPassword: &session_grpc.CheckPassword{Password: "test-password"},
			sessionID:     "session-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							idCondition,
						),
					)).
					Times(1).
					Return(nil, sessionGetErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(sessionGetErr, "DOM-qAoQrg", "failed fetching session"),
		},
		{
			testName:      "when session not found should return not found error",
			checkPassword: &session_grpc.CheckPassword{Password: "test-password"},
			sessionID:     "session-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							idCondition,
						),
					)).
					Times(1).
					Return(nil, notFoundErr)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(notFoundErr, "DOM-0XRmp8", "session not found"),
		},
		{
			testName:      "when session userID is empty should return precondition failed error",
			checkPassword: &session_grpc.CheckPassword{Password: "test-password"},
			sessionID:     "session-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							idCondition,
						),
					)).
					Times(1).
					Return(&domain.Session{
						ID:     "session-1",
						UserID: "",
					}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-hord0Z", "Errors.User.UserIDMissing"),
		},
		{
			testName:      "when retrieving user fails should return error",
			checkPassword: &session_grpc.CheckPassword{Password: "test-password"},
			sessionID:     "session-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							idCondition,
						),
					)).
					Times(1).
					Return(&domain.Session{
						ID:     "session-1",
						UserID: "user-1",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							idCondition,
						),
					)).
					Times(1).
					Return(nil, userGetErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(userGetErr, "DOM-nKD4Gq", "failed fetching user"),
		},
		{
			testName:      "when user not found should return not found error",
			checkPassword: &session_grpc.CheckPassword{Password: "test-password"},
			sessionID:     "session-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							idCondition,
						),
					)).
					Times(1).
					Return(&domain.Session{
						ID:     "session-1",
						UserID: "user-1",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							idCondition,
						),
					)).
					Times(1).
					Return(nil, notFoundErr)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(notFoundErr, "DOM-zxKosn", "Errors.User.NotFound"),
		},
		{
			testName:      "when user is not human should return precondition failed error",
			checkPassword: &session_grpc.CheckPassword{Password: "test-password"},
			sessionID:     "session-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							idCondition,
						),
					)).
					Times(1).
					Return(&domain.Session{
						ID:     "session-1",
						UserID: "user-1",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							idCondition,
						),
					)).
					Times(1).
					Return(&domain.User{
						ID:    "user-1",
						Human: nil,
					}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-ADhxAx", "user not human"),
		},
		{
			testName:      "when user is locked should return precondition failed error",
			checkPassword: &session_grpc.CheckPassword{Password: "test-password"},
			sessionID:     "session-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							idCondition,
						),
					)).
					Times(1).
					Return(&domain.Session{
						ID:     "session-1",
						UserID: "user-1",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							idCondition,
						),
					)).
					Times(1).
					Return(&domain.User{
						ID:    "user-1",
						State: domain.UserStateLocked,
						Human: &domain.HumanUser{
							Password: domain.HumanPassword{
								FailedAttempts: 5,
							},
						},
					}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(domain.NewPasswordVerificationError(5), "DOM-D804Sj", "Errors.User.Locked"),
		},
		{
			testName:      "when user password is not set should return precondition failed error",
			checkPassword: &session_grpc.CheckPassword{Password: "test-password"},
			sessionID:     "session-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							idCondition,
						),
					)).
					Times(1).
					Return(&domain.Session{
						ID:     "session-1",
						UserID: "user-1",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							idCondition,
						),
					)).
					Times(1).
					Return(&domain.User{
						ID:    "user-1",
						State: domain.UserStateActive,
						Human: &domain.HumanUser{
							Password: domain.HumanPassword{
								Password: "",
							},
						},
					}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-gklgos", "Errors.User.Password.NotSet"),
		},
		{
			testName:      "when all validations pass should return no error and set user",
			checkPassword: &session_grpc.CheckPassword{Password: "test-password"},
			sessionID:     "session-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							idCondition,
						),
					)).
					Times(1).
					Return(&domain.Session{
						ID:     "session-1",
						UserID: "user-1",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							idCondition,
						),
					)).
					Times(1).
					Return(&domain.User{
						ID:    "user-1",
						State: domain.UserStateActive,
						Human: &domain.HumanUser{
							Password: domain.HumanPassword{
								Password:       "hashed-password",
								FailedAttempts: 0,
							},
						},
					}, nil)
				return repo
			},
			expectedUser: domain.User{
				ID:    "user-1",
				State: domain.UserStateActive,
				Human: &domain.HumanUser{
					Password: domain.HumanPassword{
						Password:       "hashed-password",
						FailedAttempts: 0,
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewPasswordCheckCommand("session-1", "instance-1", nil, nil, tc.checkPassword)

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.sessionRepo != nil {
				domain.WithSessionRepo(tc.sessionRepo(ctrl))(opts)
			}
			if tc.userRepo != nil {
				domain.WithUserRepo(tc.userRepo(ctrl))(opts)
			}

			// Test
			err := cmd.Validate(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedUser, cmd.FetchedUser)
		})
	}
}

func TestPasswordCheckCommand_GetPasswordCheckAndError(t *testing.T) {
	t.Parallel()

	domain.SetPasswordHasher(&crypto.Hasher{Swapper: &passwap.Swapper{}})
	tt := []struct {
		testName          string
		inputErr          error
		fetchedUser       domain.User
		expectedCheckType domain.PasswordCheckType
		expectedError     error
	}{
		{
			testName:          "when error is nil should return CheckTypeSucceeded with no error",
			inputErr:          nil,
			fetchedUser:       domain.User{},
			expectedCheckType: &domain.CheckTypeSucceeded{},
			expectedError:     nil,
		},
		{
			testName: "when error is password mismatch should return CheckTypeFailed with invalid argument error",
			inputErr: passwap.ErrPasswordMismatch,
			fetchedUser: domain.User{
				Human: &domain.HumanUser{
					Password: domain.HumanPassword{
						FailedAttempts: 2,
					},
				},
			},
			expectedCheckType: &domain.CheckTypeFailed{},
			expectedError:     zerrors.ThrowInvalidArgument(domain.NewPasswordVerificationError(3), "DOM-3gcfDV", "Errors.User.Password.Invalid"),
		},
		{
			testName: "when error is other error should return CheckTypeFailed with internal error",
			inputErr: errors.New("some other error"),
			fetchedUser: domain.User{
				Human: &domain.HumanUser{
					Password: domain.HumanPassword{
						FailedAttempts: 0,
					},
				},
			},
			expectedCheckType: &domain.CheckTypeFailed{},
			expectedError:     zerrors.ThrowInternal(errors.New("some other error"), "DOM-xceNzI", "Errors.Internal"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			cmd := domain.NewPasswordCheckCommand("session-1", "instance-1", nil, nil, &session_grpc.CheckPassword{Password: "test"})
			cmd.FetchedUser = tc.fetchedUser

			// Test
			checkType, err := cmd.GetPasswordCheckAndError(tc.inputErr)

			// Verify
			assert.IsType(t, tc.expectedCheckType, checkType)
			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError != nil {
				checkFailed := checkType.(*domain.CheckTypeFailed)
				assert.NotZero(t, checkFailed.FailedAt)
			} else {
				checkSuccess := checkType.(*domain.CheckTypeSucceeded)
				assert.NotZero(t, checkSuccess.SucceededAt)
			}
		})
	}
}

func TestPasswordCheckCommand_GetPasswordCheckChanges(t *testing.T) {
	t.Parallel()
	listErr := errors.New("list error")
	domain.SetPasswordHasher(&crypto.Hasher{Swapper: &passwap.Swapper{}})

	tt := []struct {
		testName    string
		humanRepo   func(ctrl *gomock.Controller, checkTime time.Time) domain.HumanUserRepository
		lockoutRepo func(ctrl *gomock.Controller) domain.LockoutSettingsRepository
		updatedHash string
		checkType   domain.PasswordCheckType
		checkTime   time.Time
		fetchedUser domain.User

		expectedChanges    int
		expectedError      error
		expectedHashedPsw  string
		expectedUserLocked bool
	}{
		{
			testName: "when check type succeeded with empty hash should return only check password change",
			humanRepo: func(ctrl *gomock.Controller, checkTime time.Time) domain.HumanUserRepository {
				repo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					CheckPassword(&domain.CheckTypeSucceeded{}).
					Times(1)
				return repo
			},
			updatedHash:     "",
			checkType:       &domain.CheckTypeSucceeded{},
			expectedChanges: 1,
		},
		{
			testName: "when check type succeeded with hash should return check password and set password changes",
			humanRepo: func(ctrl *gomock.Controller, checkTime time.Time) domain.HumanUserRepository {
				repo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					CheckPassword(&domain.CheckTypeSucceeded{}).
					Times(1)
				repo.EXPECT().
					SetPassword(&domain.VerificationTypeSkipped{Value: gu.Ptr("new-hash"), VerifiedAt: checkTime}).
					Times(1)
				return repo
			},
			updatedHash:       "new-hash",
			checkType:         &domain.CheckTypeSucceeded{},
			expectedChanges:   2,
			expectedHashedPsw: "new-hash",
		},
		{
			testName: "when check type failed and no lockout policy should return internal error",
			humanRepo: func(ctrl *gomock.Controller, checkTime time.Time) domain.HumanUserRepository {
				repo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					CheckPassword(&domain.CheckTypeFailed{}).
					Times(1)
				return repo
			},
			lockoutRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				repo := domainmock.NewLockoutSettingsRepo(ctrl)

				instanceAndOrg := database.And(repo.InstanceIDCondition("instance-1"), repo.OrganizationIDCondition(gu.Ptr("org-1")))
				orgNullOrEmpty := database.Or(repo.OrganizationIDCondition(nil), repo.OrganizationIDCondition(gu.Ptr("")))
				onlyInstance := database.And(repo.InstanceIDCondition("instance-1"), orgNullOrEmpty)
				orCondition := database.Or(instanceAndOrg, onlyInstance)

				repo.EXPECT().
					List(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(orCondition)),
						dbmock.QueryOptions(database.WithOrderByAscending(repo.OrganizationIDColumn(), repo.InstanceIDColumn())),
						dbmock.QueryOptions(database.WithLimit(1)),
					).Times(1).
					Return(nil, nil)

				return repo
			},
			updatedHash: "",
			checkType:   &domain.CheckTypeFailed{},
			fetchedUser: domain.User{
				OrganizationID: "org-1",
			},
			expectedError: zerrors.ThrowInternal(domain.NewMultipleObjectsUpdatedError(1, int64(0)), "DOM-mmsrCt", "unexpected number of rows returned"),
		},
		{
			testName: "when check type failed and get lockout policy returns error should return error",
			humanRepo: func(ctrl *gomock.Controller, checkTime time.Time) domain.HumanUserRepository {
				repo := domainmock.NewMockHumanUserRepository(ctrl)
				checkFailed := &domain.CheckTypeFailed{}
				checkPasswordChange := database.NewChange(
					database.NewColumn("zitadel.humans", "password_check_failed_at"),
					checkTime,
				)
				repo.EXPECT().
					CheckPassword(checkFailed).
					Times(1).
					Return(checkPasswordChange)
				return repo
			},
			lockoutRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				repo := domainmock.NewLockoutSettingsRepo(ctrl)

				instanceAndOrg := database.And(repo.InstanceIDCondition("instance-1"), repo.OrganizationIDCondition(gu.Ptr("org-1")))
				orgNullOrEmpty := database.Or(repo.OrganizationIDCondition(nil), repo.OrganizationIDCondition(gu.Ptr("")))
				onlyInstance := database.And(repo.InstanceIDCondition("instance-1"), orgNullOrEmpty)
				orCondition := database.Or(instanceAndOrg, onlyInstance)

				repo.EXPECT().
					List(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(orCondition)),
						dbmock.QueryOptions(database.WithOrderByAscending(repo.OrganizationIDColumn(), repo.InstanceIDColumn())),
						dbmock.QueryOptions(database.WithLimit(1)),
					).Times(1).
					Return(nil, listErr)

				return repo
			},
			updatedHash:   "",
			checkType:     &domain.CheckTypeFailed{},
			expectedError: zerrors.ThrowInternal(listErr, "DOM-3B8Z6s", "failed fetching lockout settings"),
			fetchedUser: domain.User{
				OrganizationID: "org-1",
			},
		},
		{
			testName: "when check type failed with nil max attempts should return check type failed change",
			humanRepo: func(ctrl *gomock.Controller, checkTime time.Time) domain.HumanUserRepository {
				repo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					CheckPassword(&domain.CheckTypeFailed{}).
					Times(1)
				return repo
			},
			lockoutRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				repo := domainmock.NewLockoutSettingsRepo(ctrl)

				instanceAndOrg := database.And(repo.InstanceIDCondition("instance-1"), repo.OrganizationIDCondition(gu.Ptr("org-1")))
				orgNullOrEmpty := database.Or(repo.OrganizationIDCondition(nil), repo.OrganizationIDCondition(gu.Ptr("")))
				onlyInstance := database.And(repo.InstanceIDCondition("instance-1"), orgNullOrEmpty)
				orCondition := database.Or(instanceAndOrg, onlyInstance)

				setting := &domain.LockoutSettings{
					Settings:                  domain.Settings{},
					LockoutSettingsAttributes: domain.LockoutSettingsAttributes{},
				}

				repo.EXPECT().
					List(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(orCondition)),
						dbmock.QueryOptions(database.WithOrderByAscending(repo.OrganizationIDColumn(), repo.InstanceIDColumn())),
						dbmock.QueryOptions(database.WithLimit(1)),
					).Times(1).
					Return([]*domain.LockoutSettings{
						setting,
					}, nil)

				return repo
			},
			updatedHash: "",
			checkType:   &domain.CheckTypeFailed{},
			fetchedUser: domain.User{
				OrganizationID: "org-1",
			},
			expectedChanges: 1,
		},
		{
			testName: "when check type failed and user max attempts >= lockout policy max attempts should return check type and user locked changes",
			humanRepo: func(ctrl *gomock.Controller, checkTime time.Time) domain.HumanUserRepository {
				repo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					CheckPassword(&domain.CheckTypeFailed{}).
					Times(1)

				repo.EXPECT().
					SetState(domain.UserStateLocked).
					Times(1)
				return repo
			},
			lockoutRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				repo := domainmock.NewLockoutSettingsRepo(ctrl)

				instanceAndOrg := database.And(repo.InstanceIDCondition("instance-1"), repo.OrganizationIDCondition(gu.Ptr("org-1")))
				orgNullOrEmpty := database.Or(repo.OrganizationIDCondition(nil), repo.OrganizationIDCondition(gu.Ptr("")))
				onlyInstance := database.And(repo.InstanceIDCondition("instance-1"), orgNullOrEmpty)
				orCondition := database.Or(instanceAndOrg, onlyInstance)

				setting := &domain.LockoutSettings{
					Settings: domain.Settings{},
					LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
						MaxPasswordAttempts: gu.Ptr(uint64(1)),
					},
				}

				repo.EXPECT().
					List(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(orCondition)),
						dbmock.QueryOptions(database.WithOrderByAscending(repo.OrganizationIDColumn(), repo.InstanceIDColumn())),
						dbmock.QueryOptions(database.WithLimit(1)),
					).Times(1).
					Return([]*domain.LockoutSettings{
						setting,
					}, nil)

				return repo
			},
			updatedHash: "",
			checkType:   &domain.CheckTypeFailed{},
			fetchedUser: domain.User{
				OrganizationID: "org-1",
				Human:          &domain.HumanUser{Password: domain.HumanPassword{FailedAttempts: 2}},
			},
			expectedChanges:    2,
			expectedUserLocked: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewPasswordCheckCommand("session-1", "instance-1", nil, func(_, _ string) (_ string, _ error) { return "", nil }, &session_grpc.CheckPassword{Password: "test"})
			cmd.FetchedUser = tc.fetchedUser
			cmd.CheckTime = time.Now()

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.lockoutRepo != nil {
				domain.WithLockoutSettingsRepo(tc.lockoutRepo(ctrl))(opts)
			}

			// Test
			changes, err := cmd.GetPasswordCheckChanges(ctx, opts, tc.humanRepo(ctrl, cmd.CheckTime), tc.updatedHash, tc.checkType)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedHashedPsw, cmd.UpdatedHashedPsw)
			assert.Equal(t, tc.expectedUserLocked, cmd.IsUserLocked)
			assert.Len(t, changes, tc.expectedChanges)
		})
	}
}

func TestPasswordCheckCommand_Execute(t *testing.T) {
	t.Parallel()

	domain.SetPasswordHasher(&crypto.Hasher{Swapper: &passwap.Swapper{}})
	userUpdateErr := errors.New("user update error")
	sessionUpdateErr := errors.New("session update error")

	tt := []struct {
		testName      string
		humanRepo     func(ctrl *gomock.Controller) domain.HumanUserRepository
		lockoutRepo   func(ctrl *gomock.Controller) domain.LockoutSettingsRepository
		sessionRepo   func(ctrl *gomock.Controller) domain.SessionRepository
		verifier      func(password string, hash string) (string, error)
		checkPassword *session_grpc.CheckPassword
		sessionID     string
		instanceID    string
		fetchedUser   domain.User
		tarpitCalled  bool

		expectedError             error
		expectedValidated         bool
		expectedValidationSuccess bool
		expectedUserLocked        bool
	}{
		{
			testName:      "when checkPassword is nil should return no error",
			checkPassword: nil,
			verifier:      nil,
			expectedError: nil,
		},
		{
			testName:      "when user update fails should return error",
			checkPassword: &session_grpc.CheckPassword{Password: "test-password"},
			sessionID:     "session-1",
			instanceID:    "instance-1",
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					Password: domain.HumanPassword{
						Password:       "hashed-password",
						FailedAttempts: 0,
					},
				},
			},
			verifier: func(password string, hash string) (string, error) { return "", nil },
			humanRepo: func(ctrl *gomock.Controller) domain.HumanUserRepository {
				repo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					CheckPassword(gomock.Any()).
					Times(1).
					Return(nil)
				idCondition := getHumanUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					Times(1).
					Return(int64(0), userUpdateErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(userUpdateErr, "DOM-netNam", "failed updating user"),
		},
		{
			testName:      "when user not found should return not found error",
			checkPassword: &session_grpc.CheckPassword{Password: "test-password"},
			sessionID:     "session-1",
			instanceID:    "instance-1",
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					Password: domain.HumanPassword{
						Password:       "hashed-password",
						FailedAttempts: 0,
					},
				},
			},
			verifier: func(password string, hash string) (string, error) { return "", nil },
			humanRepo: func(ctrl *gomock.Controller) domain.HumanUserRepository {
				repo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					CheckPassword(gomock.Any()).
					Times(1).
					Return(nil)
				idCondition := getHumanUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					Times(1).
					Return(int64(0), nil)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(nil, "DOM-8wVrNc", "user not found"),
		},
		{
			testName:      "when user update returns multiple rows should return internal error",
			checkPassword: &session_grpc.CheckPassword{Password: "test-password"},
			sessionID:     "session-1",
			instanceID:    "instance-1",
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					Password: domain.HumanPassword{
						Password:       "hashed-password",
						FailedAttempts: 0,
					},
				},
			},
			verifier: func(password string, hash string) (string, error) { return "", nil },
			humanRepo: func(ctrl *gomock.Controller) domain.HumanUserRepository {
				repo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					CheckPassword(gomock.Any()).
					Times(1).
					Return(nil)
				idCondition := getHumanUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					Times(1).
					Return(int64(2), nil)
				return repo
			},
			expectedError: zerrors.ThrowInternal(domain.NewMultipleObjectsUpdatedError(1, 2), "DOM-D4hy9C", "unexpected number of rows updated"),
		},
		{
			testName:      "when session update fails should return error",
			checkPassword: &session_grpc.CheckPassword{Password: "test-password"},
			sessionID:     "session-1",
			instanceID:    "instance-1",
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					Password: domain.HumanPassword{
						Password:       "hashed-password",
						FailedAttempts: 0,
					},
				},
			},
			verifier: func(password string, hash string) (string, error) { return "", nil },
			humanRepo: func(ctrl *gomock.Controller) domain.HumanUserRepository {
				repo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					CheckPassword(gomock.Any()).
					Times(1).
					Return(nil)
				idCondition := getHumanUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				sessionFactorChange := getSessionPasswordFactorChange(repo, gu.Ptr(time.Now()), nil)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, sessionFactorChange).
					Times(1).
					Return(int64(0), sessionUpdateErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(sessionUpdateErr, "DOM-IZagay", "failed updating session"),
		},
		{
			testName:      "when session not found should return not found error",
			checkPassword: &session_grpc.CheckPassword{Password: "test-password"},
			sessionID:     "session-1",
			instanceID:    "instance-1",
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					Password: domain.HumanPassword{
						Password:       "hashed-password",
						FailedAttempts: 0,
					},
				},
			},
			verifier: func(password string, hash string) (string, error) { return "", nil },
			humanRepo: func(ctrl *gomock.Controller) domain.HumanUserRepository {
				repo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					CheckPassword(gomock.Any()).
					Times(1).
					Return(nil)
				idCondition := getHumanUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				sessionFactorChange := getSessionPasswordFactorChange(repo, gu.Ptr(time.Now()), nil)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, sessionFactorChange).
					Times(1).
					Return(int64(0), nil)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(nil, "DOM-H9Q59c", "session not found"),
		},
		{
			testName:      "when session update returns multiple rows should return internal error",
			checkPassword: &session_grpc.CheckPassword{Password: "test-password"},
			sessionID:     "session-1",
			instanceID:    "instance-1",
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					Password: domain.HumanPassword{
						Password:       "hashed-password",
						FailedAttempts: 0,
					},
				},
			},
			verifier: func(password string, hash string) (string, error) { return "", nil },
			humanRepo: func(ctrl *gomock.Controller) domain.HumanUserRepository {
				repo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					CheckPassword(gomock.Any()).
					Times(1).
					Return(nil)
				idCondition := getHumanUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				sessionFactorChange := getSessionPasswordFactorChange(repo, gu.Ptr(time.Now()), nil)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, sessionFactorChange).
					Times(1).
					Return(int64(2), nil)
				return repo
			},
			expectedError: zerrors.ThrowInternal(domain.NewMultipleObjectsUpdatedError(1, 2), "DOM-Tbvpy8", "unexpected number of rows updated"),
		},
		{
			testName:      "when password check succeeds should execute successfully",
			checkPassword: &session_grpc.CheckPassword{Password: "test-password"},
			sessionID:     "session-1",
			instanceID:    "instance-1",
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					Password: domain.HumanPassword{
						Password:       "hashed-password",
						FailedAttempts: 0,
					},
				},
			},
			verifier: func(password string, hash string) (string, error) { return "", nil },
			humanRepo: func(ctrl *gomock.Controller) domain.HumanUserRepository {
				repo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					CheckPassword(gomock.Any()).
					Times(1).
					Return(nil)
				idCondition := getHumanUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				sessionFactorChange := getSessionPasswordFactorChange(repo, gu.Ptr(time.Now()), nil)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, sessionFactorChange).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			expectedError:             nil,
			expectedValidated:         true,
			expectedValidationSuccess: true,
		},
		{
			testName:      "when password check fails should execute successfully with failed state",
			checkPassword: &session_grpc.CheckPassword{Password: "wrong-password"},
			sessionID:     "session-1",
			instanceID:    "instance-1",
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					Password: domain.HumanPassword{
						Password:       "hashed-password",
						FailedAttempts: 2,
					},
				},
			},
			verifier: func(password string, hash string) (string, error) {
				return "", passwap.ErrPasswordMismatch
			},
			lockoutRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				repo := domainmock.NewLockoutSettingsRepo(ctrl)

				instanceAndOrg := database.And(repo.InstanceIDCondition("instance-1"), repo.OrganizationIDCondition(gu.Ptr("org-1")))
				orgNullOrEmpty := database.Or(repo.OrganizationIDCondition(nil), repo.OrganizationIDCondition(gu.Ptr("")))
				onlyInstance := database.And(repo.InstanceIDCondition("instance-1"), orgNullOrEmpty)
				orCondition := database.Or(instanceAndOrg, onlyInstance)

				setting := &domain.LockoutSettings{
					Settings: domain.Settings{},
					LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
						MaxPasswordAttempts: gu.Ptr(uint64(1)),
					},
				}

				repo.EXPECT().
					List(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(orCondition)),
						dbmock.QueryOptions(database.WithOrderByAscending(repo.OrganizationIDColumn(), repo.InstanceIDColumn())),
						dbmock.QueryOptions(database.WithLimit(1)),
					).Times(1).
					Return([]*domain.LockoutSettings{
						setting,
					}, nil)

				return repo
			},
			humanRepo: func(ctrl *gomock.Controller) domain.HumanUserRepository {
				repo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					CheckPassword(gomock.Any()).
					Times(1).
					Return(nil)
				repo.EXPECT().
					SetState(domain.UserStateLocked).
					Times(1)
				idCondition := getHumanUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				sessionFactorChange := getSessionPasswordFactorChange(repo, nil, gu.Ptr(time.Now()))
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, sessionFactorChange).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			tarpitCalled:              true,
			expectedError:             zerrors.ThrowInvalidArgument(domain.NewPasswordVerificationError(3), "DOM-3gcfDV", "Errors.User.Password.Invalid"),
			expectedValidated:         true,
			expectedValidationSuccess: false,
			expectedUserLocked:        true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)

			tarpitCalled := false
			tarpitFunc := func(failedAttempts uint64) {
				tarpitCalled = true
			}

			cmd := domain.NewPasswordCheckCommand("session-1", "instance-1", tarpitFunc, tc.verifier, tc.checkPassword)
			cmd.FetchedUser = tc.fetchedUser

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.sessionRepo != nil {
				domain.WithSessionRepo(tc.sessionRepo(ctrl))(opts)
			}
			if tc.lockoutRepo != nil {
				domain.WithLockoutSettingsRepo(tc.lockoutRepo(ctrl))(opts)
			}
			if tc.humanRepo != nil {
				humanRepo := tc.humanRepo(ctrl)
				userRepo := domainmock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().Human().Times(1).Return(humanRepo)
				domain.WithUserRepo(userRepo)(opts)
			}

			// Test
			err := cmd.Execute(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			if tc.checkPassword != nil {
				assert.Equal(t, tc.expectedValidated, cmd.IsValidated)
				assert.Equal(t, tc.expectedValidationSuccess, cmd.IsValidationSuccessful)
				assert.Equal(t, tc.expectedUserLocked, cmd.IsUserLocked)
				assert.Equal(t, tc.tarpitCalled, tarpitCalled)
			}
		})
	}
}

func TestPasswordCheckCommand_Events(t *testing.T) {
	t.Parallel()
	domain.SetPasswordHasher(&crypto.Hasher{})
	userAgg := user.NewAggregate("user-1", "org-1").Aggregate
	sessAgg := session.NewAggregate("session-1", "instance-1").Aggregate

	tt := []struct {
		testName               string
		checkPassword          *session_grpc.CheckPassword
		isValidated            bool
		isValidationSuccessful bool
		isUserLocked           bool
		updatedHashedPsw       string
		fetchedUser            domain.User
		sessionID              string
		instanceID             string
		checkTime              time.Time
		expectedEventTypes     []eventstore.Command
	}{
		{
			testName:      "when checkPassword is nil should return nil",
			checkPassword: nil,
			isValidated:   true,
		},
		{
			testName:      "when not validated should return nil",
			checkPassword: &session_grpc.CheckPassword{Password: "test"},
			isValidated:   false,
		},
		{
			testName:               "when validation successful with no hash update should return check succeeded and password checked events",
			checkPassword:          &session_grpc.CheckPassword{Password: "test"},
			isValidated:            true,
			isValidationSuccessful: true,
			updatedHashedPsw:       "",
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
			},
			sessionID:  "session-1",
			instanceID: "instance-1",
			checkTime:  time.Now(),
			expectedEventTypes: []eventstore.Command{
				user.NewHumanPasswordCheckSucceededEvent(t.Context(), &userAgg, nil),
				session.NewPasswordCheckedEvent(t.Context(), &sessAgg, time.Now()),
			},
		},
		{
			testName:               "when validation successful with hash update should return check succeeded, hash updated, and password checked events",
			checkPassword:          &session_grpc.CheckPassword{Password: "test"},
			isValidated:            true,
			isValidationSuccessful: true,
			updatedHashedPsw:       "new-hash",
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
			},
			sessionID:  "session-1",
			instanceID: "instance-1",
			checkTime:  time.Now(),
			expectedEventTypes: []eventstore.Command{
				user.NewHumanPasswordCheckSucceededEvent(t.Context(), &userAgg, nil),
				user.NewHumanPasswordHashUpdatedEvent(t.Context(), &userAgg, "new-hash"),
				session.NewPasswordCheckedEvent(t.Context(), &sessAgg, time.Now()),
			},
		},
		{
			testName:               "when validation failed without user locked should return check failed and password checked events",
			checkPassword:          &session_grpc.CheckPassword{Password: "wrong"},
			isValidated:            true,
			isValidationSuccessful: false,
			isUserLocked:           false,
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
			},
			sessionID:  "session-1",
			instanceID: "instance-1",
			checkTime:  time.Now(),
			expectedEventTypes: []eventstore.Command{
				user.NewHumanPasswordCheckFailedEvent(t.Context(), &userAgg, nil),
				session.NewPasswordCheckedEvent(t.Context(), &sessAgg, time.Now()),
			},
		},
		{
			testName:               "when validation failed with user locked should return check failed, user locked, and password checked events",
			checkPassword:          &session_grpc.CheckPassword{Password: "wrong"},
			isValidated:            true,
			isValidationSuccessful: false,
			isUserLocked:           true,
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
			},
			sessionID:  "session-1",
			instanceID: "instance-1",
			checkTime:  time.Now(),
			expectedEventTypes: []eventstore.Command{
				user.NewHumanPasswordCheckFailedEvent(t.Context(), &userAgg, nil),
				user.NewUserLockedEvent(t.Context(), &userAgg),
				session.NewPasswordCheckedEvent(t.Context(), &sessAgg, time.Now())},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext(tc.instanceID, "", "")
			cmd := domain.NewPasswordCheckCommand(tc.sessionID, tc.instanceID, nil, nil, tc.checkPassword)
			cmd.FetchedUser = tc.fetchedUser
			cmd.IsValidated = tc.isValidated
			cmd.IsValidationSuccessful = tc.isValidationSuccessful
			cmd.IsUserLocked = tc.isUserLocked
			cmd.UpdatedHashedPsw = tc.updatedHashedPsw
			cmd.CheckTime = tc.checkTime

			// Test
			events, err := cmd.Events(ctx, &domain.InvokeOpts{})

			// Verify
			assert.NoError(t, err)
			require.Len(t, events, len(tc.expectedEventTypes))
			for i, expectedType := range tc.expectedEventTypes {
				assert.IsType(t, expectedType, events[i])
				switch expectedAssertedType := expectedType.(type) {
				case *user.HumanPasswordCheckSucceededEvent:
					actualAssertedType, ok := events[i].(*user.HumanPasswordCheckSucceededEvent)
					require.True(t, ok)
					assert.Equal(t, expectedAssertedType.AuthRequestInfo, actualAssertedType.AuthRequestInfo)
				case *user.HumanPasswordCheckFailedEvent:
					actualAssertedType, ok := events[i].(*user.HumanPasswordCheckFailedEvent)
					require.True(t, ok)
					assert.Equal(t, expectedAssertedType.AuthRequestInfo, actualAssertedType.AuthRequestInfo)
				case *user.UserLockedEvent:
					continue
				case *user.HumanPasswordHashUpdatedEvent:
					actualAssertedType, ok := events[i].(*user.HumanPasswordHashUpdatedEvent)
					require.True(t, ok)
					assert.Equal(t, expectedAssertedType.EncodedHash, actualAssertedType.EncodedHash)
				case *session.PasswordCheckedEvent:
					actualAssertedType, ok := events[i].(*session.PasswordCheckedEvent)
					require.True(t, ok)
					assert.InDelta(t, expectedAssertedType.CheckedAt.UnixMilli(), actualAssertedType.CheckedAt.UnixMilli(), 1.5)
				}
			}
		})
	}
}
