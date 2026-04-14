package domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestRecoveryCodeCheckCommand_Validate(t *testing.T) {
	t.Parallel()

	getError := errors.New("get error")

	tests := []struct {
		name        string
		sessionID   string
		instanceID  string
		check       *domain.CheckTypeRecoveryCode
		sessionRepo func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo    func(ctrl *gomock.Controller) domain.UserRepository
		wantErr     error
	}{
		{
			name: "no recovery code check to be done",
		},
		{
			name:    "empty recovery code",
			check:   &domain.CheckTypeRecoveryCode{},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-cEKxoG", "Errors.User.MFA.RecoveryCodes.Empty"),
		},
		{
			name:    "no session id",
			check:   &domain.CheckTypeRecoveryCode{RecoveryCode: "test-code"},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-hsbVyd", "Errors.Session.IDMissing"),
		},
		{
			name:      "no instance id",
			sessionID: "session-1",
			check:     &domain.CheckTypeRecoveryCode{RecoveryCode: "test-code"},
			wantErr:   zerrors.ThrowPreconditionFailed(nil, "DOM-lGIe1v", "Errors.Instance.IDMissing"),
		},
		{
			name:       "failed to get session - session not found",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckTypeRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				primaryKeyCondition := sessionRepo.PrimaryKeyCondition("instance-1", "session-1")
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(primaryKeyCondition))).
					Times(1).
					Return(nil, &database.NoRowFoundError{})
				return sessionRepo
			},
			wantErr: zerrors.ThrowNotFound(&database.NoRowFoundError{}, "DOM-2sF2kF", "Session not found"),
		},
		{
			name:       "failed to get session - database error",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckTypeRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				primaryKeyCondition := sessionRepo.PrimaryKeyCondition("instance-1", "session-1")
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(primaryKeyCondition))).
					Times(1).
					Return(nil, getError)
				return sessionRepo
			},
			wantErr: zerrors.ThrowInternal(getError, "DOM-2sF2kF", "failed fetching Session"),
		},
		{
			name:       "missing user id in session",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckTypeRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				primaryKeyCondition := sessionRepo.PrimaryKeyCondition("instance-1", "session-1")
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(primaryKeyCondition))).
					Times(1).
					Return(&domain.Session{}, nil)
				return sessionRepo
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-EaLqwq", "Errors.User.UserIDMissing"),
		},
		{
			name:       "failed to get user - user not found",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckTypeRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				primaryKeyCondition := sessionRepo.PrimaryKeyCondition("instance-1", "session-1")
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(primaryKeyCondition))).
					Times(1).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				userIDCondition := userRepo.IDCondition("user-1")
				userRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(userIDCondition))).
					Times(1).
					Return(nil, &database.NoRowFoundError{})
				return userRepo
			},
			wantErr: zerrors.ThrowNotFound(&database.NoRowFoundError{}, "DOM-Ot3qO6", "User not found"),
		},
		{
			name:       "failed to get user - database error",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckTypeRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				primaryKeyCondition := sessionRepo.PrimaryKeyCondition("instance-1", "session-1")
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(primaryKeyCondition))).
					Times(1).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				userIDCondition := userRepo.IDCondition("user-1")
				userRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(userIDCondition))).
					Times(1).
					Return(nil, getError)
				return userRepo
			},
			wantErr: zerrors.ThrowInternal(getError, "DOM-Ot3qO6", "failed fetching User"),
		},
		{
			name:       "user is locked",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckTypeRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				primaryKeyCondition := sessionRepo.PrimaryKeyCondition("instance-1", "session-1")
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(primaryKeyCondition))).
					Times(1).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				userIDCondition := userRepo.IDCondition("user-1")
				userRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(userIDCondition))).
					Times(1).
					Return(&domain.User{State: domain.UserStateLocked}, nil)
				return userRepo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-47H1Ii", "Errors.User.Locked"),
		},
		{
			name:       "user does not have recovery codes - not a human user",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckTypeRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				primaryKeyCondition := sessionRepo.PrimaryKeyCondition("instance-1", "session-1")
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(primaryKeyCondition))).
					Times(1).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				userIDCondition := userRepo.IDCondition("user-1")
				userRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(userIDCondition))).
					Times(1).
					Return(
						&domain.User{
							InstanceID:     "instance-1",
							OrganizationID: "org-1",
							State:          domain.UserStateActive,
						}, nil)
				return userRepo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-tzN2a1", "Errors.User.MFA.RecoveryCodes.NotReady"),
		},
		{
			name:       "user does not have recovery codes",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckTypeRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				primaryKeyCondition := sessionRepo.PrimaryKeyCondition("instance-1", "session-1")
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(primaryKeyCondition))).
					Times(1).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				userIDCondition := userRepo.IDCondition("user-1")
				userRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(userIDCondition))).
					Times(1).
					Return(
						&domain.User{
							InstanceID:     "instance-1",
							OrganizationID: "org-1",
							State:          domain.UserStateActive,
							Human:          &domain.HumanUser{},
						}, nil)
				return userRepo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-tzN2a1", "Errors.User.MFA.RecoveryCodes.NotReady"),
		},
		{
			name:       "valid request",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckTypeRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				primaryKeyCondition := sessionRepo.PrimaryKeyCondition("instance-1", "session-1")
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(primaryKeyCondition))).
					Times(1).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				userIDCondition := userRepo.IDCondition("user-1")
				userRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(userIDCondition))).
					Times(1).
					Return(
						&domain.User{
							InstanceID:     "instance-1",
							OrganizationID: "org-1",
							State:          domain.UserStateActive,
							Human: &domain.HumanUser{
								RecoveryCodes: &domain.HumanRecoveryCodes{
									Codes: []string{"code1", "code2", "code3"},
								},
							},
						}, nil)
				return userRepo
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tt.userRepo != nil {
				domain.WithUserRepo(tt.userRepo(ctrl))(opts)
			}
			if tt.sessionRepo != nil {
				domain.WithSessionRepo(tt.sessionRepo(ctrl))(opts)
			}

			cmd := domain.NewRecoveryCodeCheckCommand(tt.sessionID, tt.instanceID, tt.check, nil)
			got := cmd.Validate(context.Background(), opts)
			assert.ErrorIs(t, got, tt.wantErr)
		})
	}
}

