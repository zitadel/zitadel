package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user/authenticator"
	"github.com/zitadel/zitadel/internal/repository/user/schemauser"
)

type PasswordV3WriteModel struct {
	eventstore.WriteModel
	UserID string

	EncodedHash    string
	ChangeRequired bool
}

func NewPasswordV3WriteModel(resourceOwner, id string) *PasswordV3WriteModel {
	return &PasswordV3WriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   id,
			ResourceOwner: resourceOwner,
		},
		UserID: id,
	}
}

func (wm *PasswordV3WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *authenticator.PasswordCreatedEvent:
			wm.UserID = e.UserID
			wm.EncodedHash = e.EncodedHash
			wm.ChangeRequired = e.ChangeRequired
		case *authenticator.PasswordDeletedEvent:
			wm.UserID = ""
			wm.EncodedHash = ""
			wm.ChangeRequired = false
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *PasswordV3WriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(schemauser.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			authenticator.PasswordCreatedType,
			authenticator.PasswordDeletedType,
		).Builder()
}
