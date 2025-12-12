package domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/api/authz"
	old_domain "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestPasskeyCheckCommand_Validate(t *testing.T) {
	t.Parallel()
	getErr := errors.New("get error")

	tt := []struct {
		testName      string
		sessionRepo   func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo      func(ctrl *gomock.Controller) domain.UserRepository
		checkPasskey  *session_grpc.CheckWebAuthN
		expectedError error
	}{
		{
			testName:      "when checkPasskey is nil should return no error",
			checkPasskey:  nil,
			expectedError: nil,
		},
		{
			testName: "when retrieving session fails should return error",
			checkPasskey: &session_grpc.CheckWebAuthN{
				CredentialAssertionData: &structpb.Struct{},
			},
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
					Return(nil, getErr)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)

				repo.EXPECT().
					LoadPasskeys().
					Times(1).
					Return(repo)
				return repo
			},
			expectedError: zerrors.ThrowInternal(getErr, "DOM-CUnePh", "failed fetching session"),
		},
		{
			testName: "when session has no passkey challenge should return precondition failed error",
			checkPasskey: &session_grpc.CheckWebAuthN{
				CredentialAssertionData: &structpb.Struct{},
			},
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
						ID:         "session-1",
						Challenges: domain.SessionChallenges{},
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)

				repo.EXPECT().
					LoadPasskeys().
					Times(1).
					Return(repo)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-lQhNR4", "Errors.Session.WebAuthN.NoChallenge"),
		},
		{
			testName: "when session has no user ID should return precondition failed error",
			checkPasskey: &session_grpc.CheckWebAuthN{
				CredentialAssertionData: &structpb.Struct{},
			},
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
						Challenges: domain.SessionChallenges{
							&domain.SessionChallengePasskey{
								LastChallengedAt:     time.Time{},
								Challenge:            "challenge",
								AllowedCredentialIDs: [][]byte{},
								UserVerification:     0,
								RPID:                 "example.com",
							},
						},
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)

				repo.EXPECT().
					LoadPasskeys().
					Times(1).
					Return(repo)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-jy0zq7", "Errors.User.UserIDMissing"),
		},
		{
			testName: "when retrieving user fails should return error",
			checkPasskey: &session_grpc.CheckWebAuthN{
				CredentialAssertionData: &structpb.Struct{},
			},
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
						Challenges: domain.SessionChallenges{
							&domain.SessionChallengePasskey{
								Challenge:        "challenge",
								RPID:             "example.com",
								UserVerification: old_domain.UserVerificationRequirementRequired,
							},
						},
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					LoadPasskeys().
					Return(repo)
				idCondition := getUserIDCondition(repo, "user-1")

				repo.EXPECT().Human().Times(1).Return(humanRepo)

				pkeyCondition := database.NewNumberCondition(
					database.NewColumn("zitadel.humans", "passkey_type"),
					database.NumberOperationEqual,
					domain.PasskeyTypePasswordless,
				)

				humanRepo.EXPECT().
					PasskeyTypeCondition(database.NumberOperationEqual, domain.PasskeyTypePasswordless).
					Times(1).
					Return(pkeyCondition)

				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(
							idCondition,
						)),
						dbmock.QueryOptions(database.WithCondition(
							pkeyCondition)),
					).
					Times(1).
					Return(nil, getErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(getErr, "DOM-pB6Mlm", "failed fetching user"),
		},
		{
			testName: "when all validations pass should return no error",
			checkPasskey: &session_grpc.CheckWebAuthN{
				CredentialAssertionData: &structpb.Struct{},
			},
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
						Challenges: domain.SessionChallenges{
							&domain.SessionChallengePasskey{
								Challenge: "challenge",
								RPID:      "example.com",
							},
						},
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					LoadPasskeys().
					Return(repo)
				idCondition := getUserIDCondition(repo, "user-1")
				repo.EXPECT().Human().Times(1).Return(humanRepo)

				pkeyCondition := database.NewNumberCondition(
					database.NewColumn("zitadel.humans", "passkey_type"),
					database.NumberOperationEqual,
					domain.PasskeyTypePasswordless,
				)

				humanRepo.EXPECT().
					PasskeyTypeCondition(database.NumberOperationNotEqual, domain.PasskeyTypePasswordless).
					Times(1).
					Return(pkeyCondition)

				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(
							idCondition,
						)),
						dbmock.QueryOptions(database.WithCondition(
							pkeyCondition)),
					).
					Times(1).Return(&domain.User{
					ID:       "user-1",
					Username: "testuser",
					Human: &domain.HumanUser{
						DisplayName: "Test User",
						Passkeys:    []*domain.Passkey{},
					},
				}, nil)
				return repo
			},
			expectedError: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewPasskeyCheckCommand("session-1", "instance-1", tc.checkPasskey, nil)

			opts := &domain.InvokeOpts{}
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
		})
	}
}

