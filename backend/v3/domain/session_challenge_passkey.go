package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	old_domain "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

type ChallengeTypePasskey struct {
	Domain                      string
	UserVerificationRequirement old_domain.UserVerificationRequirement
}

var _ Commander = (*PasskeyChallengeCommand)(nil)
var _ Transactional = (*PasskeyChallengeCommand)(nil)

type beginLoginFn func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error)

// PasskeyChallengeCommand handles the passkey challenge creation during session creation.
// It supports both U2F and passwordless flows based on user verification requirements.
type PasskeyChallengeCommand struct {
	ChallengeTypePasskey *ChallengeTypePasskey

	SessionID  string
	InstanceID string

	user       *User
	beginLogin beginLoginFn

	challengePasskey  *SessionChallengePasskey          // the generated passkey challenge that is stored in the session.
	webAuthNChallenge *session_grpc.Challenges_WebAuthN // challenge to be set in the CreateSessionResponse
}

// NewPasskeyChallengeCommand creates a new PasskeyChallengeCommand.
func NewPasskeyChallengeCommand(
	sessionID,
	instanceID string,
	challengePasskey *ChallengeTypePasskey,
	beginLoginFn beginLoginFn,
) (*PasskeyChallengeCommand, error) {

	if beginLoginFn == nil {
		if webauthnConfig == nil {
			return nil, zerrors.ThrowInternal(nil, "DOM-jwk5Pe", "begin webauthn login function not set")
		}
		beginLoginFn = webauthnConfig.BeginWebAuthNLogin
	}
	return &PasskeyChallengeCommand{
		ChallengeTypePasskey: challengePasskey,
		SessionID:            sessionID,
		InstanceID:           instanceID,
		beginLogin:           beginLoginFn,
	}, nil
}

// RequiresTransaction implements [Transactional].
func (p *PasskeyChallengeCommand) RequiresTransaction() {}

// Validate implements [Commander].
// It retrieves the session, the user, and the user's passkeys based on the user verification requirement,
// and validates that the user is active.
func (p *PasskeyChallengeCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if p.ChallengeTypePasskey == nil {
		return nil
	}

	if p.SessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-EVo5yE", "Errors.Missing.SessionID")
	}
	if p.InstanceID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-sh8xvQ", "Errors.Missing.InstanceID")
	}

	// get session
	sessionRepo := opts.sessionRepo
	retrievedSession, err := sessionRepo.Get(
		ctx,
		opts.DB(),
		database.WithCondition(sessionRepo.PrimaryKeyCondition(p.InstanceID, p.SessionID)),
	)
	if err := handleGetError(err, "DOM-zy4hYC", objectTypeSession); err != nil {
		return err
	}
	if retrievedSession.UserID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-uVyrt2", "Errors.Missing.Session.UserID")
	}

	// retrieve user and their passkeys based on user verification requirement
	passkeyType := determinePasskeyType(p.ChallengeTypePasskey.UserVerificationRequirement)
	userRepo := opts.userRepo
	retrievedUser, err := userRepo.Get(
		ctx,
		opts.DB(),
		database.WithCondition(userRepo.PrimaryKeyCondition(p.InstanceID, retrievedSession.UserID)),
		database.WithCondition(userRepo.Human().PasskeyConditions().TypeCondition(database.TextOperationEqual, passkeyType)),
	)
	if err := handleGetError(err, "DOM-8cGMtd", objectTypeUser); err != nil {
		return err
	}

	// ensure the user is a human user
	if retrievedUser.Human == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-nd3f4", "Errors.User.NotHuman")
	}

	// ensure the user is active
	if retrievedUser.State != UserStateActive {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-bnxBdS", "Errors.User.NotFound")
	}

	p.user = retrievedUser

	return nil
}

