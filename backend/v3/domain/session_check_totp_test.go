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
	cryptomock "github.com/zitadel/zitadel/internal/crypto/mock"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestTOTPCheckCommand_Validate(t *testing.T) {
	t.Parallel()
	sessionGetErr := errors.New("session get error")
	userGetErr := errors.New("user get error")
	notFoundErr := database.NewNoRowFoundError(nil)

	tt := []struct {
		testName      string
		sessionRepo   func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo      func(ctrl *gomock.Controller) domain.UserRepository
		checkTOTP     *session_grpc.CheckTOTP
		sessionID     string
		instanceID    string
		expectedError error
		expectedUser  domain.User
	}{
		{
			testName:      "when checkTOTP is nil should return no error",
			checkTOTP:     nil,
			expectedError: nil,
		},
		{
			testName:      "when session ID is not set should return error",
			checkTOTP:     &session_grpc.CheckTOTP{Code: "123456"},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-HhGveR", "Errors.Missing.SessionID"),
		},
		{
			testName:      "when instance ID is not set should return error",
			checkTOTP:     &session_grpc.CheckTOTP{Code: "123456"},
			sessionID:     "session-1",
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-89Zi8Z", "Errors.Missing.InstanceID"),
		},
		{
			testName:   "when retrieving session fails should return error",
			checkTOTP:  &session_grpc.CheckTOTP{Code: "123456"},
			sessionID:  "session-1",
			instanceID: "instance-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).Return(nil, sessionGetErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(sessionGetErr, "DOM-e4OuhO", "Errors.Get.session"),
		},
		{
			testName:   "when session not found should return not found error",
			checkTOTP:  &session_grpc.CheckTOTP{Code: "123456"},
			sessionID:  "session-1",
			instanceID: "instance-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).Return(nil, notFoundErr)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(notFoundErr, "DOM-e4OuhO", "Errors.NotFound.session"),
		},
		{
			testName:   "when session userID is empty should return precondition failed error",
			checkTOTP:  &session_grpc.CheckTOTP{Code: "123456"},
			sessionID:  "session-1",
			instanceID: "instance-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).Return(&domain.Session{}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-hord0Z", "Errors.User.UserIDMissing"),
		},
		{
			testName:   "when retrieving user fails should return error",
			checkTOTP:  &session_grpc.CheckTOTP{Code: "123456"},
			sessionID:  "session-1",
			instanceID: "instance-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).Return(&domain.Session{UserID: "user-1"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).Return(nil, userGetErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(userGetErr, "DOM-PZvWq0", "Errors.Get.user"),
		},
		{
			testName:   "when user not found should return not found error",
			checkTOTP:  &session_grpc.CheckTOTP{Code: "123456"},
			sessionID:  "session-1",
			instanceID: "instance-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).Return(&domain.Session{UserID: "user-1"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).Return(nil, notFoundErr)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(notFoundErr, "DOM-PZvWq0", "Errors.NotFound.user"),
		},
		{
			testName:   "when user is not human should return precondition failed error",
			checkTOTP:  &session_grpc.CheckTOTP{Code: "123456"},
			sessionID:  "session-1",
			instanceID: "instance-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).Return(&domain.Session{UserID: "user-1"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).Return(&domain.User{}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-zzv1MO", "user not human"),
		},
		{
			testName:   "when user is locked should return precondition failed error",
			checkTOTP:  &session_grpc.CheckTOTP{Code: "123456"},
			sessionID:  "session-1",
			instanceID: "instance-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).Return(&domain.Session{UserID: "user-1"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).Return(&domain.User{
					State: domain.UserStateLocked,
					Human: &domain.HumanUser{},
				}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-gM4SUh", "Errors.User.Locked"),
		},
		{
			testName:   "when all validations pass should return no error and set user",
			checkTOTP:  &session_grpc.CheckTOTP{Code: "123456"},
			sessionID:  "session-1",
			instanceID: "instance-1",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).Return(&domain.Session{UserID: "user-1"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).Return(&domain.User{
					ID:    "user-1",
					State: domain.UserStateActive,
					Human: &domain.HumanUser{
						TOTP: domain.HumanTOTP{
							Check: &domain.Check{
								Code:           &crypto.CryptoValue{Crypted: []byte("123456")},
								FailedAttempts: 0,
							},
						},
					},
				}, nil)
				return repo
			},
			expectedUser: domain.User{
				ID:    "user-1",
				State: domain.UserStateActive,
				Human: &domain.HumanUser{
					TOTP: domain.HumanTOTP{Check: &domain.Check{
						Code:           &crypto.CryptoValue{Crypted: []byte("123456")},
						FailedAttempts: 0,
					},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cmd := domain.NewTOTPCheckCommand(tc.sessionID, tc.instanceID, nil, nil, nil, tc.checkTOTP)
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

			err := cmd.Validate(t.Context(), opts)
			assert.Equal(t, tc.expectedError, err)

			if tc.expectedError == nil {
				assert.Equal(t, tc.expectedUser, cmd.FetchedUser)
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
		checkTOTP          *session_grpc.CheckTOTP
		fetchedUser        domain.User
		userRepo           func(ctrl *gomock.Controller) domain.UserRepository
		sessionRepo        func(ctrl *gomock.Controller) domain.SessionRepository
		lockoutSettingRepo func(ctrl *gomock.Controller) domain.LockoutSettingsRepository
		tarpitFunc         func(uint64)
		validateFunc       func(string, string) bool
		encryptionAlgo     func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm
		expectedError      error
		expectedSuccess    bool
		expectedLocked     bool
	}{
		{
			testName:        "when checkTOTP is nil should return no error",
			checkTOTP:       nil,
			expectedError:   nil,
			expectedSuccess: false,
		},
		{
			testName:  "when TOTP verification succeeds should update user and session",
			checkTOTP: &session_grpc.CheckTOTP{Code: "123456"},
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					TOTP: domain.HumanTOTP{Check: &domain.Check{
						Code: &crypto.CryptoValue{Crypted: []byte("encrypted-secret")}, FailedAttempts: 0},
					},
				},
			},
			validateFunc: func(_, _ string) bool { return true },
			encryptionAlgo: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				mock := cryptomock.NewMockEncryptionAlgorithm(ctrl)
				mock.EXPECT().Algorithm().AnyTimes().Return("")
				mock.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{""})
				mock.EXPECT().DecryptString(gomock.Any(), gomock.Any()).AnyTimes().Return("clear txt", nil)
				return mock
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().Human().Return(humanRepo).AnyTimes()
				idCondition := getUserIDCondition(repo, "user-1")
				humanRepo.EXPECT().CheckTOTP(gomock.Any()).Times(1)
				repo.EXPECT().Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).Return(int64(1), nil)
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().SetFactor(gomock.Any()).Return(nil).AnyTimes()
				repo.EXPECT().Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).Return(int64(1), nil).AnyTimes()
				return repo
			},
			expectedSuccess: true,
		},
		{
			testName:  "when user update fails should return error",
			checkTOTP: &session_grpc.CheckTOTP{Code: "123456"},
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					TOTP: domain.HumanTOTP{Check: &domain.Check{
						Code: &crypto.CryptoValue{Crypted: []byte("encrypted-secret")}, FailedAttempts: 0},
					},
				},
			},
			validateFunc: func(_, _ string) bool { return true },
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().Human().Return(humanRepo).AnyTimes()
				humanRepo.EXPECT().CheckTOTP(gomock.Any()).Times(1)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).Return(int64(0), userUpdateErr)
				return repo
			},
			encryptionAlgo: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				mock := cryptomock.NewMockEncryptionAlgorithm(ctrl)
				mock.EXPECT().Algorithm().AnyTimes().Return("")
				mock.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{""})
				mock.EXPECT().DecryptString(gomock.Any(), gomock.Any()).AnyTimes().Return("clear txt", nil)
				return mock
			},
			expectedError: zerrors.ThrowInternal(userUpdateErr, "DOM-aoMAzO", "Errors.Update.user"),
		},
		{
			testName:     "when session update fails after successful TOTP should return error",
			checkTOTP:    &session_grpc.CheckTOTP{Code: "123456"},
			validateFunc: func(_, _ string) bool { return true },
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					TOTP: domain.HumanTOTP{Check: &domain.Check{
						Code: &crypto.CryptoValue{Crypted: []byte("encrypted-secret")}, FailedAttempts: 0},
					},
				},
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().Human().Return(humanRepo).AnyTimes()
				idCondition := getUserIDCondition(repo, "user-1")
				humanRepo.EXPECT().CheckTOTP(gomock.Any()).Times(1)
				repo.EXPECT().Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).Return(int64(1), nil)
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().SetFactor(gomock.Any()).Return(nil).AnyTimes()
				repo.EXPECT().Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).Return(int64(0), sessionUpdateErr)
				return repo
			},
			encryptionAlgo: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				mock := cryptomock.NewMockEncryptionAlgorithm(ctrl)
				mock.EXPECT().Algorithm().AnyTimes().Return("")
				mock.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{""})
				mock.EXPECT().DecryptString(gomock.Any(), gomock.Any()).AnyTimes().Return("clear txt", nil)
				return mock
			},
			expectedError: zerrors.ThrowInternal(sessionUpdateErr, "DOM-ymhCTD", "Errors.Update.session"),
		},
		{
			testName:  "when lockout policy fetch fails should return error",
			checkTOTP: &session_grpc.CheckTOTP{Code: "wrong-code"},
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					TOTP: domain.HumanTOTP{Check: &domain.Check{
						Code: &crypto.CryptoValue{Crypted: []byte("encrypted-secret")}, FailedAttempts: 0},
					},
				},
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				humanRepo.EXPECT().CheckTOTP(gomock.Any()).Times(1)
				repo.EXPECT().Human().Return(humanRepo).AnyTimes()
				getUserIDCondition(repo, "user-1")
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
			tarpitFunc: func(attempts uint64) {},
			encryptionAlgo: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				mock := cryptomock.NewMockEncryptionAlgorithm(ctrl)
				mock.EXPECT().Algorithm().AnyTimes().Return("")
				mock.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{""})
				mock.EXPECT().DecryptString(gomock.Any(), gomock.Any()).AnyTimes().Return("clear txt", nil)
				return mock
			},
			expectedError: zerrors.ThrowInternal(listErr, "DOM-3B8Z6s", "failed fetching lockout settings"),
		},
		{
			testName:  "when TOTP verification fails should update user and fail",
			checkTOTP: &session_grpc.CheckTOTP{Code: "wrong-code"},
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					TOTP: domain.HumanTOTP{Check: &domain.Check{
						Code: &crypto.CryptoValue{Crypted: []byte("encrypted-secret")}, FailedAttempts: 0},
					},
				},
			},
			encryptionAlgo: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				mock := cryptomock.NewMockEncryptionAlgorithm(ctrl)
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
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().Human().Return(humanRepo).AnyTimes()
				idCondition := getUserIDCondition(repo, "user-1")
				humanRepo.EXPECT().CheckTOTP(gomock.Any()).Times(1)
				repo.EXPECT().Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).Return(int64(0), userUpdateErr).AnyTimes()
				return repo
			},
			expectedError: zerrors.ThrowInternal(userUpdateErr, "DOM-lQLpIa", "Errors.Update.user"),
		},
		{
			testName:  "when TOTP verification fails should update user with failed check and fail on session update",
			checkTOTP: &session_grpc.CheckTOTP{Code: "wrong-code"},
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					TOTP: domain.HumanTOTP{Check: &domain.Check{
						Code: &crypto.CryptoValue{Crypted: []byte("encrypted-secret")}, FailedAttempts: 0},
					},
				},
			},
			encryptionAlgo: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				mock := cryptomock.NewMockEncryptionAlgorithm(ctrl)
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
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().Human().Return(humanRepo).AnyTimes()
				idCondition := getUserIDCondition(repo, "user-1")
				humanRepo.EXPECT().CheckTOTP(gomock.Any()).Times(1)
				repo.EXPECT().Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).Return(int64(1), nil).AnyTimes()
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().SetFactor(gomock.Any()).Return(nil).AnyTimes()
				repo.EXPECT().Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).Return(int64(0), sessionUpdateErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(sessionUpdateErr, "DOM-rSa1yU", "Errors.Update.session"),
		},
		{
			testName:  "when TOTP verification fails should update user with failed check",
			checkTOTP: &session_grpc.CheckTOTP{Code: "wrong-code"},
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					TOTP: domain.HumanTOTP{Check: &domain.Check{
						Code: &crypto.CryptoValue{Crypted: []byte("encrypted-secret")}, FailedAttempts: 0},
					},
				},
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().Human().Return(humanRepo).AnyTimes()
				idCondition := getUserIDCondition(repo, "user-1")
				humanRepo.EXPECT().CheckTOTP(gomock.Any()).Times(1)
				repo.EXPECT().Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).Return(int64(1), nil).AnyTimes()
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().SetFactor(gomock.Any()).Return(nil).AnyTimes()
				repo.EXPECT().Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).Return(int64(1), nil)
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
			tarpitFunc: func(attempts uint64) {},
			encryptionAlgo: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				mock := cryptomock.NewMockEncryptionAlgorithm(ctrl)
				mock.EXPECT().Algorithm().AnyTimes().Return("")
				mock.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{""})
				mock.EXPECT().DecryptString(gomock.Any(), gomock.Any()).AnyTimes().Return("", decryptErr)
				return mock
			},
			expectedError: zerrors.ThrowInternal(decryptErr, "DOM-Yqhggx", "failed decrypting TOTP secret"),
		},
		{
			testName:  "when TOTP verification fails and user exceeds max attempts should lock user",
			checkTOTP: &session_grpc.CheckTOTP{Code: "wrong-code"},
			fetchedUser: domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					TOTP: domain.HumanTOTP{Check: &domain.Check{
						Code:           &crypto.CryptoValue{Crypted: []byte("encrypted-secret")},
						FailedAttempts: 4,
					},
					},
				},
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().Human().Return(humanRepo).AnyTimes()
				idCondition := getUserIDCondition(repo, "user-1")
				humanRepo.EXPECT().CheckTOTP(gomock.Any()).Times(1)
				humanRepo.EXPECT().
					SetState(domain.UserStateLocked).
					Times(1)
				repo.EXPECT().Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).Return(int64(1), nil).AnyTimes()
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().SetFactor(gomock.Any()).Return(nil).AnyTimes()
				repo.EXPECT().Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).Return(int64(1), nil)
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
			tarpitFunc: func(attempts uint64) {},
			encryptionAlgo: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				mock := cryptomock.NewMockEncryptionAlgorithm(ctrl)
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var encAlgo crypto.EncryptionAlgorithm
			if tc.encryptionAlgo != nil {
				encAlgo = tc.encryptionAlgo(ctrl)
			}
			cmd := domain.NewTOTPCheckCommand("session-1", "instance-1", tc.tarpitFunc, tc.validateFunc, encAlgo, tc.checkTOTP)
			cmd.FetchedUser = tc.fetchedUser

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)

			if tc.userRepo != nil {
				domain.WithUserRepo(tc.userRepo(ctrl))(opts)
			}
			if tc.sessionRepo != nil {
				domain.WithSessionRepo(tc.sessionRepo(ctrl))(opts)
			}
			if tc.lockoutSettingRepo != nil {
				domain.WithLockoutSettingsRepo(tc.lockoutSettingRepo(ctrl))(opts)
			}

			err := cmd.Execute(t.Context(), opts)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedSuccess, cmd.IsCheckSuccessful)
			assert.Equal(t, tc.expectedLocked, cmd.IsUserLocked)
		})
	}
}

