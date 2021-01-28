package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
	"time"
)

type MachineKeyWriteModel struct {
	eventstore.WriteModel

	KeyID          string
	KeyType        domain.MachineKeyType
	ExpirationDate time.Time

	State domain.MachineKeyState
}

func NewMachineKeyWriteModel(userID, keyID, resourceOwner string) *MachineKeyWriteModel {
	return &MachineKeyWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		KeyID: keyID,
	}
}

func (wm *MachineKeyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.MachineKeyAddedEvent:
			if wm.KeyID != e.KeyID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.MachineKeyRemovedEvent:
			if wm.KeyID != e.KeyID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.UserRemovedEvent:
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *MachineKeyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.MachineKeyAddedEvent:
			wm.KeyID = e.KeyID
			wm.KeyType = e.KeyType
			wm.ExpirationDate = e.ExpirationDate
			wm.State = domain.MachineKeyStateActive
		case *user.MachineKeyRemovedEvent:
			wm.State = domain.MachineKeyStateRemoved
		case *user.UserRemovedEvent:
			wm.State = domain.MachineKeyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *MachineKeyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			user.MachineKeyAddedEventType,
			user.MachineKeyRemovedEventType,
			user.UserRemovedType)
}
