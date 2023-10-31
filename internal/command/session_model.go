package command

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
)

type WebAuthNChallengeModel struct {
	Challenge          string
	AllowedCrentialIDs [][]byte
	UserVerification   domain.UserVerificationRequirement
	RPID               string
}

type OTPCode struct {
	Code         *crypto.CryptoValue
	Expiry       time.Duration
	CreationDate time.Time
}

func (p *WebAuthNChallengeModel) WebAuthNLogin(human *domain.Human, credentialAssertionData []byte) *domain.WebAuthNLogin {
	return &domain.WebAuthNLogin{
		ObjectRoot:              human.ObjectRoot,
		CredentialAssertionData: credentialAssertionData,
		Challenge:               p.Challenge,
		AllowedCredentialIDs:    p.AllowedCrentialIDs,
		UserVerification:        p.UserVerification,
		RPID:                    p.RPID,
	}
}

type SessionWriteModel struct {
	eventstore.WriteModel

	TokenID              string
	UserID               string
	UserCheckedAt        time.Time
	PasswordCheckedAt    time.Time
	IntentCheckedAt      time.Time
	WebAuthNCheckedAt    time.Time
	TOTPCheckedAt        time.Time
	OTPSMSCheckedAt      time.Time
	OTPEmailCheckedAt    time.Time
	WebAuthNUserVerified bool
	Metadata             map[string][]byte
	State                domain.SessionState
	Expiration           time.Time

	WebAuthNChallenge     *WebAuthNChallengeModel
	OTPSMSCodeChallenge   *OTPCode
	OTPEmailCodeChallenge *OTPCode

	aggregate *eventstore.Aggregate
}

func NewSessionWriteModel(sessionID string, resourceOwner string) *SessionWriteModel {
	return &SessionWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   sessionID,
			ResourceOwner: resourceOwner,
		},
		Metadata:  make(map[string][]byte),
		aggregate: &session.NewAggregate(sessionID, resourceOwner).Aggregate,
	}
}

