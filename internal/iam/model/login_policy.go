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
	IdpProviders          []*IdpProvider
}

type IdpProvider struct {
	models.ObjectRoot
	Type        IdpProviderType
	IdpConfigID string
}

type PolicyState int32

const (
	PolicyStateActive PolicyState = iota
	PolicyStateRemoved
)

type IdpProviderType int32

const (
	IdpProviderTypeSystem IdpProviderType = iota
	IdpProviderTypeOrg
)

func (p *LoginPolicy) IsValid() bool {
	return p.ObjectRoot.AggregateID != ""
}

func (p *IdpProvider) IsValid() bool {
	return p.ObjectRoot.AggregateID != "" && p.IdpConfigID != ""
}

func GetIdpProvider(providers []*IdpProvider, id string) (int, *IdpProvider) {
	for i, m := range providers {
		if m.IdpConfigID == id {
			return i, m
		}
	}
	return -1, nil
}
