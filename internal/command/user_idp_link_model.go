package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type UserIDPLinkWriteModel struct {
	eventstore.WriteModel

	IDPConfigID    string
	ExternalUserID string
	DisplayName    string

	State domain.UserIDPLinkState
}

func NewUserIDPLinkWriteModel(userID, idpConfigID, externalUserID, resourceOwner string) *UserIDPLinkWriteModel {
	return &UserIDPLinkWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		IDPConfigID:    idpConfigID,
		ExternalUserID: externalUserID,
	}
}

func (wm *UserIDPLinkWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.UserIDPLinkAddedEvent:
			if e.IDPConfigID != wm.IDPConfigID || e.ExternalUserID != wm.ExternalUserID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.UserIDPExternalIDMigratedEvent:
			if e.IDPConfigID != wm.IDPConfigID || e.PreviousID != wm.ExternalUserID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.UserIDPLinkRemovedEvent:
			if e.IDPConfigID != wm.IDPConfigID || e.ExternalUserID != wm.ExternalUserID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.UserIDPLinkCascadeRemovedEvent:
			if e.IDPConfigID != wm.IDPConfigID || e.ExternalUserID != wm.ExternalUserID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.UserRemovedEvent:
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *UserIDPLinkWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.UserIDPLinkAddedEvent:
			wm.IDPConfigID = e.IDPConfigID
			wm.DisplayName = e.DisplayName
			wm.ExternalUserID = e.ExternalUserID
			wm.State = domain.UserIDPLinkStateActive
		case *user.UserIDPExternalIDMigratedEvent:
			wm.ExternalUserID = e.NewID
		case *user.UserIDPLinkRemovedEvent:
			wm.State = domain.UserIDPLinkStateRemoved
		case *user.UserIDPLinkCascadeRemovedEvent:
			wm.State = domain.UserIDPLinkStateRemoved
		case *user.UserRemovedEvent:
			wm.State = domain.UserIDPLinkStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *UserIDPLinkWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(user.UserIDPLinkAddedType,
			user.UserIDPExternalIDMigratedType,
			user.UserIDPLinkRemovedType,
			user.UserIDPLinkCascadeRemovedType,
			user.UserRemovedType).
		Builder()
}
