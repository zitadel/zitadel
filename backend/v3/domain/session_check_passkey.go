package domain

import (
	"bytes"
	"context"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	old_domain "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type FinishLoginFunc func(ctx context.Context, sessionData webauthn.SessionData, user webauthn.User, credentials []byte, rpID string) (*webauthn.Credential, error)

type PasskeyCheckCommand struct {
	// CheckPasskey is the assertion data for the passkey
	CheckPasskey []byte

	FinishLoginFn FinishLoginFunc

	SessionID  string
	InstanceID string

	FetchedUser    *User
	FetchedSession *Session

	// For Events()
	LastVerifiedAt time.Time
	UserVerified   bool
	PKeyID         string
	PKeySignCount  uint32
}

// NewPasskeyCheckCommand initializes a new [PasskeyCheckCommand]
//
// If finishLoginFn is nil, [webauthnConfig.FinishLoginWithNewDomainModel] will be used.
// If that is nil as well, an error will be returned.
//
// assertionData is passkey assertion required for validation
//
// The command does not implement [Transactional] due finishLoginFn that might take a long time to execute.
// So the DB transaction will be started only after finishLoginFn has been run.
func NewPasskeyCheckCommand(sessionID, instanceID string, assertionData []byte, finishLoginFn FinishLoginFunc) (*PasskeyCheckCommand, error) {
	if webauthnConfig == nil && finishLoginFn == nil {
		return nil, zerrors.ThrowInternal(nil, "DOM-bhzmHO", "no finish login function set")
	}

	pcc := &PasskeyCheckCommand{
		CheckPasskey: assertionData,
		SessionID:    sessionID,
		InstanceID:   instanceID,
	}

	if finishLoginFn != nil {
		pcc.FinishLoginFn = finishLoginFn
	} else {
		pcc.FinishLoginFn = webauthnConfig.FinishLoginWithNewDomainModel
	}

	return pcc, nil
}

// Events implements [Commander].
func (p *PasskeyCheckCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	if p.CheckPasskey == nil {
		return nil, nil
	}

	passkeyChallenge := p.FetchedSession.Challenges.GetPasskeyChallenge()

	events := make([]eventstore.Command, 2)

	sessionAgg := &session.NewAggregate(p.SessionID, p.InstanceID).Aggregate
	events[0] = session.NewWebAuthNCheckedEvent(ctx, sessionAgg, p.LastVerifiedAt, p.UserVerified)
	if passkeyChallenge.UserVerification == old_domain.UserVerificationRequirementRequired {
		events[1] = user.NewHumanPasswordlessSignCountChangedEvent(ctx, sessionAgg, p.PKeyID, p.PKeySignCount)
	} else {
		events[1] = user.NewHumanU2FSignCountChangedEvent(ctx, sessionAgg, p.PKeyID, p.PKeySignCount)
	}

	return events, nil
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

	webAuthCreds, err := p.FinishLoginFn(ctx, p.getWebAuthNSessionData(passkeyChallenge, p.FetchedUser.ID), webAuthnUsr, p.CheckPasskey, passkeyChallenge.RPID)
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
	p.PKeyID = matchingPKey.ID

	tx, txErr := opts.StartTransaction(ctx, nil)
	if txErr != nil {
		return zerrors.ThrowInternal(txErr, "DOM-sAAd3V", "failed starting transaction")
	}

	defer func() {
		if endErr := tx.End(ctx, txErr); endErr != nil {
			err = endErr
		}
	}()

	p.UserVerified = webAuthCreds.Flags.UserVerified
	p.LastVerifiedAt = time.Now()
	rowCount, err := sessionRepo.Update(ctx, tx,
		sessionRepo.IDCondition(p.SessionID),
		sessionRepo.SetFactor(&SessionFactorPasskey{LastVerifiedAt: p.LastVerifiedAt, UserVerified: p.UserVerified}),
	)
	if err := handleUpdateError(err, 1, rowCount, "DOM-Uadvap", "session"); err != nil {
		txErr = err
		return err
	}

	rowCount, err = userRepo.Update(ctx, tx,
		database.And(
			userRepo.Human().PrimaryKeyCondition(p.InstanceID, p.FetchedUser.ID),
			userRepo.Human().PasskeyConditions().IDCondition(matchingPKey.ID),
		),
		userRepo.Human().SetPasskeySignCount(webAuthCreds.Authenticator.SignCount),
	)
	if err := handleUpdateError(err, 1, rowCount, "DOM-wdwZYk", "user"); err != nil {
		txErr = err
		return err
	}
	p.PKeySignCount = webAuthCreds.Authenticator.SignCount

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

	if p.SessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-4QJa2k", "Errors.Missing.SessionID")
	}
	if p.InstanceID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-XlOhxU", "Errors.Missing.InstanceID")
	}

	sessionRepo := opts.sessionRepo
	userRepo := opts.userRepo

	p.FetchedSession, err = sessionRepo.Get(ctx, opts.DB(), database.WithCondition(sessionRepo.IDCondition(p.SessionID)))
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
		passKeyCondition = userRepo.Human().PasskeyConditions().TypeCondition(database.TextOperationEqual, PasskeyTypePasswordless)
	} else {
		passKeyCondition = userRepo.Human().PasskeyConditions().TypeCondition(database.TextOperationEqual, PasskeyTypeU2F)
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
