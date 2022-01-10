package command

import (
	"time"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/user"
)

type MachineTokenWriteModel struct {
	eventstore.WriteModel

	TokenID        string
	ExpirationDate time.Time

	State domain.MachineTokenState
}

func NewMachineTokenWriteModel(userID, tokenID, resourceOwner string) *MachineTokenWriteModel {
	return &MachineTokenWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		TokenID: tokenID,
	}
}

func (wm *MachineTokenWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.MachineTokenAddedEvent:
			if wm.TokenID != e.TokenID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.MachineTokenRemovedEvent:
			if wm.TokenID != e.TokenID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *user.UserRemovedEvent:
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *MachineTokenWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.MachineTokenAddedEvent:
			wm.TokenID = e.TokenID
			wm.ExpirationDate = e.Expiration
			wm.State = domain.MachineTokenStateActive
		case *user.MachineTokenRemovedEvent:
			wm.State = domain.MachineTokenStateRemoved
		case *user.UserRemovedEvent:
			wm.State = domain.MachineTokenStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *MachineTokenWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			user.MachineTokenAddedType,
			user.MachineTokenRemovedType,
			user.UserRemovedType).
		Builder()
}

func (wm *MachineTokenWriteModel) Exists() bool {
	return wm.State != domain.MachineTokenStateUnspecified && wm.State != domain.MachineTokenStateRemoved
}
