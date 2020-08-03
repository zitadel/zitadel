package model

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type IdpConfig struct {
	es_models.ObjectRoot
	IDPConfigID string
	Type        IdpConfigType
	Name        string
	LogoSrc     string
	State       IdpConfigState
	OIDCConfig  *OidcIdpConfig
}

type OidcIdpConfig struct {
	es_models.ObjectRoot
	IDPConfigID        string
	ClientID           string
	ClientSecret       *crypto.CryptoValue
	ClientSecretString string
	Issuer             string
	Scopes             []string
}

type IdpConfigType int32

const (
	IDPConfigTypeOIDC IdpConfigType = iota
	IDPConfigTypeSAML
)

type IdpConfigState int32

const (
	IdpConfigStateActive IdpConfigState = iota
	IdpConfigStateInactive
	IdpConfigStateRemoved
)

func NewIdpConfig(iamID, idpID string) *IdpConfig {
	return &IdpConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: iamID}, IDPConfigID: idpID}
}

func (idp *IdpConfig) IsValid(includeConfig bool) bool {
	if idp.Name == "" || idp.AggregateID == "" {
		return false
	}
	if !includeConfig {
		return true
	}
	if idp.Type == IDPConfigTypeOIDC && !idp.OIDCConfig.IsValid(true) {
		return false
	}
	return true
}

func (oi *OidcIdpConfig) IsValid(withSecret bool) bool {
	if withSecret {
		return oi.ClientID != "" && oi.Issuer != "" && oi.ClientSecretString != ""
	}
	return oi.ClientID != "" && oi.Issuer != ""
}

func (oi *OidcIdpConfig) CryptSecret(crypt crypto.Crypto) error {
	cryptedSecret, err := crypto.Crypt([]byte(oi.ClientSecretString), crypt)
	if err != nil {
		return err
	}
	oi.ClientSecret = cryptedSecret
	return nil
}
