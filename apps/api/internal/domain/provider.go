package domain

type IdentityProviderType int8

const (
	IdentityProviderTypeSystem IdentityProviderType = iota
	IdentityProviderTypeOrg

	identityProviderCount
)

type IdentityProviderState int32

const (
	IdentityProviderStateUnspecified IdentityProviderState = iota
	IdentityProviderStateActive
	IdentityProviderStateRemoved

	idpProviderState
)
