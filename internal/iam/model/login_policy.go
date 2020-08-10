package model

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type LoginPolicy struct {
	models.ObjectRoot

	State                 PolicyState
	UserLoginMustBeDomain bool
	Default               bool
	AllowUsernamePassword bool
	AllowRegister         bool
	AllowExternalIdp      bool
	IdpProviders          []*IdpProvider
}

type IdpProvider struct {
	models.ObjectRoot
	IdpProviderID   string
	Type            IdpProviderType
	IdpConfigID     string
	IdpCustomConfig *IdpConfig
}

type PolicyState int32

const (
	PolicyStateActive PolicyState = iota
	PolicyStateRemoved
)

type IdpProviderType int32

const (
	IdpProviderTypePreConfigured IdpProviderType = iota
	IdpProviderTypeCustom
)