func TestRecoveryCodeCheckCommand_Execute(t *testing.T) {
	t.Parallel()
	dbError := errors.New("db error")

	tests := []struct {
		name        string
		sessionID   string
		instanceID  string
		check       *domain.CheckTypeRecoveryCode
		verify      func(encoded, password string) (updated string, err error)
		sessionRepo func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo    func(ctrl *gomock.Controller) domain.UserRepository
		humanRepo   func(ctrl *gomock.Controller) domain.HumanUserRepository
		lockoutRepo func(ctrl *gomock.Controller) domain.LockoutSettingsRepository
		wantErr     error
	}{
		{
			name: "no recovery code check to be done",
		},
		{
			name:       "invalid recovery code - failed to update user",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "", errors.New("verify error")
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 2)

				// set up expectation to update user in Execute()
				updateHumanUserFailedExpectation(ctrl, userRepo, true, dbError, 0)

				return userRepo
			},
			lockoutRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				lockoutRepo := domainmock.NewLockoutSettingsRepo(ctrl)
				// set up expectation to get lockout settings in Execute()
				getLockoutSettingsSucceededExpectation("instance-1", "org-1", lockoutRepo)
				return lockoutRepo
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-845kaq", "Errors.User.MFA.RecoveryCodes.InvalidCode"),
		},
		{
			name:       "invalid recovery code - user not found",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "", errors.New("verify error")
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 2)

				// set up expectation to update user in Execute()
				updateHumanUserFailedExpectation(ctrl, userRepo, true, nil, 0)

				return userRepo
			},
			lockoutRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				lockoutRepo := domainmock.NewLockoutSettingsRepo(ctrl)
				// set up expectation to get lockout settings in Execute()
				getLockoutSettingsSucceededExpectation("instance-1", "org-1", lockoutRepo)
				return lockoutRepo
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-845kaq", "Errors.User.MFA.RecoveryCodes.InvalidCode"),
		},
		{
			name:       "invalid recovery code - multiple user rows updated",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "", errors.New("verify error")
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 2)

				// set up expectation to update user in Execute()
				updateHumanUserFailedExpectation(ctrl, userRepo, true, nil, 2)

				return userRepo
			},
			lockoutRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				lockoutRepo := domainmock.NewLockoutSettingsRepo(ctrl)
				// set up expectation to get lockout settings in Execute()
				getLockoutSettingsSucceededExpectation("instance-1", "org-1", lockoutRepo)
				return lockoutRepo
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-845kaq", "Errors.User.MFA.RecoveryCodes.InvalidCode"),
		},
		{
			name:       "invalid recovery code - failed to update session",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "", errors.New("verify error")
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)

				// set up expectation to update session in Execute()
				factor := sessionRepo.SetFactor(&domain.SessionFactorRecoveryCode{
					LastFailedAt: time.Now(),
				})
				updateSessionFailedExpectation(sessionRepo, factor, dbError, 0)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 0)

				// set up expectation to update user in Execute()
				humanRepo := domainmock.NewHumanRepo(ctrl)
				updateHumanUserSucceededExpectation(userRepo, humanRepo, humanRepo.IncrementRecoveryCodeFailedAttempts())

				return userRepo
			},
			lockoutRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				lockoutRepo := domainmock.NewLockoutSettingsRepo(ctrl)
				// set up expectation to get lockout settings in Execute()
				getLockoutSettingsSucceededExpectation("instance-1", "org-1", lockoutRepo)
				return lockoutRepo
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-845kaq", "Errors.User.MFA.RecoveryCodes.InvalidCode"),
		},
		{
			name:       "invalid recovery code - session not found",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "", errors.New("verify error")
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)

				// set up expectation to update session in Execute()
				factor := sessionRepo.SetFactor(&domain.SessionFactorRecoveryCode{
					LastFailedAt: time.Now(),
				})
				updateSessionFailedExpectation(sessionRepo, factor, nil, 0)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 0)

				// set up expectation to update user in Execute()
				humanRepo := domainmock.NewHumanRepo(ctrl)
				updateHumanUserSucceededExpectation(userRepo, humanRepo, humanRepo.IncrementRecoveryCodeFailedAttempts())

				return userRepo
			},
			lockoutRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				lockoutRepo := domainmock.NewLockoutSettingsRepo(ctrl)
				// set up expectation to get lockout settings in Execute()
				getLockoutSettingsSucceededExpectation("instance-1", "org-1", lockoutRepo)
				return lockoutRepo
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-845kaq", "Errors.User.MFA.RecoveryCodes.InvalidCode"),
		},
		{
			name:       "invalid recovery code - multiple session rows updated",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "", errors.New("verify error")
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)

				// set up expectation to update session in Execute()
				factor := sessionRepo.SetFactor(&domain.SessionFactorRecoveryCode{
					LastFailedAt: time.Now(),
				})
				updateSessionFailedExpectation(sessionRepo, factor, nil, 2)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 0)

				// set up expectation to update user in Execute()
				humanRepo := domainmock.NewHumanRepo(ctrl)
				updateHumanUserSucceededExpectation(userRepo, humanRepo, humanRepo.IncrementRecoveryCodeFailedAttempts())

				return userRepo
			},
			lockoutRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				lockoutRepo := domainmock.NewLockoutSettingsRepo(ctrl)
				// set up expectation to get lockout settings in Execute()
				getLockoutSettingsSucceededExpectation("instance-1", "org-1", lockoutRepo)
				return lockoutRepo
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-845kaq", "Errors.User.MFA.RecoveryCodes.InvalidCode"),
		},
		{
			name:       "invalid recovery code - update user and session with failed status and return error",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "", errors.New("verify error")
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)

				// set up expectation to update session in Execute()
				factor := sessionRepo.SetFactor(&domain.SessionFactorRecoveryCode{
					LastFailedAt: time.Now(),
				})
				updateSessionSucceededExpectation(sessionRepo, factor)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 0)

				// set up expectation to update user in Execute()
				humanRepo := domainmock.NewHumanRepo(ctrl)
				updateHumanUserSucceededExpectation(userRepo, humanRepo, humanRepo.IncrementRecoveryCodeFailedAttempts())

				return userRepo
			},
			lockoutRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				lockoutRepo := domainmock.NewLockoutSettingsRepo(ctrl)
				// set up expectation to get lockout settings in Execute()
				getLockoutSettingsSucceededExpectation("instance-1", "org-1", lockoutRepo)
				return lockoutRepo
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-845kaq", "Errors.User.MFA.RecoveryCodes.InvalidCode"),
		},
		{
			name:       "invalid recovery code - set user as locked, update user and session with failed status, and return error",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "", errors.New("verify error")
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)

				// set up expectation to update session in Execute()
				factor := sessionRepo.SetFactor(&domain.SessionFactorRecoveryCode{
					LastFailedAt: time.Now(),
				})
				updateSessionSucceededExpectation(sessionRepo, factor)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 2)

				// set up expectation to update user in Execute()
				humanRepo := domainmock.NewHumanRepo(ctrl)
				updateHumanUserSucceededExpectation(userRepo, humanRepo, humanRepo.IncrementRecoveryCodeFailedAttempts(),
					humanRepo.SetState(domain.UserStateLocked))

				return userRepo
			},
			lockoutRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				lockoutRepo := domainmock.NewLockoutSettingsRepo(ctrl)
				// set up expectation to get lockout settings in Execute()
				getLockoutSettingsSucceededExpectation("instance-1", "org-1", lockoutRepo)
				return lockoutRepo
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-845kaq", "Errors.User.MFA.RecoveryCodes.InvalidCode"),
		},
		{
			name:       "recovery code check succeeded - user update failed, return error",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "hashed-code1", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)

				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 2)

				// set up expectation to update user in Execute()
				humanRepo := domainmock.NewHumanRepo(ctrl)
				primaryKeyCondition := humanRepo.PrimaryKeyCondition("instance-1", "user-1")
				humanRepo.EXPECT().
					Update(gomock.Any(),
						gomock.Any(),
						primaryKeyCondition,
						gomock.Any(),
						humanRepo.RemoveRecoveryCode("hashed-code1"),
					).Times(1).
					Return(int64(0), dbError)
				userRepo.EXPECT().Human().Times(1).Return(humanRepo)
				return userRepo
			},
			wantErr: zerrors.ThrowInternal(dbError, "DOM-XGf3Tk", "user update failed"),
		},
		{
			name:       "recovery code check succeeded - user not found, return error",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "hashed-code1", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)

				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 2)

				// set up expectation to update user in Execute()
				humanRepo := domainmock.NewHumanRepo(ctrl)
				primaryKeyCondition := humanRepo.PrimaryKeyCondition("instance-1", "user-1")
				humanRepo.EXPECT().
					Update(gomock.Any(),
						gomock.Any(),
						primaryKeyCondition,
						gomock.Any(),
						humanRepo.RemoveRecoveryCode("hashed-code1"),
					).Times(1).
					Return(int64(0), nil)
				userRepo.EXPECT().Human().Times(1).Return(humanRepo)
				return userRepo
			},
			wantErr: zerrors.ThrowNotFound(nil, "DOM-hQu5ns", "Errors.User.NotFound"),
		},
		{
			name:       "recovery code check succeeded - multiple user rows updated, return error",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "hashed-code1", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)

				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 2)

				// set up expectation to update user in Execute()
				humanRepo := domainmock.NewHumanRepo(ctrl)
				primaryKeyCondition := humanRepo.PrimaryKeyCondition("instance-1", "user-1")
				humanRepo.EXPECT().
					Update(gomock.Any(),
						gomock.Any(),
						primaryKeyCondition,
						gomock.Any(),
						humanRepo.RemoveRecoveryCode("hashed-code1"),
					).Times(1).
					Return(int64(2), nil)
				userRepo.EXPECT().Human().Times(1).Return(humanRepo)
				return userRepo
			},
			wantErr: zerrors.ThrowInternal(domain.NewMultipleObjectsUpdatedError(1, 2), "DOM-EWjTOH", "Errors.Internal"),
		},
		{
			name:       "recovery code check succeeded - session update failed, return error",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "hashed-code1", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)

				// set up expectation to update session in Execute()
				factor := sessionRepo.SetFactor(&domain.SessionFactorRecoveryCode{
					LastFailedAt: time.Now(),
				})
				updateSessionFailedExpectation(sessionRepo, factor, dbError, 0)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 2)

				// set up expectation to update user in Execute()
				humanRepo := domainmock.NewHumanRepo(ctrl)
				primaryKeyCondition := humanRepo.PrimaryKeyCondition("instance-1", "user-1")
				humanRepo.EXPECT().
					Update(gomock.Any(),
						gomock.Any(),
						primaryKeyCondition,
						gomock.Any(),
						humanRepo.RemoveRecoveryCode("hashed-code1"),
					).Times(1).
					Return(int64(1), nil)
				userRepo.EXPECT().Human().Times(1).Return(humanRepo)

				return userRepo
			},
			wantErr: zerrors.ThrowInternal(dbError, "DOM-fhR4N3", "session update failed"),
		},
		{
			name:       "recovery code check succeeded - session not found, return error",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "hashed-code1", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)

				// set up expectation to update session in Execute()
				factor := sessionRepo.SetFactor(&domain.SessionFactorRecoveryCode{
					LastFailedAt: time.Now(),
				})
				updateSessionFailedExpectation(sessionRepo, factor, nil, 0)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 2)

				// set up expectation to update user in Execute()
				humanRepo := domainmock.NewHumanRepo(ctrl)
				primaryKeyCondition := humanRepo.PrimaryKeyCondition("instance-1", "user-1")
				humanRepo.EXPECT().
					Update(gomock.Any(),
						gomock.Any(),
						primaryKeyCondition,
						gomock.Any(),
						humanRepo.RemoveRecoveryCode("hashed-code1"),
					).Times(1).
					Return(int64(1), nil)
				userRepo.EXPECT().Human().Times(1).Return(humanRepo)

				return userRepo
			},
			wantErr: zerrors.ThrowNotFound(nil, "DOM-Frnwxt", "Errors.Session.NotFound"),
		},
		{
			name:       "recovery code check succeeded - multiple session rows updated, return error",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "hashed-code1", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)

				// set up expectation to update session in Execute()
				factor := sessionRepo.SetFactor(&domain.SessionFactorRecoveryCode{
					LastFailedAt: time.Now(),
				})
				updateSessionFailedExpectation(sessionRepo, factor, nil, 2)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 0)

				// set up expectation to update user in Execute()
				humanRepo := domainmock.NewHumanRepo(ctrl)
				primaryKeyCondition := humanRepo.PrimaryKeyCondition("instance-1", "user-1")
				humanRepo.EXPECT().
					Update(gomock.Any(),
						gomock.Any(),
						primaryKeyCondition,
						gomock.Any(),
						humanRepo.RemoveRecoveryCode("hashed-code1"),
					).Times(1).
					Return(int64(1), nil)
				userRepo.EXPECT().Human().Times(1).Return(humanRepo)

				return userRepo
			},
			wantErr: zerrors.ThrowInternal(domain.NewMultipleObjectsUpdatedError(1, 2), "DOM-gYp8tG", "unexpected number of rows updateCount"),
		},
		{
			name:       "recovery code check succeeded - user and session updated successfully",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "hashed-code1", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)

				// set up expectation to update session in Execute()
				factor := sessionRepo.SetFactor(&domain.SessionFactorRecoveryCode{
					LastVerifiedAt: time.Now(),
				})
				updateSessionSucceededExpectation(sessionRepo, factor)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 2)

				// set up expectation to update user in Execute()
				humanRepo := domainmock.NewHumanRepo(ctrl)
				primaryKeyCondition := humanRepo.PrimaryKeyCondition("instance-1", "user-1")
				humanRepo.EXPECT().
					Update(gomock.Any(),
						gomock.Any(),
						primaryKeyCondition,
						gomock.Any(),
						humanRepo.RemoveRecoveryCode("hashed-code1"),
					).Times(1).
					Return(int64(1), nil)
				userRepo.EXPECT().Human().Times(1).Return(humanRepo)

				return userRepo
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tt.userRepo != nil {
				domain.WithUserRepo(tt.userRepo(ctrl))(opts)
			}
			if tt.sessionRepo != nil {
				domain.WithSessionRepo(tt.sessionRepo(ctrl))(opts)
			}
			if tt.lockoutRepo != nil {
				domain.WithLockoutSettingsRepo(tt.lockoutRepo(ctrl))(opts)
			}

			ctx := context.Background()
			cmd := domain.NewRecoveryCodeCheckCommand(tt.sessionID, tt.instanceID, tt.check, tt.verify)

			// to fetch/validate session and user before calling Execute
			err := cmd.Validate(ctx, opts)
			assert.NoError(t, err)

			got := cmd.Execute(ctx, opts)
			assert.ErrorIs(t, got, tt.wantErr)
		})
	}
}

