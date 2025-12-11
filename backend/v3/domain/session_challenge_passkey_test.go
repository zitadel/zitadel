package domain_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
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
	legacy_domain "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestPasskeyChallengeCommand_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                    string
		sessionID               string
		instanceID              string
		RequestChallengePasskey *session_grpc.RequestChallenges_WebAuthN
		userRepo                func(ctrl *gomock.Controller) domain.UserRepository
		sessionRepo             func(ctrl *gomock.Controller) domain.SessionRepository
		wantErr                 error
		wantUser                *domain.User
	}{
		{
			name:                    "no request passkey challenge",
			RequestChallengePasskey: nil,
			wantErr:                 nil,
		},
		{
			name:                    "no session id",
			RequestChallengePasskey: &session_grpc.RequestChallenges_WebAuthN{},
			sessionID:               "",
			wantErr:                 zerrors.ThrowPreconditionFailed(nil, "DOM-EVo5yE", "missing session id"),
		},
		{
			name:                    "no session id",
			RequestChallengePasskey: &session_grpc.RequestChallenges_WebAuthN{},
			sessionID:               "session-id",
			instanceID:              "",
			wantErr:                 zerrors.ThrowPreconditionFailed(nil, "DOM-sh8xvQ", "missing instance id"),
		},
		{
			name:                    "failed to fetch session",
			RequestChallengePasskey: &session_grpc.RequestChallenges_WebAuthN{},
			sessionID:               "session-id",
			instanceID:              "instance-id",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					AnyTimes().
					Return(nil, assert.AnError)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-zy4hYC", "failed fetching session"),
		},
		{
			name:                    "session not found",
			RequestChallengePasskey: &session_grpc.RequestChallenges_WebAuthN{},
			sessionID:               "session-id",
			instanceID:              "instance-id",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					AnyTimes().
					Return(nil, new(database.NoRowFoundError))
				return repo
			},
			wantErr: zerrors.ThrowNotFound(new(database.NoRowFoundError), "DOM-zy4hYC", "session not found"),
		},
		{
			name:                    "session without user id",
			RequestChallengePasskey: &session_grpc.RequestChallenges_WebAuthN{},
			sessionID:               "session-id",
			instanceID:              "instance-id",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					AnyTimes().
					Return(&domain.Session{
						UserID: "",
					}, nil)
				return repo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-uVyrt2", "missing user id in session"),
		},
		{
			name:                    "failed to fetch user",
			RequestChallengePasskey: &session_grpc.RequestChallenges_WebAuthN{},
			sessionID:               "session-id",
			instanceID:              "instance-id",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					AnyTimes().
					Return(&domain.Session{
						UserID: "user-id",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					LoadPasskeys().
					AnyTimes().
					Return(repo)
				repo.EXPECT().
					Human().
					AnyTimes().
					Return(humanRepo)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								getUserIDCondition(repo, "user-id"),
							),
						),
						dbmock.QueryOptions(
							database.WithCondition(
								getPasskeyTypeCondition(humanRepo, domain.PasskeyTypeU2F),
							),
						)).
					AnyTimes().
					Return(nil, assert.AnError)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-8cGMtd", "failed fetching user"),
		},
		{
			name:                    "user not found",
			RequestChallengePasskey: &session_grpc.RequestChallenges_WebAuthN{},
			sessionID:               "session-id",
			instanceID:              "instance-id",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					AnyTimes().
					Return(&domain.Session{
						UserID: "user-id",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					LoadPasskeys().
					AnyTimes().
					Return(repo)
				repo.EXPECT().
					Human().
					AnyTimes().
					Return(humanRepo)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getUserIDCondition(repo, "user-id"),
						),
					),
						dbmock.QueryOptions(
							database.WithCondition(
								getPasskeyTypeCondition(humanRepo, domain.PasskeyTypeU2F),
							),
						)).
					AnyTimes().
					Return(nil, new(database.NoRowFoundError))
				return repo
			},
			wantErr: zerrors.ThrowNotFound(new(database.NoRowFoundError), "DOM-8cGMtd", "user not found"),
		},
		{
			name:                    "user not active",
			RequestChallengePasskey: &session_grpc.RequestChallenges_WebAuthN{},
			sessionID:               "session-id",
			instanceID:              "instance-id",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					AnyTimes().
					Return(&domain.Session{
						UserID: "user-id",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					LoadPasskeys().
					AnyTimes().
					Return(repo)
				repo.EXPECT().
					Human().
					AnyTimes().
					Return(humanRepo)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getUserIDCondition(repo, "user-id"),
						),
					),
						dbmock.QueryOptions(
							database.WithCondition(
								getPasskeyTypeCondition(humanRepo, domain.PasskeyTypeU2F),
							),
						)).
					AnyTimes().
					Return(&domain.User{
						State: domain.UserStateInactive,
					}, nil)
				return repo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-bnxBdS", "user not active"),
		},
		{
			name:                    "valid request passkey challenge",
			RequestChallengePasskey: &session_grpc.RequestChallenges_WebAuthN{},
			sessionID:               "session-id",
			instanceID:              "instance-id",
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					AnyTimes().
					Return(&domain.Session{
						UserID: "user-id",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					LoadPasskeys().
					AnyTimes().
					Return(repo)
				repo.EXPECT().
					Human().
					AnyTimes().
					Return(humanRepo)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getUserIDCondition(repo, "user-id"),
						),
					),
						dbmock.QueryOptions(
							database.WithCondition(
								getPasskeyTypeCondition(humanRepo, domain.PasskeyTypeU2F),
							),
						)).
					AnyTimes().
					Return(&domain.User{
						ID:             "user-id",
						InstanceID:     "instance-id",
						OrganizationID: "organization-id",
						Username:       "username",
						Human: &domain.HumanUser{
							Passkeys: []*domain.Passkey{
								{
									ID:        "passkey-id",
									PublicKey: []byte("public-key"),
									Name:      "My Passkey",
									Type:      domain.PasskeyTypePasswordless,
								},
							},
							Email: domain.HumanEmail{
								Address: "user@example.com",
							},
						},
						State: domain.UserStateActive,
					}, nil)
				return repo
			},
			wantUser: &domain.User{
				ID:             "user-id",
				InstanceID:     "instance-id",
				OrganizationID: "organization-id",
				Username:       "username",
				Human: &domain.HumanUser{
					Passkeys: []*domain.Passkey{
						{
							ID:        "passkey-id",
							PublicKey: []byte("public-key"),
							Name:      "My Passkey",
							Type:      domain.PasskeyTypePasswordless,
						},
					},
					Email: domain.HumanEmail{
						Address: "user@example.com",
					},
				},
				State: domain.UserStateActive,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewPasskeyChallengeCommand(
				tt.sessionID,
				tt.instanceID,
				tt.RequestChallengePasskey,
				nil,
			)
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
			err := cmd.Validate(ctx, opts)
			assert.Equal(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.wantUser, cmd.User)
			}
		})
	}
}

