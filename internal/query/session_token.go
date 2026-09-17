package query

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// SessionTokenReadModel holds the state needed to verify a session token.
type SessionTokenReadModel struct {
	eventstore.ReadModel

	UserID                string
	UserResourceOwner     string
	PasswordCheckedAt     time.Time
	IntentCheckedAt       time.Time
	WebAuthNCheckedAt     time.Time
	WebAuthNUserVerified  bool
	TOTPCheckedAt         time.Time
	OTPSMSCheckedAt       time.Time
	OTPEmailCheckedAt     time.Time
	RecoveryCodeCheckedAt time.Time
	TokenID               string
	State                 domain.SessionState
	Expiration            time.Time

	userCheckedPosition     decimal.Decimal
	passwordCheckedPosition decimal.Decimal
}

func newSessionTokenReadModel(sessionID string) *SessionTokenReadModel {
	return &SessionTokenReadModel{
		ReadModel: eventstore.ReadModel{
			AggregateID: sessionID,
		},
	}
}

func (wm *SessionTokenReadModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *session.AddedEvent:
			wm.State = domain.SessionStateActive
		case *session.UserCheckedEvent:
			wm.UserID = e.UserID
			wm.UserResourceOwner = e.UserResourceOwner
			wm.userCheckedPosition = e.Position()
		case *session.PasswordCheckedEvent:
			wm.PasswordCheckedAt = e.CheckedAt
			wm.passwordCheckedPosition = e.Position()
		case *session.IntentCheckedEvent:
			wm.IntentCheckedAt = e.CheckedAt
		case *session.WebAuthNCheckedEvent:
			wm.WebAuthNCheckedAt = e.CheckedAt
			wm.WebAuthNUserVerified = e.UserVerified
		case *session.TOTPCheckedEvent:
			wm.TOTPCheckedAt = e.CheckedAt
		case *session.OTPSMSCheckedEvent:
			wm.OTPSMSCheckedAt = e.CheckedAt
		case *session.OTPEmailCheckedEvent:
			wm.OTPEmailCheckedAt = e.CheckedAt
		case *session.RecoveryCodeCheckedEvent:
			wm.RecoveryCodeCheckedAt = e.CheckedAt
		case *session.TokenSetEvent:
			wm.TokenID = e.TokenID
		case *session.LifetimeSetEvent:
			wm.Expiration = e.CreationDate().Add(e.Lifetime)
		case *session.TerminateEvent:
			wm.State = domain.SessionStateTerminated
		}
	}
	return wm.ReadModel.Reduce()
}

func (wm *SessionTokenReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(session.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			session.AddedType,
			session.UserCheckedType,
			session.PasswordCheckedType,
			session.IntentCheckedType,
			session.WebAuthNCheckedType,
			session.TOTPCheckedType,
			session.OTPSMSCheckedType,
			session.OTPEmailCheckedType,
			session.RecoveryCodeCheckedType,
			session.TokenSetType,
			session.LifetimeSetType,
			session.TerminateType,
		).
		Builder()
}

// AuthMethodTypes returns the [domain.UserAuthMethodType] of all succeeded checks of the session.
func (wm *SessionTokenReadModel) AuthMethodTypes() []domain.UserAuthMethodType {
	types := make([]domain.UserAuthMethodType, 0, 7)
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
	if !wm.RecoveryCodeCheckedAt.IsZero() {
		types = append(types, domain.UserAuthMethodTypeRecoveryCode)
	}
	return types
}

// ActiveSessionByToken verifies the session token and returns the state of the session, based on the eventstore.
// The session is regarded as not existing if it was terminated,
// or if its user was deactivated, locked or removed, or the user's organization was deactivated or removed after the user was checked.
// A password check is dropped if the user's password was changed after it.
func (q *Queries) ActiveSessionByToken(ctx context.Context, sessionID, sessionToken string) (model *SessionTokenReadModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	model = newSessionTokenReadModel(sessionID)
	if err = q.eventstore.FilterToQueryReducer(ctx, model); err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Oof3a", "Errors.Internal")
	}
	if model.State != domain.SessionStateActive {
		return nil, zerrors.ThrowNotFound(nil, "QUERY-ohP5e", "Errors.Session.NotExisting")
	}
	if err = q.sessionTokenVerifier(ctx, sessionToken, model.AggregateID, model.TokenID); err != nil {
		return nil, zerrors.ThrowPermissionDenied(nil, "QUERY-ieV6o", "Errors.PermissionDenied")
	}
	if model.UserID == "" {
		return model, nil
	}
	invalidation := &sessionInvalidationModel{
		userID:                  model.UserID,
		userResourceOwner:       model.UserResourceOwner,
		userCheckedPosition:     model.userCheckedPosition,
		passwordCheckedPosition: model.passwordCheckedPosition,
		passwordChecked:         !model.PasswordCheckedAt.IsZero(),
	}
	if err = q.eventstore.FilterToQueryReducer(ctx, invalidation); err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Ahm7i", "Errors.Internal")
	}
	if invalidation.terminated {
		return nil, zerrors.ThrowNotFound(nil, "QUERY-Cie5v", "Errors.Session.NotExisting")
	}
	if invalidation.passwordChanged {
		model.PasswordCheckedAt = time.Time{}
	}
	return model, nil
}

// sessionInvalidationModel searches for user and organization events,
// which invalidate a session or its password check after they occurred.
type sessionInvalidationModel struct {
	eventstore.ReadModel

	userID                  string
	userResourceOwner       string
	userCheckedPosition     decimal.Decimal
	passwordCheckedPosition decimal.Decimal
	passwordChecked         bool

	terminated      bool
	passwordChanged bool
}

func (m *sessionInvalidationModel) Reduce() error {
	for _, event := range m.Events {
		switch event.Type() {
		case user.HumanPasswordChangedType:
			m.passwordChanged = true
		case user.UserDeactivatedType,
			user.UserLockedType,
			user.UserRemovedType,
			org.OrgDeactivatedEventType,
			org.OrgRemovedEventType:
			m.terminated = true
		}
	}
	return m.ReadModel.Reduce()
}

func (m *sessionInvalidationModel) Query() *eventstore.SearchQueryBuilder {
	builder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(m.userID).
		EventTypes(
			user.UserDeactivatedType,
			user.UserLockedType,
			user.UserRemovedType,
		).
		PositionAfter(m.userCheckedPosition).
		Builder()
	if m.userResourceOwner != "" {
		builder = builder.AddQuery().
			AggregateTypes(org.AggregateType).
			AggregateIDs(m.userResourceOwner).
			EventTypes(
				org.OrgDeactivatedEventType,
				org.OrgRemovedEventType,
			).
			PositionAfter(m.userCheckedPosition).
			Builder()
	}
	if m.passwordChecked {
		builder = builder.AddQuery().
			AggregateTypes(user.AggregateType).
			AggregateIDs(m.userID).
			EventTypes(
				user.HumanPasswordChangedType,
			).
			PositionAfter(m.passwordCheckedPosition).
			Builder()
	}
	return builder
}
