package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user/authenticator"
)

type UsernameV3WriteModel struct {
	eventstore.WriteModel
	UserID        string
	Username      string
	IsOrgSpecific bool
}

func NewUsernameV3WriteModel(resourceOwner, userID, id string) *UsernameV3WriteModel {
	return &UsernameV3WriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   id,
			ResourceOwner: resourceOwner,
		},
		UserID: userID,
	}
}

func (wm *UsernameV3WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *authenticator.UsernameCreatedEvent:
			if e.UserID != wm.UserID {
				continue
			}
			wm.UserID = e.UserID
			wm.Username = e.Username
			wm.IsOrgSpecific = e.IsOrgSpecific
		case *authenticator.UsernameDeletedEvent:
			wm.UserID = ""
			wm.Username = ""
			wm.IsOrgSpecific = false
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *UsernameV3WriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(authenticator.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			authenticator.UsernameCreatedType,
			authenticator.UsernameDeletedType,
		).Builder()
}