func TestPasskeyChallengeCommand_Events(t *testing.T) {
	t.Parallel()
	ctx := authz.NewMockContext("instance-id", "", "")
	challengedAt := time.Now()
	tests := []struct {
		name                    string
		requestChallengePasskey *session_grpc.RequestChallenges_WebAuthN
		challengePasskey        *domain.SessionChallengePasskey
		wantErr                 error
		wantEvent               eventstore.Command
	}{
		{
			name:                    "no request passkey challenge",
			requestChallengePasskey: nil,
			wantErr:                 zerrors.ThrowInternal(nil, "DOM-MALUxr", "failed to push WebAuthN challenged event"),
			wantEvent:               nil,
		},
		{
			name:                    "valid request passkey challenge",
			requestChallengePasskey: &session_grpc.RequestChallenges_WebAuthN{},
			challengePasskey: &domain.SessionChallengePasskey{
				Challenge:        "challenge",
				LastChallengedAt: challengedAt,
				RPID:             "rpID",
				UserVerification: legacy_domain.UserVerificationRequirementPreferred,
			},
			wantEvent: session.NewWebAuthNChallengedEvent(ctx,
				&session.NewAggregate("session-id", "instance-id").Aggregate,
				"challenge",
				nil,
				legacy_domain.UserVerificationRequirementPreferred,
				"rpID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cmd := domain.NewPasskeyChallengeCommand(
				"session-id",
				"instance-id",
				&session_grpc.RequestChallenges_WebAuthN{},
				nil,
			)
			cmd.ChallengePasskey = tt.challengePasskey
			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			events, err := cmd.Events(ctx, opts)
			assert.Equal(t, tt.wantErr, err)
			if tt.wantEvent != nil {
				require.Len(t, events, 1)
				assert.Equal(t, tt.wantEvent, events[0])
			}
		})
	}
}

