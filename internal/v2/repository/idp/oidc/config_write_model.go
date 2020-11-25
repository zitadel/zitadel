package oidc

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
)

type ConfigWriteModel struct {
	eventstore.WriteModel

	IDPConfigID  string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Issuer       string
	Scopes       []string

	IDPDisplayNameMapping MappingField
	UserNameMapping       MappingField
}

func (wm *ConfigWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *ConfigAddedEvent:
			wm.IDPConfigID = e.IDPConfigID
			wm.ClientID = e.ClientID
			wm.ClientSecret = e.ClientSecret
			wm.Issuer = e.Issuer
			wm.Scopes = e.Scopes
			wm.IDPDisplayNameMapping = e.IDPDisplayNameMapping
			wm.UserNameMapping = e.UserNameMapping
		case *ConfigChangedEvent:
			if e.ClientID != "" {
				wm.ClientID = e.ClientID
			}
			if e.Issuer != "" {
				wm.Issuer = e.Issuer
			}
			if len(e.Scopes) > 0 {
				wm.Scopes = e.Scopes
			}
			if e.IDPDisplayNameMapping.Valid() {
				wm.IDPDisplayNameMapping = e.IDPDisplayNameMapping
			}
			if e.UserNameMapping.Valid() {
				wm.UserNameMapping = e.UserNameMapping
			}
		}
	}

	return wm.WriteModel.Reduce()
}