func TestRecoveryCodeCheckCommand_Events(t *testing.T) {
	t.Parallel()
	updateErr := errors.New("update error")

	tests := []struct {
		name        string
		sessionID   string
		instanceID  string
		check       *domain.CheckTypeRecoveryCode
		verify      func(encoded, password string) (updated string, err error)
		sessionRepo func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo    func(ctrl *gomock.Controller) domain.UserRepository
		humanRepo   func(ctrl *gomock.Controller) domain.HumanUserRepository
		lockoutRepo func(ctrl *gomock.Controller) domain.LockoutSettingsRepository
		wantErr     error
		wantEvents  []eventstore.Command
	}{
		{
			name: "no check needed",
		},
		{
			// relevant events: HumanRecoveryCodeCheckSucceededEvent, NewRecoveryCodeCheckedEvent
			name:       "recovery code check succeeded - return relevant events",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "hashed-code1", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)

				// set up expectation to update session in Execute()
				factor := sessionRepo.SetFactor(&domain.SessionFactorRecoveryCode{
					LastVerifiedAt: time.Now(),
				})
				updateSessionSucceededExpectation(sessionRepo, factor)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 0)

				// set up expectation to update user in Execute()
				humanRepo := domainmock.NewHumanRepo(ctrl)
				primaryKeyCondition := humanRepo.PrimaryKeyCondition("instance-1", "user-1")
				humanRepo.EXPECT().
					Update(gomock.Any(),
						gomock.Any(),
						primaryKeyCondition,
						gomock.Any(),
						humanRepo.RemoveRecoveryCode("hashed-code1"),
					).Times(1).
					Return(int64(1), nil)
				userRepo.EXPECT().Human().Times(1).Return(humanRepo)

				return userRepo
			},
			wantEvents: []eventstore.Command{
				user.NewHumanRecoveryCodeCheckSucceededEvent(
					t.Context(),
					&user.NewAggregate("user-1", "org-1").Aggregate,
					"hashed-code1",
					nil,
				),
				session.NewRecoveryCodeCheckedEvent(
					t.Context(),
					&session.NewAggregate("session-1", "instance-1").Aggregate,
					time.Now(),
				),
			},
		},
		{
			// relevant events: HumanRecoveryCodeCheckFailedEvent;
			// the user isn't locked despite exceeding failed attempts as the recovery code is valid
			name:       "recovery code check succeeded but update failed - return error and relevant events",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "hashed-code1", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 2)

				// set up expectation to update user in Execute()
				humanRepo := domainmock.NewHumanRepo(ctrl)
				primaryKeyCondition := humanRepo.PrimaryKeyCondition("instance-1", "user-1")
				humanRepo.EXPECT().
					Update(gomock.Any(),
						gomock.Any(),
						primaryKeyCondition,
						gomock.Any(),
						humanRepo.RemoveRecoveryCode("hashed-code1"),
					).Times(1).
					Return(int64(0), updateErr)
				userRepo.EXPECT().Human().Times(1).Return(humanRepo)

				return userRepo
			},
			wantEvents: []eventstore.Command{
				user.NewHumanRecoveryCodeCheckFailedEvent(
					t.Context(),
					&user.NewAggregate("user-1", "org-1").Aggregate,
					nil,
				),
			},
			wantErr: zerrors.ThrowInternal(updateErr, "DOM-XGf3Tk", "user update failed"),
		},
		{
			// relevant events: HumanRecoveryCodeCheckFailedEvent and UserLockedEvent
			name:       "invalid recovery code and max OTP attempts reached - return relevant events and error",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "", errors.New("verify error")
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)

				// set up expectation to update session in Execute()
				factor := sessionRepo.SetFactor(&domain.SessionFactorRecoveryCode{
					LastFailedAt: time.Now(),
				})
				updateSessionSucceededExpectation(sessionRepo, factor)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 2) // 2 failed attempts, with max set to 3

				// set up expectation to update user in Execute()
				humanRepo := domainmock.NewHumanRepo(ctrl)
				updateHumanUserSucceededExpectation(
					userRepo,
					humanRepo,
					humanRepo.IncrementRecoveryCodeFailedAttempts(),
					humanRepo.SetState(domain.UserStateLocked),
				)

				return userRepo
			},
			lockoutRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				lockoutRepo := domainmock.NewLockoutSettingsRepo(ctrl)
				// set up expectation to get lockout settings in Execute()
				getLockoutSettingsSucceededExpectation("instance-1", "org-1", lockoutRepo)
				return lockoutRepo
			},
			wantEvents: []eventstore.Command{
				user.NewHumanRecoveryCodeCheckFailedEvent(
					t.Context(),
					&user.NewAggregate("user-1", "org-1").Aggregate,
					nil,
				),
				user.NewUserLockedEvent(
					t.Context(),
					&user.NewAggregate("user-1", "org-1").Aggregate,
				),
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-845kaq", "Errors.User.MFA.RecoveryCodes.InvalidCode"),
		},
		{
			// relevant events: HumanRecoveryCodeCheckFailedEvent
			name:       "invalid recovery code - return relevant events and error",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "", errors.New("verify error")
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)

				// set up expectation to update session in Execute()
				factor := sessionRepo.SetFactor(&domain.SessionFactorRecoveryCode{
					LastFailedAt: time.Now(),
				})
				updateSessionSucceededExpectation(sessionRepo, factor)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 1)

				// set up expectation to update user in Execute()
				humanRepo := domainmock.NewHumanRepo(ctrl)
				updateHumanUserSucceededExpectation(
					userRepo,
					humanRepo, humanRepo.IncrementRecoveryCodeFailedAttempts(),
				)

				return userRepo
			},
			lockoutRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				lockoutRepo := domainmock.NewLockoutSettingsRepo(ctrl)
				// set up expectation to get lockout settings in Execute()
				getLockoutSettingsSucceededExpectation("instance-1", "org-1", lockoutRepo)
				return lockoutRepo
			},
			wantEvents: []eventstore.Command{
				user.NewHumanRecoveryCodeCheckFailedEvent(
					t.Context(),
					&user.NewAggregate("user-1", "org-1").Aggregate,
					nil,
				),
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-845kaq", "Errors.User.MFA.RecoveryCodes.InvalidCode"),
		},
		{
			// relevant events: HumanRecoveryCodeCheckFailedEvent
			name:       "invalid recovery code and update failed - return relevant events and error",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "test-code",
			},
			verify: func(encoded, password string) (updated string, err error) {
				return "", errors.New("verify error")
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				// set up expectation to get session in Validate()
				getSessionSucceededExpectation(sessionRepo)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				// set up expectation to get recovery codes in Validate()
				getUserSucceededExpectation(userRepo, 1)

				// set up expectation to update user in Execute()
				updateHumanUserFailedExpectation(
					ctrl,
					userRepo,
					false,
					updateErr,
					0,
				)

				return userRepo
			},
			lockoutRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				lockoutRepo := domainmock.NewLockoutSettingsRepo(ctrl)
				// set up expectation to get lockout settings in Execute()
				getLockoutSettingsSucceededExpectation("instance-1", "org-1", lockoutRepo)
				return lockoutRepo
			},
			wantEvents: []eventstore.Command{
				user.NewHumanRecoveryCodeCheckFailedEvent(
					t.Context(),
					&user.NewAggregate("user-1", "org-1").Aggregate,
					nil,
				),
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-845kaq", "Errors.User.MFA.RecoveryCodes.InvalidCode"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tt.userRepo != nil {
				domain.WithUserRepo(tt.userRepo(ctrl))(opts)
			}
			if tt.sessionRepo != nil {
				domain.WithSessionRepo(tt.sessionRepo(ctrl))(opts)
			}
			if tt.lockoutRepo != nil {
				domain.WithLockoutSettingsRepo(tt.lockoutRepo(ctrl))(opts)
			}

			ctx := context.Background()
			cmd := domain.NewRecoveryCodeCheckCommand(tt.sessionID, tt.instanceID, tt.check, tt.verify)

			// to fetch/validate session and user before calling Execute
			err := cmd.Validate(ctx, opts)
			assert.NoError(t, err)

			// to update user/session before calling Events
			err = cmd.Execute(ctx, opts)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			got, gotErr := cmd.Events(ctx, opts)
			assert.NoError(t, gotErr)
			assertEvents(t, got, tt.wantEvents)
		})
	}
}

