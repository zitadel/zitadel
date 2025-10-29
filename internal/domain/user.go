package domain

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

func (s UserState) Exists() bool {
	return s != UserStateUnspecified && s != UserStateDeleted
}

func (s UserState) IsEnabled() bool {
	return s == UserStateActive || s == UserStateInitial
}

type UserType int32

const (
	UserTypeUnspecified UserType = iota
	UserTypeHuman
	UserTypeMachine
	userTypeCount
)

type UserAuthMethodType int32

const (
	UserAuthMethodTypeUnspecified UserAuthMethodType = iota
	UserAuthMethodTypeTOTP
	UserAuthMethodTypeU2F
	UserAuthMethodTypePasswordless
	UserAuthMethodTypePassword
	UserAuthMethodTypeIDP
	UserAuthMethodTypeOTPSMS
	UserAuthMethodTypeOTPEmail
	UserAuthMethodTypeOTP // generic OTP when parsing AMR from OIDC
	UserAuthMethodTypePrivateKey
	userAuthMethodTypeCount
)

// HasMFA checks whether the user authenticated with multiple auth factors.
// This can either be true if the list contains a [UserAuthMethodType] which by itself is MFA (e.g. [UserAuthMethodTypePasswordless])
// or if multiple factors were used (e.g. [UserAuthMethodTypePassword] and [UserAuthMethodTypeU2F])
func HasMFA(methods []UserAuthMethodType) bool {
	var factors int
	for _, method := range methods {
		switch method {
		case UserAuthMethodTypePasswordless:
			return true
		case UserAuthMethodTypePassword,
			UserAuthMethodTypeU2F,
			UserAuthMethodTypeTOTP,
			UserAuthMethodTypeOTPSMS,
			UserAuthMethodTypeOTPEmail,
			UserAuthMethodTypeIDP,
			UserAuthMethodTypeOTP,
			UserAuthMethodTypePrivateKey:
			factors++
		case UserAuthMethodTypeUnspecified,
			userAuthMethodTypeCount:
			// ignore
		}
	}
	return factors > 1
}

// Has2FA checks whether the auth factors provided are a second factor and will return true if at least one is.
func Has2FA(methods []UserAuthMethodType) bool {
	var factors int
	for _, method := range methods {
		switch method {
		case
			UserAuthMethodTypeU2F,
			UserAuthMethodTypeTOTP,
			UserAuthMethodTypeOTPSMS,
			UserAuthMethodTypeOTPEmail,
			UserAuthMethodTypeOTP:
			factors++
		case UserAuthMethodTypeUnspecified,
			UserAuthMethodTypePassword,
			UserAuthMethodTypePasswordless,
			UserAuthMethodTypeIDP,
			UserAuthMethodTypePrivateKey,
			userAuthMethodTypeCount:
			// ignore
		}
	}
	return factors > 0
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

// AuthMethodToSecondFactor maps user auth methods to their corresponding second factor types
func AuthMethodToSecondFactor(method UserAuthMethodType) SecondFactorType {
	switch method {
	case UserAuthMethodTypeTOTP:
		return SecondFactorTypeTOTP
	case UserAuthMethodTypeU2F:
		return SecondFactorTypeU2F
	case UserAuthMethodTypeOTPSMS:
		return SecondFactorTypeOTPSMS
	case UserAuthMethodTypeOTPEmail:
		return SecondFactorTypeOTPEmail
	case UserAuthMethodTypeOTP:
		return SecondFactorTypeOTPSMS
	default:
		// First-factor methods: password, IDP, passwordless, private key
		return 0
	}
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
