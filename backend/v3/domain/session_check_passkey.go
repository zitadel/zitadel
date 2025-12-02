package domain

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	old_domain "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

type FinishLoginFunc func(ctx context.Context, sessionData webauthn.SessionData, user webauthn.User, credentials []byte, rpID string) (*webauthn.Credential, error)

type PasskeyCheckCommand struct {
	CheckPasskey *session_grpc.CheckWebAuthN

	finishLoginFn FinishLoginFunc

	sessionID  string
	instanceID string

	fetchedUser    *User
	fetchedSession *Session

	// For Events()
	lastVeriedAt  time.Time
	userVerified  bool
	pkeyID        string
	pkeySignCount uint32
}

func NewPasskeyCheckCommand(sessionID, instanceID string, request *session_grpc.CheckWebAuthN, finishLoginFn FinishLoginFunc) *PasskeyCheckCommand {
	pcc := &PasskeyCheckCommand{
		CheckPasskey:  request,
		sessionID:     sessionID,
		instanceID:    instanceID,
		finishLoginFn: webauthnConfig.FinishLoginWithNewDomainModel,
	}

	if finishLoginFn != nil {
		pcc.finishLoginFn = finishLoginFn
	}

	return pcc
}

// RequiresTransaction implements [Transactional].
func (p *PasskeyCheckCommand) RequiresTransaction() {}

// Events implements [Commander].
func (p *PasskeyCheckCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	if p.CheckPasskey == nil {
		return nil, nil
	}

	passkeyChallenge := p.fetchedSession.Challenges.GetPasskeyChallenge()

	toReturn := make([]eventstore.Command, 2)

	sessionAgg := &session.NewAggregate(p.sessionID, p.instanceID).Aggregate
	toReturn[0] = session.NewWebAuthNCheckedEvent(ctx, sessionAgg, p.lastVeriedAt, p.userVerified)
	if passkeyChallenge.UserVerification == old_domain.UserVerificationRequirementRequired {
		toReturn[1] = user.NewHumanPasswordlessSignCountChangedEvent(ctx, sessionAgg, p.pkeyID, p.pkeySignCount)
	} else {
		toReturn[1] = user.NewHumanU2FSignCountChangedEvent(ctx, sessionAgg, p.pkeyID, p.pkeySignCount)
	}

	return toReturn, nil
}

// Execute implements [Commander].
func (p *PasskeyCheckCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	if p.CheckPasskey == nil {
		return nil
	}

	sessionRepo := opts.sessionRepo
	userRepo := opts.userRepo

	passkeyChallenge := p.fetchedSession.Challenges.GetPasskeyChallenge()

	var selectedPkeys []*Passkey
	if passkeyChallenge.UserVerification == old_domain.UserVerificationRequirementRequired {
		selectedPkeys = p.fetchedUser.Human.Passkeys.GetPasskeysOfType([]PasskeyType{PasskeyTypePasswordless})
	} else {
		selectedPkeys = p.fetchedUser.Human.Passkeys.GetPasskeysOfType([]PasskeyType{PasskeyTypeU2F, PasskeyTypeUnspecified})
	}

	webAuthnUsr := &webAuthNUser{
		userID:      p.fetchedUser.ID,
		username:    p.fetchedUser.Username,
		displayName: p.fetchedUser.Human.DisplayName,
		creds:       PasskeysToCredentials(ctx, selectedPkeys, passkeyChallenge.RPID),
	}
	credentialAssertionData, err := json.Marshal(p.CheckPasskey.GetCredentialAssertionData())
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-asd", "Errors.Internal")
	}

	webAuthCreds, err := p.finishLoginFn(ctx, p.getWebAuthNSessionData(passkeyChallenge, p.fetchedUser.ID), webAuthnUsr, credentialAssertionData, passkeyChallenge.RPID)
	if err != nil && (webAuthCreds == nil || webAuthCreds.ID == nil) {
		return err
	}

	var matchingPKey *Passkey
	for _, pkey := range selectedPkeys {
		if bytes.Equal(pkey.KeyID, webAuthCreds.ID) {
			matchingPKey = pkey
			break
		}
	}
	if matchingPKey == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-asd", "Errors.User.WebAuthN.NotFound")
	}

	p.userVerified = webAuthCreds.Flags.UserVerified
	p.lastVeriedAt = time.Now()
	rowCount, err := sessionRepo.Update(ctx, opts.DB(),
		sessionRepo.IDCondition(p.sessionID),
		sessionRepo.SetFactor(&SessionFactorPasskey{LastVerifiedAt: p.lastVeriedAt, UserVerified: p.userVerified}),
	)
	if err := handleUpdateError(err, 1, rowCount, "DOM-asd", "session"); err != nil {
		return err
	}

	rowCount, err = userRepo.Update(ctx, opts.DB(),
		userRepo.Human().PasskeyIDCondition(matchingPKey.ID),
		userRepo.Human().SetPasskeySignCount(webAuthCreds.Authenticator.SignCount),
	)
	if err := handleUpdateError(err, 1, rowCount, "DOM-asd", "user"); err != nil {
		return err
	}

	return nil
}

// String implements [Commander].
func (p *PasskeyCheckCommand) String() string {
	return "PasskeyCheckCommand"
}

// Validate implements [Commander].
func (p *PasskeyCheckCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if p.CheckPasskey == nil {
		return nil
	}

	sessionRepo := opts.sessionRepo
	userRepo := opts.userRepo.LoadPasskeys()

	p.fetchedSession, err = sessionRepo.Get(ctx, opts.DB(), database.WithCondition(sessionRepo.IDCondition(p.sessionID)))
	if err := handleGetError(err, "DOM-asd", "session"); err != nil {
		return err
	}

	passkeyChallenge := p.fetchedSession.Challenges.GetPasskeyChallenge()
	if passkeyChallenge == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-asd", "Errors.Session.WebAuthN.NoChallenge")
	}

	if p.fetchedSession.UserID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-asd", "Errors.User.UserIDMissing")
	}

	p.fetchedUser, err = userRepo.Get(ctx, opts.DB(), database.WithCondition(userRepo.IDCondition(p.fetchedSession.UserID)))
	if err := handleGetError(err, "DOM-asd", "user"); err != nil {
		return err
	}

	return nil
}

func (p *PasskeyCheckCommand) getWebAuthNSessionData(sessionChallenge *SessionChallengePasskey, userID string) webauthn.SessionData {
	return webauthn.SessionData{
		Challenge:            sessionChallenge.Challenge,
		UserID:               []byte(userID),
		AllowedCredentialIDs: sessionChallenge.AllowedCredentialIDs,
		UserVerification:     UserVerificationFromDomain(sessionChallenge.UserVerification),
	}
}

var _ Commander = (*PasskeyCheckCommand)(nil)
var _ Transactional = (*PasskeyCheckCommand)(nil)
