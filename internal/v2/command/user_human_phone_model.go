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

	State domain.PhoneState
}

func NewHumanPhoneWriteModel(userID, resourceOwner string) *HumanPhoneWriteModel {
	return &HumanPhoneWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanPhoneWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanAddedEvent, *user.HumanRegisteredEvent:
			wm.AppendEvents(e)
		case *user.HumanPhoneChangedEvent:
			wm.AppendEvents(e)
		case *user.HumanPhoneVerifiedEvent:
			wm.AppendEvents(e)
		case *user.HumanPhoneRemovedEvent:
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
			if e.PhoneNumber != "" {
				wm.Phone = e.PhoneNumber
				wm.State = domain.PhoneStateActive
			}
		case *user.HumanRegisteredEvent:
			if e.PhoneNumber != "" {
				wm.Phone = e.PhoneNumber
				wm.State = domain.PhoneStateActive
			}
		case *user.HumanPhoneChangedEvent:
			wm.Phone = e.PhoneNumber
			wm.IsPhoneVerified = false
			wm.State = domain.PhoneStateActive
		case *user.HumanPhoneVerifiedEvent:
			wm.IsPhoneVerified = true
		case *user.HumanPhoneRemovedEvent:
			wm.State = domain.PhoneStateRemoved
		case *user.UserRemovedEvent:
			wm.State = domain.PhoneStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanPhoneWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner)
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
