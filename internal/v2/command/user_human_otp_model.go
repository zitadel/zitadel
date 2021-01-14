package command

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

type HumanOTPWriteModel struct {
	eventstore.WriteModel

	MFAState domain.MFAState
	Secret   *crypto.CryptoValue
	OTPState domain.OTPState
}

func NewHumanOTPWriteModel(userID, resourceOwner string) *HumanOTPWriteModel {
	return &HumanOTPWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanOTPWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
}

func (wm *HumanOTPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanOTPAddedEvent:
			wm.Secret = e.Secret
			wm.MFAState = domain.MFAStateNotReady
			wm.OTPState = domain.OTPStateActive
		case *user.HumanOTPVerifiedEvent:
			wm.MFAState = domain.MFAStateReady
		case *user.HumanOTPRemovedEvent:
			wm.OTPState = domain.OTPStateRemoved
		case *user.UserRemovedEvent:
			wm.OTPState = domain.OTPStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanOTPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}
