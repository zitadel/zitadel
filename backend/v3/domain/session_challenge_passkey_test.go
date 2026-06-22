package domain_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	"github.com/zitadel/zitadel/internal/api/authz"
	old_domain "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestNewPasskeyChallengeCommand(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                 string
		sessionID            string
		instanceID           string
		challengeTypePasskey *domain.ChallengeTypePasskey
		webAuthNBeginLogin   func(
			ctx context.Context,
			user webauthn.User,
			rpID string,
			userVerification protocol.UserVerificationRequirement,
		) (sessionData *webauthn.SessionData,
			cred []byte,
			relyingPartyID string,
			err error)
		wantCmd *domain.PasskeyChallengeCommand
		wantErr error
	}{
		{
			name:       "missing begin webauthn login function",
			sessionID:  "session-1",
			instanceID: "instance-1",
			challengeTypePasskey: &domain.ChallengeTypePasskey{
				Domain:                      "example.com",
				UserVerificationRequirement: old_domain.UserVerificationRequirementRequired,
			},
			wantErr: zerrors.ThrowInternal(nil, "DOM-jwk5Pe", "begin webauthn login function not set"),
		},
		{
			name:       "valid passkey challenge command",
			sessionID:  "session-1",
			instanceID: "instance-1",
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, []byte("assertion data"), "example.com", nil
			},
			challengeTypePasskey: &domain.ChallengeTypePasskey{
				Domain:                      "example.com",
				UserVerificationRequirement: old_domain.UserVerificationRequirementRequired,
			},
			wantCmd: &domain.PasskeyChallengeCommand{
				InstanceID: "instance-1",
				SessionID:  "session-1",
				ChallengeTypePasskey: &domain.ChallengeTypePasskey{
					Domain:                      "example.com",
					UserVerificationRequirement: old_domain.UserVerificationRequirementRequired,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := domain.NewPasskeyChallengeCommand(tt.sessionID, tt.instanceID, tt.challengeTypePasskey, tt.webAuthNBeginLogin)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantCmd.SessionID, cmd.SessionID)
			assert.Equal(t, tt.wantCmd.InstanceID, cmd.InstanceID)
			assert.Equal(t, tt.wantCmd.ChallengeTypePasskey, cmd.ChallengeTypePasskey)
		})
	}
}

