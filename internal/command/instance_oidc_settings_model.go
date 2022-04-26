package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceOIDCSettingsWriteModel struct {
	eventstore.WriteModel

	AccessTokenLifetime        time.Duration
	IdTokenLifetime            time.Duration
	RefreshTokenIdleExpiration time.Duration
	RefreshTokenExpiration     time.Duration
	State                      domain.OIDCSettingsState
}

func NewInstanceOIDCSettingsWriteModel(ctx context.Context) *InstanceOIDCSettingsWriteModel {
	return &InstanceOIDCSettingsWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   authz.GetInstance(ctx).InstanceID(),
			ResourceOwner: authz.GetInstance(ctx).InstanceID(),
		},
	}
}

func (wm *InstanceOIDCSettingsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *instance.OIDCSettingsAddedEvent:
			wm.AccessTokenLifetime = e.AccessTokenLifetime
			wm.IdTokenLifetime = e.IdTokenLifetime
			wm.RefreshTokenIdleExpiration = e.RefreshTokenIdleExpiration
			wm.RefreshTokenExpiration = e.RefreshTokenExpiration
			wm.State = domain.OIDCSettingsStateActive
		case *instance.OIDCSettingsChangedEvent:
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

func (wm *InstanceOIDCSettingsWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.OIDCSettingsAddedEventType,
			instance.OIDCSettingsChangedEventType).
		Builder()
}

func (wm *InstanceOIDCSettingsWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	accessTokenLifetime,
	idTokenLifetime,
	refreshTokenIdleExpiration,
	refreshTokenExpiration time.Duration,
) (*instance.OIDCSettingsChangedEvent, bool, error) {
	changes := make([]instance.OIDCSettingsChanges, 0, 4)
	var err error

	if wm.AccessTokenLifetime != accessTokenLifetime {
		changes = append(changes, instance.ChangeOIDCSettingsAccessTokenLifetime(accessTokenLifetime))
	}
	if wm.IdTokenLifetime != idTokenLifetime {
		changes = append(changes, instance.ChangeOIDCSettingsIdTokenLifetime(idTokenLifetime))
	}
	if wm.RefreshTokenIdleExpiration != refreshTokenIdleExpiration {
		changes = append(changes, instance.ChangeOIDCSettingsRefreshTokenIdleExpiration(refreshTokenIdleExpiration))
	}
	if wm.RefreshTokenExpiration != refreshTokenExpiration {
		changes = append(changes, instance.ChangeOIDCSettingsRefreshTokenExpiration(refreshTokenExpiration))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := instance.NewOIDCSettingsChangeEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
