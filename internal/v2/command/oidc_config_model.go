package command

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/idpconfig"
)

type OIDCConfigWriteModel struct {
	eventstore.WriteModel

	IDPConfigID  string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Issuer       string
	Scopes       []string

	IDPDisplayNameMapping domain.OIDCMappingField
	UserNameMapping       domain.OIDCMappingField
	State                 domain.IDPConfigState
}

func (wm *OIDCConfigWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idpconfig.OIDCConfigAddedEvent:
			wm.reduceConfigAddedEvent(e)
		case *idpconfig.OIDCConfigChangedEvent:
			wm.reduceConfigChangedEvent(e)
		case *idpconfig.IDPConfigDeactivatedEvent:
			wm.State = domain.IDPConfigStateInactive
		case *idpconfig.IDPConfigReactivatedEvent:
			wm.State = domain.IDPConfigStateActive
		case *idpconfig.IDPConfigRemovedEvent:
			wm.State = domain.IDPConfigStateRemoved
		}
	}

	return wm.WriteModel.Reduce()
}

func (wm *OIDCConfigWriteModel) reduceConfigAddedEvent(e *idpconfig.OIDCConfigAddedEvent) {
	wm.IDPConfigID = e.IDPConfigID
	wm.ClientID = e.ClientID
	wm.ClientSecret = e.ClientSecret
	wm.Issuer = e.Issuer
	wm.Scopes = e.Scopes
	wm.IDPDisplayNameMapping = e.IDPDisplayNameMapping
	wm.UserNameMapping = e.UserNameMapping
	wm.State = domain.IDPConfigStateActive
}

func (wm *OIDCConfigWriteModel) reduceConfigChangedEvent(e *idpconfig.OIDCConfigChangedEvent) {
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
