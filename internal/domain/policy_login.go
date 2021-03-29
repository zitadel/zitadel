package domain

import "github.com/caos/zitadel/internal/eventstore/v1/models"

type LoginPolicy struct {
	models.ObjectRoot

	Default               bool
	AllowUsernamePassword bool
	AllowRegister         bool
	AllowExternalIDP      bool
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

	Name          string
	StylingType   IDPConfigStylingType
	IDPConfigType IDPConfigType
	IDPState      IDPConfigState
}

func (p IDPProvider) IsValid() bool {
	return p.IDPConfigID != ""
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

func (p *LoginPolicy) HasSecondFactors() bool {
	return len(p.SecondFactors) > 0
}

func (p *LoginPolicy) HasMultiFactors() bool {
	return len(p.MultiFactors) > 0
}
