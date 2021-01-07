package command

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

type HumanOTPWriteModel struct {
	eventstore.WriteModel

	Secret *crypto.CryptoValue

	UserState domain.UserState
}

func NewHumanOTPWriteModel(userID string) *HumanOTPWriteModel {
	return &HumanOTPWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: userID,
		},
	}
}

func (wm *HumanOTPWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanOTPAddedEvent:
			wm.AppendEvents(e)
		case *user.HumanOTPRemovedEvent:
			wm.AppendEvents(e)
		case *user.HumanAddedEvent, *user.HumanRegisteredEvent:
			wm.AppendEvents(e)
		case *user.UserRemovedEvent:
			wm.AppendEvents(e)
		}
	}
}

func (wm *HumanOTPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.UserState = domain.UserStateActive
		case *user.HumanRegisteredEvent:
			wm.UserState = domain.UserStateActive
		case *user.HumanOTPAddedEvent:
			wm.Secret = e.Secret
		case *user.HumanOTPRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanOTPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID)
}
