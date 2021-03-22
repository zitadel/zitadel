package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"
	"time"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/user"
)

type HumanInitCodeWriteModel struct {
	eventstore.WriteModel

	Email           string
	IsEmailVerified bool

	Code             *crypto.CryptoValue
	CodeCreationDate time.Time
	CodeExpiry       time.Duration

	UserState domain.UserState
}

func NewHumanInitCodeWriteModel(userID, resourceOwner string) *HumanInitCodeWriteModel {
	return &HumanInitCodeWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanInitCodeWriteModel) Reduce() error {
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
			wm.Code = nil
			if wm.UserState == domain.UserStateInitial {
				wm.UserState = domain.UserStateActive
			}
		case *user.HumanInitialCodeAddedEvent:
			wm.Code = e.Code
			wm.CodeCreationDate = e.CreationDate()
			wm.CodeExpiry = e.Expiry
			wm.UserState = domain.UserStateInitial
		case *user.HumanInitializedCheckSucceededEvent:
			wm.Code = nil
			wm.UserState = domain.UserStateActive
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanInitCodeWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(user.HumanAddedType,
			user.HumanRegisteredType,
			user.HumanEmailChangedType,
			user.HumanEmailVerifiedType,
			user.HumanInitialCodeAddedType,
			user.HumanInitializedCheckSucceededType,
			user.UserRemovedType)
	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func (wm *HumanInitCodeWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	email string,
) (*user.HumanEmailChangedEvent, bool) {
	changedEvent := user.NewHumanEmailChangedEvent(ctx, aggregate, email)
	return changedEvent, wm.Email != email
}