func (wm *SessionWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *session.AddedEvent:
			wm.reduceAdded(e)
		case *session.UserCheckedEvent:
			wm.reduceUserChecked(e)
		case *session.PasswordCheckedEvent:
			wm.reducePasswordChecked(e)
		case *session.IntentCheckedEvent:
			wm.reduceIntentChecked(e)
		case *session.WebAuthNChallengedEvent:
			wm.reduceWebAuthNChallenged(e)
		case *session.WebAuthNCheckedEvent:
			wm.reduceWebAuthNChecked(e)
		case *session.TOTPCheckedEvent:
			wm.reduceTOTPChecked(e)
		case *session.OTPSMSChallengedEvent:
			wm.reduceOTPSMSChallenged(e)
		case *session.OTPSMSCheckedEvent:
			wm.reduceOTPSMSChecked(e)
		case *session.OTPEmailChallengedEvent:
			wm.reduceOTPEmailChallenged(e)
		case *session.OTPEmailCheckedEvent:
			wm.reduceOTPEmailChecked(e)
		case *session.TokenSetEvent:
			wm.reduceTokenSet(e)
		case *session.LifetimeSetEvent:
			wm.reduceLifetimeSet(e)
		case *session.TerminateEvent:
			wm.reduceTerminate()
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *SessionWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(session.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			session.AddedType,
			session.UserCheckedType,
			session.PasswordCheckedType,
			session.IntentCheckedType,
			session.WebAuthNChallengedType,
			session.WebAuthNCheckedType,
			session.TOTPCheckedType,
			session.OTPSMSChallengedType,
			session.OTPSMSCheckedType,
			session.OTPEmailChallengedType,
			session.OTPEmailCheckedType,
			session.TokenSetType,
			session.MetadataSetType,
			session.LifetimeSetType,
			session.TerminateType,
		).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func (wm *SessionWriteModel) reduceAdded(e *session.AddedEvent) {
	wm.State = domain.SessionStateActive
}

func (wm *SessionWriteModel) reduceUserChecked(e *session.UserCheckedEvent) {
	wm.UserID = e.UserID
	wm.UserCheckedAt = e.CheckedAt
}

func (wm *SessionWriteModel) reducePasswordChecked(e *session.PasswordCheckedEvent) {
	wm.PasswordCheckedAt = e.CheckedAt
}

func (wm *SessionWriteModel) reduceIntentChecked(e *session.IntentCheckedEvent) {
	wm.IntentCheckedAt = e.CheckedAt
}

func (wm *SessionWriteModel) reduceWebAuthNChallenged(e *session.WebAuthNChallengedEvent) {
	wm.WebAuthNChallenge = &WebAuthNChallengeModel{
		Challenge:          e.Challenge,
		AllowedCrentialIDs: e.AllowedCrentialIDs,
		UserVerification:   e.UserVerification,
		RPID:               e.RPID,
	}
}

func (wm *SessionWriteModel) reduceWebAuthNChecked(e *session.WebAuthNCheckedEvent) {
	wm.WebAuthNChallenge = nil
	wm.WebAuthNCheckedAt = e.CheckedAt
	wm.WebAuthNUserVerified = e.UserVerified
}

func (wm *SessionWriteModel) reduceTOTPChecked(e *session.TOTPCheckedEvent) {
	wm.TOTPCheckedAt = e.CheckedAt
}

func (wm *SessionWriteModel) reduceOTPSMSChallenged(e *session.OTPSMSChallengedEvent) {
	wm.OTPSMSCodeChallenge = &OTPCode{
		Code:         e.Code,
		Expiry:       e.Expiry,
		CreationDate: e.CreationDate(),
	}
}

func (wm *SessionWriteModel) reduceOTPSMSChecked(e *session.OTPSMSCheckedEvent) {
	wm.OTPSMSCodeChallenge = nil
	wm.OTPSMSCheckedAt = e.CheckedAt
}

func (wm *SessionWriteModel) reduceOTPEmailChallenged(e *session.OTPEmailChallengedEvent) {
	wm.OTPEmailCodeChallenge = &OTPCode{
		Code:         e.Code,
		Expiry:       e.Expiry,
		CreationDate: e.CreationDate(),
	}
}

func (wm *SessionWriteModel) reduceOTPEmailChecked(e *session.OTPEmailCheckedEvent) {
	wm.OTPEmailCodeChallenge = nil
	wm.OTPEmailCheckedAt = e.CheckedAt
}

func (wm *SessionWriteModel) reduceTokenSet(e *session.TokenSetEvent) {
	wm.TokenID = e.TokenID
}

func (wm *SessionWriteModel) reduceLifetimeSet(e *session.LifetimeSetEvent) {
	wm.Expiration = e.CreationDate().Add(e.Lifetime)
}

func (wm *SessionWriteModel) reduceTerminate() {
	wm.State = domain.SessionStateTerminated
}

// AuthenticationTime returns the time the user authenticated using the latest time of all checks
func (wm *SessionWriteModel) AuthenticationTime() time.Time {
	var authTime time.Time
	for _, check := range []time.Time{
		wm.PasswordCheckedAt,
		wm.WebAuthNCheckedAt,
		wm.TOTPCheckedAt,
		wm.IntentCheckedAt,
		wm.OTPSMSCheckedAt,
		wm.OTPEmailCheckedAt,
	} {
		if check.After(authTime) {
			authTime = check
		}
	}
	return authTime
}

// AuthMethodTypes returns a list of UserAuthMethodTypes based on succeeded checks
func (wm *SessionWriteModel) AuthMethodTypes() []domain.UserAuthMethodType {
	types := make([]domain.UserAuthMethodType, 0, domain.UserAuthMethodTypeIDP)
	if !wm.PasswordCheckedAt.IsZero() {
		types = append(types, domain.UserAuthMethodTypePassword)
	}
	if !wm.WebAuthNCheckedAt.IsZero() {
		if wm.WebAuthNUserVerified {
			types = append(types, domain.UserAuthMethodTypePasswordless)
		} else {
			types = append(types, domain.UserAuthMethodTypeU2F)
		}
	}
	if !wm.IntentCheckedAt.IsZero() {
		types = append(types, domain.UserAuthMethodTypeIDP)
	}
	if !wm.TOTPCheckedAt.IsZero() {
		types = append(types, domain.UserAuthMethodTypeTOTP)
	}
	if !wm.OTPSMSCheckedAt.IsZero() {
		types = append(types, domain.UserAuthMethodTypeOTPSMS)
	}
	if !wm.OTPEmailCheckedAt.IsZero() {
		types = append(types, domain.UserAuthMethodTypeOTPEmail)
	}
	return types
}

// CheckNotInvalidated checks that the session was not invalidated either manually ([session.TerminateType])
// or automatically (expired).
func (wm *SessionWriteModel) CheckNotInvalidated() error {
	if wm.State == domain.SessionStateTerminated {
		return errors.ThrowPreconditionFailed(nil, "COMMAND-Hewfq", "Errors.Session.Terminated")
	}
	if !wm.Expiration.IsZero() && wm.Expiration.Before(time.Now()) {
		return errors.ThrowPreconditionFailed(nil, "COMMAND-Hkl3d", "Errors.Session.Expired")
	}
	return nil
}

// CheckIsActive checks that the session was not invalidated ([CheckNotInvalidated]) and actually already exists.
func (wm *SessionWriteModel) CheckIsActive() error {
	if wm.State == domain.SessionStateUnspecified {
		return errors.ThrowPreconditionFailed(nil, "COMMAND-Flk38", "Errors.Session.NotExisting")
	}
	return wm.CheckNotInvalidated()
}
