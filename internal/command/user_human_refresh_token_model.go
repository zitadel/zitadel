package command

import (
	"time"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/user"
)

type HumanRefreshTokenWriteModel struct {
	eventstore.WriteModel

	TokenID      string
	RefreshToken string

	UserState      domain.UserState
	IdleExpiration time.Time
	Expiration     time.Time
}

func NewHumanRefreshTokenWriteModel(userID, resourceOwner, tokenID string) *HumanRefreshTokenWriteModel {
	return &HumanRefreshTokenWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		TokenID: tokenID,
	}
}

func (wm *HumanRefreshTokenWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanRefreshTokenAddedEvent:
			if wm.TokenID != e.TokenID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.HumanRefreshTokenRenewedEvent:
			if wm.TokenID != e.TokenID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.HumanRefreshTokenRemovedEvent:
			if wm.TokenID != e.TokenID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *HumanRefreshTokenWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanRefreshTokenAddedEvent:
			wm.TokenID = e.TokenID
			wm.RefreshToken = e.TokenID
			wm.IdleExpiration = e.CreationDate().Add(e.IdleExpiration)
			wm.Expiration = e.CreationDate().Add(e.Expiration)
			wm.UserState = domain.UserStateActive
		case *user.HumanRefreshTokenRenewedEvent:
			if wm.UserState == domain.UserStateActive {
				wm.RefreshToken = e.RefreshToken
			}
			wm.RefreshToken = e.RefreshToken
			wm.IdleExpiration = e.CreationDate().Add(e.IdleExpiration)
		case *user.HumanRefreshTokenRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanRefreshTokenWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			user.HumanRefreshTokenAddedType,
			user.HumanRefreshTokenRenewedType,
			user.HumanRefreshTokenRemovedType,
			user.UserRemovedType).
		SearchQueryBuilder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}
