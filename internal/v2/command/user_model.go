package command

import (
	"strings"

	caos_errors "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
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

func (wm *UserWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
}

//TODO: Compute OTPState? initial/active
func (wm *UserWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.UserName = e.UserName
			wm.UserState = domain.UserStateInitial
		case *user.HumanRegisteredEvent:
			wm.UserName = e.UserName
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
		AggregateIDs(wm.AggregateID)
	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func UserAggregateFromWriteModel(wm *eventstore.WriteModel) *user.Aggregate {
	return &user.Aggregate{
		Aggregate: *eventstore.AggregateFromWriteModel(wm, user.AggregateType, user.AggregateVersion),
	}
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
