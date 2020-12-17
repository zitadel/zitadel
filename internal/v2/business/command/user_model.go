package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/business/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

type UserWriteModel struct {
	eventstore.WriteModel

	UserName  string
	UserState domain.UserState
}

func NewUserWriteModel(userID string) *UserWriteModel {
	return &UserWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: userID,
		},
	}
}

func (wm *UserWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanEmailChangedEvent:
			wm.AppendEvents(e)
		case *user.HumanEmailVerifiedEvent:
			wm.AppendEvents(e)
		case *user.HumanAddedEvent, *user.HumanRegisteredEvent:
			wm.AppendEvents(e)
		case *user.MachineAddedEvent:
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

//TODO: Compute State? initial/active
func (wm *UserWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.UserName = e.UserName
			wm.UserState = domain.UserStateInitial
		case *user.HumanRegisteredEvent:
			wm.UserName = e.UserName
			wm.UserState = domain.UserStateInitial

		case *user.MachineAddedEvent:
			wm.UserName = e.UserName
			wm.UserState = domain.UserStateActive
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

func (wm *UserWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID)
}

func UserAggregateFromWriteModel(wm *eventstore.WriteModel) *user.Aggregate {
	return &user.Aggregate{
		Aggregate: *eventstore.AggregateFromWriteModel(wm, user.AggregateType, user.AggregateVersion),
	}
}
