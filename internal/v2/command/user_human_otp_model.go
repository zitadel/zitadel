package command

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

type HumanOTPWriteModel struct {
	eventstore.WriteModel

	State  domain.MFAState
	Secret *crypto.CryptoValue
}

func NewHumanOTPWriteModel(userID, resourceOwner string) *HumanOTPWriteModel {
	return &HumanOTPWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanOTPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanOTPAddedEvent:
			wm.Secret = e.Secret
			wm.State = domain.MFAStateNotReady
		case *user.HumanOTPVerifiedEvent:
			wm.State = domain.MFAStateReady
		case *user.HumanOTPRemovedEvent:
			wm.State = domain.MFAStateRemoved
		case *user.UserRemovedEvent:
			wm.State = domain.MFAStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanOTPWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(user.HumanMFAOTPAddedType,
			user.HumanMFAOTPVerifiedType,
			user.HumanMFAOTPRemovedType,
			user.UserRemovedType)
	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}