func TestTOTPCheckCommand_Events(t *testing.T) {
	t.Parallel()

	sessionAgg := session.NewAggregate("session-1", "instance-1").Aggregate
	userAgg := user.NewAggregate("user-1", "org-1").Aggregate

	tt := []struct {
		testName          string
		checkTOTP         *session_grpc.CheckTOTP
		sessionID         string
		instanceID        string
		fetchedUser       domain.User
		checkedAt         time.Time
		isCheckSuccessful bool
		isUserLocked      bool
		expectedEvents    []eventstore.Command
	}{
		{
			testName:       "when checkTOTP is nil should return no events",
			checkTOTP:      nil,
			expectedEvents: []eventstore.Command{},
		},
		{
			testName:          "when check is successful should emit user succeeded and session totp checked events",
			checkTOTP:         &session_grpc.CheckTOTP{},
			checkedAt:         time.Now(),
			isCheckSuccessful: true,

			fetchedUser: domain.User{ID: "user-1", OrganizationID: "org-1"},
			expectedEvents: []eventstore.Command{
				user.NewHumanOTPCheckSucceededEvent(t.Context(), &userAgg, nil),
				session.NewTOTPCheckedEvent(t.Context(), &sessionAgg, time.Now()),
			},
		},
		{
			testName:  "when check is unsuccessful should emit user failed and session totp checked events",
			checkTOTP: &session_grpc.CheckTOTP{},
			checkedAt: time.Now(),

			fetchedUser: domain.User{ID: "user-1", OrganizationID: "org-1"},
			expectedEvents: []eventstore.Command{
				user.NewHumanOTPCheckFailedEvent(t.Context(), &userAgg, nil),
				session.NewTOTPCheckedEvent(t.Context(), &sessionAgg, time.Now()),
			},
		},
		{
			testName:     "when check is unsuccessful and user is locked should emit user failed, user locked and session totp checked events",
			checkTOTP:    &session_grpc.CheckTOTP{},
			checkedAt:    time.Now(),
			isUserLocked: true,

			fetchedUser: domain.User{ID: "user-1", OrganizationID: "org-1"},
			expectedEvents: []eventstore.Command{
				user.NewHumanOTPCheckFailedEvent(t.Context(), &userAgg, nil),
				user.NewUserLockedEvent(t.Context(), &userAgg),
				session.NewTOTPCheckedEvent(t.Context(), &sessionAgg, time.Now()),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			cmd := domain.NewTOTPCheckCommand("session-1", "instance-1", nil, nil, nil, tc.checkTOTP)
			cmd.FetchedUser = tc.fetchedUser
			cmd.CheckedAt = tc.checkedAt
			cmd.IsCheckSuccessful = tc.isCheckSuccessful
			cmd.IsUserLocked = tc.isUserLocked

			// Test
			events, err := cmd.Events(ctx, &domain.InvokeOpts{})

			// Verify
			assert.NoError(t, err)
			assert.Len(t, events, len(tc.expectedEvents))
			for i, expectedType := range tc.expectedEvents {
				assert.IsType(t, expectedType, events[i])
				switch expectedAssertedType := expectedType.(type) {
				case *session.TOTPCheckedEvent:
					actualAssertedType, ok := events[i].(*session.TOTPCheckedEvent)
					require.True(t, ok)
					assert.InDelta(t, expectedAssertedType.CheckedAt.UnixMilli(), actualAssertedType.CheckedAt.UnixMilli(), 1.5)
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
