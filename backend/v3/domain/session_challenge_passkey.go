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

const (
	expectedUpdatedRows = 1
	objectTypeSession   = "session"
	objectTypeUser      = "user"
)

var _ Commander = (*PasskeyChallengeCommand)(nil)
var _ Transactional = (*PasskeyChallengeCommand)(nil)

type beginLoginFn func(ctx context.Context, user webauthn.User, rpID string, userVerification protocol.UserVerificationRequirement) (sessionData *webauthn.SessionData, cred []byte, relyingPartyID string, err error)

// PasskeyChallengeCommand handles the passkey challenge creation during session creation.
// It supports both U2F and passwordless flows based on user verification requirements.
type PasskeyChallengeCommand struct {
	RequestChallengePasskey *session_grpc.RequestChallenges_WebAuthN

	sessionID  string
	instanceID string
	beginLogin beginLoginFn

	User    *User
	Session *Session

	ChallengePasskey  *SessionChallengePasskey          // the generated passkey challenge that is stored in the session.
	WebAuthNChallenge *session_grpc.Challenges_WebAuthN // challenge to be set in the CreateSessionResponse
}

// NewPasskeyChallengeCommand creates a new PasskeyChallengeCommand.
func NewPasskeyChallengeCommand(
	sessionID,
	instanceID string,
	requestChallengePasskey *session_grpc.RequestChallenges_WebAuthN,
	beginLoginFn beginLoginFn,
) *PasskeyChallengeCommand {

	if beginLoginFn == nil {
		beginLoginFn = webauthnConfig.BeginWebAuthNLogin
	}
	return &PasskeyChallengeCommand{
		RequestChallengePasskey: requestChallengePasskey,
		sessionID:               sessionID,
		instanceID:              instanceID,
		beginLogin:              beginLoginFn,
	}
}

// RequiresTransaction implements [Transactional].
func (p *PasskeyChallengeCommand) RequiresTransaction() {}

// Validate implements [Commander].
// It retrieves the session, the user, and the user's passkeys based on the user verification requirement,
// and validates that the user is active.
func (p *PasskeyChallengeCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if p.RequestChallengePasskey == nil {
		return nil
	}

	if p.sessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-EVo5yE", "missing session id")
	}
	if p.instanceID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-sh8xvQ", "missing instance id")
	}

	// get session
	sessionRepo := opts.sessionRepo
	session, err := sessionRepo.Get(ctx, opts.DB(), database.WithCondition(sessionRepo.IDCondition(p.sessionID)))
	if err := handleGetError(err, "DOM-zy4hYC", objectTypeSession); err != nil {
		return err
	}
	if session.UserID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-uVyrt2", "missing user id in session")
	}

	// retrieve user and their passkeys based on user verification requirement
	passkeyType := determinePasskeyType(p.RequestChallengePasskey.GetUserVerificationRequirement())
	userRepo := opts.userRepo.LoadPasskeys()
	user, err := userRepo.Get(
		ctx,
		opts.DB(),
		database.WithCondition(userRepo.IDCondition(session.UserID)),
		database.WithCondition(userRepo.Human().PasskeyTypeCondition(database.NumberOperationEqual, passkeyType)),
	)
	if err := handleGetError(err, "DOM-8cGMtd", objectTypeUser); err != nil {
		return err
	}

	// ensure user is active
	if user.State != UserStateActive {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-bnxBdS", "user not active")
	}

	p.Session = session
	p.User = user

	return nil
}