func assertEvents(t *testing.T, got, want []eventstore.Command) {
	require.Len(t, got, len(want))
	for i, event := range want {
		assert.IsType(t, event, got[i])
		switch expectedEvent := event.(type) {
		case *user.HumanRecoveryCodeCheckSucceededEvent:
			gotEvent, ok := got[i].(*user.HumanRecoveryCodeCheckSucceededEvent)
			require.True(t, ok)
			assert.Equal(t, expectedEvent, gotEvent)
		case *session.RecoveryCodeCheckedEvent:
			gotEvent, ok := got[i].(*session.RecoveryCodeCheckedEvent)
			require.True(t, ok)
			assert.NotZero(t, gotEvent.CheckedAt)
			assert.Equal(t, expectedEvent.Aggregate().ID, gotEvent.Aggregate().ID)
			assert.Equal(t, expectedEvent.Aggregate().ResourceOwner, gotEvent.Aggregate().ResourceOwner)
		case *user.HumanRecoveryCodeCheckFailedEvent:
			gotEvent, ok := got[i].(*user.HumanRecoveryCodeCheckFailedEvent)
			require.True(t, ok)
			assert.Equal(t, expectedEvent, gotEvent)
		case *user.UserLockedEvent:
			gotEvent, ok := got[i].(*user.UserLockedEvent)
			require.True(t, ok)
			assert.Equal(t, expectedEvent, gotEvent)
		}
	}

}

