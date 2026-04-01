package domain_test

import (
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
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestTOTPCheckCommand_Validate(t *testing.T) {
	t.Parallel()
	sessionGetErr := errors.New("session get error")
	userGetErr := errors.New("user get error")
	notFoundErr := database.NewNoRowFoundError(nil)
	now := time.Now()

	tt := []struct {
		testName    string
		sessionRepo func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo    func(ctrl *gomock.Controller) domain.UserRepository
		cmd         *domain.TOTPCheckCommand

		expectedError error
		expectedUser  domain.User
	}{
		{
			testName:      "when checkTOTP is nil should return no error",
			cmd:           &domain.TOTPCheckCommand{},
			expectedError: nil,
		},
		{
			testName:      "when session ID is not set should return error",
			cmd:           &domain.TOTPCheckCommand{CheckTOTP: &domain.CheckTOTPType{}},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-ZNWO80", "Errors.Missing.SessionID"),
		},
		{
			testName:      "when instance ID is not set should return error",
			cmd:           &domain.TOTPCheckCommand{SessionID: "session-1", CheckTOTP: &domain.CheckTOTPType{}},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-47G8S3", "Errors.Missing.InstanceID"),
		},
		{
			testName: "when retrieving session fails should return error",
			cmd:      &domain.TOTPCheckCommand{SessionID: "session-1", InstanceID: "instance-1", CheckTOTP: &domain.CheckTOTPType{Code: "123456"}},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(nil, sessionGetErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(sessionGetErr, "DOM-e4OuhO", "failed fetching session"),
		},
		{
			testName: "when session not found should return not found error",
			cmd:      &domain.TOTPCheckCommand{SessionID: "session-1", InstanceID: "instance-1", CheckTOTP: &domain.CheckTOTPType{Code: "123456"}},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(nil, notFoundErr)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(notFoundErr, "DOM-e4OuhO", "session not found"),
		},
		{
			testName: "when session userID is empty should return precondition failed error",
			cmd:      &domain.TOTPCheckCommand{SessionID: "session-1", InstanceID: "instance-1", CheckTOTP: &domain.CheckTOTPType{Code: "123456"}},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-hord0Z", "Errors.User.UserIDMissing"),
		},
		{
			testName: "when retrieving user fails should return error",
			cmd:      &domain.TOTPCheckCommand{SessionID: "session-1", InstanceID: "instance-1", CheckTOTP: &domain.CheckTOTPType{Code: "123456"}},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(nil, userGetErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(userGetErr, "DOM-PZvWq0", "failed fetching user"),
		},
		{
			testName: "when user not found should return not found error",
			cmd:      &domain.TOTPCheckCommand{SessionID: "session-1", InstanceID: "instance-1", CheckTOTP: &domain.CheckTOTPType{Code: "123456"}},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(nil, notFoundErr)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(notFoundErr, "DOM-PZvWq0", "user not found"),
		},
		{
			testName: "when user is not human should return precondition failed error",
			cmd:      &domain.TOTPCheckCommand{SessionID: "session-1", InstanceID: "instance-1", CheckTOTP: &domain.CheckTOTPType{Code: "123456"}},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-zzv1MO", "Errors.User.NotHuman"),
		},
		{
			testName: "when user has no TOTP should return precondition failed error",
			cmd:      &domain.TOTPCheckCommand{SessionID: "session-1", InstanceID: "instance-1", CheckTOTP: &domain.CheckTOTPType{Code: "123456"}},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{
						State: domain.UserStateLocked,
						Human: &domain.HumanUser{},
					}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-V6Av2a", "Errors.User.NoTOTP"),
		},
		{
			testName: "when user TOTP has no secret set should return precondition failed error",
			cmd:      &domain.TOTPCheckCommand{SessionID: "session-1", InstanceID: "instance-1", CheckTOTP: &domain.CheckTOTPType{Code: "123456"}},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{
						State: domain.UserStateLocked,
						Human: &domain.HumanUser{TOTP: &domain.HumanTOTP{}},
					}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-b44CWR", "Errors.User.NoTOTPSecret"),
		},
		{
			testName: "when TOTP is not successfully checked should return precondition error",
			cmd:      &domain.TOTPCheckCommand{SessionID: "session-1", InstanceID: "instance-1", CheckTOTP: &domain.CheckTOTPType{Code: "123456"}},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{
						State: domain.UserStateLocked,
						Human: &domain.HumanUser{TOTP: &domain.HumanTOTP{Secret: &crypto.CryptoValue{}}},
					}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-0g4ZAU", "Errors.User.MFA.OTP.NotReady"),
		},
		{
			testName: "when user is locked should return precondition failed error",
			cmd:      &domain.TOTPCheckCommand{SessionID: "session-1", InstanceID: "instance-1", CheckTOTP: &domain.CheckTOTPType{Code: "123456"}},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{
						State: domain.UserStateLocked,
						Human: &domain.HumanUser{TOTP: &domain.HumanTOTP{Secret: &crypto.CryptoValue{}, VerifiedAt: now}},
					}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-gM4SUh", "Errors.User.Locked"),
		},
		{
			testName: "when all validations pass should return no error and set user",
			cmd:      &domain.TOTPCheckCommand{SessionID: "session-1", InstanceID: "instance-1", CheckTOTP: &domain.CheckTOTPType{Code: "123456"}},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{
						ID:    "user-1",
						State: domain.UserStateActive,
						Human: &domain.HumanUser{
							TOTP: &domain.HumanTOTP{
								VerifiedAt:     now,
								Secret:         &crypto.CryptoValue{Crypted: []byte("123456")},
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
					TOTP: &domain.HumanTOTP{
						VerifiedAt:     now,
						Secret:         &crypto.CryptoValue{Crypted: []byte("123456")},
						FailedAttempts: 0,
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

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

			err := tc.cmd.Validate(t.Context(), opts)
			assert.ErrorIs(t, err, tc.expectedError)

			if tc.expectedError == nil {
				assert.Equal(t, tc.expectedUser, tc.cmd.FetchedUser)
			}
		})
	}
}

func TestTOTPCheckCommand_Execute(t *testing.T) {
	t.Parallel()

	userUpdateErr := errors.New("user update error")
	sessionUpdateErr := errors.New("session update error")
	listErr := errors.New("list error")
	decryptErr := errors.New("decrypt error")

	tt := []struct {
		testName           string
		cmd                *domain.TOTPCheckCommand
		encryptionAlgo     func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm
		userRepo           func(ctrl *gomock.Controller) domain.HumanUserRepository
		sessionRepo        func(ctrl *gomock.Controller) domain.SessionRepository
		lockoutSettingRepo func(ctrl *gomock.Controller) domain.LockoutSettingsRepository
		expectedError      error
		expectedSuccess    bool
		expectedLocked     bool
	}{
		{
			testName:        "when checkTOTP is nil should return no error",
			cmd:             &domain.TOTPCheckCommand{},
			expectedError:   nil,
			expectedSuccess: false,
		},
		{
			testName: "when TOTP verification succeeds should update user and session",
			cmd: &domain.TOTPCheckCommand{
				CheckTOTP:  &domain.CheckTOTPType{Code: "123456"},
				InstanceID: "instance-1",
				SessionID:  "session-1",
				FetchedUser: domain.User{
					ID:             "user-1",
					OrganizationID: "org-1",
					Human: &domain.HumanUser{
						TOTP: &domain.HumanTOTP{
							Secret: &crypto.CryptoValue{Crypted: []byte("encrypted-secret")}, FailedAttempts: 0},
					},
				},
				ValidateFunc: func(_, _ string) bool { return true },
			},
			encryptionAlgo: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				mock := crypto.NewMockEncryptionAlgorithm(ctrl)
				mock.EXPECT().Algorithm().AnyTimes().Return("")
				mock.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{""})
				mock.EXPECT().DecryptString(gomock.Any(), gomock.Any()).AnyTimes().Return("clear txt", nil)
				return mock
			},
			userRepo: func(ctrl *gomock.Controller) domain.HumanUserRepository {
				repo := domainmock.NewHumanRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "user-1")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, repo.SetLastSuccessfulTOTPCheck(time.Now())).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "session-1")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, repo.SetFactor(&domain.SessionFactorTOTP{LastVerifiedAt: time.Now()})).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			expectedSuccess: true,
		},
		{
			testName: "when user update fails should return error",
			cmd: &domain.TOTPCheckCommand{
				CheckTOTP:  &domain.CheckTOTPType{Code: "123456"},
				InstanceID: "instance-1",
				SessionID:  "session-1",
				FetchedUser: domain.User{
					ID:             "user-1",
					OrganizationID: "org-1",
					Human: &domain.HumanUser{
						TOTP: &domain.HumanTOTP{
							Secret: &crypto.CryptoValue{Crypted: []byte("encrypted-secret")}, FailedAttempts: 0},
					},
				},
				ValidateFunc: func(_, _ string) bool { return true },
			},
			userRepo: func(ctrl *gomock.Controller) domain.HumanUserRepository {
				repo := domainmock.NewHumanRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "user-1")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, repo.SetLastSuccessfulTOTPCheck(time.Now())).
					Times(1).
					Return(int64(0), userUpdateErr)
				return repo
			},
			encryptionAlgo: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				mock := crypto.NewMockEncryptionAlgorithm(ctrl)
				mock.EXPECT().Algorithm().AnyTimes().Return("")
				mock.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{""})
				mock.EXPECT().DecryptString(gomock.Any(), gomock.Any()).AnyTimes().Return("clear txt", nil)
				return mock
			},
			expectedError: zerrors.ThrowInternal(userUpdateErr, "DOM-aoMAzO", "failed updating user"),
		},
		{
			testName: "when session update fails after successful TOTP should return error",
			cmd: &domain.TOTPCheckCommand{
				CheckTOTP:  &domain.CheckTOTPType{Code: "123456"},
				InstanceID: "instance-1",
				SessionID:  "session-1",
				FetchedUser: domain.User{
					ID:             "user-1",
					OrganizationID: "org-1",
					Human: &domain.HumanUser{
						TOTP: &domain.HumanTOTP{
							Secret: &crypto.CryptoValue{Crypted: []byte("encrypted-secret")}, FailedAttempts: 0},
					},
				},
				ValidateFunc: func(_, _ string) bool { return true },
			},
			userRepo: func(ctrl *gomock.Controller) domain.HumanUserRepository {
				repo := domainmock.NewHumanRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "user-1")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, repo.SetLastSuccessfulTOTPCheck(time.Now())).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "session-1")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, repo.SetFactor(&domain.SessionFactorTOTP{LastVerifiedAt: time.Now()})).
					Times(1).
					Return(int64(0), sessionUpdateErr)
				return repo
			},
			encryptionAlgo: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				mock := crypto.NewMockEncryptionAlgorithm(ctrl)
				mock.EXPECT().Algorithm().AnyTimes().Return("")
				mock.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{""})
				mock.EXPECT().DecryptString(gomock.Any(), gomock.Any()).AnyTimes().Return("clear txt", nil)
				return mock
			},
			expectedError: zerrors.ThrowInternal(sessionUpdateErr, "DOM-ymhCTD", "failed updating session"),
		},
		{
			testName: "when lockout policy fetch fails should return error",
			cmd: &domain.TOTPCheckCommand{
				CheckTOTP:  &domain.CheckTOTPType{Code: "wrong-code"},
				InstanceID: "instance-1",
				SessionID:  "session-1",
				FetchedUser: domain.User{
					ID:             "user-1",
					OrganizationID: "org-1",
					Human: &domain.HumanUser{
						TOTP: &domain.HumanTOTP{
							Secret: &crypto.CryptoValue{Crypted: []byte("encrypted-secret")}, FailedAttempts: 0},
					},
				},
				ValidateFunc: func(_, _ string) bool { return false },
			},
			userRepo: func(ctrl *gomock.Controller) domain.HumanUserRepository {
				repo := domainmock.NewHumanRepo(ctrl)
				return repo
			},
			lockoutSettingRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				repo := domainmock.NewLockoutSettingsRepo(ctrl)
				instanceAndOrg := database.And(repo.InstanceIDCondition("instance-1"), repo.OrganizationIDCondition(gu.Ptr("org-1")))
				orgNullOrEmpty := database.Or(repo.OrganizationIDCondition(nil), repo.OrganizationIDCondition(gu.Ptr("")))
				onlyInstance := database.And(repo.InstanceIDCondition("instance-1"), orgNullOrEmpty)
				conds := database.WithCondition(database.Or(instanceAndOrg, onlyInstance))

				repo.EXPECT().
					List(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(conds),
						dbmock.QueryOptions(database.WithOrderByAscending(repo.OrganizationIDColumn(), repo.InstanceIDColumn())),
						dbmock.QueryOptions(database.WithLimit(1)),
					).Times(1).
					Return(nil, listErr)
				return repo
			},
			encryptionAlgo: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				mock := crypto.NewMockEncryptionAlgorithm(ctrl)
				mock.EXPECT().Algorithm().AnyTimes().Return("")
				mock.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{""})
				mock.EXPECT().DecryptString(gomock.Any(), gomock.Any()).AnyTimes().Return("clear txt", nil)
				return mock
			},
			expectedError: zerrors.ThrowInternal(listErr, "DOM-3B8Z6s", "failed fetching lockout settings"),
		},
		{
			testName: "when TOTP verification fails should update user and fail",
			cmd: &domain.TOTPCheckCommand{
				CheckTOTP:  &domain.CheckTOTPType{Code: "wrong-code"},
				InstanceID: "instance-1",
				SessionID:  "session-1",
				FetchedUser: domain.User{
					ID:             "user-1",
					OrganizationID: "org-1",
					Human: &domain.HumanUser{
						TOTP: &domain.HumanTOTP{
							Secret: &crypto.CryptoValue{Crypted: []byte("encrypted-secret")}, FailedAttempts: 0},
					},
				},
				ValidateFunc: func(_, _ string) bool { return false },
			},
			encryptionAlgo: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				mock := crypto.NewMockEncryptionAlgorithm(ctrl)
				mock.EXPECT().Algorithm().AnyTimes().Return("")
				mock.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{""})
				mock.EXPECT().DecryptString(gomock.Any(), gomock.Any()).AnyTimes().Return("", decryptErr)
				return mock
			},
			lockoutSettingRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				repo := domainmock.NewLockoutSettingsRepo(ctrl)
				instanceAndOrg := database.And(repo.InstanceIDCondition("instance-1"), repo.OrganizationIDCondition(gu.Ptr("org-1")))
				orgNullOrEmpty := database.Or(repo.OrganizationIDCondition(nil), repo.OrganizationIDCondition(gu.Ptr("")))
				onlyInstance := database.And(repo.InstanceIDCondition("instance-1"), orgNullOrEmpty)
				conds := database.WithCondition(database.Or(instanceAndOrg, onlyInstance))

				repo.EXPECT().
					List(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(conds),
						dbmock.QueryOptions(database.WithOrderByAscending(repo.OrganizationIDColumn(), repo.InstanceIDColumn())),
						dbmock.QueryOptions(database.WithLimit(1)),
					).Times(1).
					Return([]*domain.LockoutSettings{
						{LockoutSettingsAttributes: domain.LockoutSettingsAttributes{MaxOTPAttempts: gu.Ptr(uint64(5))}},
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.HumanUserRepository {
				repo := domainmock.NewHumanRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "user-1")
				changes := database.Changes{
					repo.IncrementTOTPFailedAttempts(),
				}
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						idCondition,
						changes,
					).
					Times(1).
					Return(int64(0), userUpdateErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(userUpdateErr, "DOM-lQLpIa", "failed updating user"),
		},
		{
			testName: "when TOTP verification fails should update user with failed check and fail on session update",
			cmd: &domain.TOTPCheckCommand{
				CheckTOTP:  &domain.CheckTOTPType{Code: "wrong-code"},
				InstanceID: "instance-1",
				SessionID:  "session-1",
				FetchedUser: domain.User{
					ID:             "user-1",
					OrganizationID: "org-1",
					Human: &domain.HumanUser{
						TOTP: &domain.HumanTOTP{
							Secret: &crypto.CryptoValue{Crypted: []byte("encrypted-secret")}, FailedAttempts: 0},
					},
				},
				ValidateFunc: func(_, _ string) bool { return false },
			},
			encryptionAlgo: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				mock := crypto.NewMockEncryptionAlgorithm(ctrl)
				mock.EXPECT().Algorithm().AnyTimes().Return("")
				mock.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{""})
				mock.EXPECT().DecryptString(gomock.Any(), gomock.Any()).AnyTimes().Return("", decryptErr)
				return mock
			},
			lockoutSettingRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				repo := domainmock.NewLockoutSettingsRepo(ctrl)
				instanceAndOrg := database.And(repo.InstanceIDCondition("instance-1"), repo.OrganizationIDCondition(gu.Ptr("org-1")))
				orgNullOrEmpty := database.Or(repo.OrganizationIDCondition(nil), repo.OrganizationIDCondition(gu.Ptr("")))
				onlyInstance := database.And(repo.InstanceIDCondition("instance-1"), orgNullOrEmpty)
				conds := database.WithCondition(database.Or(instanceAndOrg, onlyInstance))

				repo.EXPECT().
					List(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(conds),
						dbmock.QueryOptions(database.WithOrderByAscending(repo.OrganizationIDColumn(), repo.InstanceIDColumn())),
						dbmock.QueryOptions(database.WithLimit(1)),
					).Times(1).
					Return([]*domain.LockoutSettings{
						{
							Settings: domain.Settings{},
							LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
								MaxOTPAttempts: gu.Ptr(uint64(5)),
							},
						},
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.HumanUserRepository {
				repo := domainmock.NewHumanRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "user-1")
				changes := database.Changes{
					repo.IncrementTOTPFailedAttempts(),
				}
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						idCondition,
						changes,
					).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "session-1")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, repo.SetFactor(&domain.SessionFactorTOTP{LastVerifiedAt: time.Now()})).
					Times(1).
					Return(int64(0), sessionUpdateErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(sessionUpdateErr, "DOM-rSa1yU", "failed updating session"),
		},
		{
			testName: "when TOTP verification fails should update user with failed check",
			cmd: &domain.TOTPCheckCommand{
				CheckTOTP:  &domain.CheckTOTPType{Code: "wrong-code"},
				InstanceID: "instance-1",
				SessionID:  "session-1",
				FetchedUser: domain.User{
					ID:             "user-1",
					OrganizationID: "org-1",
					Human: &domain.HumanUser{
						TOTP: &domain.HumanTOTP{
							Secret: &crypto.CryptoValue{Crypted: []byte("encrypted-secret")}, FailedAttempts: 0},
					},
				},
				TarpitFunc:   func(_ uint64) {},
				ValidateFunc: func(_, _ string) bool { return false },
			},
			userRepo: func(ctrl *gomock.Controller) domain.HumanUserRepository {
				repo := domainmock.NewHumanRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "user-1")
				changes := database.Changes{
					repo.IncrementTOTPFailedAttempts(),
				}
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						idCondition,
						changes,
					).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "session-1")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, repo.SetFactor(&domain.SessionFactorTOTP{LastVerifiedAt: time.Now()})).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			lockoutSettingRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				repo := domainmock.NewLockoutSettingsRepo(ctrl)
				instanceAndOrg := database.And(repo.InstanceIDCondition("instance-1"), repo.OrganizationIDCondition(gu.Ptr("org-1")))
				orgNullOrEmpty := database.Or(repo.OrganizationIDCondition(nil), repo.OrganizationIDCondition(gu.Ptr("")))
				onlyInstance := database.And(repo.InstanceIDCondition("instance-1"), orgNullOrEmpty)
				conds := database.WithCondition(database.Or(instanceAndOrg, onlyInstance))

				repo.EXPECT().
					List(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(conds),
						dbmock.QueryOptions(database.WithOrderByAscending(repo.OrganizationIDColumn(), repo.InstanceIDColumn())),
						dbmock.QueryOptions(database.WithLimit(1)),
					).Times(1).
					Return([]*domain.LockoutSettings{
						{
							Settings: domain.Settings{},
							LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
								MaxOTPAttempts: gu.Ptr(uint64(5)),
							},
						},
					}, nil)
				return repo
			},
			encryptionAlgo: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				mock := crypto.NewMockEncryptionAlgorithm(ctrl)
				mock.EXPECT().Algorithm().AnyTimes().Return("")
				mock.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{""})
				mock.EXPECT().DecryptString(gomock.Any(), gomock.Any()).AnyTimes().Return("", decryptErr)
				return mock
			},
			expectedError: zerrors.ThrowInternal(decryptErr, "DOM-Yqhggx", "Errors.TOTP.FailedToDecryptSecret"),
		},
		{
			testName: "when TOTP verification fails and user exceeds max attempts should lock user",
			cmd: &domain.TOTPCheckCommand{
				CheckTOTP:  &domain.CheckTOTPType{Code: "wrong-code"},
				InstanceID: "instance-1",
				SessionID:  "session-1",
				FetchedUser: domain.User{
					ID:             "user-1",
					OrganizationID: "org-1",
					Human: &domain.HumanUser{
						TOTP: &domain.HumanTOTP{
							Secret: &crypto.CryptoValue{Crypted: []byte("encrypted-secret")}, FailedAttempts: 4},
					},
				},
				TarpitFunc:   func(_ uint64) {},
				ValidateFunc: func(_, _ string) bool { return false },
			},
			userRepo: func(ctrl *gomock.Controller) domain.HumanUserRepository {
				repo := domainmock.NewHumanRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "user-1")
				changes := database.Changes{
					repo.IncrementTOTPFailedAttempts(),
					repo.SetState(domain.UserStateLocked),
				}
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						idCondition,
						changes,
					).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.PrimaryKeyCondition("instance-1", "session-1")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, repo.SetFactor(&domain.SessionFactorTOTP{LastVerifiedAt: time.Now()})).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			lockoutSettingRepo: func(ctrl *gomock.Controller) domain.LockoutSettingsRepository {
				repo := domainmock.NewLockoutSettingsRepo(ctrl)
				instanceAndOrg := database.And(repo.InstanceIDCondition("instance-1"), repo.OrganizationIDCondition(gu.Ptr("org-1")))
				orgNullOrEmpty := database.Or(repo.OrganizationIDCondition(nil), repo.OrganizationIDCondition(gu.Ptr("")))
				onlyInstance := database.And(repo.InstanceIDCondition("instance-1"), orgNullOrEmpty)
				conds := database.WithCondition(database.Or(instanceAndOrg, onlyInstance))

				repo.EXPECT().
					List(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(conds),
						dbmock.QueryOptions(database.WithOrderByAscending(repo.OrganizationIDColumn(), repo.InstanceIDColumn())),
						dbmock.QueryOptions(database.WithLimit(1)),
					).Times(1).
					Return([]*domain.LockoutSettings{
						{
							Settings: domain.Settings{},
							LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
								MaxOTPAttempts: gu.Ptr(uint64(5)),
							},
						},
					}, nil)
				return repo
			},
			encryptionAlgo: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				mock := crypto.NewMockEncryptionAlgorithm(ctrl)
				mock.EXPECT().Algorithm().AnyTimes().Return("")
				mock.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{""})
				mock.EXPECT().DecryptString(gomock.Any(), gomock.Any()).AnyTimes().Return("clear txt", nil)
				return mock
			},
			expectedError:  zerrors.ThrowInvalidArgument(nil, "DOM-o5cVir", "Errors.User.MFA.OTP.InvalidCode"),
			expectedLocked: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			var encAlgo crypto.EncryptionAlgorithm
			if tc.encryptionAlgo != nil {
				encAlgo = tc.encryptionAlgo(ctrl)
				tc.cmd.EncryptionAlgorithm = encAlgo
			}

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)

			if tc.userRepo != nil {
				userRepo := domainmock.NewUserRepo(ctrl)
				humanRepo := tc.userRepo(ctrl)
				userRepo.EXPECT().Human().Times(1).Return(humanRepo)
				domain.WithUserRepo(userRepo)(opts)
			}
			if tc.sessionRepo != nil {
				domain.WithSessionRepo(tc.sessionRepo(ctrl))(opts)
			}
			if tc.lockoutSettingRepo != nil {
				domain.WithLockoutSettingsRepo(tc.lockoutSettingRepo(ctrl))(opts)
			}

			err := tc.cmd.Execute(t.Context(), opts)

			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedSuccess, tc.cmd.IsCheckSuccessful)
			assert.Equal(t, tc.expectedLocked, tc.cmd.IsUserLocked)
		})
	}
}

