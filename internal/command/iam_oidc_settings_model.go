package command

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMOIDCSettingsWriteModel struct {
	eventstore.WriteModel

	AccessTokenLifetime        time.Duration
	IdTokenLifetime            time.Duration
	RefreshTokenIdleExpiration time.Duration
	RefreshTokenExpiration     time.Duration
	State                      domain.OIDCSettingsState
}

func NewIAMOIDCSettingsWriteModel() *IAMOIDCSettingsWriteModel {
	return &IAMOIDCSettingsWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   domain.IAMID,
			ResourceOwner: domain.IAMID,
		},
	}
}

func (wm *IAMOIDCSettingsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *iam.OIDCSettingsAddedEvent:
			wm.AccessTokenLifetime = e.AccessTokenLifetime
			wm.IdTokenLifetime = e.IdTokenLifetime
			wm.RefreshTokenIdleExpiration = e.RefreshTokenIdleExpiration
			wm.RefreshTokenExpiration = e.RefreshTokenExpiration
			wm.State = domain.OIDCSettingsStateActive
		case *iam.OIDCSettingsChangedEvent:
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

func (wm *IAMOIDCSettingsWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			iam.OIDCSettingsAddedEventType,
			iam.OIDCSettingsChangedEventType).
		Builder()
}

func (wm *IAMOIDCSettingsWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	accessTokenLifetime,
	idTokenLifetime,
	refreshTokenIdleExpiration,
	refreshTokenExpiration time.Duration,
) (*iam.OIDCSettingsChangedEvent, bool, error) {
	changes := make([]iam.OIDCSettingsChanges, 0, 4)
	var err error

	if wm.AccessTokenLifetime != accessTokenLifetime {
		changes = append(changes, iam.ChangeOIDCSettingsAccessTokenLifetime(accessTokenLifetime))
	}
	if wm.IdTokenLifetime != idTokenLifetime {
		changes = append(changes, iam.ChangeOIDCSettingsIdTokenLifetime(idTokenLifetime))
	}
	if wm.RefreshTokenIdleExpiration != refreshTokenIdleExpiration {
		changes = append(changes, iam.ChangeOIDCSettingsRefreshTokenIdleExpiration(refreshTokenIdleExpiration))
	}
	if wm.RefreshTokenExpiration != refreshTokenExpiration {
		changes = append(changes, iam.ChangeOIDCSettingsRefreshTokenExpiration(refreshTokenExpiration))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := iam.NewOIDCSettingsChangeEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