func TestPasskeyCheckCommand_Execute(t *testing.T) {
	t.Parallel()
	finishLoginErr := errors.New("finish login error")
	sessionUpdateErr := errors.New("session update error")
	userUpdateErr := errors.New("user update error")

	tt := []struct {
		testName       string
		sessionRepo    func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo       func(ctrl *gomock.Controller) domain.UserRepository
		finishLoginFn  func(ctx context.Context, sessionData webauthn.SessionData, user webauthn.User, credentials []byte, rpID string) (*webauthn.Credential, error)
		checkPasskey   *session_grpc.CheckWebAuthN
		fetchedUser    *domain.User
		fetchedSession *domain.Session

		expectedError error
	}{
		{
			testName:      "when checkPasskey is nil should return no error",
			checkPasskey:  nil,
			expectedError: nil,
		},
		{
			testName: "when finish login fails should return error",
			checkPasskey: &session_grpc.CheckWebAuthN{
				CredentialAssertionData: &structpb.Struct{},
			},
			fetchedSession: &domain.Session{
				ID:     "session-1",
				UserID: "user-1",
				Challenges: domain.SessionChallenges{
					&domain.SessionChallengePasskey{
						Challenge: "challenge",
						RPID:      "example.com",
					},
				},
			},
			fetchedUser: &domain.User{
				ID:       "user-1",
				Username: "testuser",
				Human: &domain.HumanUser{
					DisplayName: "Test User",
					Passkeys:    []*domain.Passkey{},
				},
			},
			finishLoginFn: func(ctx context.Context, sessionData webauthn.SessionData, user webauthn.User, credentials []byte, rpID string) (*webauthn.Credential, error) {
				return nil, finishLoginErr
			},
			expectedError: finishLoginErr,
		},
		{
			testName: "when passkey not found should return precondition failed error",
			checkPasskey: &session_grpc.CheckWebAuthN{
				CredentialAssertionData: &structpb.Struct{},
			},
			fetchedSession: &domain.Session{
				ID:     "session-1",
				UserID: "user-1",
				Challenges: domain.SessionChallenges{
					&domain.SessionChallengePasskey{
						Challenge: "challenge",
						RPID:      "example.com",
					},
				},
			},
			fetchedUser: &domain.User{
				ID:       "user-1",
				Username: "testuser",
				Human: &domain.HumanUser{
					DisplayName: "Test User",
					Passkeys: []*domain.Passkey{
						{
							ID:    "pkey-1",
							KeyID: []byte("key-id-1"),
						},
					},
				},
			},
			finishLoginFn: func(ctx context.Context, sessionData webauthn.SessionData, user webauthn.User, credentials []byte, rpID string) (*webauthn.Credential, error) {
				return &webauthn.Credential{
					ID: []byte("non-matching-key-id"),
				}, nil
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-uuxodH", "Errors.User.WebAuthN.NotFound"),
		},
		{
			testName: "when session update fails should return error",
			checkPasskey: &session_grpc.CheckWebAuthN{
				CredentialAssertionData: &structpb.Struct{},
			},
			fetchedSession: &domain.Session{
				ID:     "session-1",
				UserID: "user-1",
				Challenges: domain.SessionChallenges{
					&domain.SessionChallengePasskey{
						Challenge: "challenge",
						RPID:      "example.com",
					},
				},
			},
			fetchedUser: &domain.User{
				ID:       "user-1",
				Username: "testuser",
				Human: &domain.HumanUser{
					DisplayName: "Test User",
					Passkeys: []*domain.Passkey{
						{
							ID:    "pkey-1",
							KeyID: []byte("key-id-1"),
						},
					},
				},
			},
			finishLoginFn: func(ctx context.Context, sessionData webauthn.SessionData, user webauthn.User, credentials []byte, rpID string) (*webauthn.Credential, error) {
				return &webauthn.Credential{
					ID: []byte("key-id-1"),
					Flags: webauthn.CredentialFlags{
						UserVerified: true,
					},
					Authenticator: webauthn.Authenticator{
						SignCount: 5,
					},
				}, nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				factorChange := getSessionPasskeyFactorChange(repo, time.Now(), true)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						idCondition,
						factorChange).
					Times(1).
					Return(int64(0), sessionUpdateErr)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				return repo
			},
			expectedError: zerrors.ThrowInternal(sessionUpdateErr, "DOM-Uadvap", "failed updating session"),
		},
		{
			testName: "when session not found should return not found error",
			checkPasskey: &session_grpc.CheckWebAuthN{
				CredentialAssertionData: &structpb.Struct{},
			},
			fetchedSession: &domain.Session{
				ID:     "session-1",
				UserID: "user-1",
				Challenges: domain.SessionChallenges{
					&domain.SessionChallengePasskey{
						Challenge: "challenge",
						RPID:      "example.com",
					},
				},
			},
			fetchedUser: &domain.User{
				ID:       "user-1",
				Username: "testuser",
				Human: &domain.HumanUser{
					DisplayName: "Test User",
					Passkeys: []*domain.Passkey{
						{
							ID:    "pkey-1",
							KeyID: []byte("key-id-1"),
						},
					},
				},
			},
			finishLoginFn: func(ctx context.Context, sessionData webauthn.SessionData, user webauthn.User, credentials []byte, rpID string) (*webauthn.Credential, error) {
				return &webauthn.Credential{
					ID: []byte("key-id-1"),
					Flags: webauthn.CredentialFlags{
						UserVerified: true,
					},
					Authenticator: webauthn.Authenticator{
						SignCount: 5,
					},
				}, nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				factorChange := getSessionPasskeyFactorChange(repo, time.Now(), true)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						idCondition,
						factorChange).
					Times(1).
					Return(int64(0), nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(nil, "DOM-Uadvap", "session not found"),
		},
		{
			testName: "when user update fails should return error",
			checkPasskey: &session_grpc.CheckWebAuthN{
				CredentialAssertionData: &structpb.Struct{},
			},
			fetchedSession: &domain.Session{
				ID:     "session-1",
				UserID: "user-1",
				Challenges: domain.SessionChallenges{
					&domain.SessionChallengePasskey{
						Challenge: "challenge",
						RPID:      "example.com",
					},
				},
			},
			fetchedUser: &domain.User{
				ID:       "user-1",
				Username: "testuser",
				Human: &domain.HumanUser{
					DisplayName: "Test User",
					Passkeys: []*domain.Passkey{
						{
							ID:    "pkey-1",
							KeyID: []byte("key-id-1"),
						},
					},
				},
			},
			finishLoginFn: func(ctx context.Context, sessionData webauthn.SessionData, user webauthn.User, credentials []byte, rpID string) (*webauthn.Credential, error) {
				return &webauthn.Credential{
					ID: []byte("key-id-1"),
					Flags: webauthn.CredentialFlags{
						UserVerified: true,
					},
					Authenticator: webauthn.Authenticator{
						SignCount: 5,
					},
				}, nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				factorChange := getSessionPasskeyFactorChange(repo, time.Now(), true)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						idCondition,
						factorChange).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().Human().AnyTimes().Return(humanRepo)
				pkeyCondition := getPasskeyCondition(humanRepo, "pkey-1")
				pkeySignCountChange := getHumanPasskeySignCount(humanRepo, 5)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						pkeyCondition,
						pkeySignCountChange).
					Times(1).
					Return(int64(0), userUpdateErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(userUpdateErr, "DOM-wdwZYk", "failed updating user"),
		},
		{
			testName: "when execute succeeds should return no error",
			checkPasskey: &session_grpc.CheckWebAuthN{
				CredentialAssertionData: &structpb.Struct{},
			},
			fetchedSession: &domain.Session{
				ID:     "session-1",
				UserID: "user-1",
				Challenges: domain.SessionChallenges{
					&domain.SessionChallengePasskey{
						Challenge: "challenge",
						RPID:      "example.com",
					},
				},
			},
			fetchedUser: &domain.User{
				ID:       "user-1",
				Username: "testuser",
				Human: &domain.HumanUser{
					DisplayName: "Test User",
					Passkeys: []*domain.Passkey{
						{
							ID:    "pkey-1",
							KeyID: []byte("key-id-1"),
						},
					},
				},
			},
			finishLoginFn: func(ctx context.Context, sessionData webauthn.SessionData, user webauthn.User, credentials []byte, rpID string) (*webauthn.Credential, error) {
				return &webauthn.Credential{
					ID: []byte("key-id-1"),
					Flags: webauthn.CredentialFlags{
						UserVerified: true,
					},
					Authenticator: webauthn.Authenticator{
						SignCount: 5,
					},
				}, nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				idCondition := getSessionIDCondition(repo, "session-1")
				factorChange := getSessionPasskeyFactorChange(repo, time.Now(), true)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						idCondition,
						factorChange).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().Human().AnyTimes().Return(humanRepo)
				pkeyCondition := getPasskeyCondition(humanRepo, "pkey-1")
				pkeySignCountChange := getHumanPasskeySignCount(humanRepo, 5)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						pkeyCondition,
						pkeySignCountChange).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			expectedError: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewPasskeyCheckCommand("session-1", "instance-1", tc.checkPasskey, tc.finishLoginFn)
			cmd.FetchedSession = tc.fetchedSession
			cmd.FetchedUser = tc.fetchedUser

			opts := &domain.InvokeOpts{}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.sessionRepo != nil {
				domain.WithSessionRepo(tc.sessionRepo(ctrl))(opts)
			}
			if tc.userRepo != nil {
				domain.WithUserRepo(tc.userRepo(ctrl))(opts)
			}

			// Test
			err := cmd.Execute(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			if tc.checkPasskey != nil && tc.expectedError == nil {
				assert.NotZero(t, cmd.LastVeriedAt)
			}
		})
	}
}

