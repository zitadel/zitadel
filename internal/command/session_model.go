package command

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
)

type PasskeyChallengeModel struct {
	Challenge          string
	AllowedCrentialIDs [][]byte
	UserVerification   domain.UserVerificationRequirement
	RPID               string
}

func (p *PasskeyChallengeModel) WebAuthNLogin(human *domain.Human, credentialAssertionData []byte) (*domain.WebAuthNLogin, error) {
	if p == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Ioqu5", "Errors.Session.Passkey.NoChallenge")
	}
	return &domain.WebAuthNLogin{
		ObjectRoot:              human.ObjectRoot,
		CredentialAssertionData: credentialAssertionData,
		Challenge:               p.Challenge,
		AllowedCredentialIDs:    p.AllowedCrentialIDs,
		UserVerification:        p.UserVerification,
		RPID:                    p.RPID,
	}, nil
}

type SessionWriteModel struct {
	eventstore.WriteModel

	TokenID           string
	UserID            string
	UserCheckedAt     time.Time
	PasswordCheckedAt time.Time
	IntentCheckedAt   time.Time
	PasskeyCheckedAt  time.Time
	Metadata          map[string][]byte
	Domain            string
	State             domain.SessionState

	PasskeyChallenge *PasskeyChallengeModel

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
		case *session.PasskeyChallengedEvent:
			wm.reducePasskeyChallenged(e)
		case *session.PasskeyCheckedEvent:
			wm.reducePasskeyChecked(e)
		case *session.TokenSetEvent:
			wm.reduceTokenSet(e)
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
			session.PasskeyChallengedType,
			session.PasskeyCheckedType,
			session.TokenSetType,
			session.MetadataSetType,
			session.TerminateType,
		).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func (wm *SessionWriteModel) reduceAdded(e *session.AddedEvent) {
	wm.Domain = e.Domain
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

func (wm *SessionWriteModel) reducePasskeyChallenged(e *session.PasskeyChallengedEvent) {
	wm.PasskeyChallenge = &PasskeyChallengeModel{
		Challenge:          e.Challenge,
		AllowedCrentialIDs: e.AllowedCrentialIDs,
		UserVerification:   e.UserVerification,
		RPID:               wm.Domain,
	}
}

func (wm *SessionWriteModel) reducePasskeyChecked(e *session.PasskeyCheckedEvent) {
	wm.PasskeyChallenge = nil
	wm.PasskeyCheckedAt = e.CheckedAt
}

func (wm *SessionWriteModel) reduceTokenSet(e *session.TokenSetEvent) {
	wm.TokenID = e.TokenID
}

func (wm *SessionWriteModel) reduceTerminate() {
	wm.State = domain.SessionStateTerminated
}

// AuthenticationTime returns the time the user authenticated using the latest time of all checks
func (wm *SessionWriteModel) AuthenticationTime() time.Time {
	var authTime time.Time
	for _, check := range []time.Time{
		wm.PasswordCheckedAt,
		wm.PasskeyCheckedAt,
		wm.IntentCheckedAt,
		// TODO: add U2F and OTP check https://github.com/zitadel/zitadel/issues/5477
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
	if !wm.PasskeyCheckedAt.IsZero() {
		types = append(types, domain.UserAuthMethodTypePasswordless)
	}
	if !wm.IntentCheckedAt.IsZero() {
		types = append(types, domain.UserAuthMethodTypeIDP)
	}
	// TODO: add checks with https://github.com/zitadel/zitadel/issues/5477
	/*
		if !wm.TOTPCheckedAt.IsZero() {
			types = append(types, domain.UserAuthMethodTypeOTP)
		}
		if !wm.U2FCheckedAt.IsZero() {
			types = append(types, domain.UserAuthMethodTypeU2F)
		}
	*/
	return types
}
