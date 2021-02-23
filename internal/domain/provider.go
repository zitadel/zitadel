package domain

type IdentityProviderType int8

const (
	IdentityProviderTypeSystem IdentityProviderType = iota
	IdentityProviderTypeOrg

	identityProviderCount
)

func (f IdentityProviderType) Valid() bool {
	return f >= 0 && f < identityProviderCount
}

type IdentityProviderState int32

const (
	IdentityProviderStateUnspecified IdentityProviderState = iota
	IdentityProviderStateActive
	IdentityProviderStateRemoved

	idpProviderState
)

func (s IdentityProviderState) Valid() bool {
	return s >= 0 && s < idpProviderState
}
