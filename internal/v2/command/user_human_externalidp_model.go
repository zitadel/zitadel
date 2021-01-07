package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

type HumanExternalIDPWriteModel struct {
	eventstore.WriteModel

	IDPConfigID    string
	ExternalUserID string
	DisplayName    string

	UserState domain.UserState
}

func NewHumanExternalIDPWriteModel(userID, idpConfigID, externalUserID string) *HumanExternalIDPWriteModel {
	return &HumanExternalIDPWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: userID,
		},
		IDPConfigID:    idpConfigID,
		ExternalUserID: externalUserID,
	}
}

func (wm *HumanExternalIDPWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanExternalIDPAddedEvent:
			if wm.IDPConfigID == e.IDPConfigID && wm.ExternalUserID == e.UserID {
				wm.AppendEvents(e)
			}
		case *user.HumanExternalIDPRemovedEvent:
			if wm.IDPConfigID == e.IDPConfigID {
				wm.AppendEvents(e)
			}
		case *user.HumanExternalIDPCascadeRemovedEvent:
			if wm.IDPConfigID == e.IDPConfigID {
				wm.AppendEvents(e)
			}
		case *user.HumanAddedEvent, *user.HumanRegisteredEvent:
			wm.AppendEvents(e)
		case *user.UserRemovedEvent:
			wm.AppendEvents(e)
		}
	}
}

func (wm *HumanExternalIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.UserState = domain.UserStateActive
		case *user.HumanRegisteredEvent:
			wm.UserState = domain.UserStateActive
		case *user.HumanExternalIDPAddedEvent:
			wm.IDPConfigID = e.IDPConfigID
			wm.DisplayName = e.DisplayName
		case *user.HumanExternalIDPRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		case *user.HumanExternalIDPCascadeRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanExternalIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID)
}
