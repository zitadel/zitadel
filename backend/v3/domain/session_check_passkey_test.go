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
)

func TestPasskeyCheckCommand_Validate(t *testing.T) {
	t.Parallel()
	getErr := errors.New("get error")

	tt := []struct {
		testName      string
		sessionRepo   func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo      func(ctrl *gomock.Controller) domain.UserRepository
		cmd           *domain.PasskeyCheckCommand
		expectedError error
	}{
		{
			testName:      "when checkPasskey is nil should return no error",
			cmd:           &domain.PasskeyCheckCommand{},
			expectedError: nil,
		},
		{
			testName:      "when sessionID is not set should return error",
			cmd:           &domain.PasskeyCheckCommand{CheckPasskey: []byte{}},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-4QJa2k", "Errors.Missing.SessionID"),
		},
		{
			testName:      "when instanceID is not set should return error",
			cmd:           &domain.PasskeyCheckCommand{CheckPasskey: []byte{}, SessionID: "session-1"},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-XlOhxU", "Errors.Missing.InstanceID"),
		},
		{
			testName: "when retrieving session fails should return error",
			cmd:      &domain.PasskeyCheckCommand{CheckPasskey: []byte{}, SessionID: "session-1", InstanceID: "instance-1"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.IDCondition("session-1")
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
				repo := domainmock.NewUserRepo(ctrl)
				return repo
			},
			expectedError: zerrors.ThrowInternal(getErr, "DOM-CUnePh", "failed fetching session"),
		},
		{
			testName: "when session has no passkey challenge should return precondition failed error",
			cmd:      &domain.PasskeyCheckCommand{CheckPasskey: []byte{}, SessionID: "session-1", InstanceID: "instance-1"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.IDCondition("session-1")
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
				repo := domainmock.NewUserRepo(ctrl)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-lQhNR4", "Errors.Session.WebAuthN.NoChallenge"),
		},
		{
			testName: "when session has no user ID should return precondition failed error",
			cmd:      &domain.PasskeyCheckCommand{CheckPasskey: []byte{}, SessionID: "session-1", InstanceID: "instance-1"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)

				idCondition := repo.IDCondition("session-1")
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
				repo := domainmock.NewUserRepo(ctrl)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-jy0zq7", "Errors.User.UserIDMissing"),
		},
		{
			testName: "when retrieving user fails should return error",
			cmd:      &domain.PasskeyCheckCommand{CheckPasskey: []byte{}, SessionID: "session-1", InstanceID: "instance-1"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.IDCondition("session-1")
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
				repo := domainmock.NewUserRepo(ctrl)
				humanRepo := domainmock.NewHumanRepo(ctrl)
				repo.EXPECT().Human().Times(1).Return(humanRepo)

				idCondition := repo.IDCondition("user-1")

				pkeyCondition := humanRepo.
					PasskeyConditions().
					TypeCondition(database.TextOperationEqual, domain.PasskeyTypePasswordless)

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
			cmd:      &domain.PasskeyCheckCommand{CheckPasskey: []byte{}, SessionID: "session-1", InstanceID: "instance-1"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.IDCondition("session-1")
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
				repo := domainmock.NewUserRepo(ctrl)
				humanRepo := domainmock.NewHumanRepo(ctrl)

				repo.EXPECT().Human().Times(1).Return(humanRepo)
				idCondition := repo.IDCondition("user-1")

				pkeyCondition := humanRepo.
					PasskeyConditions().
					TypeCondition(database.TextOperationEqual, domain.PasskeyTypeU2F)

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
			ctx := authz.NewMockContext(tc.cmd.InstanceID, "", "")
			ctrl := gomock.NewController(t)

			opts := &domain.InvokeOpts{}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.sessionRepo != nil {
				domain.WithSessionRepo(tc.sessionRepo(ctrl))(opts)
			}
			if tc.userRepo != nil {
				domain.WithUserRepo(tc.userRepo(ctrl))(opts)
			}

			// Test
			err := tc.cmd.Validate(ctx, opts)

			// Verify
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestPasskeyCheckCommand_Execute(t *testing.T) {
	t.Parallel()
	finishLoginErr := errors.New("finish login error")
	sessionUpdateErr := errors.New("session update error")
	userUpdateErr := errors.New("user update error")

	finishLoginErrFn := func(ctx context.Context, sessionData webauthn.SessionData, user webauthn.User, credentials []byte, rpID string) (*webauthn.Credential, error) {
		return nil, finishLoginErr
	}
	finishLoginNonMatchingPasskeyFn := func(ctx context.Context, sessionData webauthn.SessionData, user webauthn.User, credentials []byte, rpID string) (*webauthn.Credential, error) {
		return &webauthn.Credential{ID: []byte("non-matching-key-id")}, nil
	}
	finishLoginOKFn := func(ctx context.Context, sessionData webauthn.SessionData, user webauthn.User, credentials []byte, rpID string) (*webauthn.Credential, error) {
		return &webauthn.Credential{
			ID: []byte("key-id-1"),
			Flags: webauthn.CredentialFlags{
				UserVerified: true,
			},
			Authenticator: webauthn.Authenticator{
				SignCount: 5,
			},
		}, nil
	}

	tt := []struct {
		testName      string
		sessionRepo   func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo      func(ctrl *gomock.Controller) domain.UserRepository
		finishLoginFn func(ctx context.Context, sessionData webauthn.SessionData, user webauthn.User, credentials []byte, rpID string) (*webauthn.Credential, error)
		cmd           *domain.PasskeyCheckCommand

		expectedError            error
		expectedPasskeySignCount uint32
		expectedPasskeyID        string
		expectedUserVerified     bool
	}{
		{
			testName:      "when checkPasskey is nil should return no error",
			cmd:           &domain.PasskeyCheckCommand{},
			expectedError: nil,
		},
		{
			testName: "when finish login fails should return error",
			cmd: &domain.PasskeyCheckCommand{
				CheckPasskey: []byte{},
				SessionID:    "session-1",
				InstanceID:   "instance-1",
				FetchedUser: &domain.User{
					ID:       "user-1",
					Username: "testuser",
					Human:    &domain.HumanUser{DisplayName: "Test User", Passkeys: []*domain.Passkey{}},
				},
				FetchedSession: &domain.Session{
					ID:     "session-1",
					UserID: "user-1",
					Challenges: domain.SessionChallenges{
						&domain.SessionChallengePasskey{Challenge: "challenge", RPID: "example.com"},
					},
				},
				FinishLoginFn: finishLoginErrFn,
			},
			expectedError: finishLoginErr,
		},
		{
			testName: "when passkey not found should return precondition failed error",
			cmd: &domain.PasskeyCheckCommand{
				CheckPasskey: []byte{},
				SessionID:    "session-1",
				InstanceID:   "instance-1",
				FetchedUser: &domain.User{
					ID:       "user-1",
					Username: "testuser",
					Human: &domain.HumanUser{
						DisplayName: "Test User",
						Passkeys: []*domain.Passkey{
							{ID: "pkey-1", KeyID: []byte("key-id-1")},
						},
					},
				},
				FetchedSession: &domain.Session{
					ID:     "session-1",
					UserID: "user-1",
					Challenges: domain.SessionChallenges{
						&domain.SessionChallengePasskey{Challenge: "challenge", RPID: "example.com"},
					},
				},
				FinishLoginFn: finishLoginNonMatchingPasskeyFn,
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-uuxodH", "Errors.User.WebAuthN.NotFound"),
		},
		{
			testName: "when session update fails should return error",
			cmd: &domain.PasskeyCheckCommand{
				CheckPasskey: []byte{},
				SessionID:    "session-1",
				InstanceID:   "instance-1",
				FetchedUser: &domain.User{
					ID:       "user-1",
					Username: "testuser",
					Human: &domain.HumanUser{
						DisplayName: "Test User",
						Passkeys: []*domain.Passkey{
							{ID: "pkey-1", KeyID: []byte("key-id-1")},
						},
					},
				},
				FetchedSession: &domain.Session{
					ID:     "session-1",
					UserID: "user-1",
					Challenges: domain.SessionChallenges{
						&domain.SessionChallengePasskey{Challenge: "challenge", RPID: "example.com"},
					},
				},
				FinishLoginFn: finishLoginOKFn,
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.IDCondition("session-1")
				factorChange := repo.SetFactor(&domain.SessionFactorPasskey{LastVerifiedAt: time.Now(), UserVerified: true})
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						idCondition,
						factorChange).
					Times(1).
					Return(int64(0), sessionUpdateErr)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				return repo
			},
			expectedError:        zerrors.ThrowInternal(sessionUpdateErr, "DOM-Uadvap", "failed updating session"),
			expectedPasskeyID:    "pkey-1",
			expectedUserVerified: true,
		},
		{
			testName: "when session not found should return not found error",
			cmd: &domain.PasskeyCheckCommand{
				CheckPasskey: []byte{},
				SessionID:    "session-1",
				InstanceID:   "instance-1",
				FetchedUser: &domain.User{
					ID:       "user-1",
					Username: "testuser",
					Human: &domain.HumanUser{
						DisplayName: "Test User",
						Passkeys: []*domain.Passkey{
							{ID: "pkey-1", KeyID: []byte("key-id-1")},
						},
					},
				},
				FetchedSession: &domain.Session{
					ID:     "session-1",
					UserID: "user-1",
					Challenges: domain.SessionChallenges{
						&domain.SessionChallengePasskey{Challenge: "challenge", RPID: "example.com"},
					},
				},
				FinishLoginFn: finishLoginOKFn,
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.IDCondition("session-1")
				factorChange := repo.SetFactor(&domain.SessionFactorPasskey{LastVerifiedAt: time.Now(), UserVerified: true})
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						idCondition,
						factorChange).
					Times(1).
					Return(int64(0), nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				return repo
			},
			expectedError:        zerrors.ThrowNotFound(nil, "DOM-Uadvap", "session not found"),
			expectedPasskeyID:    "pkey-1",
			expectedUserVerified: true,
		},
		{
			testName: "when user update fails should return error",
			cmd: &domain.PasskeyCheckCommand{
				CheckPasskey: []byte{},
				SessionID:    "session-1",
				InstanceID:   "instance-1",
				FetchedUser: &domain.User{
					ID:       "user-1",
					Username: "testuser",
					Human: &domain.HumanUser{
						DisplayName: "Test User",
						Passkeys: []*domain.Passkey{
							{ID: "pkey-1", KeyID: []byte("key-id-1")},
						},
					},
				},
				FetchedSession: &domain.Session{
					ID:     "session-1",
					UserID: "user-1",
					Challenges: domain.SessionChallenges{
						&domain.SessionChallengePasskey{Challenge: "challenge", RPID: "example.com"},
					},
				},
				FinishLoginFn: finishLoginOKFn,
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.IDCondition("session-1")
				factorChange := repo.SetFactor(&domain.SessionFactorPasskey{LastVerifiedAt: time.Now(), UserVerified: true})
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						idCondition,
						factorChange).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				humanRepo := domainmock.NewHumanRepo(ctrl)
				repo.EXPECT().Human().Times(3).Return(humanRepo)
				updateConditions := database.And(
					humanRepo.PrimaryKeyCondition("instance-1", "user-1"),
					humanRepo.PasskeyConditions().IDCondition("pkey-1"),
				)
				pkeySignCountChange := humanRepo.SetPasskeySignCount(5)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						updateConditions,
						pkeySignCountChange).
					Times(1).
					Return(int64(0), userUpdateErr)
				return repo
			},
			expectedError:        zerrors.ThrowInternal(userUpdateErr, "DOM-wdwZYk", "failed updating user"),
			expectedPasskeyID:    "pkey-1",
			expectedUserVerified: true,
		},
		{
			testName: "when execute succeeds should return no error",
			cmd: &domain.PasskeyCheckCommand{
				CheckPasskey: []byte{},
				SessionID:    "session-1",
				InstanceID:   "instance-1",
				FetchedUser: &domain.User{
					ID:       "user-1",
					Username: "testuser",
					Human: &domain.HumanUser{
						DisplayName: "Test User",
						Passkeys: []*domain.Passkey{
							{ID: "pkey-1", KeyID: []byte("key-id-1")},
						},
					},
				},
				FetchedSession: &domain.Session{
					ID:     "session-1",
					UserID: "user-1",
					Challenges: domain.SessionChallenges{
						&domain.SessionChallengePasskey{Challenge: "challenge", RPID: "example.com"},
					},
				},
				FinishLoginFn: finishLoginOKFn,
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				idCondition := repo.IDCondition("session-1")
				factorChange := repo.SetFactor(&domain.SessionFactorPasskey{LastVerifiedAt: time.Now(), UserVerified: true})
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						idCondition,
						factorChange).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				humanRepo := domainmock.NewHumanRepo(ctrl)
				repo.EXPECT().Human().Times(3).Return(humanRepo)
				updateConditions := database.And(
					humanRepo.PrimaryKeyCondition("instance-1", "user-1"),
					humanRepo.PasskeyConditions().IDCondition("pkey-1"),
				)
				pkeySignCountChange := humanRepo.SetPasskeySignCount(5)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(),
						updateConditions,
						pkeySignCountChange).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			expectedError:            nil,
			expectedPasskeyID:        "pkey-1",
			expectedPasskeySignCount: 5,
			expectedUserVerified:     true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext(tc.cmd.InstanceID, "", "")
			ctrl := gomock.NewController(t)

			opts := &domain.InvokeOpts{}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.sessionRepo != nil {
				domain.WithSessionRepo(tc.sessionRepo(ctrl))(opts)
			}
			if tc.userRepo != nil {
				domain.WithUserRepo(tc.userRepo(ctrl))(opts)
			}

			// Test
			err := tc.cmd.Execute(ctx, opts)

			// Verify
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedPasskeyID, tc.cmd.PKeyID)
			assert.Equal(t, tc.expectedPasskeySignCount, tc.cmd.PKeySignCount)
			assert.Equal(t, tc.expectedUserVerified, tc.cmd.UserVerified)
			if tc.cmd.CheckPasskey != nil && tc.expectedError == nil {
				assert.NotZero(t, tc.cmd.LastVerifiedAt)
			}
		})
	}
}

func TestPasskeyCheckCommand_Events(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName string
		cmd      *domain.PasskeyCheckCommand

		expectedEvents []eventstore.Command
	}{
		{
			testName:       "when checkPasskey is nil should return no events",
			cmd:            &domain.PasskeyCheckCommand{},
			expectedEvents: []eventstore.Command{},
		},
		{
			testName: "when user verification is required should return WebAuthNCheckedEvent and PasswordlessSignCountChangedEvent",
			cmd: &domain.PasskeyCheckCommand{
				CheckPasskey: []byte{},
				SessionID:    "session-1",
				InstanceID:   "instance-1",
				FetchedUser: &domain.User{
					ID:       "user-1",
					Username: "testuser",
					Human: &domain.HumanUser{
						DisplayName: "Test User",
						Passkeys: []*domain.Passkey{
							{ID: "pkey-1", KeyID: []byte("key-id-1")},
						},
					},
				},
				FetchedSession: &domain.Session{
					ID:     "session-1",
					UserID: "user-1",
					Challenges: domain.SessionChallenges{
						&domain.SessionChallengePasskey{Challenge: "challenge", RPID: "example.com", UserVerification: old_domain.UserVerificationRequirementRequired},
					},
				},
				LastVerifiedAt: time.Now(),
				UserVerified:   true,
				PKeyID:         "pkey-1",
				PKeySignCount:  5,
			},
			expectedEvents: []eventstore.Command{
				session.NewWebAuthNCheckedEvent(t.Context(), nil, time.Now(), true),
				user.NewHumanPasswordlessSignCountChangedEvent(t.Context(), nil, "pkey-1", 5),
			},
		},
		{
			testName: "when user verification is not required should return WebAuthNCheckedEvent and U2FSignCountChangedEvent",
			cmd: &domain.PasskeyCheckCommand{
				CheckPasskey: []byte{},
				SessionID:    "session-1",
				InstanceID:   "instance-1",
				FetchedUser: &domain.User{
					ID:       "user-1",
					Username: "testuser",
					Human: &domain.HumanUser{
						DisplayName: "Test User",
						Passkeys: []*domain.Passkey{
							{ID: "pkey-2", KeyID: []byte("key-id-2")},
						},
					},
				},
				FetchedSession: &domain.Session{
					ID:     "session-1",
					UserID: "user-1",
					Challenges: domain.SessionChallenges{
						&domain.SessionChallengePasskey{Challenge: "challenge", RPID: "example.com", UserVerification: old_domain.UserVerificationRequirementPreferred},
					},
				},
				LastVerifiedAt: time.Now(),
				UserVerified:   false,
				PKeyID:         "pkey-2",
				PKeySignCount:  10,
			},
			expectedEvents: []eventstore.Command{
				session.NewWebAuthNCheckedEvent(t.Context(), nil, time.Now(), false),
				user.NewHumanU2FSignCountChangedEvent(t.Context(), nil, "pkey-2", 10),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Given
			ctx := authz.NewMockContext(tc.cmd.InstanceID, "", "")

			opts := &domain.InvokeOpts{}

			// Test
			events, err := tc.cmd.Events(ctx, opts)

			// Verify
			assert.NoError(t, err)
			require.Len(t, events, len(tc.expectedEvents))
			for i, expectedType := range tc.expectedEvents {
				assert.IsType(t, expectedType, events[i])
				switch expectedAssertedType := expectedType.(type) {
				case *session.WebAuthNCheckedEvent:
					actualAssertedType, ok := events[i].(*session.WebAuthNCheckedEvent)
					require.True(t, ok)
					assert.Equal(t, tc.cmd.LastVerifiedAt, actualAssertedType.CheckedAt)
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