// Execute implements [Commander].
// It begins the WebAuthN login process and updates the session with the passkey challenge.
func (p *PasskeyChallengeCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	if p.ChallengeTypePasskey == nil {
		return nil
	}
	// begin webauthn login
	sessionData, credentialAssertionData, rpID, err := p.beginWebAuthNLogin(ctx, p.ChallengeTypePasskey.UserVerificationRequirement)
	if err != nil {
		return err
	}
	// to set the challenge for CreateSessionResponse
	webAuthNChallenge := &session_grpc.Challenges_WebAuthN{
		PublicKeyCredentialRequestOptions: new(structpb.Struct),
	}
	if err = json.Unmarshal(credentialAssertionData, webAuthNChallenge.PublicKeyCredentialRequestOptions); err != nil {
		return zerrors.ThrowInternal(err, "DOM-liSCA4", "Errors.Unmarshal")
	}

	// update the session with the passkey challenge
	challengePasskey := &SessionChallengePasskey{
		Challenge:            sessionData.Challenge,
		AllowedCredentialIDs: sessionData.AllowedCredentialIDs,
		UserVerification:     p.ChallengeTypePasskey.UserVerificationRequirement,
		RPID:                 rpID,
		LastChallengedAt:     time.Now(),
	}
	sessionRepo := opts.sessionRepo
	updated, err := sessionRepo.Update(
		ctx,
		opts.DB(),
		sessionRepo.PrimaryKeyCondition(p.InstanceID, p.SessionID),
		sessionRepo.SetChallenge(challengePasskey),
	)
	if err := handleUpdateError(err, expectedUpdatedRows, updated, "DOM-yd3f4", objectTypeSession); err != nil {
		return err
	}
	p.webAuthNChallenge = webAuthNChallenge
	p.challengePasskey = challengePasskey

	return nil
}

// Events implements [Commander].
// It creates a WebAuthN challenged event if a passkey challenge was requested and created.
func (p *PasskeyChallengeCommand) Events(ctx context.Context, _ *InvokeOpts) ([]eventstore.Command, error) {
	if p.ChallengeTypePasskey == nil {
		return nil, nil
	}
	return []eventstore.Command{
		session.NewWebAuthNChallengedEvent(
			ctx,
			&session.NewAggregate(p.SessionID, p.InstanceID).Aggregate,
			p.challengePasskey.Challenge,
			p.challengePasskey.AllowedCredentialIDs,
			p.challengePasskey.UserVerification,
			p.challengePasskey.RPID,
		),
	}, nil
}

func (p *PasskeyChallengeCommand) GetWebAuthNChallenge() *session_grpc.Challenges_WebAuthN {
	return p.webAuthNChallenge
}

func (p *PasskeyChallengeCommand) String() string {
	return "PasskeyChallengeCommand"
}

// beginWebAuthNLogin starts the WebAuthN login process for the user with the specified user verification requirement.
func (p *PasskeyChallengeCommand) beginWebAuthNLogin(ctx context.Context, userVerificationDomain old_domain.UserVerificationRequirement) (*webauthn.SessionData, []byte, string, error) {
	webUser := &webAuthNUser{
		userID:      p.user.ID,
		username:    p.user.Username,
		displayName: p.user.Human.DisplayName,
		creds:       PasskeysToCredentials(ctx, p.user.Human.Passkeys, p.ChallengeTypePasskey.Domain),
	}
	sessionData, credentialAssertionData, rpID, err := p.beginLogin(
		ctx,
		webUser,
		p.ChallengeTypePasskey.Domain,
		UserVerificationFromDomain(userVerificationDomain),
	)
	if err != nil {
		return nil, nil, "", err
	}
	return sessionData, credentialAssertionData, rpID, nil
}

// determinePasskeyType determines the passkey type (U2F or Passwordless) based on the user verification requirement.
func determinePasskeyType(userVerificationRequirement old_domain.UserVerificationRequirement) PasskeyType {
	passkeyType := PasskeyTypeU2F
	if userVerificationRequirement == old_domain.UserVerificationRequirementRequired {
		passkeyType = PasskeyTypePasswordless
	}
	return passkeyType
}
