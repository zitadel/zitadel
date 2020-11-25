package oidc

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
)

type ConfigReadModel struct {
	eventstore.ReadModel

	IDPConfigID           string
	ClientID              string
	ClientSecret          *crypto.CryptoValue
	ClientSecretString    string
	Issuer                string
	Scopes                []string
	IDPDisplayNameMapping MappingField
	UsernameMapping       MappingField
}

func (rm *ConfigReadModel) AppendEvents(events ...eventstore.EventReader) {
	rm.ReadModel.AppendEvents(events...)
}

func (rm *ConfigReadModel) Reduce() error {
	return nil
}

type MappingField int32

const (
	OIDCMappingFieldUnspecified MappingField = iota
	OIDCMappingFieldPreferredLoginName
	OIDCMappingFieldEmail
)
