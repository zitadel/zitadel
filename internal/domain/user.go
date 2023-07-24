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

func (s UserState) NotDisabled() bool {
	return s == UserStateActive || s == UserStateInitial
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

type UserAuthMethodType int32

const (
	UserAuthMethodTypeUnspecified UserAuthMethodType = iota
	UserAuthMethodTypeOTP
	UserAuthMethodTypeU2F
	UserAuthMethodTypePasswordless
	UserAuthMethodTypePassword
	UserAuthMethodTypeIDP
	userAuthMethodTypeCount
)

func (f UserAuthMethodType) Valid() bool {
	return f >= 0 && f < userAuthMethodTypeCount
}

// HasMFA checks whether the user authenticated with multiple auth factors.
// This can either be true if the list contains a [UserAuthMethodType] which by itself is MFA (e.g. [UserAuthMethodTypePasswordless])
// or if multiple factors were used (e.g. [UserAuthMethodTypePassword] and [UserAuthMethodTypeU2F])
func HasMFA(methods []UserAuthMethodType) bool {
	var factors int
	for _, method := range methods {
		switch method {
		case UserAuthMethodTypePassword:
			factors++
		case UserAuthMethodTypePasswordless:
			return true
		case UserAuthMethodTypeU2F:
			factors++
		case UserAuthMethodTypeOTP:
			factors++
		case UserAuthMethodTypeIDP:
			factors++
		case UserAuthMethodTypeUnspecified,
			userAuthMethodTypeCount:
			// ignore
		}
	}
	return factors > 1
}

// RequiresMFA checks whether the user requires to authenticate with multiple auth factors based on the LoginPolicy and the authentication type.
// Internal authentication will require MFA if either option is activated.
// External authentication will only require MFA if it's forced generally and not local only.
func RequiresMFA(forceMFA, forceMFALocalOnly, isInternalLogin bool) bool {
	if isInternalLogin {
		return forceMFA || forceMFALocalOnly
	}
	return forceMFA && !forceMFALocalOnly
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
