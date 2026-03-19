package domain_test

import (
	"encoding/base64"
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
	"github.com/zitadel/zitadel/internal/repository/idpintent"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestIDPIntentCheckCommand_Events(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name           string
		command        *domain.IDPIntentCheckCommand
		expectedEvents []eventstore.Command
	}{
		{
			name: "nil CheckIntent returns nil events",
			command: &domain.IDPIntentCheckCommand{
				CheckIntent:     nil,
				IsCheckComplete: true,
			},
		},
		{
			name: "isCheckComplete false returns nil events",
			command: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{},
			},
		},
		{
			name: "valid command returns events",
			command: &domain.IDPIntentCheckCommand{
				CheckIntent:        &domain.CheckIDPIntentType{ID: "intent-123"},
				SessionID:          "session-456",
				InstanceID:         "instance-789",
				IsCheckComplete:    true,
				IntentLastVerified: time.Now(),
			},

			expectedEvents: []eventstore.Command{
				session.NewIntentCheckedEvent(t.Context(), &session.NewAggregate("session-456", "instance-789").Aggregate, time.Now()),
				idpintent.NewConsumedEvent(t.Context(), &idpintent.NewAggregate("intent-123", "").Aggregate),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Test
			events, err := tc.command.Events(t.Context(), &domain.InvokeOpts{})

			// Verify
			assert.NoError(t, err)

			require.Len(t, events, len(tc.expectedEvents))

			for i, expectedType := range tc.expectedEvents {
				require.IsType(t, expectedType, events[i])
				switch expectedType.(type) {
				case *session.IntentCheckedEvent:
					actualAssertedType, ok := events[i].(*session.IntentCheckedEvent)
					require.True(t, ok)
					assert.Equal(t, tc.command.IntentLastVerified, actualAssertedType.CheckedAt)
				case *idpintent.ConsumedEvent:
					_, ok := events[i].(*idpintent.ConsumedEvent)
					require.True(t, ok)
				}
			}
		})
	}
}

