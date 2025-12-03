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

type PasskeyChallengeCommand struct {
	RequestChallengePasskey *session_grpc.RequestChallenges_WebAuthN

	SessionID  string
	InstanceID string
	BeginLogin beginLoginFn

	User             *User
	Session          *Session
	ChallengePasskey *SessionChallengePasskey

	WebAuthNChallenge *session_grpc.Challenges_WebAuthN // todo: review
}

func NewPasskeyChallengeCommand(
	sessionID,
	instanceID string,
	requestChallengePasskey *session_grpc.RequestChallenges_WebAuthN,
	beginLoginFn beginLoginFn,
) *PasskeyChallengeCommand {
	passkeyChallengeCmd := &PasskeyChallengeCommand{
		RequestChallengePasskey: requestChallengePasskey,
		SessionID:               sessionID,
		InstanceID:              instanceID,
		BeginLogin:              beginLoginFn,
	}
	if beginLoginFn == nil {
		passkeyChallengeCmd.BeginLogin = webauthnConfig.BeginWebAuthNLogin
	}
	return passkeyChallengeCmd
}

func (p *PasskeyChallengeCommand) RequiresTransaction() {}

func (p *PasskeyChallengeCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if p.RequestChallengePasskey == nil {
		return nil
	}
	// get session
	sessionRepo := opts.sessionRepo
	p.Session, err = sessionRepo.Get(ctx, opts.DB(), database.WithCondition(sessionRepo.IDCondition(p.SessionID)))
	if err := handleGetError(err, "DOM-zy4hYC", objectTypeSession); err != nil {
		return err
	}
	if p.Session.UserID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-uVyrt2", "missing user id in session")
	}

	// retrieve user and passkeys based on user verification requirement
	passkeyType := PasskeyTypeU2F
	userVerification := UserVerificationRequirementToDomain(p.RequestChallengePasskey.GetUserVerificationRequirement())
	if userVerification == old_domain.UserVerificationRequirementRequired {
		passkeyType = PasskeyTypePasswordless
	}
	userRepo := opts.userRepo.LoadPasskeys()
	p.User, err = userRepo.Get(
		ctx,
		opts.DB(),
		database.WithCondition(userRepo.IDCondition(p.Session.UserID)),
		database.WithCondition(userRepo.Human().PasskeyTypeCondition(database.NumberOperationEqual, passkeyType)),
	)
	if err := handleGetError(err, "DOM-8cGMtd", objectTypeUser); err != nil {
		return err
	}

	// ensure user is active
	if p.User.State != UserStateActive {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-bnxBdS", "user not active")
	}

	return nil
}

func (p *PasskeyChallengeCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	if p.RequestChallengePasskey == nil {
		return nil
	}
	// begin webauthn login
	sessionRepo := opts.sessionRepo
	challenge := &session_grpc.Challenges_WebAuthN{
		PublicKeyCredentialRequestOptions: new(structpb.Struct),
	}
	if p.User.Human == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-nd3f4", "user is not a human user")
	}
	webUser := &webAuthNUser{
		userID:      p.User.ID,
		username:    p.User.Username,
		displayName: p.User.Human.DisplayName,
		creds:       PasskeysToCredentials(ctx, p.User.Human.Passkeys, p.RequestChallengePasskey.GetDomain()),
	}
	userVerification := UserVerificationFromDomain(UserVerificationRequirementToDomain(p.RequestChallengePasskey.GetUserVerificationRequirement()))
	sessionData, credentialAssertionData, rpID, err := p.BeginLogin(
		ctx,
		webUser,
		p.RequestChallengePasskey.GetDomain(),
		userVerification,
	)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-Fy333Q", "failed to begin webauthn login")
	}
	// todo: where is this challenge used? set in the CreateSessionResponse in the current implementation
	if err = json.Unmarshal(credentialAssertionData, challenge.PublicKeyCredentialRequestOptions); err != nil {
		return zerrors.ThrowInternal(nil, "DOM-liSCA4", "failed to unmarshal credential assertion data")
	}

	// todo: review
	p.WebAuthNChallenge = challenge

	// set challenge in session
	p.ChallengePasskey = &SessionChallengePasskey{
		Challenge:            sessionData.Challenge,
		AllowedCredentialIDs: sessionData.AllowedCredentialIDs,
		UserVerification:     UserVerificationRequirementToDomain(p.RequestChallengePasskey.GetUserVerificationRequirement()),
		RPID:                 rpID,
		LastChallengedAt:     time.Now(),
	}
	updated, err := sessionRepo.Update(
		ctx,
		opts.DB(),
		sessionRepo.IDCondition(p.Session.ID),
		sessionRepo.SetChallenge(p.ChallengePasskey),
	)
	if err := handleUpdateError(err, expectedUpdatedRows, updated, "DOM-yd3f4", objectTypeSession); err != nil {
		return err
	}
	return nil
}

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
			&session.NewAggregate(p.SessionID, p.InstanceID).Aggregate,
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
