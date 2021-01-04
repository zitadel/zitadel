package domain

import "github.com/caos/zitadel/internal/eventstore/models"

type LoginPolicy struct {
	models.ObjectRoot

	Default               bool
	AllowUsernamePassword bool
	AllowRegister         bool
	AllowExternalIdp      bool
	IDPProviders          []*IDPProvider
	ForceMFA              bool
	SecondFactors         []SecondFactorType
	MultiFactors          []MultiFactorType
	PasswordlessType      PasswordlessType
}

type IDPProvider struct {
	models.ObjectRoot
	Type        IdentityProviderType
	IDPConfigID string
}

type PasswordlessType int32

const (
	PasswordlessTypeNotAllowed PasswordlessType = iota
	PasswordlessTypeAllowed

	passwordlessCount
)

func (f PasswordlessType) Valid() bool {
	return f >= 0 && f < passwordlessCount
}
