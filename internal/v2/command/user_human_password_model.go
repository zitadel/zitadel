package command

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
	"time"
)

type HumanPasswordWriteModel struct {
	eventstore.WriteModel

	Secret               *crypto.CryptoValue
	SecretChangeRequired bool

	Code             *crypto.CryptoValue
	CodeCreationDate time.Time
	CodeExpiry       time.Duration

	UserState domain.UserState
}

func NewHumanPasswordWriteModel(userID, resourceOwner string) *HumanPasswordWriteModel {
	return &HumanPasswordWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanPasswordWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
}

func (wm *HumanPasswordWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.Secret = e.Secret
			wm.SecretChangeRequired = e.ChangeRequired
			wm.UserState = domain.UserStateInitial
		case *user.HumanRegisteredEvent:
			wm.Secret = e.Secret
			wm.SecretChangeRequired = e.ChangeRequired
			wm.UserState = domain.UserStateActive
		case *user.HumanPasswordChangedEvent:
			wm.Secret = e.Secret
			wm.SecretChangeRequired = e.ChangeRequired
			wm.Code = nil
		case *user.HumanPasswordCodeAddedEvent:
			wm.Code = e.Code
			wm.CodeCreationDate = e.CreationDate()
			wm.CodeExpiry = e.Expiry
		case *user.HumanEmailVerifiedEvent:
			if wm.UserState == domain.UserStateInitial {
				wm.UserState = domain.UserStateActive
			}
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanPasswordWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}
