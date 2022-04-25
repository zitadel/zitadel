package models

type AuthRequestInt interface {
	GetID() string
	GetApplicationID() string
	GetRelayState() string
	GetNameID() string
	GetAccessConsumerServiceURL() string
	GetBindingType() string
	GetAuthRequestID() string
	GetCode() string
	GetIssuer() string
	GetIssuerName() string
	GetDestination() string
	GetUserID() string
	GetUserName() string
	Done() bool
}

type AttributeSetter interface {
	SetEmail(string)
	SetFullName(string)
	SetGivenName(string)
	SetSurname(string)
	SetUserID(string)
	SetUsername(string)
}
