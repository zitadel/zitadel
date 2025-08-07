package user

import "github.com/zitadel/zitadel/internal/domain"

type SearchQuery_ResourceOwner struct {
	ResourceOwner *ResourceOwnerQuery
}

func (SearchQuery_ResourceOwner) isSearchQuery_Query() {}

type ResourceOwnerQuery struct {
	OrgID string
}

type UserType = isUser_Type

type MembershipType = isMembership_Type

func (s UserState) ToDomain() domain.UserState {
	switch s {
	case UserState_USER_STATE_UNSPECIFIED:
		return domain.UserStateUnspecified
	case UserState_USER_STATE_ACTIVE:
		return domain.UserStateActive
	case UserState_USER_STATE_INACTIVE:
		return domain.UserStateInactive
	case UserState_USER_STATE_DELETED:
		return domain.UserStateDeleted
	case UserState_USER_STATE_LOCKED:
		return domain.UserStateLocked
	case UserState_USER_STATE_SUSPEND:
		return domain.UserStateSuspend
	case UserState_USER_STATE_INITIAL:
		return domain.UserStateInitial
	default:
		return domain.UserStateUnspecified
	}
}

func (t Type) ToDomain() domain.UserType {
	switch t {
	case Type_TYPE_UNSPECIFIED:
		return domain.UserTypeUnspecified
	case Type_TYPE_HUMAN:
		return domain.UserTypeHuman
	case Type_TYPE_MACHINE:
		return domain.UserTypeMachine
	default:
		return domain.UserTypeUnspecified
	}
}
