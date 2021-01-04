package model

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type LoginPolicy struct {
	models.ObjectRoot

	State                 PolicyState
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
	Type        IDPProviderType
	IDPConfigID string
}

type PolicyState int32

const (
	PolicyStateActive PolicyState = iota
	PolicyStateRemoved
)

type IDPProviderType int32

const (
	IDPProviderTypeSystem IDPProviderType = iota
	IDPProviderTypeOrg
)

type SecondFactorType int32

const (
	SecondFactorTypeUnspecified SecondFactorType = iota
	SecondFactorTypeOTP
	SecondFactorTypeU2F
)

type MultiFactorType int32

const (
	MultiFactorTypeUnspecified MultiFactorType = iota
	MultiFactorTypeU2FWithPIN
)

type PasswordlessType int32

const (
	PasswordlessTypeNotAllowed PasswordlessType = iota
	PasswordlessTypeAllowed
)

func (p *LoginPolicy) IsValid() bool {
	return p.ObjectRoot.AggregateID != ""
}

func (p *IDPProvider) IsValid() bool {
	return p.ObjectRoot.AggregateID != "" && p.IDPConfigID != ""
}

func (p *LoginPolicy) GetIdpProvider(id string) (int, *IDPProvider) {
	for i, m := range p.IDPProviders {
		if m.IDPConfigID == id {
			return i, m
		}
	}
	return -1, nil
}

func (p *LoginPolicy) GetSecondFactor(mfaType SecondFactorType) (int, SecondFactorType) {
	for i, m := range p.SecondFactors {
		if m == mfaType {
			return i, m
		}
	}
	return -1, 0
}

func (p *LoginPolicy) GetMultiFactor(mfaType MultiFactorType) (int, MultiFactorType) {
	for i, m := range p.MultiFactors {
		if m == mfaType {
			return i, m
		}
	}
	return -1, 0
}