func TestPasskeyCheckCommand_Events(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName           string
		checkPasskey       *session_grpc.CheckWebAuthN
		fetchedSession     *domain.Session
		fetchedUser        *domain.User
		lastVerifiedAt     time.Time
		userVerified       bool
		pkeyID             string
		pkeySignCount      uint32
		expectedEventCount int
		expectedEvents     []eventstore.Command
	}{
		{
			testName:           "when checkPasskey is nil should return no events",
			checkPasskey:       nil,
			expectedEventCount: 0,
			expectedEvents:     []eventstore.Command{},
		},
		{
			testName:       "when user verification is required should return WebAuthNCheckedEvent and PasswordlessSignCountChangedEvent",
			checkPasskey:   &session_grpc.CheckWebAuthN{},
			lastVerifiedAt: time.Now(),
			userVerified:   true,
			pkeyID:         "pkey-1",
			pkeySignCount:  5,
			fetchedSession: &domain.Session{
				ID:     "session-1",
				UserID: "user-1",
				Challenges: domain.SessionChallenges{
					&domain.SessionChallengePasskey{
						Challenge:        "challenge",
						RPID:             "example.com",
						UserVerification: old_domain.UserVerificationRequirementRequired,
					},
				},
			},
			fetchedUser: &domain.User{
				ID:       "user-1",
				Username: "testuser",
				Human: &domain.HumanUser{
					DisplayName: "Test User",
					Passkeys: []*domain.Passkey{
						{
							ID:    "pkey-1",
							KeyID: []byte("key-id-1"),
						},
					},
				},
			},
			expectedEventCount: 2,
			expectedEvents: []eventstore.Command{
				session.NewWebAuthNCheckedEvent(t.Context(), nil, time.Now(), true),
				user.NewHumanPasswordlessSignCountChangedEvent(t.Context(), nil, "pkey-1", 5),
			},
		},
		{
			testName:       "when user verification is not required should return WebAuthNCheckedEvent and U2FSignCountChangedEvent",
			checkPasskey:   &session_grpc.CheckWebAuthN{},
			lastVerifiedAt: time.Now(),
			userVerified:   false,
			pkeyID:         "pkey-2",
			pkeySignCount:  10,
			fetchedSession: &domain.Session{
				ID:     "session-1",
				UserID: "user-1",
				Challenges: domain.SessionChallenges{
					&domain.SessionChallengePasskey{
						Challenge:        "challenge",
						RPID:             "example.com",
						UserVerification: old_domain.UserVerificationRequirementPreferred,
					},
				},
			},
			fetchedUser: &domain.User{
				ID:       "user-1",
				Username: "testuser",
				Human: &domain.HumanUser{
					DisplayName: "Test User",
					Passkeys: []*domain.Passkey{
						{
							ID:    "pkey-2",
							KeyID: []byte("key-id-2"),
						},
					},
				},
			},
			expectedEventCount: 2,
			expectedEvents: []eventstore.Command{
				session.NewWebAuthNCheckedEvent(t.Context(), nil, time.Now(), false),
				user.NewHumanU2FSignCountChangedEvent(t.Context(), nil, "pkey-2", 10),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			cmd := domain.NewPasskeyCheckCommand("session-1", "instance-1", tc.checkPasskey, nil)
			cmd.FetchedSession = tc.fetchedSession
			cmd.FetchedUser = tc.fetchedUser
			cmd.LastVeriedAt = tc.lastVerifiedAt
			cmd.UserVerified = tc.userVerified
			cmd.PKeyID = tc.pkeyID
			cmd.PKeySignCount = tc.pkeySignCount

			opts := &domain.InvokeOpts{}

			// Test
			events, err := cmd.Events(ctx, opts)

			// Verify
			assert.NoError(t, err)
			assert.Len(t, events, tc.expectedEventCount)
			for i, expectedType := range tc.expectedEvents {
				assert.IsType(t, expectedType, events[i])
				switch expectedAssertedType := expectedType.(type) {
				case *session.WebAuthNCheckedEvent:
					actualAssertedType, ok := events[i].(*session.WebAuthNCheckedEvent)
					require.True(t, ok)
					assert.InDelta(t, expectedAssertedType.CheckedAt.UnixMilli(), actualAssertedType.CheckedAt.UnixMilli(), 1.5)
					assert.Equal(t, expectedAssertedType.UserVerified, actualAssertedType.UserVerified)
				case *user.HumanPasswordlessSignCountChangedEvent:
					actualAssertedType, ok := events[i].(*user.HumanPasswordlessSignCountChangedEvent)
					require.True(t, ok)
					assert.Equal(t, expectedAssertedType.WebAuthNTokenID, actualAssertedType.WebAuthNTokenID)
					assert.Equal(t, expectedAssertedType.SignCount, actualAssertedType.SignCount)
				case *user.HumanU2FSignCountChangedEvent:
					actualAssertedType, ok := events[i].(*user.HumanU2FSignCountChangedEvent)
					require.True(t, ok)
					assert.Equal(t, expectedAssertedType.WebAuthNTokenID, actualAssertedType.WebAuthNTokenID)
					assert.Equal(t, expectedAssertedType.SignCount, actualAssertedType.SignCount)
				}
			}
		})
	}
}
