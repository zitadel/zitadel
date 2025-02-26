package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type HumanPhoneWriteModel struct {
	eventstore.WriteModel

	Phone           domain.PhoneNumber
	IsPhoneVerified bool

	Code             *crypto.CryptoValue
	CodeCreationDate time.Time
	CodeExpiry       time.Duration
	GeneratorID      string
	VerificationID   string

	State     domain.PhoneState
	UserState domain.UserState
}

func NewHumanPhoneWriteModel(userID, resourceOwner string) *HumanPhoneWriteModel {
	return &HumanPhoneWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
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
			wm.UserState = domain.UserStateActive
		case *user.HumanRegisteredEvent:
			if e.PhoneNumber != "" {
				wm.Phone = e.PhoneNumber
				wm.State = domain.PhoneStateActive
			}
			wm.UserState = domain.UserStateActive
		case *user.HumanInitialCodeAddedEvent:
			wm.UserState = domain.UserStateInitial
		case *user.HumanInitializedCheckSucceededEvent:
			wm.UserState = domain.UserStateActive
		case *user.HumanPhoneChangedEvent:
			wm.Phone = e.PhoneNumber
			wm.IsPhoneVerified = false
			wm.State = domain.PhoneStateActive
			wm.Code = nil
		case *user.HumanPhoneVerifiedEvent:
			wm.IsPhoneVerified = true
			wm.Code = nil
		case *user.HumanPhoneCodeAddedEvent:
			wm.Code = e.Code
			wm.CodeCreationDate = e.CreationDate()
			wm.CodeExpiry = e.Expiry
			wm.GeneratorID = e.GeneratorID
		case *user.HumanPhoneCodeSentEvent:
			wm.GeneratorID = e.GeneratorInfo.GetID()
			wm.VerificationID = e.GeneratorInfo.GetVerificationID()
		case *user.HumanPhoneRemovedEvent:
			wm.State = domain.PhoneStateRemoved
			wm.IsPhoneVerified = false
			wm.Phone = ""
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
			wm.IsPhoneVerified = false
			wm.Phone = ""
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanPhoneWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(user.HumanAddedType,
			user.HumanRegisteredType,
			user.HumanInitialCodeAddedType,
			user.HumanInitializedCheckSucceededType,
			user.HumanPhoneChangedType,
			user.HumanPhoneVerifiedType,
			user.HumanPhoneCodeAddedType,
			user.HumanPhoneCodeSentType,
			user.HumanPhoneRemovedType,
			user.UserRemovedType,
			user.UserV1AddedType,
			user.UserV1RegisteredType,
			user.UserV1InitialCodeAddedType,
			user.UserV1InitializedCheckSucceededType,
			user.UserV1PhoneCodeAddedType,
			user.UserV1PhoneChangedType,
			user.UserV1PhoneVerifiedType,
			user.UserV1PhoneRemovedType).
		Builder()
}

func (wm *HumanPhoneWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	phone domain.PhoneNumber,
) (*user.HumanPhoneChangedEvent, bool) {
	changedEvent := user.NewHumanPhoneChangedEvent(ctx, aggregate, phone)
	return changedEvent, phone != wm.Phone
}
