package command

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type PersonalAccessTokenWriteModel struct {
	eventstore.WriteModel

	TokenID        string
	ExpirationDate time.Time

	State domain.PersonalAccessTokenState
}

func NewPersonalAccessTokenWriteModel(userID, tokenID, resourceOwner string) *PersonalAccessTokenWriteModel {
	return &PersonalAccessTokenWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		TokenID: tokenID,
	}
}

func (wm *PersonalAccessTokenWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.PersonalAccessTokenAddedEvent:
			if wm.TokenID != e.TokenID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.PersonalAccessTokenRemovedEvent:
			if wm.TokenID != e.TokenID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.UserRemovedEvent:
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *PersonalAccessTokenWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.PersonalAccessTokenAddedEvent:
			wm.TokenID = e.TokenID
			wm.ExpirationDate = e.Expiration
			wm.State = domain.PersonalAccessTokenStateActive
		case *user.PersonalAccessTokenRemovedEvent:
			wm.State = domain.PersonalAccessTokenStateRemoved
		case *user.UserRemovedEvent:
			wm.State = domain.PersonalAccessTokenStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *PersonalAccessTokenWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			user.PersonalAccessTokenAddedType,
			user.PersonalAccessTokenRemovedType,
			user.UserRemovedType).
		Builder()
}

func (wm *PersonalAccessTokenWriteModel) Exists() bool {
	return wm.State != domain.PersonalAccessTokenStateUnspecified && wm.State != domain.PersonalAccessTokenStateRemoved
}