// Execute implements [Commander].
// It begins the WebAuthN login process and updates the session with the passkey challenge.
func (p *PasskeyChallengeCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	if p.RequestChallengePasskey == nil {
		return nil
	}
	// begin webauthn login
	userVerificationDomain := UserVerificationRequirementToDomain(p.RequestChallengePasskey.GetUserVerificationRequirement())
	sessionData, credentialAssertionData, rpID, err := p.beginWebAuthNLogin(ctx, userVerificationDomain)
	if err != nil {
		return err
	}
	// set the challenge for CreateSessionResponse
	webAuthNChallenge := &session_grpc.Challenges_WebAuthN{
		PublicKeyCredentialRequestOptions: new(structpb.Struct),
	}
	if err = json.Unmarshal(credentialAssertionData, webAuthNChallenge.PublicKeyCredentialRequestOptions); err != nil {
		return zerrors.ThrowInternal(nil, "DOM-liSCA4", "failed to unmarshal credential assertion data")
	}

	// update the session with the passkey challenge
	challengePasskey := &SessionChallengePasskey{
		Challenge:            sessionData.Challenge,
		AllowedCredentialIDs: sessionData.AllowedCredentialIDs,
		UserVerification:     userVerificationDomain,
		RPID:                 rpID,
		LastChallengedAt:     time.Now(),
	}
	sessionRepo := opts.sessionRepo
	updated, err := sessionRepo.Update(
		ctx,
		opts.DB(),
		sessionRepo.IDCondition(p.Session.ID),
		sessionRepo.SetChallenge(challengePasskey),
	)
	if err := handleUpdateError(err, expectedUpdatedRows, updated, "DOM-yd3f4", objectTypeSession); err != nil {
		return err
	}
	p.WebAuthNChallenge = webAuthNChallenge
	p.ChallengePasskey = challengePasskey
	return nil
}

// Events implements [Commander].
// It creates a WebAuthN challenged event if a passkey challenge was requested and created.
func (p *PasskeyChallengeCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	if p.RequestChallengePasskey == nil {
		return nil, nil
	}
	if p.ChallengePasskey == nil {
		return nil, zerrors.ThrowInternal(nil, "DOM-MALUxr", "failed to push WebAuthN challenged event")
	}
	return []eventstore.Command{
		session.NewWebAuthNChallengedEvent(
			ctx,
			&session.NewAggregate(p.sessionID, p.instanceID).Aggregate,
			p.ChallengePasskey.Challenge,
			p.ChallengePasskey.AllowedCredentialIDs,
			p.ChallengePasskey.UserVerification,
			p.ChallengePasskey.RPID,
		),
	}, nil
}

func (p *PasskeyChallengeCommand) String() string {
	return "PasskeyChallengeCommand"
}

// beginWebAuthNLogin starts the WebAuthN login process for the user with the specified user verification requirement.
func (p *PasskeyChallengeCommand) beginWebAuthNLogin(ctx context.Context, userVerificationDomain old_domain.UserVerificationRequirement) (*webauthn.SessionData, []byte, string, error) {
	if p.User.Human == nil {
		return nil, nil, "", zerrors.ThrowPreconditionFailed(nil, "DOM-nd3f4", "user is not a human user")
	}
	webUser := &webAuthNUser{
		userID:      p.User.ID,
		username:    p.User.Username,
		displayName: p.User.Human.DisplayName,
		creds:       PasskeysToCredentials(ctx, p.User.Human.Passkeys, p.RequestChallengePasskey.GetDomain()),
	}
	sessionData, credentialAssertionData, rpID, err := p.beginLogin(
		ctx,
		webUser,
		p.RequestChallengePasskey.GetDomain(),
		UserVerificationFromDomain(userVerificationDomain),
	)
	if err != nil {
		return nil, nil, "", zerrors.ThrowInternal(err, "DOM-Fy333Q", "failed to begin webauthn login")
	}
	return sessionData, credentialAssertionData, rpID, nil
}

// determinePasskeyType determines the passkey type (U2F or Passwordless) based on the user verification requirement.
func determinePasskeyType(pbUserVerificationRequirement session_grpc.UserVerificationRequirement) PasskeyType {
	passkeyType := PasskeyTypeU2F
	userVerification := UserVerificationRequirementToDomain(pbUserVerificationRequirement)
	if userVerification == old_domain.UserVerificationRequirementRequired {
		passkeyType = PasskeyTypePasswordless
	}
	return passkeyType
}