func TestPasskeyChallengeCommand_Execute(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                    string
		requestChallengePasskey *session_grpc.RequestChallenges_WebAuthN
		sessionRepo             func(ctrl *gomock.Controller) domain.SessionRepository
		user                    *domain.User
		session                 *domain.Session
		webAuthNBeginLogin      func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error)
		wantErr                 error
	}{
		{
			name:                    "no request passkey challenge",
			requestChallengePasskey: nil,
			wantErr:                 nil,
		},
		{
			name:                    "failed to begin webauthn login",
			requestChallengePasskey: &session_grpc.RequestChallenges_WebAuthN{},
			user: &domain.User{
				Human: &domain.HumanUser{},
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return nil, nil, "", assert.AnError
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-Fy333Q", "failed to begin webauthn login"),
		},
		{
			name:                    "invalid credential assertion data",
			requestChallengePasskey: &session_grpc.RequestChallenges_WebAuthN{},
			user: &domain.User{
				Human: &domain.HumanUser{},
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "challenge",
				}, []byte("invalid"), "rpID", nil
			},
			wantErr: zerrors.ThrowInternal(nil, "DOM-liSCA4", "failed to unmarshal credential assertion data"),
		},
		{
			name: "failed to update session",
			requestChallengePasskey: &session_grpc.RequestChallenges_WebAuthN{
				UserVerificationRequirement: 2,
			},
			user: &domain.User{
				Human: &domain.HumanUser{
					Passkeys: []*domain.Passkey{
						{
							ID:    "passkey-1",
							KeyID: []byte("key-id-1"),
						},
						{
							ID:    "passkey-2",
							KeyID: []byte("key-id-2"),
						},
					},
				},
			},
			session: &domain.Session{
				ID: "session-id",
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, []byte(`{"publicKey":{"challenge":"Y2hhbGxlbmdl","timeout":60000,"rpId":"example.com","allowCredentials":[{"type":"public-key","id":"cGFzc2tleS1pZA"}],"userVerification":"preferred"}}`), "example.com", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				expectedPasskeyChallenge := &domain.SessionChallengePasskey{
					Challenge:        "Y2hhbGxlbmdl",
					RPID:             "example.com",
					UserVerification: legacy_domain.UserVerificationRequirementPreferred,
				}
				repo.EXPECT().
					SetChallenge(gomock.Any()).
					AnyTimes().
					DoAndReturn(assertPasskeyChallengeChange(t, expectedPasskeyChallenge))
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					AnyTimes().
					Return(int64(0), assert.AnError)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-yd3f4", "failed updating session"),
		},
		{
			name: "failed to update session - no rows updated",
			requestChallengePasskey: &session_grpc.RequestChallenges_WebAuthN{
				UserVerificationRequirement: 2, // to set userVerification to "preferred"
			},
			user: &domain.User{
				Human: &domain.HumanUser{
					Passkeys: []*domain.Passkey{
						{
							ID:    "passkey-1",
							KeyID: []byte("key-id-1"),
						},
						{
							ID:    "passkey-2",
							KeyID: []byte("key-id-2"),
						},
					},
				},
			},
			session: &domain.Session{
				ID: "session-id",
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, []byte(`{"publicKey":{"challenge":"Y2hhbGxlbmdl","timeout":60000,"rpId":"example.com","allowCredentials":[{"type":"public-key","id":"cGFzc2tleS1pZA"}],"userVerification":"preferred"}}`), "example.com", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				expectedPasskeyChallenge := &domain.SessionChallengePasskey{
					Challenge:        "Y2hhbGxlbmdl",
					RPID:             "example.com",
					UserVerification: legacy_domain.UserVerificationRequirementPreferred,
				}
				repo.EXPECT().
					SetChallenge(gomock.Any()).
					AnyTimes().
					DoAndReturn(assertPasskeyChallengeChange(t, expectedPasskeyChallenge))
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					AnyTimes().
					Return(int64(0), nil)
				return repo
			},
			wantErr: zerrors.ThrowNotFound(nil, "DOM-yd3f4", "session not found"),
		},
		{
			name: "failed to update session - more than 1 row updated",
			requestChallengePasskey: &session_grpc.RequestChallenges_WebAuthN{
				UserVerificationRequirement: 2, // to set userVerification to "preferred"
			},
			user: &domain.User{
				Human: &domain.HumanUser{
					Passkeys: []*domain.Passkey{
						{
							ID:    "passkey-1",
							KeyID: []byte("key-id-1"),
						},
						{
							ID:    "passkey-2",
							KeyID: []byte("key-id-2"),
						},
					},
				},
			},
			session: &domain.Session{
				ID: "session-id",
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, []byte(`{"publicKey":{"challenge":"Y2hhbGxlbmdl","timeout":60000,"rpId":"example.com","allowCredentials":[{"type":"public-key","id":"cGFzc2tleS1pZA"}],"userVerification":"preferred"}}`), "example.com", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				expectedPasskeyChallenge := &domain.SessionChallengePasskey{
					Challenge:        "Y2hhbGxlbmdl",
					RPID:             "example.com",
					UserVerification: legacy_domain.UserVerificationRequirementPreferred,
				}
				repo.EXPECT().
					SetChallenge(gomock.Any()).
					AnyTimes().
					DoAndReturn(assertPasskeyChallengeChange(t, expectedPasskeyChallenge))
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					AnyTimes().
					Return(int64(2), nil)
				return repo
			},
			wantErr: zerrors.ThrowInternal(domain.NewMultipleObjectsUpdatedError(1, 2), "DOM-yd3f4", "unexpected number of rows updated"),
		},
		{
			name: "session updated successfully with the passkey challenge",
			requestChallengePasskey: &session_grpc.RequestChallenges_WebAuthN{
				UserVerificationRequirement: 2, // to set userVerification to "preferred"
			},
			user: &domain.User{
				Human: &domain.HumanUser{
					Passkeys: []*domain.Passkey{
						{
							ID:    "passkey-1",
							KeyID: []byte("key-id-1"),
						},
						{
							ID:    "passkey-2",
							KeyID: []byte("key-id-2"),
						},
					},
				},
			},
			session: &domain.Session{
				ID: "session-id",
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, []byte(`{"publicKey":{"challenge":"Y2hhbGxlbmdl","timeout":60000,"rpId":"example.com","allowCredentials":[{"type":"public-key","id":"cGFzc2tleS1pZA"}],"userVerification":"preferred"}}`), "example.com", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				expectedPasskeyChallenge := &domain.SessionChallengePasskey{
					Challenge:        "Y2hhbGxlbmdl",
					RPID:             "example.com",
					UserVerification: legacy_domain.UserVerificationRequirementPreferred,
				}
				repo.EXPECT().
					SetChallenge(gomock.Any()).
					AnyTimes().
					DoAndReturn(assertPasskeyChallengeChange(t, expectedPasskeyChallenge))
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					AnyTimes().
					Return(int64(1), nil)
				return repo
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := authz.NewMockContext("instance-id", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewPasskeyChallengeCommand(
				"session-id",
				"instance-id",
				tt.requestChallengePasskey,
				tt.webAuthNBeginLogin,
			)
			cmd.User = tt.user
			cmd.Session = tt.session
			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			if tt.sessionRepo != nil {
				domain.WithSessionRepo(tt.sessionRepo(ctrl))(opts)
			}
			err := cmd.Execute(ctx, opts)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func assertPasskeyChallengeChange(t *testing.T, expectedPasskeyChallenge *domain.SessionChallengePasskey) func(challenge *domain.SessionChallengePasskey) database.Change {
	return func(challenge *domain.SessionChallengePasskey) database.Change {
		assert.Equal(t, expectedPasskeyChallenge.Challenge, challenge.Challenge)
		assert.Equal(t, expectedPasskeyChallenge.RPID, challenge.RPID)
		assert.Equal(t, expectedPasskeyChallenge.UserVerification, challenge.UserVerification)
		return database.NewChanges(
			database.NewChange(
				database.NewColumn("zitadel.sessions", "passkey_challenge_challenge"), challenge.Challenge,
			),
			database.NewChange(
				database.NewColumn("zitadel.sessions", "passkey_challenge_rp_id"), challenge.RPID,
			),
			database.NewChange(
				database.NewColumn("zitadel.sessions", "passkey_challenge_user_verification"), challenge.UserVerification,
			),
		)
	}
}
