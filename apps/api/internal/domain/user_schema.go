package domain

type UserSchemaState int32

const (
	UserSchemaStateUnspecified UserSchemaState = iota
	UserSchemaStateActive
	UserSchemaStateInactive
	UserSchemaStateDeleted
	userSchemaStateCount
)

type AuthenticatorType int32

const (
	AuthenticatorTypeUnspecified AuthenticatorType = iota
	AuthenticatorTypeUsername
	AuthenticatorTypePassword
	AuthenticatorTypeWebAuthN
	AuthenticatorTypeTOTP
	AuthenticatorTypeOTPEmail
	AuthenticatorTypeOTPSMS
	AuthenticatorTypeAuthenticationKey
	AuthenticatorTypeIdentityProvider
	authenticatorTypeCount
)
