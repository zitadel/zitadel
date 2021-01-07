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

	State domain.ExternalIDPState
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
		case *user.UserRemovedEvent:
			wm.AppendEvents(e)
		}
	}
}

func (wm *HumanExternalIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanExternalIDPAddedEvent:
			wm.IDPConfigID = e.IDPConfigID
			wm.DisplayName = e.DisplayName
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
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID)
}
