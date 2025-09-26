package command

import (
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type OIDCSessionWriteModel struct {
	eventstore.WriteModel

	UserID                     string
	UserResourceOwner          string
	PreferredLanguage          *language.Tag
	SessionID                  string
	ClientID                   string
	Audience                   []string
	Scope                      []string
	AuthMethods                []domain.UserAuthMethodType
	AuthTime                   time.Time
	Nonce                      string
	UserAgent                  *domain.UserAgent
	State                      domain.OIDCSessionState
	AccessTokenID              string
	AccessTokenCreation        time.Time
	AccessTokenExpiration      time.Time
	AccessTokenReason          domain.TokenReason
	AccessTokenActor           *domain.TokenActor
	RefreshTokenID             string
	RefreshToken               string
	RefreshTokenExpiration     time.Time
	RefreshTokenIdleExpiration time.Time

	aggregate *eventstore.Aggregate
}

func NewOIDCSessionWriteModel(id string, resourceOwner string) *OIDCSessionWriteModel {
	return &OIDCSessionWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   id,
			ResourceOwner: resourceOwner,
		},
		aggregate: &oidcsession.NewAggregate(id, resourceOwner).Aggregate,
	}
}

func (wm *OIDCSessionWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *oidcsession.AddedEvent:
			wm.reduceAdded(e)
		case *oidcsession.TerminatedEvent:
			wm.reduceTerminated(e)
		case *oidcsession.AccessTokenAddedEvent:
			wm.reduceAccessTokenAdded(e)
		case *oidcsession.AccessTokenRevokedEvent:
			wm.reduceAccessTokenRevoked(e)
		case *oidcsession.RefreshTokenAddedEvent:
			wm.reduceRefreshTokenAdded(e)
		case *oidcsession.RefreshTokenRenewedEvent:
			wm.reduceRefreshTokenRenewed(e)
		case *oidcsession.RefreshTokenRevokedEvent:
			wm.reduceRefreshTokenRevoked(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *OIDCSessionWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(oidcsession.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			oidcsession.AddedType,
			oidcsession.TerminatedType,
			oidcsession.AccessTokenAddedType,
			oidcsession.RefreshTokenAddedType,
			oidcsession.RefreshTokenRenewedType,
			oidcsession.RefreshTokenRevokedType,
		).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func (wm *OIDCSessionWriteModel) reduceAdded(e *oidcsession.AddedEvent) {
	wm.UserID = e.UserID
	wm.UserResourceOwner = e.UserResourceOwner
	wm.SessionID = e.SessionID
	wm.ClientID = e.ClientID
	wm.Audience = e.Audience
	wm.Scope = e.Scope
	wm.AuthMethods = e.AuthMethods
	wm.AuthTime = e.AuthTime
	wm.Nonce = e.Nonce
	wm.PreferredLanguage = e.PreferredLanguage
	wm.UserAgent = e.UserAgent
	wm.State = domain.OIDCSessionStateActive
	// the write model might be initialized without resource owner,
	// so update the aggregate
	if wm.ResourceOwner == "" {
		wm.aggregate = &oidcsession.NewAggregate(wm.AggregateID, e.Aggregate().ResourceOwner).Aggregate
	}
}

func (wm *OIDCSessionWriteModel) reduceTerminated(e *oidcsession.TerminatedEvent) {
	wm.State = domain.OIDCSessionStateTerminated
}

func (wm *OIDCSessionWriteModel) reduceAccessTokenAdded(e *oidcsession.AccessTokenAddedEvent) {
	wm.AccessTokenID = e.ID
	wm.AccessTokenExpiration = e.CreationDate().Add(e.Lifetime)
	wm.AccessTokenReason = e.Reason
	wm.AccessTokenActor = e.Actor
}

func (wm *OIDCSessionWriteModel) reduceAccessTokenRevoked(e *oidcsession.AccessTokenRevokedEvent) {
	wm.AccessTokenID = ""
	wm.AccessTokenExpiration = e.CreationDate()
}

func (wm *OIDCSessionWriteModel) reduceRefreshTokenAdded(e *oidcsession.RefreshTokenAddedEvent) {
	wm.RefreshTokenID = e.ID
	wm.RefreshTokenExpiration = e.CreationDate().Add(e.Lifetime)
	wm.RefreshTokenIdleExpiration = e.CreationDate().Add(e.IdleLifetime)
}

func (wm *OIDCSessionWriteModel) reduceRefreshTokenRenewed(e *oidcsession.RefreshTokenRenewedEvent) {
	wm.RefreshTokenID = e.ID
	wm.RefreshTokenIdleExpiration = e.CreationDate().Add(e.IdleLifetime)
}

func (wm *OIDCSessionWriteModel) reduceRefreshTokenRevoked(e *oidcsession.RefreshTokenRevokedEvent) {
	wm.RefreshTokenID = ""
	wm.RefreshTokenExpiration = e.CreationDate()
	wm.RefreshTokenIdleExpiration = e.CreationDate()
	wm.AccessTokenID = ""
	wm.AccessTokenExpiration = e.CreationDate()
}

func (wm *OIDCSessionWriteModel) CheckRefreshToken(refreshTokenID string) error {
	if wm.State != domain.OIDCSessionStateActive {
		return zerrors.ThrowPreconditionFailed(nil, "OIDCS-s3hjk", "Errors.OIDCSession.RefreshTokenInvalid")
	}
	if wm.RefreshTokenID != refreshTokenID {
		return zerrors.ThrowPreconditionFailed(nil, "OIDCS-28ubl", "Errors.OIDCSession.RefreshTokenInvalid")
	}
	now := time.Now()
	if wm.RefreshTokenExpiration.Before(now) || wm.RefreshTokenIdleExpiration.Before(now) {
		return zerrors.ThrowPreconditionFailed(nil, "OIDCS-3jt2w", "Errors.OIDCSession.RefreshTokenInvalid")
	}
	return nil
}

func (wm *OIDCSessionWriteModel) CheckAccessToken(accessTokenID string) error {
	if wm.State != domain.OIDCSessionStateActive {
		return zerrors.ThrowPreconditionFailed(nil, "OIDCS-KL2pk", "Errors.OIDCSession.Token.Invalid")
	}
	if wm.AccessTokenID != accessTokenID {
		return zerrors.ThrowPreconditionFailed(nil, "OIDCS-JLKW2", "Errors.OIDCSession.Token.Invalid")
	}
	if wm.AccessTokenExpiration.Before(time.Now()) {
		return zerrors.ThrowPreconditionFailed(nil, "OIDCS-3j3md", "Errors.OIDCSession.Token.Invalid")
	}
	return nil
}

func (wm *OIDCSessionWriteModel) CheckClient(clientID string) error {
	for _, aud := range wm.Audience {
		if aud == clientID {
			return nil
		}
	}
	return zerrors.ThrowPreconditionFailed(nil, "OIDCS-SKjl3", "Errors.OIDCSession.InvalidClient")
}

func (wm *OIDCSessionWriteModel) OIDCRefreshTokenID(refreshTokenID string) string {
	return wm.AggregateID + TokenDelimiter + refreshTokenID
}
