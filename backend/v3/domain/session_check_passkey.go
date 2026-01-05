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

	FetchedUser    *User
	FetchedSession *Session

	// For Events()
	LastVeriedAt  time.Time
	UserVerified  bool
	PKeyID        string
	PKeySignCount uint32
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

	passkeyChallenge := p.FetchedSession.Challenges.GetPasskeyChallenge()

	toReturn := make([]eventstore.Command, 2)

	sessionAgg := &session.NewAggregate(p.sessionID, p.instanceID).Aggregate
	toReturn[0] = session.NewWebAuthNCheckedEvent(ctx, sessionAgg, p.LastVeriedAt, p.UserVerified)
	if passkeyChallenge.UserVerification == old_domain.UserVerificationRequirementRequired {
		toReturn[1] = user.NewHumanPasswordlessSignCountChangedEvent(ctx, sessionAgg, p.PKeyID, p.PKeySignCount)
	} else {
		toReturn[1] = user.NewHumanU2FSignCountChangedEvent(ctx, sessionAgg, p.PKeyID, p.PKeySignCount)
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

	passkeyChallenge := p.FetchedSession.Challenges.GetPasskeyChallenge()

	userPKeys := p.FetchedUser.Human.Passkeys

	webAuthnUsr := &webAuthNUser{
		userID:      p.FetchedUser.ID,
		username:    p.FetchedUser.Username,
		displayName: p.FetchedUser.Human.DisplayName,
		creds:       PasskeysToCredentials(ctx, userPKeys, passkeyChallenge.RPID),
	}
	credentialAssertionData, err := json.Marshal(p.CheckPasskey.GetCredentialAssertionData())
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-I84Iyp", "Errors.Internal")
	}

	webAuthCreds, err := p.finishLoginFn(ctx, p.getWebAuthNSessionData(passkeyChallenge, p.FetchedUser.ID), webAuthnUsr, credentialAssertionData, passkeyChallenge.RPID)
	if err != nil && (webAuthCreds == nil || webAuthCreds.ID == nil) {
		return err
	}

	var matchingPKey *Passkey
	for _, pkey := range userPKeys {
		if bytes.Equal(pkey.KeyID, webAuthCreds.ID) {
			matchingPKey = pkey
			break
		}
	}
	if matchingPKey == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-uuxodH", "Errors.User.WebAuthN.NotFound")
	}

	p.UserVerified = webAuthCreds.Flags.UserVerified
	p.LastVeriedAt = time.Now()
	rowCount, err := sessionRepo.Update(ctx, opts.DB(),
		sessionRepo.IDCondition(p.sessionID),
		sessionRepo.SetFactor(&SessionFactorPasskey{LastVerifiedAt: p.LastVeriedAt, UserVerified: p.UserVerified}),
	)
	if err := handleUpdateError(err, 1, rowCount, "DOM-Uadvap", "session"); err != nil {
		return err
	}

	rowCount, err = userRepo.Update(ctx, opts.DB(),
		userRepo.Human().PasskeyIDCondition(matchingPKey.ID),
		userRepo.Human().SetPasskeySignCount(webAuthCreds.Authenticator.SignCount),
	)
	if err := handleUpdateError(err, 1, rowCount, "DOM-wdwZYk", "user"); err != nil {
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

	if p.sessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-4QJa2k", "Errors.Missing.SessionID")
	}
	if p.instanceID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-XlOhxU", "Errors.Missing.InstanceID")
	}

	sessionRepo := opts.sessionRepo
	userRepo := opts.userRepo.LoadPasskeys()

	p.FetchedSession, err = sessionRepo.Get(ctx, opts.DB(), database.WithCondition(sessionRepo.IDCondition(p.sessionID)))
	if err := handleGetError(err, "DOM-CUnePh", "session"); err != nil {
		return err
	}

	passkeyChallenge := p.FetchedSession.Challenges.GetPasskeyChallenge()
	if passkeyChallenge == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-lQhNR4", "Errors.Session.WebAuthN.NoChallenge")
	}

	if p.FetchedSession.UserID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-jy0zq7", "Errors.User.UserIDMissing")
	}

	var passKeyCondition database.Condition
	if passkeyChallenge.UserVerification == old_domain.UserVerificationRequirementRequired {
		passKeyCondition = userRepo.Human().PasskeyTypeCondition(database.NumberOperationEqual, PasskeyTypePasswordless)
	} else {
		passKeyCondition = userRepo.Human().PasskeyTypeCondition(database.NumberOperationNotEqual, PasskeyTypePasswordless)
	}

	p.FetchedUser, err = userRepo.Get(ctx, opts.DB(),
		database.WithCondition(userRepo.IDCondition(p.FetchedSession.UserID)),
		database.WithCondition(passKeyCondition),
	)
	if err := handleGetError(err, "DOM-pB6Mlm", "user"); err != nil {
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
