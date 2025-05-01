package command

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
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

func (wm *MachineKeyWriteModel) AppendEvents(events ...eventstore.Event) {
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
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			user.MachineKeyAddedEventType,
			user.MachineKeyRemovedEventType,
			user.UserRemovedType).
		Builder()
}

func (wm *MachineKeyWriteModel) Exists() bool {
	return wm.State != domain.MachineKeyStateUnspecified && wm.State != domain.MachineKeyStateRemoved
}
