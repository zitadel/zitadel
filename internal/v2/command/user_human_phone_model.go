package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

type HumanPhoneWriteModel struct {
	eventstore.WriteModel

	Phone           string
	IsPhoneVerified bool

	UserState domain.UserState
}

func NewHumanPhoneWriteModel(userID string) *HumanPhoneWriteModel {
	return &HumanPhoneWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: userID,
		},
	}
}

func (wm *HumanPhoneWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanPhoneChangedEvent:
			wm.AppendEvents(e)
		case *user.HumanPhoneVerifiedEvent:
			wm.AppendEvents(e)
		case *user.HumanAddedEvent, *user.HumanRegisteredEvent:
			wm.AppendEvents(e)
		case *user.UserRemovedEvent:
			wm.AppendEvents(e)
		}
	}
}

func (wm *HumanPhoneWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.Phone = e.PhoneNumber
			wm.UserState = domain.UserStateActive
		case *user.HumanRegisteredEvent:
			wm.Phone = e.PhoneNumber
			wm.UserState = domain.UserStateActive
		case *user.HumanPhoneChangedEvent:
			wm.Phone = e.PhoneNumber
			wm.IsPhoneVerified = false
		case *user.HumanPhoneVerifiedEvent:
			wm.IsPhoneVerified = true
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanPhoneWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID)
}

func (wm *HumanPhoneWriteModel) NewChangedEvent(
	ctx context.Context,
	phone string,
) (*user.HumanPhoneChangedEvent, bool) {
	hasChanged := false
	changedEvent := user.NewHumanPhoneChangedEvent(ctx)
	if wm.Phone != phone {
		hasChanged = true
		changedEvent.PhoneNumber = phone
	}
	return changedEvent, hasChanged
}
