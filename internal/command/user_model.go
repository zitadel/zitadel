package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type UserWriteModel struct {
	eventstore.WriteModel

	UserName  string
	IDPLinks  []*domain.UserIDPLink
	UserState domain.UserState
	UserType  domain.UserType
}

func NewUserWriteModel(userID, resourceOwner string) *UserWriteModel {
	return &UserWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		IDPLinks: make([]*domain.UserIDPLink, 0),
	}
}

func (wm *UserWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.UserName = e.UserName
			wm.UserState = domain.UserStateActive
			wm.UserType = domain.UserTypeHuman
		case *user.HumanRegisteredEvent:
			wm.UserName = e.UserName
			wm.UserState = domain.UserStateActive
			wm.UserType = domain.UserTypeHuman
		case *user.HumanInitialCodeAddedEvent:
			wm.UserState = domain.UserStateInitial
		case *user.HumanInitializedCheckSucceededEvent:
			wm.UserState = domain.UserStateActive
		case *user.UserIDPLinkAddedEvent:
			wm.IDPLinks = append(wm.IDPLinks, &domain.UserIDPLink{IDPConfigID: e.IDPConfigID, ExternalUserID: e.ExternalUserID})
		case *user.UserIDPLinkRemovedEvent:
			idx, _ := wm.IDPLinkByID(e.IDPConfigID, e.ExternalUserID)
			if idx < 0 {
				continue
			}
			copy(wm.IDPLinks[idx:], wm.IDPLinks[idx+1:])
			wm.IDPLinks[len(wm.IDPLinks)-1] = nil
			wm.IDPLinks = wm.IDPLinks[:len(wm.IDPLinks)-1]
		case *user.UserIDPLinkCascadeRemovedEvent:
			idx, _ := wm.IDPLinkByID(e.IDPConfigID, e.ExternalUserID)
			if idx < 0 {
				continue
			}
			copy(wm.IDPLinks[idx:], wm.IDPLinks[idx+1:])
			wm.IDPLinks[len(wm.IDPLinks)-1] = nil
			wm.IDPLinks = wm.IDPLinks[:len(wm.IDPLinks)-1]
		case *user.MachineAddedEvent:
			wm.UserName = e.UserName
			wm.UserState = domain.UserStateActive
			wm.UserType = domain.UserTypeMachine
		case *user.UsernameChangedEvent:
			wm.UserName = e.UserName
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
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *UserWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			user.HumanAddedType,
			user.HumanRegisteredType,
			user.HumanInitializedCheckSucceededType,
			user.UserIDPLinkAddedType,
			user.UserIDPLinkRemovedType,
			user.UserIDPLinkCascadeRemovedType,
			user.MachineAddedEventType,
			user.UserUserNameChangedType,
			user.MachineChangedEventType,
			user.UserLockedType,
			user.UserUnlockedType,
			user.UserDeactivatedType,
			user.UserReactivatedType,
			user.UserRemovedType,
			user.UserV1AddedType,
			user.UserV1RegisteredType,
			user.UserV1InitializedCheckSucceededType).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func UserAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, user.AggregateType, user.AggregateVersion)
}

func isUserStateExists(state domain.UserState) bool {
	return !hasUserState(state, domain.UserStateDeleted, domain.UserStateUnspecified)
}

func isUserStateInactive(state domain.UserState) bool {
	return hasUserState(state, domain.UserStateInactive)
}

func isUserStateInitial(state domain.UserState) bool {
	return hasUserState(state, domain.UserStateInitial)
}

func hasUserState(check domain.UserState, states ...domain.UserState) bool {
	for _, state := range states {
		if check == state {
			return true
		}
	}
	return false
}

func (wm *UserWriteModel) IDPLinkByID(idpID, externalUserID string) (idx int, idp *domain.UserIDPLink) {
	for idx, idp = range wm.IDPLinks {
		if idp.IDPConfigID == idpID && idp.ExternalUserID == externalUserID {
			return idx, idp
		}
	}
	return -1, nil
}
