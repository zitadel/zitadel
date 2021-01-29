package domain

import "github.com/caos/zitadel/internal/eventstore/models"

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
	if p.SecondFactors == nil || len(p.SecondFactors) == 0 {
		return false
	}
	return true
}

func (p *LoginPolicy) HasMultiFactors() bool {
	if p.MultiFactors == nil || len(p.MultiFactors) == 0 {
		return false
	}
	return true
}
