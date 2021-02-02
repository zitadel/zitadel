package command

import (
	"context"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
	"time"
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

func (wm *HumanInitCodeWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
}

func (wm *HumanInitCodeWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.Email = e.EmailAddress
			wm.UserState = domain.UserStateInitial
		case *user.HumanRegisteredEvent:
			wm.Email = e.EmailAddress
			wm.UserState = domain.UserStateInitial
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
		case *user.HumanInitializedCheckSucceededEvent:
			wm.Code = nil
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanInitCodeWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID)
	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func (wm *HumanInitCodeWriteModel) NewChangedEvent(
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
