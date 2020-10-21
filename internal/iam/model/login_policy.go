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
	SoftwareMFAs          []SoftwareMFAType
	HardwareMFAs          []HardwareMFAType
}

type IDPProvider struct {
	models.ObjectRoot
	Type        IDPProviderType
	IdpConfigID string
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

type SoftwareMFAType int32

const (
	SoftwareMFATypeUnspecified SoftwareMFAType = iota
	SoftwareMFATypeOTP
)

type HardwareMFAType int32

const (
	HardwareMFATypeUnspecified HardwareMFAType = iota
	HardwareMFATypeU2F
)

func (p *LoginPolicy) IsValid() bool {
	return p.ObjectRoot.AggregateID != ""
}

func (p *IDPProvider) IsValid() bool {
	return p.ObjectRoot.AggregateID != "" && p.IdpConfigID != ""
}

func (p *LoginPolicy) GetIdpProvider(id string) (int, *IDPProvider) {
	for i, m := range p.IDPProviders {
		if m.IdpConfigID == id {
			return i, m
		}
	}
	return -1, nil
}

func (p *LoginPolicy) GetSoftwareMFA(mfaType SoftwareMFAType) (int, SoftwareMFAType) {
	for i, m := range p.SoftwareMFAs {
		if m == mfaType {
			return i, m
		}
	}
	return -1, 0
}

func (p *LoginPolicy) GetHardwareMFA(mfaType HardwareMFAType) (int, HardwareMFAType) {
	for i, m := range p.HardwareMFAs {
		if m == mfaType {
			return i, m
		}
	}
	return -1, 0
}
