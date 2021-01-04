package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

type HumanEmailWriteModel struct {
	eventstore.WriteModel

	Email           string
	IsEmailVerified bool

	UserState domain.UserState
}

func NewHumanEmailWriteModel(userID string) *HumanEmailWriteModel {
	return &HumanEmailWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: userID,
		},
	}
}

func (wm *HumanEmailWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanEmailChangedEvent:
			wm.AppendEvents(e)
		case *user.HumanEmailVerifiedEvent:
			wm.AppendEvents(e)
		case *user.HumanAddedEvent, *user.HumanRegisteredEvent:
			wm.AppendEvents(e)
		case *user.UserRemovedEvent:
			wm.AppendEvents(e)
		}
	}
}

func (wm *HumanEmailWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.Email = e.EmailAddress
			wm.UserState = domain.UserStateActive
		case *user.HumanRegisteredEvent:
			wm.Email = e.EmailAddress
			wm.UserState = domain.UserStateActive
		case *user.HumanEmailChangedEvent:
			wm.Email = e.EmailAddress
			wm.IsEmailVerified = false
		case *user.HumanEmailVerifiedEvent:
			wm.IsEmailVerified = true
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanEmailWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID)
}

func (wm *HumanEmailWriteModel) NewChangedEvent(
	ctx context.Context,
	email string,
) (*user.HumanEmailChangedEvent, bool) {
	hasChanged := false
	changedEvent := user.NewHumanEmailChangedEvent(ctx)
	if wm.Email != email {
		hasChanged = true
		changedEvent.EmailAddress = email
	}
	return changedEvent, hasChanged
}
