package query

import (
	"context"
	"strings"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type OIDCSessionAccessTokenReadModel struct {
	eventstore.WriteModel

	UserID                string
	SessionID             string
	ClientID              string
	Audience              []string
	Scope                 []string
	AuthMethods           []domain.UserAuthMethodType
	AuthTime              time.Time
	State                 domain.OIDCSessionState
	AccessTokenID         string
	AccessTokenCreation   time.Time
	AccessTokenExpiration time.Time
	Reason                domain.TokenReason
	Actor                 *domain.TokenActor
}

func newOIDCSessionAccessTokenReadModel(id string) *OIDCSessionAccessTokenReadModel {
	return &OIDCSessionAccessTokenReadModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: id,
		},
	}
}

func (wm *OIDCSessionAccessTokenReadModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *oidcsession.AddedEvent:
			wm.reduceAdded(e)
		case *oidcsession.AccessTokenAddedEvent:
			wm.reduceAccessTokenAdded(e)
		case *oidcsession.AccessTokenRevokedEvent,
			*oidcsession.RefreshTokenRevokedEvent:
			wm.reduceTokenRevoked(event)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *OIDCSessionAccessTokenReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AwaitOpenTransactions().
		AllowTimeTravel().
		AddQuery().
		AggregateTypes(oidcsession.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			oidcsession.AddedType,
			oidcsession.AccessTokenAddedType,
			oidcsession.AccessTokenRevokedType,
			oidcsession.RefreshTokenRevokedType,
		).
		Builder()
}

func (wm *OIDCSessionAccessTokenReadModel) reduceAdded(e *oidcsession.AddedEvent) {
	wm.UserID = e.UserID
	wm.SessionID = e.SessionID
	wm.ClientID = e.ClientID
	wm.Audience = e.Audience
	wm.Scope = e.Scope
	wm.AuthMethods = e.AuthMethods
	wm.AuthTime = e.AuthTime
	wm.State = domain.OIDCSessionStateActive
}

func (wm *OIDCSessionAccessTokenReadModel) reduceAccessTokenAdded(e *oidcsession.AccessTokenAddedEvent) {
	wm.AccessTokenID = e.ID
	wm.AccessTokenCreation = e.CreationDate()
	wm.AccessTokenExpiration = e.CreationDate().Add(e.Lifetime)
	wm.Reason = e.Reason
	wm.Actor = e.Actor
}

func (wm *OIDCSessionAccessTokenReadModel) reduceTokenRevoked(e eventstore.Event) {
	wm.AccessTokenID = ""
	wm.AccessTokenExpiration = e.CreatedAt()
}

// ActiveAccessTokenByToken will check if the token is active by retrieving the OIDCSession events from the eventstore.
// Refreshed or expired tokens will return an error as well as if the underlying sessions has been terminated.
func (q *Queries) ActiveAccessTokenByToken(ctx context.Context, token string) (model *OIDCSessionAccessTokenReadModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	split := strings.Split(token, "-")
	if len(split) != 2 {
		return nil, zerrors.ThrowPermissionDenied(nil, "QUERY-LJK2W", "Errors.OIDCSession.Token.Invalid")
	}
	model, err = q.accessTokenByOIDCSessionAndTokenID(ctx, split[0], split[1])
	if err != nil {
		return nil, err
	}
	if !model.AccessTokenExpiration.After(time.Now()) {
		return nil, zerrors.ThrowPermissionDenied(nil, "QUERY-SAF3rf", "Errors.OIDCSession.Token.Expired")
	}
	if err = q.checkSessionNotTerminatedAfter(ctx, model.SessionID, model.UserID, model.AccessTokenCreation); err != nil {
		return nil, err
	}
	return model, nil
}

func (q *Queries) accessTokenByOIDCSessionAndTokenID(ctx context.Context, oidcSessionID, tokenID string) (model *OIDCSessionAccessTokenReadModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	model = newOIDCSessionAccessTokenReadModel(oidcSessionID)
	if err = q.eventstore.FilterToQueryReducer(ctx, model); err != nil {
		return nil, zerrors.ThrowPermissionDenied(err, "QUERY-ASfe2", "Errors.OIDCSession.Token.Invalid")
	}
	if model.AccessTokenID != tokenID {
		return nil, zerrors.ThrowPermissionDenied(nil, "QUERY-M2u9w", "Errors.OIDCSession.Token.Invalid")
	}
	return model, nil
}

// checkSessionNotTerminatedAfter checks if a [session.TerminateType] event occurred after a certain time
// and will return an error if so.
func (q *Queries) checkSessionNotTerminatedAfter(ctx context.Context, sessionID, userID string, creation time.Time) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	model := &sessionTerminatedModel{
		sessionID: sessionID,
		creation:  creation,
		userID:    userID,
	}
	err = q.eventstore.FilterToQueryReducer(ctx, model)
	if err != nil {
		return zerrors.ThrowPermissionDenied(err, "QUERY-SJ642", "Errors.Internal")
	}

	if model.terminated {
		return zerrors.ThrowPermissionDenied(nil, "QUERY-IJL3H", "Errors.OIDCSession.Token.Invalid")
	}
	return nil
}

type sessionTerminatedModel struct {
	creation  time.Time
	sessionID string
	userID    string

	events     int
	terminated bool
}

func (s *sessionTerminatedModel) Reduce() error {
	s.terminated = s.events > 0
	return nil
}

func (s *sessionTerminatedModel) AppendEvents(events ...eventstore.Event) {
	s.events += len(events)
}

func (s *sessionTerminatedModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AwaitOpenTransactions().
		CreationDateAfter(s.creation).
		AddQuery().
		AggregateTypes(session.AggregateType).
		AggregateIDs(s.sessionID).
		EventTypes(
			session.TerminateType,
		).
		Builder()
	if s.userID == "" {
		return query
	}
	return query.
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(s.userID).
		EventTypes(
			user.UserDeactivatedType,
			user.UserLockedType,
			user.UserRemovedType,
		).
		Builder()
}