func TestPasskeyChallengeCommand_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		sessionID          string
		instanceID         string
		webAuthNBeginLogin func(
			ctx context.Context,
			user webauthn.User,
			rpID string,
			userVerification protocol.UserVerificationRequirement,
		) (sessionData *webauthn.SessionData,
			cred []byte,
			relyingPartyID string,
			err error)
		challengeTypePasskey *domain.ChallengeTypePasskey
		userRepo             func(ctrl *gomock.Controller) domain.UserRepository
		sessionRepo          func(ctrl *gomock.Controller) domain.SessionRepository
		wantErr              error
	}{
		{
			name: "no request passkey challenge",
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, []byte("assertion data"), "example.com", nil
			},
		},
		{
			name:                 "no session id",
			challengeTypePasskey: &domain.ChallengeTypePasskey{},
			sessionID:            "",
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, []byte("assertion data"), "example.com", nil
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-EVo5yE", "Errors.Missing.SessionID"),
		},
		{
			name:                 "no instance id",
			challengeTypePasskey: &domain.ChallengeTypePasskey{},
			sessionID:            "session-1",
			instanceID:           "",
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, []byte("assertion data"), "example.com", nil
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-sh8xvQ", "Errors.Missing.InstanceID"),
		},
		{
			name:                 "failed to fetch session",
			challengeTypePasskey: &domain.ChallengeTypePasskey{},
			sessionID:            "session-1",
			instanceID:           "instance-1",
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, []byte("assertion data"), "example.com", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							repo.PrimaryKeyCondition("instance-1", "session-1"),
						),
					)).
					Times(1).
					Return(nil, assert.AnError)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-zy4hYC", "failed fetching Session"),
		},
		{
			name:                 "session not found",
			challengeTypePasskey: &domain.ChallengeTypePasskey{},
			sessionID:            "session-1",
			instanceID:           "instance-1",
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, []byte("assertion data"), "example.com", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							repo.PrimaryKeyCondition("instance-1", "session-1"),
						),
					)).
					Times(1).
					Return(nil, new(database.NoRowFoundError))
				return repo
			},
			wantErr: zerrors.ThrowNotFound(new(database.NoRowFoundError), "DOM-zy4hYC", "Session not found"),
		},
		{
			name:                 "session without user id",
			challengeTypePasskey: &domain.ChallengeTypePasskey{},
			sessionID:            "session-1",
			instanceID:           "instance-1",
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, []byte("assertion data"), "example.com", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							repo.PrimaryKeyCondition("instance-1", "session-1"),
						),
					)).
					Times(1).
					Return(&domain.Session{
						UserID: "",
					}, nil)
				return repo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-uVyrt2", "Errors.Missing.Session.UserID"),
		},
		{
			name:                 "failed to fetch user",
			challengeTypePasskey: &domain.ChallengeTypePasskey{},
			sessionID:            "session-1",
			instanceID:           "instance-1",
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, []byte("assertion data"), "example.com", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)

				humanRepo := domainmock.NewHumanRepo(ctrl)
				repo.EXPECT().Human().Times(1).Return(humanRepo)

				repo.EXPECT().
					Get(gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								repo.PrimaryKeyCondition("instance-1", "user-1"),
							),
						),
						dbmock.QueryOptions(
							database.WithCondition(
								humanRepo.PasskeyConditions().TypeCondition(database.TextOperationEqual, domain.PasskeyTypeU2F),
							),
						),
					).
					Times(1).
					Return(nil, assert.AnError)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-8cGMtd", "failed fetching User"),
		},
		{
			name:                 "user not found",
			challengeTypePasskey: &domain.ChallengeTypePasskey{},
			sessionID:            "session-1",
			instanceID:           "instance-1",
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, []byte("assertion data"), "example.com", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				humanRepo := domainmock.NewHumanRepo(ctrl)
				repo.EXPECT().Human().Times(1).Return(humanRepo)

				repo.EXPECT().
					Get(gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								repo.PrimaryKeyCondition("instance-1", "user-1"),
							),
						),
						dbmock.QueryOptions(
							database.WithCondition(
								humanRepo.PasskeyConditions().TypeCondition(database.TextOperationEqual, domain.PasskeyTypeU2F),
							),
						),
					).
					Times(1).
					Return(nil, new(database.NoRowFoundError))
				return repo
			},
			wantErr: zerrors.ThrowNotFound(new(database.NoRowFoundError), "DOM-8cGMtd", "User not found"),
		},
		{
			name:                 "user not human",
			challengeTypePasskey: &domain.ChallengeTypePasskey{},
			sessionID:            "session-1",
			instanceID:           "instance-1",
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, []byte("assertion data"), "example.com", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				humanRepo := domainmock.NewHumanRepo(ctrl)
				repo.EXPECT().Human().Times(1).Return(humanRepo)

				repo.EXPECT().
					Get(gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								repo.PrimaryKeyCondition("instance-1", "user-1"),
							),
						),
						dbmock.QueryOptions(
							database.WithCondition(
								humanRepo.PasskeyConditions().TypeCondition(database.TextOperationEqual, domain.PasskeyTypeU2F),
							),
						),
					).
					Times(1).
					Return(&domain.User{}, nil)
				return repo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-nd3f4", "Errors.User.NotHuman"),
		},
		{
			name:                 "user not active",
			challengeTypePasskey: &domain.ChallengeTypePasskey{},
			sessionID:            "session-1",
			instanceID:           "instance-1",
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, []byte("assertion data"), "example.com", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				humanRepo := domainmock.NewHumanRepo(ctrl)
				repo.EXPECT().Human().Times(1).Return(humanRepo)

				repo.EXPECT().
					Get(gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								repo.PrimaryKeyCondition("instance-1", "user-1"),
							),
						),
						dbmock.QueryOptions(
							database.WithCondition(
								humanRepo.PasskeyConditions().TypeCondition(database.TextOperationEqual, domain.PasskeyTypeU2F),
							),
						),
					).
					Times(1).
					Return(&domain.User{
						State: domain.UserStateInactive,
						Human: &domain.HumanUser{},
					}, nil)
				return repo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-bnxBdS", "Errors.User.NotFound"),
		},
		{
			name:                 "valid request passkey challenge",
			challengeTypePasskey: &domain.ChallengeTypePasskey{},
			sessionID:            "session-1",
			instanceID:           "instance-1",
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, []byte("assertion data"), "example.com", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				return getUserByPasskeyType(ctrl, domain.PasskeyTypeU2F)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd, err := domain.NewPasskeyChallengeCommand(
				tt.sessionID,
				tt.instanceID,
				tt.challengeTypePasskey,
				tt.webAuthNBeginLogin,
			)
			require.NoError(t, err)

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			if tt.userRepo != nil {
				domain.WithUserRepo(tt.userRepo(ctrl))(opts)
			}
			if tt.sessionRepo != nil {
				domain.WithSessionRepo(tt.sessionRepo(ctrl))(opts)
			}
			err = cmd.Validate(ctx, opts)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestPasskeyChallengeCommand_Execute(t *testing.T) {
	t.Parallel()
	challengedAt := time.Now()
	credentialAssertionDataRequired := []byte(`{"publicKey":{"challenge":"Y2hhbGxlbmdl","timeout":60000,"rpId":"example.com","allowCredentials":[{"type":"public-key","id":"cGFzc2tleS1pZA"}],"userVerification":"required"}}`)
	credentialAssertionDataPreferred := []byte(`{"publicKey":{"challenge":"Y2hhbGxlbmdl","timeout":60000,"rpId":"example.com","allowCredentials":[{"type":"public-key","id":"cGFzc2tleS1pZA"}],"userVerification":"preferred"}}`)

	tests := []struct {
		name                 string
		challengeTypePasskey *domain.ChallengeTypePasskey
		sessionRepo          func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo             func(ctrl *gomock.Controller) domain.UserRepository
		webAuthNBeginLogin   func(
			ctx context.Context,
			user webauthn.User,
			rpID string,
			userVerification protocol.UserVerificationRequirement,
		) (sessionData *webauthn.SessionData,
			cred []byte,
			relyingPartyID string,
			err error)
		wantErr               error
		wantWebAuthNChallenge *session_grpc.Challenges_WebAuthN
	}{
		{
			name: "no request passkey challenge",
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, credentialAssertionDataPreferred, "example.com", nil
			},
		},
		{
			name: "failed to begin webauthn login",
			challengeTypePasskey: &domain.ChallengeTypePasskey{
				UserVerificationRequirement: old_domain.UserVerificationRequirementPreferred,
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return nil, nil, "", assert.AnError
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				return getUserByPasskeyType(ctrl, domain.PasskeyTypeU2F)
			},
			wantErr: assert.AnError,
		},
		{
			name:                 "invalid credential assertion data",
			challengeTypePasskey: &domain.ChallengeTypePasskey{},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "challenge",
				}, []byte("invalid"), "rpID", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				return getUserByPasskeyType(ctrl, domain.PasskeyTypeU2F)
			},
			wantErr: zerrors.ThrowInternal(nil, "DOM-liSCA4", "Errors.Unmarshal"),
		},
		{
			name: "failed to update session",
			challengeTypePasskey: &domain.ChallengeTypePasskey{
				UserVerificationRequirement: old_domain.UserVerificationRequirementPreferred,
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, credentialAssertionDataPreferred, "example.com", nil
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				return getUserByPasskeyType(ctrl, domain.PasskeyTypeU2F)
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)

				challenge := repo.SetChallenge(&domain.SessionChallengePasskey{
					Challenge:        "Y2hhbGxlbmdl",
					RPID:             "example.com",
					UserVerification: old_domain.UserVerificationRequirementPreferred,
					LastChallengedAt: time.Now(),
				})
				updateSessionFailedExpectation(repo, challenge, assert.AnError, 0)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-yd3f4", "failed updating Session"),
		},
		{
			name: "failed to update session - no rows updated",
			challengeTypePasskey: &domain.ChallengeTypePasskey{
				UserVerificationRequirement: old_domain.UserVerificationRequirementPreferred,
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, credentialAssertionDataPreferred, "example.com", nil
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				return getUserByPasskeyType(ctrl, domain.PasskeyTypeU2F)
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)

				challenge := repo.SetChallenge(&domain.SessionChallengePasskey{
					Challenge:        "Y2hhbGxlbmdl",
					RPID:             "example.com",
					UserVerification: old_domain.UserVerificationRequirementPreferred,
					LastChallengedAt: challengedAt,
				})
				updateSessionFailedExpectation(repo, challenge, nil, 0)
				return repo
			},
			wantErr: zerrors.ThrowNotFound(nil, "DOM-yd3f4", "Session not found"),
		},
		{
			name: "failed to update session - more than 1 row updated",
			challengeTypePasskey: &domain.ChallengeTypePasskey{
				UserVerificationRequirement: old_domain.UserVerificationRequirementPreferred,
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, credentialAssertionDataPreferred, "example.com", nil
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				return getUserByPasskeyType(ctrl, domain.PasskeyTypeU2F)
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)

				challenge := repo.SetChallenge(&domain.SessionChallengePasskey{
					Challenge:        "Y2hhbGxlbmdl",
					RPID:             "example.com",
					UserVerification: old_domain.UserVerificationRequirementPreferred,
					LastChallengedAt: challengedAt,
				})
				updateSessionFailedExpectation(repo, challenge, nil, 2)
				return repo
			},
			wantErr: zerrors.ThrowInternal(nil, "DOM-yd3f4", "unexpected number of rows updated"),
		},
		{
			name: "session updated successfully with the passkey challenge - required user verification",
			challengeTypePasskey: &domain.ChallengeTypePasskey{
				UserVerificationRequirement: old_domain.UserVerificationRequirementRequired,
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, credentialAssertionDataRequired, "example.com", nil
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				return getUserByPasskeyType(ctrl, domain.PasskeyTypePasswordless)
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				challenge := repo.SetChallenge(&domain.SessionChallengePasskey{
					Challenge:        "Y2hhbGxlbmdl",
					RPID:             "example.com",
					UserVerification: old_domain.UserVerificationRequirementRequired,
					LastChallengedAt: time.Now(),
				})
				updateSessionSucceededExpectation(repo, challenge)
				return repo
			},
			wantWebAuthNChallenge: getChallengePasskey(t, credentialAssertionDataRequired),
		},
		{
			name: "session updated successfully with the passkey challenge - preferred user verification",
			challengeTypePasskey: &domain.ChallengeTypePasskey{
				UserVerificationRequirement: old_domain.UserVerificationRequirementPreferred,
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, credentialAssertionDataPreferred, "example.com", nil
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				return getUserByPasskeyType(ctrl, domain.PasskeyTypeU2F)
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				challenge := repo.SetChallenge(&domain.SessionChallengePasskey{
					Challenge:        "Y2hhbGxlbmdl",
					RPID:             "example.com",
					UserVerification: old_domain.UserVerificationRequirementPreferred,
					LastChallengedAt: time.Now(),
				})
				updateSessionSucceededExpectation(repo, challenge)
				return repo
			},
			wantWebAuthNChallenge: getChallengePasskey(t, credentialAssertionDataPreferred),
		},
		{
			name: "session updated successfully with the passkey challenge - unspecified user verification",
			challengeTypePasskey: &domain.ChallengeTypePasskey{
				UserVerificationRequirement: old_domain.UserVerificationRequirementUnspecified,
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, credentialAssertionDataPreferred, "example.com", nil
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				return getUserByPasskeyType(ctrl, domain.PasskeyTypeU2F)
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				challenge := repo.SetChallenge(&domain.SessionChallengePasskey{
					Challenge:        "Y2hhbGxlbmdl",
					RPID:             "example.com",
					UserVerification: old_domain.UserVerificationRequirementUnspecified,
					LastChallengedAt: time.Now(),
				})
				updateSessionSucceededExpectation(repo, challenge)
				return repo
			},
			wantWebAuthNChallenge: getChallengePasskey(t, credentialAssertionDataPreferred),
		},
		{
			name: "session updated successfully with the passkey challenge - discouraged user verification",
			challengeTypePasskey: &domain.ChallengeTypePasskey{
				UserVerificationRequirement: old_domain.UserVerificationRequirementDiscouraged,
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, credentialAssertionDataPreferred, "example.com", nil
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				return getUserByPasskeyType(ctrl, domain.PasskeyTypeU2F)
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				challenge := repo.SetChallenge(&domain.SessionChallengePasskey{
					Challenge:        "Y2hhbGxlbmdl",
					RPID:             "example.com",
					UserVerification: old_domain.UserVerificationRequirementDiscouraged,
					LastChallengedAt: time.Now(),
				})
				updateSessionSucceededExpectation(repo, challenge)
				return repo
			},
			wantWebAuthNChallenge: getChallengePasskey(t, credentialAssertionDataPreferred),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd, err := domain.NewPasskeyChallengeCommand(
				"session-1",
				"instance-1",
				tt.challengeTypePasskey,
				tt.webAuthNBeginLogin,
			)
			require.NoError(t, err)

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			if tt.sessionRepo != nil {
				domain.WithSessionRepo(tt.sessionRepo(ctrl))(opts)
			}
			if tt.userRepo != nil {
				domain.WithUserRepo(tt.userRepo(ctrl))(opts)
			}

			// to fetch/validate session and user before calling Execute
			err = cmd.Validate(ctx, opts)
			require.NoError(t, err)

			err = cmd.Execute(ctx, opts)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantWebAuthNChallenge, cmd.GetWebAuthNChallenge())
		})
	}
}

