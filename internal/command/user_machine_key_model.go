package command

import (
	"github.com/caos/zitadel/internal/eventstore"
	"time"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/user"
)

type MachineKeyWriteModel struct {
	eventstore.WriteModel

	KeyID          string
	KeyType        domain.AuthNKeyType
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
			//TODO: adlerhurst we should decide who should handle the correct event appending
			// IMO in this append events we should only get events with the correct keyID
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

func (wm *MachineKeyWriteModel) Exists() bool {
	return wm.State != domain.MachineKeyStateUnspecified && wm.State != domain.MachineKeyStateRemoved
}