func updateHumanUserSucceededExpectation(userRepo *domainmock.UserRepo, humanRepo *domainmock.HumanRepo, userUpdates ...database.Change) {
	primaryKeyCondition := humanRepo.PrimaryKeyCondition("instance-1", "user-1")
	humanRepo.EXPECT().
		Update(gomock.Any(), gomock.Any(), primaryKeyCondition, userUpdates).
		Times(1).
		Return(int64(1), nil)
	userRepo.EXPECT().Human().Times(1).Return(humanRepo)
}

func updateHumanUserFailedExpectation(ctrl *gomock.Controller, userRepo *domainmock.UserRepo, lock bool, err error, updateCount int64) {
	humanRepo := domainmock.NewHumanRepo(ctrl)
	primaryKeyCondition := humanRepo.PrimaryKeyCondition("instance-1", "user-1")

	userUpdates := make([]database.Change, 0, 2)
	userUpdates = append(userUpdates, humanRepo.IncrementRecoveryCodeFailedAttempts())

	if lock {
		userUpdates = append(userUpdates, humanRepo.SetState(domain.UserStateLocked))
	}
	userRepo.EXPECT().Human().Times(1).Return(humanRepo)

	humanRepo.EXPECT().
		Update(gomock.Any(), gomock.Any(), primaryKeyCondition, userUpdates).
		Times(1).
		Return(updateCount, err)
}