func TestTOTPCheckCommand_Events(t *testing.T) {
	t.Parallel()

	sessionAgg := session.NewAggregate("session-1", "instance-1").Aggregate
	userAgg := user.NewAggregate("user-1", "org-1").Aggregate

	tt := []struct {
		testName       string
		cmd            *domain.TOTPCheckCommand
		expectedEvents []eventstore.Command
	}{
		{
			testName:       "when checkTOTP is nil should return no events",
			cmd:            &domain.TOTPCheckCommand{},
			expectedEvents: []eventstore.Command{},
		},
		{
			testName: "when check is successful should emit user succeeded and session totp checked events",
			cmd: &domain.TOTPCheckCommand{
				CheckTOTP:         &domain.CheckTOTPType{},
				SessionID:         "session-1",
				InstanceID:        "instance-1",
				FetchedUser:       domain.User{ID: "user-1", OrganizationID: "org-1"},
				IsCheckSuccessful: true,
				CheckedAt:         time.Now(),
			},

			expectedEvents: []eventstore.Command{
				user.NewHumanOTPCheckSucceededEvent(t.Context(), &userAgg, nil),
				session.NewTOTPCheckedEvent(t.Context(), &sessionAgg, time.Now()),
			},
		},
		{
			testName: "when check is unsuccessful should emit user failed event",
			cmd: &domain.TOTPCheckCommand{
				CheckTOTP:   &domain.CheckTOTPType{},
				SessionID:   "session-1",
				InstanceID:  "instance-1",
				FetchedUser: domain.User{ID: "user-1", OrganizationID: "org-1"},
				CheckedAt:   time.Now(),
			},
			expectedEvents: []eventstore.Command{
				user.NewHumanOTPCheckFailedEvent(t.Context(), &userAgg, nil),
			},
		},
		{
			testName: "when check is unsuccessful and user is locked should emit user failed and user locked events",
			cmd: &domain.TOTPCheckCommand{
				CheckTOTP:    &domain.CheckTOTPType{},
				SessionID:    "session-1",
				InstanceID:   "instance-1",
				FetchedUser:  domain.User{ID: "user-1", OrganizationID: "org-1"},
				IsUserLocked: true,
				CheckedAt:    time.Now(),
			},
			expectedEvents: []eventstore.Command{
				user.NewHumanOTPCheckFailedEvent(t.Context(), &userAgg, nil),
				user.NewUserLockedEvent(t.Context(), &userAgg),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")

			// Test
			events, err := tc.cmd.Events(ctx, &domain.InvokeOpts{})

			// Verify
			assert.NoError(t, err)
			assert.Len(t, events, len(tc.expectedEvents))
			for i, expectedType := range tc.expectedEvents {
				assert.IsType(t, expectedType, events[i])
				switch expectedType.(type) {
				case *session.TOTPCheckedEvent:
					actualAssertedType, ok := events[i].(*session.TOTPCheckedEvent)
					require.True(t, ok)
					assert.Equal(t, tc.cmd.CheckedAt, actualAssertedType.CheckedAt)
				case *user.HumanOTPCheckSucceededEvent:
					_, ok := events[i].(*user.HumanOTPCheckSucceededEvent)
					require.True(t, ok)
				case *user.HumanOTPCheckFailedEvent:
					_, ok := events[i].(*user.HumanOTPCheckFailedEvent)
					require.True(t, ok)
				}
			}
		})
	}
}
