package command

import (
	"github.com/caos/zitadel/internal/eventstore"
	"strings"

	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/user"
)

type UserWriteModel struct {
	eventstore.WriteModel

	UserName  string
	UserState domain.UserState
}

func NewUserWriteModel(userID, resourceOwner string) *UserWriteModel {
	return &UserWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *UserWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.UserName = e.UserName
			wm.UserState = domain.UserStateActive
		case *user.HumanRegisteredEvent:
			wm.UserName = e.UserName
			wm.UserState = domain.UserStateActive
		case *user.HumanInitialCodeAddedEvent:
			wm.UserState = domain.UserStateInitial
		case *user.HumanInitializedCheckSucceededEvent:
			wm.UserState = domain.UserStateActive
		case *user.MachineAddedEvent:
			wm.UserName = e.UserName
			wm.UserState = domain.UserStateActive
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
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			user.HumanAddedType,
			user.HumanRegisteredType,
			user.HumanInitializedCheckSucceededType,
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
			user.UserV1InitializedCheckSucceededType)
	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func UserAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, user.AggregateType, user.AggregateVersion)
}

func CheckOrgIAMPolicyForUserName(userName string, policy *domain.OrgIAMPolicy) error {
	if policy == nil {
		return caos_errors.ThrowPreconditionFailed(nil, "COMMAND-3Mb9s", "Errors.Users.OrgIamPolicyNil")
	}
	if policy.UserLoginMustBeDomain && strings.Contains(userName, "@") {
		return caos_errors.ThrowPreconditionFailed(nil, "COMMAND-4M9vs", "Errors.User.EmailAsUsernameNotAllowed")
	}
	return nil
}

func isUserStateExists(state domain.UserState) bool {
	return !hasUserState(state, domain.UserStateDeleted, domain.UserStateUnspecified)
}

func isUserStateInactive(state domain.UserState) bool {
	return hasUserState(state, domain.UserStateInactive)
}

func hasUserState(check domain.UserState, states ...domain.UserState) bool {
	for _, state := range states {
		if check == state {
			return true
		}
	}
	return false
}
