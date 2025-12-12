package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestOTPCheckCommand_Validate(t *testing.T) {
	t.Parallel()
	sessionGetErr := errors.New("session get error")
	notFoundErr := database.NewNoRowFoundError(nil)
	userGetErr := errors.New("user get error")

	tt := []struct {
		testName      string
		sessionRepo   func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo      func(ctrl *gomock.Controller) domain.UserRepository
		checkOTP      *session_grpc.CheckOTP
		sessionID     string
		requestType   domain.OTPRequestType
		expectedError error
		expectedUser  domain.User
	}{
		{
			testName:      "when checkOTP is nil should return no error",
			checkOTP:      nil,
			expectedError: nil,
		},
		{
			testName:    "when retrieving session fails should return error",
			checkOTP:    &session_grpc.CheckOTP{Code: "123456"},
			sessionID:   "session-1",
			requestType: domain.OTPSMSRequestType,
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(nil, sessionGetErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(sessionGetErr, "DOM-eppPwQ", "failed fetching session"),
		},
		{
			testName:    "when session not found should return not found error",
			checkOTP:    &session_grpc.CheckOTP{Code: "123456"},
			sessionID:   "session-1",
			requestType: domain.OTPSMSRequestType,
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(nil, notFoundErr)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(notFoundErr, "DOM-eppPwQ", "session not found"),
		},
		{
			testName:    "when retrieving user fails should return error",
			checkOTP:    &session_grpc.CheckOTP{Code: "123456"},
			sessionID:   "session-1",
			requestType: domain.OTPSMSRequestType,
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(nil, userGetErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(userGetErr, "DOM-TxDSma", "failed fetching user"),
		},
		{
			testName:    "when user not found should return not found error",
			checkOTP:    &session_grpc.CheckOTP{Code: "123456"},
			sessionID:   "session-1",
			requestType: domain.OTPSMSRequestType,
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(nil, notFoundErr)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(notFoundErr, "DOM-TxDSma", "user not found"),
		},
		{
			testName:    "when user is not human should return precondition failed error",
			checkOTP:    &session_grpc.CheckOTP{Code: "123456"},
			sessionID:   "session-1",
			requestType: domain.OTPSMSRequestType,
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{Human: nil}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-pBmqRN", "user not human"),
		},
		{
			testName:    "when checkOTP code is empty should return invalid argument error",
			checkOTP:    &session_grpc.CheckOTP{Code: ""},
			sessionID:   "session-1",
			requestType: domain.OTPSMSRequestType,
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{Human: &domain.HumanUser{}}, nil)
				return repo
			},
			expectedError: zerrors.ThrowInvalidArgument(nil, "DOM-u7KQi4", "Errors.User.Code.Empty"),
		},
		{
			testName:    "when SMS OTP challenge is nil should return precondition failed error",
			checkOTP:    &session_grpc.CheckOTP{Code: "123456"},
			sessionID:   "session-1",
			requestType: domain.OTPSMSRequestType,
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{
						UserID:     "user-1",
						Challenges: domain.SessionChallenges{},
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{Human: &domain.HumanUser{}}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-UpslUc", "no OTP SMS challenge set"),
		},
		{
			testName:    "when SMS OTP user phone is nil should return precondition failed error",
			checkOTP:    &session_grpc.CheckOTP{Code: "123456"},
			sessionID:   "session-1",
			requestType: domain.OTPSMSRequestType,
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{
						UserID: "user-1",
						Challenges: domain.SessionChallenges{
							&domain.SessionChallengeOTPSMS{},
						},
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{Human: &domain.HumanUser{Phone: nil}}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-fzSWTO", "no phone set"),
		},
		{
			testName:    "when SMS OTP not enabled should return precondition failed error",
			checkOTP:    &session_grpc.CheckOTP{Code: "123456"},
			sessionID:   "session-1",
			requestType: domain.OTPSMSRequestType,
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{
						UserID: "user-1",
						Challenges: domain.SessionChallenges{
							&domain.SessionChallengeOTPSMS{},
						},
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{
						Human: &domain.HumanUser{
							Phone: &domain.HumanPhone{
								OTP: domain.OTP{EnabledAt: time.Time{}},
							},
						},
					}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-iJZ4jp", "Errors.User.MFA.OTP.NotReady"),
		},
		{
			testName:    "when SMS OTP code and generator ID both missing should return precondition failed error",
			checkOTP:    &session_grpc.CheckOTP{Code: "123456"},
			sessionID:   "session-1",
			requestType: domain.OTPSMSRequestType,
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{
						UserID: "user-1",
						Challenges: domain.SessionChallenges{
							&domain.SessionChallengeOTPSMS{
								Code:        nil,
								GeneratorID: "",
							},
						},
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{
						Human: &domain.HumanUser{
							Phone: &domain.HumanPhone{
								OTP: domain.OTP{EnabledAt: time.Now()},
							},
						},
					}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-tAK4Cc", "Errors.User.Code.NotFound"),
		},
		{
			testName:    "when email OTP challenge is nil should return precondition failed error",
			checkOTP:    &session_grpc.CheckOTP{Code: "123456"},
			sessionID:   "session-1",
			requestType: domain.OTPEmailRequestType,
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{
						UserID:     "user-1",
						Challenges: domain.SessionChallenges{},
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{Human: &domain.HumanUser{}}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-2DlM76", "no OTP Email challenge set"),
		},
		{
			testName:    "when email OTP not enabled should return precondition failed error",
			checkOTP:    &session_grpc.CheckOTP{Code: "123456"},
			sessionID:   "session-1",
			requestType: domain.OTPEmailRequestType,
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{
						UserID: "user-1",
						Challenges: domain.SessionChallenges{
							&domain.SessionChallengeOTPEmail{},
						},
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{
						Human: &domain.HumanUser{
							Email: domain.HumanEmail{
								OTP: domain.OTP{EnabledAt: time.Time{}},
							},
						},
					}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-2uf0SY", "Errors.User.MFA.OTP.NotReady"),
		},
		{
			testName:    "when email OTP code is nil should return precondition failed error",
			checkOTP:    &session_grpc.CheckOTP{Code: "123456"},
			sessionID:   "session-1",
			requestType: domain.OTPEmailRequestType,
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{
						UserID: "user-1",
						Challenges: domain.SessionChallenges{
							&domain.SessionChallengeOTPEmail{Code: nil},
						},
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{
						Human: &domain.HumanUser{
							Email: domain.HumanEmail{
								OTP: domain.OTP{EnabledAt: time.Now()},
							},
						},
					}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-RegOgD", "Errors.User.Code.NotFound"),
		},
		{
			testName:    "when SMS OTP validation passes should return no error and set user",
			checkOTP:    &session_grpc.CheckOTP{Code: "123456"},
			sessionID:   "session-1",
			requestType: domain.OTPSMSRequestType,
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{
						UserID: "user-1",
						Challenges: domain.SessionChallenges{
							&domain.SessionChallengeOTPSMS{
								Code:        &crypto.CryptoValue{},
								GeneratorID: "",
							},
						},
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{
						ID: "user-1",
						Human: &domain.HumanUser{
							Phone: &domain.HumanPhone{
								OTP: domain.OTP{EnabledAt: time.Now()},
							},
						},
					}, nil)
				return repo
			},
			expectedUser: domain.User{
				ID: "user-1",
				Human: &domain.HumanUser{
					Phone: &domain.HumanPhone{
						OTP: domain.OTP{EnabledAt: time.Now()},
					},
				},
			},
		},
		{
			testName:    "when email OTP validation passes should return no error and set user",
			checkOTP:    &session_grpc.CheckOTP{Code: "123456"},
			sessionID:   "session-1",
			requestType: domain.OTPEmailRequestType,
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{
						UserID: "user-1",
						Challenges: domain.SessionChallenges{
							&domain.SessionChallengeOTPEmail{
								Code: &crypto.CryptoValue{},
							},
						},
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{
						ID: "user-1",
						Human: &domain.HumanUser{
							Email: domain.HumanEmail{
								OTP: domain.OTP{EnabledAt: time.Now()},
							},
						},
					}, nil)
				return repo
			},
			expectedUser: domain.User{
				ID: "user-1",
				Human: &domain.HumanUser{
					Email: domain.HumanEmail{
						OTP: domain.OTP{EnabledAt: time.Now()},
					},
				},
			},
		},
		{
			testName:    "when request type unknown should return invalid argument error",
			checkOTP:    &session_grpc.CheckOTP{Code: "123456"},
			sessionID:   "session-1",
			requestType: 99,
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.Session{
						UserID: "user-1",
						Challenges: domain.SessionChallenges{
							&domain.SessionChallengeOTPEmail{
								Code: &crypto.CryptoValue{},
							},
						},
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(idCondition))).
					Return(&domain.User{
						ID: "user-1",
						Human: &domain.HumanUser{
							Email: domain.HumanEmail{
								OTP: domain.OTP{EnabledAt: time.Now()},
							},
						},
					}, nil)
				return repo
			},
			expectedError: zerrors.ThrowInvalidArgument(nil, "DOM-oeUlih", "invalid OTP request type"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

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

			cmd := domain.NewOTPCheckCommand(tc.sessionID, "instance-1", nil, nil, nil, nil, tc.checkOTP, tc.requestType)
			err := cmd.Validate(t.Context(), opts)

			assert.Equal(t, tc.expectedError, err)
			if tc.checkOTP != nil && tc.expectedError == nil {
				assert.Equal(t, tc.expectedUser.ID, cmd.FetchedUser.ID)
			}
		})
	}
}
func TestOTPCheckCommand_Events(t *testing.T) {
	t.Parallel()
	userAgg := user.NewAggregate("user-1", "org-1").Aggregate
	sessAgg := session.NewAggregate("session-1", "instance-1").Aggregate

	tt := []struct {
		testName           string
		checkOTP           *session_grpc.CheckOTP
		isSMSCheckSucc     bool
		isEmailCheckSucc   bool
		isUserLocked       bool
		requestType        domain.OTPRequestType
		expectedEventTypes []eventstore.Command
	}{
		{
			testName: "when checkOTP is nil should return no events",
			checkOTP: nil,
		},
		{
			testName:         "when SMS check succeeded should return user and session events",
			checkOTP:         &session_grpc.CheckOTP{Code: "123456"},
			isSMSCheckSucc:   true,
			isEmailCheckSucc: false,
			isUserLocked:     false,
			requestType:      domain.OTPSMSRequestType,
			expectedEventTypes: []eventstore.Command{
				user.NewHumanOTPSMSCheckSucceededEvent(t.Context(), &userAgg, nil),
				session.NewOTPSMSCheckedEvent(t.Context(), &sessAgg, time.Now()),
			},
		},
		{
			testName:         "when SMS check failed should return user and session events",
			checkOTP:         &session_grpc.CheckOTP{Code: "123456"},
			isSMSCheckSucc:   false,
			isEmailCheckSucc: false,
			isUserLocked:     false,
			requestType:      domain.OTPSMSRequestType,
			expectedEventTypes: []eventstore.Command{
				user.NewHumanOTPSMSCheckFailedEvent(t.Context(), &userAgg, nil),
				session.NewOTPSMSCheckedEvent(t.Context(), &sessAgg, time.Now()),
			},
		},
		{
			testName:         "when SMS check failed and user locked should return three events",
			checkOTP:         &session_grpc.CheckOTP{Code: "123456"},
			isSMSCheckSucc:   false,
			isEmailCheckSucc: false,
			isUserLocked:     true,
			requestType:      domain.OTPSMSRequestType,
			expectedEventTypes: []eventstore.Command{
				user.NewHumanOTPSMSCheckFailedEvent(t.Context(), &userAgg, nil),
				user.NewUserLockedEvent(t.Context(), &userAgg),
				session.NewOTPSMSCheckedEvent(t.Context(), &sessAgg, time.Now()),
			},
		},
		{
			testName:         "when Email check succeeded should return user and session events",
			checkOTP:         &session_grpc.CheckOTP{Code: "123456"},
			isSMSCheckSucc:   false,
			isEmailCheckSucc: true,
			isUserLocked:     false,
			requestType:      domain.OTPEmailRequestType,
			expectedEventTypes: []eventstore.Command{
				user.NewHumanOTPEmailCheckSucceededEvent(t.Context(), &userAgg, nil),
				session.NewOTPEmailCheckedEvent(t.Context(), &sessAgg, time.Now()),
			},
		},
		{
			testName:         "when Email check failed should return user and session events",
			checkOTP:         &session_grpc.CheckOTP{Code: "123456"},
			isSMSCheckSucc:   false,
			isEmailCheckSucc: false,
			isUserLocked:     false,
			requestType:      domain.OTPEmailRequestType,
			expectedEventTypes: []eventstore.Command{
				user.NewHumanOTPEmailCheckFailedEvent(t.Context(), &userAgg, nil),
				session.NewOTPEmailCheckedEvent(t.Context(), &sessAgg, time.Now()),
			},
		},
		{
			testName:         "when Email check failed and user locked should return three events",
			checkOTP:         &session_grpc.CheckOTP{Code: "123456"},
			isSMSCheckSucc:   false,
			isEmailCheckSucc: false,
			isUserLocked:     true,
			requestType:      domain.OTPEmailRequestType,
			expectedEventTypes: []eventstore.Command{
				user.NewHumanOTPEmailCheckFailedEvent(t.Context(), &userAgg, nil),
				user.NewUserLockedEvent(t.Context(), &userAgg),
				session.NewOTPEmailCheckedEvent(t.Context(), &sessAgg, time.Now()),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Given
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)

			cmd := domain.NewOTPCheckCommand("session-1", "instance-1", nil, nil, nil, nil, tc.checkOTP, tc.requestType)
			cmd.IsSMSCheckSucceeded = tc.isSMSCheckSucc
			cmd.IsEmailCheckSucceeded = tc.isEmailCheckSucc
			cmd.IsUserLocked = tc.isUserLocked
			cmd.FetchedUser = &domain.User{
				ID:             "user-1",
				OrganizationID: "org-1",
				Human: &domain.HumanUser{
					Phone: &domain.HumanPhone{},
					Email: domain.HumanEmail{},
				},
			}
			cmd.CheckTime = time.Now()

			// Test
			events, err := cmd.Events(t.Context(), opts)

			endTime := time.Now()

			// Verify
			assert.NoError(t, err)
			require.Len(t, events, len(tc.expectedEventTypes))
			for i, expectedType := range tc.expectedEventTypes {
				assert.IsType(t, expectedType, events[i])
				switch expectedAssertedType := expectedType.(type) {
				case *user.HumanOTPSMSCheckSucceededEvent:
					actualAssertedType, ok := events[i].(*user.HumanOTPSMSCheckSucceededEvent)
					require.True(t, ok)
					assert.Equal(t, expectedAssertedType.AuthRequestInfo, actualAssertedType.AuthRequestInfo)
				case *user.HumanOTPSMSCheckFailedEvent:
					actualAssertedType, ok := events[i].(*user.HumanOTPSMSCheckFailedEvent)
					require.True(t, ok)
					assert.Equal(t, expectedAssertedType.AuthRequestInfo, actualAssertedType.AuthRequestInfo)
				case *user.UserLockedEvent:
					continue
				case *user.HumanOTPEmailCheckSucceededEvent:
					actualAssertedType, ok := events[i].(*user.HumanOTPEmailCheckSucceededEvent)
					require.True(t, ok)
					assert.Equal(t, expectedAssertedType.AuthRequestInfo, actualAssertedType.AuthRequestInfo)
				case *user.HumanOTPEmailCheckFailedEvent:
					actualAssertedType, ok := events[i].(*user.HumanOTPEmailCheckFailedEvent)
					require.True(t, ok)
					assert.Equal(t, expectedAssertedType.AuthRequestInfo, actualAssertedType.AuthRequestInfo)
				case *session.OTPSMSCheckedEvent:
					actualAssertedType, ok := events[i].(*session.OTPSMSCheckedEvent)
					require.True(t, ok)
					assert.WithinRange(t, actualAssertedType.CheckedAt, expectedAssertedType.CheckedAt, endTime)
				case *session.OTPEmailCheckedEvent:
					actualAssertedType, ok := events[i].(*session.OTPEmailCheckedEvent)
					require.True(t, ok)
					assert.WithinRange(t, actualAssertedType.CheckedAt, expectedAssertedType.CheckedAt, endTime)
				}
			}
		})
	}
}

// TODO(IAM-Marco): Write tests for Execute()
