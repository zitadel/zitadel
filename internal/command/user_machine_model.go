package command

import (
	"context"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type MachineWriteModel struct {
	eventstore.WriteModel

	UserName string

	Name            string
	Description     string
	UserState       domain.UserState
	AccessTokenType domain.OIDCTokenType
	HashedSecret    string
}

func NewMachineWriteModel(userID, resourceOwner string) *MachineWriteModel {
	return &MachineWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *MachineWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.MachineAddedEvent:
			wm.UserName = e.UserName
			wm.Name = e.Name
			wm.Description = e.Description
			wm.AccessTokenType = e.AccessTokenType
			wm.UserState = domain.UserStateActive
		case *user.UsernameChangedEvent:
			wm.UserName = e.UserName
		case *user.MachineChangedEvent:
			if e.Name != nil {
				wm.Name = *e.Name
			}
			if e.Description != nil {
				wm.Description = *e.Description
			}
			if e.AccessTokenType != nil {
				wm.AccessTokenType = *e.AccessTokenType
			}
		case *user.UserLockedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateLocked
			}
		case *user.UserUnlockedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateActive
			}
		case *user.UserDeactivatedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateInactive
			}
		case *user.UserReactivatedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateActive
			}
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		case *user.MachineSecretSetEvent:
			wm.HashedSecret = crypto.SecretOrEncodedHash(e.ClientSecret, e.HashedSecret)
		case *user.MachineSecretRemovedEvent:
			wm.HashedSecret = ""
		case *user.MachineSecretHashUpdatedEvent:
			wm.HashedSecret = e.HashedSecret
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *MachineWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(user.MachineAddedEventType,
			user.UserUserNameChangedType,
			user.MachineChangedEventType,
			user.UserLockedType,
			user.UserUnlockedType,
			user.UserDeactivatedType,
			user.UserReactivatedType,
			user.UserRemovedType,
			user.MachineSecretSetType,
			user.MachineSecretRemovedType,
			user.MachineSecretHashUpdatedType,
		).Builder()
}

func (wm *MachineWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name,
	description string,
	accessTokenType domain.OIDCTokenType,
) (*user.MachineChangedEvent, bool) {
	changes := make([]user.MachineChanges, 0)

	if wm.Name != name {
		changes = append(changes, user.ChangeName(name))
	}
	if wm.Description != description {
		changes = append(changes, user.ChangeDescription(description))
	}
	if wm.AccessTokenType != accessTokenType {
		changes = append(changes, user.ChangeAccessTokenType(accessTokenType))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changeEvent := user.NewMachineChangedEvent(ctx, aggregate, changes)
	return changeEvent, true
}
