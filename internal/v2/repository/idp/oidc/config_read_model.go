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
	Issuer                string
	Scopes                []string
	IDPDisplayNameMapping MappingField
	UserNameMapping       MappingField
}

func (rm *ConfigReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *ConfigAddedEvent:
			rm.IDPConfigID = e.IDPConfigID
			rm.ClientID = e.ClientID
			rm.ClientSecret = e.ClientSecret
			rm.Issuer = e.Issuer
			rm.Scopes = e.Scopes
			rm.IDPDisplayNameMapping = e.IDPDisplayNameMapping
			rm.UserNameMapping = e.UserNameMapping
		case *ConfigChangedEvent:
			if e.ClientID != "" {
				rm.ClientID = e.ClientID
			}
			if e.Issuer != "" {
				rm.Issuer = e.Issuer
			}
			if len(e.Scopes) > 0 {
				rm.Scopes = e.Scopes
			}
			if e.IDPDisplayNameMapping.Valid() {
				rm.IDPDisplayNameMapping = e.IDPDisplayNameMapping
			}
			if e.UserNameMapping.Valid() {
				rm.UserNameMapping = e.UserNameMapping
			}
		}
	}

	return rm.ReadModel.Reduce()
}
