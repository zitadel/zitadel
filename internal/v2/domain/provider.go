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
