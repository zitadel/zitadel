package command

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
)

type OIDCSessionWriteModel struct {
	eventstore.WriteModel

	UserID                     string
	SessionID                  string
	ClientID                   string
	Audience                   []string
	Scope                      []string
	AuthMethods                []domain.UserAuthMethodType
	AuthTime                   time.Time
	State                      domain.OIDCSessionState
	AccessTokenCreation        time.Time
	AccessTokenExpiration      time.Time
	RefreshTokenID             string
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
		case *oidcsession.AccessTokenAddedEvent:
			wm.reduceAccessTokenAdded(e)
		case *oidcsession.RefreshTokenAddedEvent:
			wm.reduceRefreshTokenAdded(e)
		case *oidcsession.RefreshTokenRenewedEvent:
			wm.reduceRefreshTokenRenewed(e)
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
			oidcsession.AccessTokenAddedType,
			oidcsession.RefreshTokenAddedType,
			oidcsession.RefreshTokenRenewedType,
		).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func (wm *OIDCSessionWriteModel) reduceAdded(e *oidcsession.AddedEvent) {
	wm.UserID = e.UserID
	wm.SessionID = e.SessionID
	wm.ClientID = e.ClientID
	wm.Audience = e.Audience
	wm.Scope = e.Scope
	wm.AuthMethods = e.AuthMethods
	wm.AuthTime = e.AuthTime
	wm.State = domain.OIDCSessionStateActive
	if wm.ResourceOwner == "" {
		wm.ResourceOwner = e.Aggregate().ResourceOwner
		wm.aggregate = &oidcsession.NewAggregate(wm.AggregateID, e.Aggregate().ResourceOwner).Aggregate
	}
}

func (wm *OIDCSessionWriteModel) reduceAccessTokenAdded(e *oidcsession.AccessTokenAddedEvent) {
	wm.AccessTokenExpiration = e.CreationDate().Add(e.Lifetime)
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

func (wm *OIDCSessionWriteModel) CheckRefreshToken(refreshTokenID string) error {
	if wm.State != domain.OIDCSessionStateActive {
		return caos_errs.ThrowPreconditionFailed(nil, "OIDCS-s3hjk", "Errors.OIDCSession.RefreshTokenInvalid")
	}
	if wm.RefreshTokenID != refreshTokenID {
		return caos_errs.ThrowPreconditionFailed(nil, "OIDCS-28ubl", "Errors.OIDCSession.RefreshTokenInvalid")
	}
	now := time.Now()
	if wm.RefreshTokenExpiration.Before(now) || wm.RefreshTokenIdleExpiration.Before(now) {
		return caos_errs.ThrowPreconditionFailed(nil, "OIDCS-3jt2w", "Errors.OIDCSession.RefreshTokenInvalid")
	}
	return nil
}