func getUserSucceededExpectation(userRepo *domainmock.UserRepo, failedAttempts uint8) {
	userIDCondition := userRepo.IDCondition("user-1")
	userRepo.EXPECT().
		Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(userIDCondition))).
		Times(1).
		Return(
			&domain.User{
				ID:             "user-1",
				InstanceID:     "instance-1",
				OrganizationID: "org-1",
				State:          domain.UserStateActive,
				Human: &domain.HumanUser{
					RecoveryCodes: &domain.HumanRecoveryCodes{
						Codes:          []string{"hashed-code1", "hashed-code2", "hashed-code3"},
						FailedAttempts: failedAttempts,
					},
				},
			}, nil)
}

func updateSessionSucceededExpectation(sessionRepo *domainmock.SessionRepo, change database.Change) {
	primaryKeyCondition := sessionRepo.PrimaryKeyCondition("instance-1", "session-1")
	sessionRepo.EXPECT().
		Update(gomock.Any(), gomock.Any(), primaryKeyCondition, change).
		Times(1).
		Return(int64(1), nil)
}

func updateSessionFailedExpectation(sessionRepo *domainmock.SessionRepo, change database.Change, err error, updateCount int64) {
	primaryKeyCondition := sessionRepo.PrimaryKeyCondition("instance-1", "session-1")
	sessionRepo.EXPECT().
		Update(gomock.Any(), gomock.Any(), primaryKeyCondition, change).
		Times(1).
		Return(updateCount, err)
}