func TestIDPIntentCheckCommand_Validate(t *testing.T) {
	t.Parallel()
	getErr := errors.New("get error")
	notFoundErr := database.NewNoRowFoundError(nil)
	decryptionError := errors.New("decryption error")
	invalidTokenEncoding := base64.RawURLEncoding.EncodeToString([]byte("invalid-token"))
	validTokenEncoding := base64.RawURLEncoding.EncodeToString([]byte("token-123"))

	tt := []struct {
		testName            string
		sessionRepo         func(ctrl *gomock.Controller) domain.SessionRepository
		idpIntentRepo       func(ctrl *gomock.Controller) domain.IDPIntentRepository
		userRepo            func(ctrl *gomock.Controller) domain.UserRepository
		encryptionAlgorithm func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm
		cmd                 *domain.IDPIntentCheckCommand
		expectedError       error
		expectedUser        domain.User
	}{
		{
			testName:      "when CheckIntent is nil should return no error",
			cmd:           &domain.IDPIntentCheckCommand{},
			expectedError: nil,
		},
		{
			testName: "when session ID is not set should return error",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: validTokenEncoding},
				SessionID:   "",
				InstanceID:  "instance-1",
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-5Y8pb4", "Errors.Missing.SessionID"),
		},
		{
			testName: "when instance ID is not set should return error",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: validTokenEncoding},
				SessionID:   "session-1",
				InstanceID:  "",
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-Q4YFIq", "Errors.Missing.InstanceID"),
		},
		{
			testName: "when retrieving session fails should return error",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: validTokenEncoding},
				SessionID:   "session-1",
				InstanceID:  "instance-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "session-1")))).
					Return(nil, getErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(getErr, "DOM-EhIgey", "failed fetching session"),
		},
		{
			testName: "when session not found should return error",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: validTokenEncoding},
				SessionID:   "session-1",
				InstanceID:  "instance-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "session-1")))).
					Return(nil, notFoundErr)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(notFoundErr, "DOM-EhIgey", "session not found"),
		},
		{
			testName: "when session userID is empty should return error",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: validTokenEncoding},
				SessionID:   "session-1",
				InstanceID:  "instance-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "session-1")))).
					Return(&domain.Session{ID: "session-1", UserID: ""}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-IJcVkV", "Errors.User.UserIDMissing"),
		},
		{
			testName: "when token verification fails should return error",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: invalidTokenEncoding},
				SessionID:   "session-1",
				InstanceID:  "instance-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "session-1")))).
					Return(&domain.Session{ID: "session-1", UserID: "user-1"}, nil)
				return repo
			},
			encryptionAlgorithm: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				cryptoAlg := crypto.NewMockEncryptionAlgorithm(ctrl)
				cryptoAlg.EXPECT().
					EncryptionKeyID().
					Times(1)
				cryptoAlg.EXPECT().
					DecryptString(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", decryptionError)
				return cryptoAlg
			},
			expectedError: zerrors.ThrowPermissionDenied(decryptionError, "CRYPTO-Sf4gt", "Errors.Intent.InvalidToken"),
		},
		{
			testName: "when retrieving intent fails should return error",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: validTokenEncoding},
				SessionID:   "session-1",
				InstanceID:  "instance-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "session-1")))).
					Return(&domain.Session{ID: "session-1", UserID: "user-1"}, nil)
				return repo
			},
			encryptionAlgorithm: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				cryptoAlg := crypto.NewMockEncryptionAlgorithm(ctrl)
				cryptoAlg.EXPECT().
					EncryptionKeyID().
					Times(1)
				cryptoAlg.EXPECT().
					DecryptString(gomock.Any(), gomock.Any()).
					Times(1).
					Return("intent-123", nil)
				return cryptoAlg
			},
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "intent-123")))).
					Return(nil, getErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(getErr, "DOM-5XkWJV", "failed fetching intent"),
		},
		{
			testName: "when intent not found should return error",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: validTokenEncoding},
				SessionID:   "session-1",
				InstanceID:  "instance-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "session-1")))).
					Return(&domain.Session{ID: "session-1", UserID: "user-1"}, nil)
				return repo
			},
			encryptionAlgorithm: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				cryptoAlg := crypto.NewMockEncryptionAlgorithm(ctrl)
				cryptoAlg.EXPECT().
					EncryptionKeyID().
					Times(1)
				cryptoAlg.EXPECT().
					DecryptString(gomock.Any(), gomock.Any()).
					Times(1).
					Return("intent-123", nil)
				return cryptoAlg
			},
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "intent-123")))).
					Return(nil, notFoundErr)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(notFoundErr, "DOM-5XkWJV", "intent not found"),
		},
		{
			testName: "when intent state is not succeeded should return error",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: validTokenEncoding},
				SessionID:   "session-1",
				InstanceID:  "instance-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "session-1")))).
					Return(&domain.Session{ID: "session-1", UserID: "user-1"}, nil)
				return repo
			},
			encryptionAlgorithm: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				cryptoAlg := crypto.NewMockEncryptionAlgorithm(ctrl)
				cryptoAlg.EXPECT().
					EncryptionKeyID().
					Times(1)
				cryptoAlg.EXPECT().
					DecryptString(gomock.Any(), gomock.Any()).
					Times(1).
					Return("intent-123", nil)
				return cryptoAlg
			},
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "intent-123")))).
					Return(&domain.IDPIntent{ID: "intent-123", State: domain.IDPIntentStateFailed, ExpiresAt: gu.Ptr(time.Now().Add(time.Hour))}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-0UKHku", "Errors.Intent.NotSucceeded"),
		},
		{
			testName: "when intent is expired should return error",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: validTokenEncoding},
				SessionID:   "session-1",
				InstanceID:  "instance-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "session-1")))).
					Return(&domain.Session{ID: "session-1", UserID: "user-1"}, nil)
				return repo
			},
			encryptionAlgorithm: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				cryptoAlg := crypto.NewMockEncryptionAlgorithm(ctrl)
				cryptoAlg.EXPECT().
					EncryptionKeyID().
					Times(1)
				cryptoAlg.EXPECT().
					DecryptString(gomock.Any(), gomock.Any()).
					Times(1).
					Return("intent-123", nil)
				return cryptoAlg
			},
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "intent-123")))).
					Return(&domain.IDPIntent{ID: "intent-123", State: domain.IDPIntentStateSucceeded, ExpiresAt: gu.Ptr(time.Now().Add(-time.Hour))}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-kDR1XK", "Errors.Intent.Expired"),
		},
		{
			testName: "when intent expiration is nil should return error",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: validTokenEncoding},
				SessionID:   "session-1",
				InstanceID:  "instance-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "session-1")))).
					Return(&domain.Session{ID: "session-1", UserID: "user-1"}, nil)
				return repo
			},
			encryptionAlgorithm: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				cryptoAlg := crypto.NewMockEncryptionAlgorithm(ctrl)
				cryptoAlg.EXPECT().
					EncryptionKeyID().
					Times(1)
				cryptoAlg.EXPECT().
					DecryptString(gomock.Any(), gomock.Any()).
					Times(1).
					Return("intent-123", nil)
				return cryptoAlg
			},
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "intent-123")))).
					Return(&domain.IDPIntent{ID: "intent-123", State: domain.IDPIntentStateSucceeded, ExpiresAt: nil}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-kDR1XK", "Errors.Intent.Expired"),
		},
		{
			testName: "when retrieving user fails should return error",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: validTokenEncoding},
				SessionID:   "session-1",
				InstanceID:  "instance-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "session-1")))).
					Return(&domain.Session{ID: "session-1", UserID: "user-1"}, nil)
				return repo
			},
			encryptionAlgorithm: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				cryptoAlg := crypto.NewMockEncryptionAlgorithm(ctrl)
				cryptoAlg.EXPECT().
					EncryptionKeyID().
					Times(1)
				cryptoAlg.EXPECT().
					DecryptString(gomock.Any(), gomock.Any()).
					Times(1).
					Return("intent-123", nil)
				return cryptoAlg
			},
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "intent-123")))).
					Return(&domain.IDPIntent{ID: "intent-123", State: domain.IDPIntentStateSucceeded, ExpiresAt: gu.Ptr(time.Now().Add(time.Hour)), UserID: ""}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "user-1")))).
					Return(nil, getErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(getErr, "DOM-Vnx2G9", "failed fetching user"),
		},
		{
			testName: "when user is not found should return not found error",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: validTokenEncoding},
				SessionID:   "session-1",
				InstanceID:  "instance-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "session-1")))).
					Return(&domain.Session{ID: "session-1", UserID: "user-1"}, nil)
				return repo
			},
			encryptionAlgorithm: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				cryptoAlg := crypto.NewMockEncryptionAlgorithm(ctrl)
				cryptoAlg.EXPECT().
					EncryptionKeyID().
					Times(1)
				cryptoAlg.EXPECT().
					DecryptString(gomock.Any(), gomock.Any()).
					Times(1).
					Return("intent-123", nil)
				return cryptoAlg
			},
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "intent-123")))).
					Return(&domain.IDPIntent{ID: "intent-123", State: domain.IDPIntentStateSucceeded, ExpiresAt: gu.Ptr(time.Now().Add(time.Hour)), UserID: ""}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "user-1")))).
					Return(nil, notFoundErr)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(notFoundErr, "DOM-Vnx2G9", "user not found"),
		},
		{
			testName: "when user is not human should return error",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: validTokenEncoding},
				SessionID:   "session-1",
				InstanceID:  "instance-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "session-1")))).
					Return(&domain.Session{ID: "session-1", UserID: "user-1"}, nil)
				return repo
			},
			encryptionAlgorithm: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				cryptoAlg := crypto.NewMockEncryptionAlgorithm(ctrl)
				cryptoAlg.EXPECT().
					EncryptionKeyID().
					Times(1)
				cryptoAlg.EXPECT().
					DecryptString(gomock.Any(), gomock.Any()).
					Times(1).
					Return("intent-123", nil)
				return cryptoAlg
			},
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "intent-123")))).
					Return(&domain.IDPIntent{ID: "intent-123", State: domain.IDPIntentStateSucceeded, ExpiresAt: gu.Ptr(time.Now().Add(time.Hour)), UserID: ""}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "user-1")))).
					Return(&domain.User{ID: "user-1", State: domain.UserStateActive, Human: nil}, nil)
				return repo
			},
			expectedError: zerrors.ThrowInternal(nil, "DOM-FkX5lZ", "user not human"),
		},
		{
			testName: "when intent has user and it matches session user should return no error and set fetched user",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: validTokenEncoding},
				SessionID:   "session-1",
				InstanceID:  "instance-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "session-1")))).
					Return(&domain.Session{ID: "session-1", UserID: "user-1"}, nil)
				return repo
			},
			encryptionAlgorithm: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				cryptoAlg := crypto.NewMockEncryptionAlgorithm(ctrl)
				cryptoAlg.EXPECT().
					EncryptionKeyID().
					Times(1)
				cryptoAlg.EXPECT().
					DecryptString(gomock.Any(), gomock.Any()).
					Times(1).
					Return("intent-123", nil)
				return cryptoAlg
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "user-1")))).
					Return(&domain.User{
						ID:    "user-1",
						State: domain.UserStateActive,
						Human: &domain.HumanUser{IdentityProviderLinks: []*domain.IdentityProviderLink{{ProviderID: "idp-123"}}},
					}, nil)
				return repo
			},
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "intent-123")))).
					Return(&domain.IDPIntent{ID: "intent-123", State: domain.IDPIntentStateSucceeded, ExpiresAt: gu.Ptr(time.Now().Add(time.Hour)), UserID: "user-1"}, nil)
				return repo
			},
			expectedError: nil,
			expectedUser: domain.User{
				ID:    "user-1",
				State: domain.UserStateActive,
				Human: &domain.HumanUser{
					IdentityProviderLinks: []*domain.IdentityProviderLink{
						{ProviderID: "idp-123"},
					},
				},
			},
		},
		{
			testName: "when intent has user but it does not match session user should return error",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: validTokenEncoding},
				SessionID:   "session-1",
				InstanceID:  "instance-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "session-1")))).
					Return(&domain.Session{ID: "session-1", UserID: "user-1"}, nil)
				return repo
			},
			encryptionAlgorithm: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				cryptoAlg := crypto.NewMockEncryptionAlgorithm(ctrl)
				cryptoAlg.EXPECT().
					EncryptionKeyID().
					Times(1)
				cryptoAlg.EXPECT().
					DecryptString(gomock.Any(), gomock.Any()).
					Times(1).
					Return("intent-123", nil)
				return cryptoAlg
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "user-1")))).
					Return(&domain.User{
						ID:    "user-1",
						State: domain.UserStateActive,
						Human: &domain.HumanUser{IdentityProviderLinks: []*domain.IdentityProviderLink{{ProviderID: "idp-123"}}},
					}, nil)
				return repo
			},
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "intent-123")))).
					Return(&domain.IDPIntent{ID: "intent-123", State: domain.IDPIntentStateSucceeded, ExpiresAt: gu.Ptr(time.Now().Add(time.Hour)), UserID: "user-2"}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-FLdnLH", "Errors.Intent.OtherUser"),
		},
		{
			testName: "when user has no matching IDP link should return error",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: validTokenEncoding},
				SessionID:   "session-1",
				InstanceID:  "instance-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "session-1")))).
					Return(&domain.Session{ID: "session-1", UserID: "user-1"}, nil)
				return repo
			},
			encryptionAlgorithm: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				cryptoAlg := crypto.NewMockEncryptionAlgorithm(ctrl)
				cryptoAlg.EXPECT().
					EncryptionKeyID().
					Times(1)
				cryptoAlg.EXPECT().
					DecryptString(gomock.Any(), gomock.Any()).
					Times(1).
					Return("intent-123", nil)
				return cryptoAlg
			},
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "intent-123")))).
					Return(&domain.IDPIntent{ID: "intent-123", State: domain.IDPIntentStateSucceeded, ExpiresAt: gu.Ptr(time.Now().Add(time.Hour)), UserID: "", IDPID: "idp-123"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "user-1")))).
					Return(&domain.User{
						ID:    "user-1",
						State: domain.UserStateActive,
						Human: &domain.HumanUser{IdentityProviderLinks: []*domain.IdentityProviderLink{{ProviderID: "idp-456"}}},
					}, nil)
				return repo
			},
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-XuNkt7", "Errors.Intent.OtherUser"),
		},
		{
			testName: "when all validations pass should return no error",
			cmd: &domain.IDPIntentCheckCommand{
				CheckIntent: &domain.CheckIDPIntentType{ID: "intent-123", Token: validTokenEncoding},
				SessionID:   "session-1",
				InstanceID:  "instance-1",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "session-1")))).
					Return(&domain.Session{ID: "session-1", UserID: "user-1"}, nil)
				return repo
			},
			encryptionAlgorithm: func(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
				cryptoAlg := crypto.NewMockEncryptionAlgorithm(ctrl)
				cryptoAlg.EXPECT().
					EncryptionKeyID().
					Times(1)
				cryptoAlg.EXPECT().
					DecryptString(gomock.Any(), gomock.Any()).
					Times(1).
					Return("intent-123", nil)
				return cryptoAlg
			},
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "intent-123")))).
					Return(&domain.IDPIntent{ID: "intent-123", State: domain.IDPIntentStateSucceeded, ExpiresAt: gu.Ptr(time.Now().Add(time.Hour)), UserID: "", IDPID: "idp-123"}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(database.WithCondition(repo.PrimaryKeyCondition("instance-1", "user-1")))).
					Return(&domain.User{
						ID:    "user-1",
						State: domain.UserStateActive,
						Human: &domain.HumanUser{
							IdentityProviderLinks: []*domain.IdentityProviderLink{
								{ProviderID: "idp-123"},
							},
						},
					}, nil)
				return repo
			},
			expectedError: nil,
			expectedUser: domain.User{
				ID:    "user-1",
				State: domain.UserStateActive,
				Human: &domain.HumanUser{
					IdentityProviderLinks: []*domain.IdentityProviderLink{
						{ProviderID: "idp-123"},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctrl := gomock.NewController(t)
			ctx := authz.NewMockContext(tc.cmd.InstanceID, "", "")

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
			if tc.idpIntentRepo != nil {
				domain.WithIDPIntentRepo(tc.idpIntentRepo(ctrl))(opts)
			}
			if tc.encryptionAlgorithm != nil {
				tc.cmd.EncAlgo = tc.encryptionAlgorithm(ctrl)
			}

			// Test
			err := tc.cmd.Validate(ctx, opts)

			// Verify
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedUser, tc.cmd.FetchedUser)
		})
	}
}