func TestPasskeyChallengeCommand_Events(t *testing.T) {
	t.Parallel()
	ctx := authz.NewMockContext("instance-1", "", "")
	challengedAt := time.Now()

	tests := []struct {
		name                 string
		challengeTypePasskey *domain.ChallengeTypePasskey
		sessionRepo          func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo             func(ctrl *gomock.Controller) domain.UserRepository
		webAuthNBeginLogin   func(
			ctx context.Context,
			user webauthn.User,
			rpID string,
			userVerification protocol.UserVerificationRequirement,
		) (sessionData *webauthn.SessionData,
			cred []byte,
			relyingPartyID string,
			err error)
		wantEvent eventstore.Command
		wantErr   error
	}{
		{
			name: "no request passkey challenge",
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{
					Challenge: "Y2hhbGxlbmdl",
				}, []byte("assertion data"), "example.com", nil
			},
		},
		{
			name: "valid request passkey challenge - preferred user verification",
			challengeTypePasskey: &domain.ChallengeTypePasskey{
				UserVerificationRequirement: old_domain.UserVerificationRequirementPreferred,
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{Challenge: "Y2hhbGxlbmdl"},
					[]byte(`{"publicKey":{"challenge":"Y2hhbGxlbmdl","timeout":60000,"rpId":"example.com","allowCredentials":[{"type":"public-key","id":"cGFzc2tleS1pZA"}],"userVerification":"preferred"}}`),
					"example.com",
					nil
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				return getUserByPasskeyType(ctrl, domain.PasskeyTypeU2F)
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				challenge := repo.SetChallenge(&domain.SessionChallengePasskey{
					Challenge:        "Y2hhbGxlbmdl",
					RPID:             "example.com",
					UserVerification: old_domain.UserVerificationRequirementPreferred,
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challenge)
				return repo
			},
			wantEvent: session.NewWebAuthNChallengedEvent(ctx,
				&session.NewAggregate("session-1", "instance-1").Aggregate,
				"Y2hhbGxlbmdl",
				nil,
				old_domain.UserVerificationRequirementPreferred,
				"example.com"),
		},
		{
			name: "valid request passkey challenge - required user verification",
			challengeTypePasskey: &domain.ChallengeTypePasskey{
				UserVerificationRequirement: old_domain.UserVerificationRequirementRequired,
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{Challenge: "Y2hhbGxlbmdl"},
					[]byte(`{"publicKey":{"challenge":"Y2hhbGxlbmdl","timeout":60000,"rpId":"example.com","allowCredentials":[{"type":"public-key","id":"cGFzc2tleS1pZA"}],"userVerification":"required"}}`),
					"example.com",
					nil
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				return getUserByPasskeyType(ctrl, domain.PasskeyTypePasswordless)
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				challenge := repo.SetChallenge(&domain.SessionChallengePasskey{
					Challenge:        "Y2hhbGxlbmdl",
					RPID:             "example.com",
					UserVerification: old_domain.UserVerificationRequirementRequired,
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challenge)
				return repo
			},
			wantEvent: session.NewWebAuthNChallengedEvent(ctx,
				&session.NewAggregate("session-1", "instance-1").Aggregate,
				"Y2hhbGxlbmdl",
				nil,
				old_domain.UserVerificationRequirementRequired,
				"example.com"),
		},
		{
			name: "valid request passkey challenge - unspecified user verification",
			challengeTypePasskey: &domain.ChallengeTypePasskey{
				UserVerificationRequirement: old_domain.UserVerificationRequirementUnspecified,
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{Challenge: "Y2hhbGxlbmdl"},
					[]byte(`{"publicKey":{"challenge":"Y2hhbGxlbmdl","timeout":60000,"rpId":"example.com","allowCredentials":[{"type":"public-key","id":"cGFzc2tleS1pZA"}],"userVerification":"preferred"}}`),
					"example.com",
					nil
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				return getUserByPasskeyType(ctrl, domain.PasskeyTypeU2F)
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				challenge := repo.SetChallenge(&domain.SessionChallengePasskey{
					Challenge:        "Y2hhbGxlbmdl",
					RPID:             "example.com",
					UserVerification: old_domain.UserVerificationRequirementUnspecified,
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challenge)
				return repo
			},
			wantEvent: session.NewWebAuthNChallengedEvent(ctx,
				&session.NewAggregate("session-1", "instance-1").Aggregate,
				"Y2hhbGxlbmdl",
				nil,
				old_domain.UserVerificationRequirementUnspecified,
				"example.com"),
		},
		{
			name: "valid request passkey challenge - discouraged user verification",
			challengeTypePasskey: &domain.ChallengeTypePasskey{
				UserVerificationRequirement: old_domain.UserVerificationRequirementDiscouraged,
			},
			webAuthNBeginLogin: func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error) {
				return &webauthn.SessionData{Challenge: "Y2hhbGxlbmdl"},
					[]byte(`{"publicKey":{"challenge":"Y2hhbGxlbmdl","timeout":60000,"rpId":"example.com","allowCredentials":[{"type":"public-key","id":"cGFzc2tleS1pZA"}],"userVerification":"preferred"}}`),
					"example.com",
					nil
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				return getUserByPasskeyType(ctrl, domain.PasskeyTypeU2F)
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				challenge := repo.SetChallenge(&domain.SessionChallengePasskey{
					Challenge:        "Y2hhbGxlbmdl",
					RPID:             "example.com",
					UserVerification: old_domain.UserVerificationRequirementDiscouraged,
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challenge)
				return repo
			},
			wantEvent: session.NewWebAuthNChallengedEvent(ctx,
				&session.NewAggregate("session-1", "instance-1").Aggregate,
				"Y2hhbGxlbmdl",
				nil,
				old_domain.UserVerificationRequirementDiscouraged,
				"example.com"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd, err := domain.NewPasskeyChallengeCommand(
				"session-1",
				"instance-1",
				tt.challengeTypePasskey,
				tt.webAuthNBeginLogin,
			)
			require.NoError(t, err)

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			if tt.sessionRepo != nil {
				domain.WithSessionRepo(tt.sessionRepo(ctrl))(opts)
			}
			if tt.userRepo != nil {
				domain.WithUserRepo(tt.userRepo(ctrl))(opts)
			}

			// to fetch/validate session and user before calling Execute
			err = cmd.Validate(ctx, opts)
			require.NoError(t, err)

			// to update session before calling Execute
			err = cmd.Execute(ctx, opts)
			require.NoError(t, err)

			events, err := cmd.Events(ctx, opts)
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantEvent != nil {
				require.Len(t, events, 1)
				assert.Equal(t, tt.wantEvent, events[0])
				return
			}
			assert.Empty(t, events)
		})
	}
}

func getChallengePasskey(t *testing.T, data []byte) *session_grpc.Challenges_WebAuthN {
	webAuthNChallenge := &session_grpc.Challenges_WebAuthN{
		PublicKeyCredentialRequestOptions: new(structpb.Struct),
	}
	err := json.Unmarshal(data, webAuthNChallenge.PublicKeyCredentialRequestOptions)
	require.NoError(t, err)

	return webAuthNChallenge
}

func getUserByPasskeyType(ctrl *gomock.Controller, passkeyType domain.PasskeyType) *domainmock.UserRepo {
	repo := domainmock.NewUserRepo(ctrl)
	humanRepo := domainmock.NewHumanRepo(ctrl)
	repo.EXPECT().Human().Times(1).Return(humanRepo)

	repo.EXPECT().
		Get(gomock.Any(),
			gomock.Any(),
			dbmock.QueryOptions(
				database.WithCondition(
					repo.PrimaryKeyCondition("instance-1", "user-1"),
				),
			),
			dbmock.QueryOptions(
				database.WithCondition(
					humanRepo.PasskeyConditions().TypeCondition(database.TextOperationEqual, passkeyType),
				),
			),
		).
		Times(1).
		Return(&domain.User{
			ID:             "user-id",
			InstanceID:     "instance-1",
			OrganizationID: "organization-id",
			Username:       "username",
			Human: &domain.HumanUser{
				Passkeys: []*domain.Passkey{
					{
						ID:         "passkey-id",
						PublicKey:  []byte("public-key"),
						Name:       "My Passkey",
						Type:       passkeyType,
						KeyID:      []byte("key-id"),
						VerifiedAt: time.Now().Add(-2 * time.Hour),
					},
				},
				Email: domain.HumanEmail{
					Address: "user@example.com",
				},
			},
			State: domain.UserStateActive,
		}, nil)
	return repo
}
