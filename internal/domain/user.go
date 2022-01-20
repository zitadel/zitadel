package domain

type User interface {
	GetUsername() string
	GetState() UserState
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

func (s UserState) Exists() bool {
	return s != UserStateUnspecified && s != UserStateDeleted
}

type UserType int32

const (
	UserTypeUnspecified UserType = iota
	UserTypeHuman
	UserTypeMachine
	userTypeCount
)

func (f UserType) Valid() bool {
	return f >= 0 && f < userTypeCount
}

type PersonalAccessTokenState int32

const (
	PersonalAccessTokenStateUnspecified PersonalAccessTokenState = iota
	PersonalAccessTokenStateActive
	PersonalAccessTokenStateRemoved

	personalAccessTokenStateCount
)

func (f PersonalAccessTokenState) Valid() bool {
	return f >= 0 && f < personalAccessTokenStateCount
}
