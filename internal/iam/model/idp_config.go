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
	State       IDPConfigState
	OIDCConfig  *OIDCIDPConfig
}

type OIDCIDPConfig struct {
	es_models.ObjectRoot
	IDPConfigID        string
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

func (idp *IDPConfig) IsValid(includeConfig bool) bool {
	if idp.Name == "" || idp.AggregateID == "" {
		return false
	}
	if !includeConfig {
		return true
	}
	if idp.Type == IDPConfigTypeOIDC && !idp.OIDCConfig.IsValid() {
		return false
	}
	return true
}

func (oi *OIDCIDPConfig) IsValid() bool {
	return oi.ClientID != "" && oi.ClientSecretString != "" && oi.Issuer != ""
}

func (oi *OIDCIDPConfig) CryptSecret(crypt crypto.Crypto) error {
	cryptedSecret, err := crypto.Crypt([]byte(oi.ClientSecretString), crypt)
	if err != nil {
		return err
	}
	oi.ClientSecret = cryptedSecret
	return nil
}
