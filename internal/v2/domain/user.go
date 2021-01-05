package domain

import es_models "github.com/caos/zitadel/internal/eventstore/models"

type User struct {
	es_models.ObjectRoot
	State    UserState
	UserName string

	*Human
	*Machine
}

type UserState int32

const (
	UserStateUnspecified UserState = iota
	UserStateActive
	UserStateInactive
	UserStateDeleted
	UserStateLocked
	UserStateSuspend
	UserStateInitial

	userStateCount
)

func (f UserState) Valid() bool {
	return f >= 0 && f < userStateCount
}

func (u *User) IsValid() bool {
	if u.Human == nil && u.Machine == nil || u.UserName == "" {
		return false
	}
	if u.Human != nil {
		return u.Human.IsValid()
	}
	return u.Machine.IsValid()
}
