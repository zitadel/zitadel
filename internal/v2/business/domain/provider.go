package domain

type IdentityProviderType int8

const (
	IdentityProviderTypeSystem IdentityProviderType = iota
	IdentityProviderTypeOrg

	typeCount
)

func (f IdentityProviderType) Valid() bool {
	return f >= 0 && f < typeCount
}
