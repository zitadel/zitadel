package command

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMOIDCConfigWriteModel struct {
	eventstore.WriteModel

	AccessTokenLifetime        time.Duration
	IdTokenLifetime            time.Duration
	RefreshTokenIdleExpiration time.Duration
	RefreshTokenExpiration     time.Duration
	State                      domain.OIDCConfigState
}

func NewIAMOIDCConfigWriteModel() *IAMOIDCConfigWriteModel {
	return &IAMOIDCConfigWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   domain.IAMID,
			ResourceOwner: domain.IAMID,
		},
	}
}

func (wm *IAMOIDCConfigWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *iam.OIDCConfigAddedEvent:
			wm.AccessTokenLifetime = e.AccessTokenLifetime
			wm.IdTokenLifetime = e.IdTokenLifetime
			wm.RefreshTokenIdleExpiration = e.RefreshTokenIdleExpiration
			wm.RefreshTokenExpiration = e.RefreshTokenExpiration
			wm.State = domain.OIDCConfigStateActive
		case *iam.OIDCConfigChangedEvent:
			if e.AccessTokenLifetime != nil {
				wm.AccessTokenLifetime = *e.AccessTokenLifetime
			}
			if e.IdTokenLifetime != nil {
				wm.IdTokenLifetime = *e.IdTokenLifetime
			}
			if e.RefreshTokenIdleExpiration != nil {
				wm.RefreshTokenIdleExpiration = *e.RefreshTokenIdleExpiration
			}
			if e.RefreshTokenExpiration != nil {
				wm.RefreshTokenExpiration = *e.RefreshTokenExpiration
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *IAMOIDCConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			iam.OIDCConfigAddedEventType,
			iam.OIDCConfigChangedEventType).
		Builder()
}

func (wm *IAMOIDCConfigWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	accessTokenLifetime,
	idTokenLifetime,
	refreshTokenIdleExpiration,
	refreshTokenExpiration time.Duration,
) (*iam.OIDCConfigChangedEvent, bool, error) {
	changes := make([]iam.OIDCConfigChanges, 0)
	var err error

	if wm.AccessTokenLifetime != accessTokenLifetime {
		changes = append(changes, iam.ChangeOIDCConfigAccessTokenLifetime(accessTokenLifetime))
	}
	if wm.IdTokenLifetime != idTokenLifetime {
		changes = append(changes, iam.ChangeOIDCConfigIdTokenLifetime(idTokenLifetime))
	}
	if wm.RefreshTokenIdleExpiration != refreshTokenIdleExpiration {
		changes = append(changes, iam.ChangeOIDCConfigRefreshTokenIdleExpiration(refreshTokenIdleExpiration))
	}
	if wm.RefreshTokenExpiration != refreshTokenExpiration {
		changes = append(changes, iam.ChangeOIDCConfigRefreshTokenExpiration(refreshTokenExpiration))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := iam.NewOIDCConfigChangeEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
