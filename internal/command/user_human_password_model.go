package command

import (
	"github.com/caos/zitadel/internal/eventstore"
	"time"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/user"
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

func (wm *HumanPasswordWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.Secret = e.Secret
			wm.SecretChangeRequired = e.ChangeRequired
			wm.UserState = domain.UserStateActive
		case *user.HumanRegisteredEvent:
			wm.Secret = e.Secret
			wm.SecretChangeRequired = e.ChangeRequired
			wm.UserState = domain.UserStateActive
		case *user.HumanInitialCodeAddedEvent:
			wm.UserState = domain.UserStateInitial
		case *user.HumanInitializedCheckSucceededEvent:
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
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(user.HumanAddedType,
			user.HumanRegisteredType,
			user.HumanInitialCodeAddedType,
			user.HumanInitializedCheckSucceededType,
			user.HumanPasswordChangedType,
			user.HumanPasswordCodeAddedType,
			user.HumanEmailVerifiedType,
			user.UserRemovedType,
			user.UserV1AddedType,
			user.UserV1RegisteredType,
			user.UserV1InitialCodeAddedType,
			user.UserV1InitializedCheckSucceededType,
			user.UserV1PasswordChangedType,
			user.UserV1PasswordCodeAddedType,
			user.UserV1EmailVerifiedType)
	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}