func TestIDPIntentCheckCommand_Execute(t *testing.T) {
	t.Parallel()
	deleteErr := errors.New("delete error")
	updateErr := errors.New("update error")

	tt := []struct {
		testName       string
		idpIntentRepo  func(ctrl *gomock.Controller) domain.IDPIntentRepository
		sessionRepo    func(ctrl *gomock.Controller) domain.SessionRepository
		cmd            *domain.IDPIntentCheckCommand
		expectedError  error
		expectComplete bool
	}{
		{
			testName:       "when CheckIntent is nil should return no error and not mark complete",
			cmd:            &domain.IDPIntentCheckCommand{},
			expectedError:  nil,
			expectComplete: false,
		},
		{
			testName: "when intent deletion fails should return error",
			cmd:      domain.NewIDPIntentCheckCommand(&domain.CheckIDPIntentType{ID: "intent-123"}, "session-1", "instance-1", nil),
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Delete(gomock.Any(), gomock.Any(), repo.PrimaryKeyCondition("instance-1", "intent-123")).
					Return(int64(0), deleteErr)
				return repo
			},
			expectedError:  zerrors.ThrowInternal(deleteErr, "DOM-j1s5Eu", "failed deleting IDP intent"),
			expectComplete: false,
		},
		{
			testName: "when intent deletion returns unexpected row count should return error",
			cmd:      domain.NewIDPIntentCheckCommand(&domain.CheckIDPIntentType{ID: "intent-123"}, "session-1", "instance-1", nil),
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Delete(gomock.Any(), gomock.Any(), repo.PrimaryKeyCondition("instance-1", "intent-123")).
					Return(int64(0), nil)
				return repo
			},
			expectedError:  zerrors.ThrowInternal(domain.NewRowsReturnedMismatchError(1, 0), "DOM-3CBpdB", "unexpected number of rows deleted"),
			expectComplete: false,
		},
		{
			testName: "when session update fails should return internal error",
			cmd:      domain.NewIDPIntentCheckCommand(&domain.CheckIDPIntentType{ID: "intent-123"}, "session-1", "instance-1", nil),
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Delete(gomock.Any(), gomock.Any(), repo.PrimaryKeyCondition("instance-1", "intent-123")).
					Return(int64(1), nil)
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), repo.PrimaryKeyCondition("instance-1", "session-1"), repo.SetFactor(&domain.SessionFactorIdentityProviderIntent{LastVerifiedAt: time.Now()})).
					Times(1).
					Return(int64(0), updateErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(updateErr, "DOM-pec0al", "failed updating session"),
		},
		{
			testName: "when session update returns no rows should return not found error",
			cmd:      domain.NewIDPIntentCheckCommand(&domain.CheckIDPIntentType{ID: "intent-123"}, "session-1", "instance-1", nil),
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Delete(gomock.Any(), gomock.Any(), repo.PrimaryKeyCondition("instance-1", "intent-123")).
					Return(int64(1), nil)
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), repo.PrimaryKeyCondition("instance-1", "session-1"), repo.SetFactor(&domain.SessionFactorIdentityProviderIntent{LastVerifiedAt: time.Now()})).
					Times(1).
					Return(int64(0), nil)
				return repo
			},
			expectedError: zerrors.ThrowNotFound(nil, "DOM-CopO4e", "session not found"),
		},
		{
			testName: "when session update returns no too many rows should return internal error",
			cmd:      domain.NewIDPIntentCheckCommand(&domain.CheckIDPIntentType{ID: "intent-123"}, "session-1", "instance-1", nil),
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Delete(gomock.Any(), gomock.Any(), repo.PrimaryKeyCondition("instance-1", "intent-123")).
					Return(int64(1), nil)
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), repo.PrimaryKeyCondition("instance-1", "session-1"), repo.SetFactor(&domain.SessionFactorIdentityProviderIntent{LastVerifiedAt: time.Now()})).
					Times(1).
					Return(int64(2), nil)
				return repo
			},
			expectedError: zerrors.ThrowInternal(domain.NewMultipleObjectsUpdatedError(1, 2), "DOM-mlbibw", "unexpected number of rows updated"),
		},
		{
			testName: "when session update returns 1 row updated should return no error",
			cmd:      domain.NewIDPIntentCheckCommand(&domain.CheckIDPIntentType{ID: "intent-123"}, "session-1", "instance-1", nil),
			idpIntentRepo: func(ctrl *gomock.Controller) domain.IDPIntentRepository {
				repo := domainmock.NewIDPIntentRepo(ctrl)
				repo.EXPECT().
					Delete(gomock.Any(), gomock.Any(), repo.PrimaryKeyCondition("instance-1", "intent-123")).
					Return(int64(1), nil)
				return repo
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), repo.PrimaryKeyCondition("instance-1", "session-1"), repo.SetFactor(&domain.SessionFactorIdentityProviderIntent{LastVerifiedAt: time.Now()})).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			expectComplete: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctrl := gomock.NewController(t)
			ctx := authz.NewMockContext(tc.cmd.InstanceID, "", "")

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)

			if tc.idpIntentRepo != nil {
				domain.WithIDPIntentRepo(tc.idpIntentRepo(ctrl))(opts)
			}
			if tc.sessionRepo != nil {
				domain.WithSessionRepo(tc.sessionRepo(ctrl))(opts)
			}

			// Test
			err := tc.cmd.Execute(ctx, opts)

			// Verify
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectComplete, tc.cmd.IsCheckComplete)
			if tc.expectedError == nil && tc.expectComplete {
				assert.NotZero(t, tc.cmd.IntentLastVerified)
			}
		})
	}
}