func getSessionSucceededExpectation(sessionRepo *domainmock.SessionRepo) {
	primaryKeyCondition := sessionRepo.PrimaryKeyCondition("instance-1", "session-1")
	sessionRepo.EXPECT().
		Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(primaryKeyCondition))).
		Times(1).
		Return(&domain.Session{
			UserID:     "user-1",
			InstanceID: "instance-1",
			ID:         "session-1",
		}, nil)
}

func getLockoutSettingsSucceededExpectation(instanceID, orgID string, lockoutRepo *domainmock.LockoutSettingsRepo) {
	instanceAndOrg := database.And(lockoutRepo.InstanceIDCondition(instanceID), lockoutRepo.OrganizationIDCondition(gu.Ptr(orgID)))
	orgNullOrEmpty := database.Or(lockoutRepo.OrganizationIDCondition(nil), lockoutRepo.OrganizationIDCondition(gu.Ptr("")))
	onlyInstance := database.And(lockoutRepo.InstanceIDCondition(instanceID), orgNullOrEmpty)
	orCondition := database.Or(instanceAndOrg, onlyInstance)
	lockoutRepo.EXPECT().
		List(gomock.Any(), gomock.Any(),
			dbmock.QueryOptions(database.WithCondition(orCondition)),
			dbmock.QueryOptions(database.WithOrderByAscending(lockoutRepo.OrganizationIDColumn(), lockoutRepo.InstanceIDColumn())),
			dbmock.QueryOptions(database.WithLimit(1)),
		).Times(1).Return(
		[]*domain.LockoutSettings{
			{
				LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
					MaxOTPAttempts: gu.Ptr(uint64(3)),
				},
			},
		}, nil)
}
