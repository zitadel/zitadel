package model

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type IDPConfig struct {
	es_models.ObjectRoot
	IDPConfigID string
	Type        IDPConfigType
	Name        string
	LogoSrc     string
	OIDCConfig  OIDCIDPConfig
}

type OIDCIDPConfig struct {
	ClientID           string
	ClientSecret       *crypto.CryptoValue
	ClientSecretString string
	Issuer             string
	Scopes             []string
}

type IDPConfigType int32

const (
	IDPConfigTypeOIDC IDPConfigType = iota
	IDPConfigTypeSAML
)

type IDPConfigState int32

const (
	IDPConfigStateActive IDPConfigState = iota
	IDPConfigStateInactive
	IDPConfigStateRemoved
)
