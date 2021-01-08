package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

type MachineWriteModel struct {
	eventstore.WriteModel

	UserName string

	Name        string
	Description string
	UserState   domain.UserState
}

func NewMachineWriteModel(userID string) *MachineWriteModel {
	return &MachineWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: userID,
		},
	}
}

func (wm *MachineWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.MachineAddedEvent:
			wm.AppendEvents(e)
		case *user.UsernameChangedEvent:
			wm.AppendEvents(e)
		case *user.MachineChangedEvent:
			wm.AppendEvents(e)
		case *user.UserDeactivatedEvent:
			wm.AppendEvents(e)
		case *user.UserReactivatedEvent:
			wm.AppendEvents(e)
		case *user.UserLockedEvent:
			wm.AppendEvents(e)
		case *user.UserUnlockedEvent:
			wm.AppendEvents(e)
		case *user.UserRemovedEvent:
			wm.AppendEvents(e)
		}
	}
}

//TODO: Compute OTPState? initial/active
func (wm *MachineWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.MachineAddedEvent:
			wm.UserName = e.UserName
			wm.Name = e.Name
			wm.Description = e.Description
			wm.UserState = domain.UserStateActive
		case *user.UsernameChangedEvent:
			wm.UserName = e.UserName
		case *user.MachineChangedEvent:
			if e.Name != nil {
				wm.Name = *e.Name
			}
			if e.Description != nil {
				wm.Description = *e.Description
			}
		case *user.UserLockedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateLocked
			}
		case *user.UserUnlockedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateActive
			}
		case *user.UserDeactivatedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateInactive
			}
		case *user.UserReactivatedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateActive
			}
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *MachineWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID)
}

func (wm *MachineWriteModel) NewChangedEvent(
	ctx context.Context,
	name,
	description string,
) (*user.MachineChangedEvent, bool) {
	hasChanged := false
	changedEvent := user.NewMachineChangedEvent(ctx)
	if wm.Name != name {
		hasChanged = true
		changedEvent.Name = &name
	}
	if wm.Description != description {
		hasChanged = true
		changedEvent.Description = &description
	}
	return changedEvent, hasChanged
}
