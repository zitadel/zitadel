package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/user"
)

type HumanExternalIDPWriteModel struct {
	eventstore.WriteModel

	IDPConfigID    string
	ExternalUserID string
	DisplayName    string

	State domain.ExternalIDPState
}

func NewHumanExternalIDPWriteModel(userID, idpConfigID, externalUserID, resourceOwner string) *HumanExternalIDPWriteModel {
	return &HumanExternalIDPWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		IDPConfigID:    idpConfigID,
		ExternalUserID: externalUserID,
	}
}

func (wm *HumanExternalIDPWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanExternalIDPAddedEvent:
			if e.IDPConfigID != wm.IDPConfigID && e.ExternalUserID != wm.ExternalUserID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.HumanExternalIDPRemovedEvent:
			if e.IDPConfigID != wm.IDPConfigID && e.ExternalUserID != wm.ExternalUserID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.HumanExternalIDPCascadeRemovedEvent:
			if e.IDPConfigID != wm.IDPConfigID && e.ExternalUserID != wm.ExternalUserID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.UserRemovedEvent:
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *HumanExternalIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanExternalIDPAddedEvent:
			wm.IDPConfigID = e.IDPConfigID
			wm.DisplayName = e.DisplayName
			wm.ExternalUserID = e.ExternalUserID
			wm.State = domain.ExternalIDPStateActive
		case *user.HumanExternalIDPRemovedEvent:
			wm.State = domain.ExternalIDPStateRemoved
		case *user.HumanExternalIDPCascadeRemovedEvent:
			wm.State = domain.ExternalIDPStateRemoved
		case *user.UserRemovedEvent:
			wm.State = domain.ExternalIDPStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanExternalIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(user.HumanExternalIDPAddedType,
			user.HumanExternalIDPRemovedType,
			user.HumanExternalIDPCascadeRemovedType,
			user.UserRemovedType).
		Builder()
}
