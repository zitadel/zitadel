package query

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/idpconfig"
)

type OIDCConfigReadModel struct {
	eventstore.ReadModel

	IDPConfigID           string
	ClientID              string
	ClientSecret          *crypto.CryptoValue
	Issuer                string
	Scopes                []string
	IDPDisplayNameMapping domain.OIDCMappingField
	UserNameMapping       domain.OIDCMappingField
}

func (rm *OIDCConfigReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *idpconfig.OIDCConfigAddedEvent:
			rm.reduceConfigAddedEvent(e)
		case *idpconfig.OIDCConfigChangedEvent:
			rm.reduceConfigChangedEvent(e)
		}
	}

	return rm.ReadModel.Reduce()
}

func (rm *OIDCConfigReadModel) reduceConfigAddedEvent(e *idpconfig.OIDCConfigAddedEvent) {
	rm.IDPConfigID = e.IDPConfigID
	rm.ClientID = e.ClientID
	rm.ClientSecret = e.ClientSecret
	rm.Issuer = e.Issuer
	rm.Scopes = e.Scopes
	rm.IDPDisplayNameMapping = e.IDPDisplayNameMapping
	rm.UserNameMapping = e.UserNameMapping
}

func (rm *OIDCConfigReadModel) reduceConfigChangedEvent(e *idpconfig.OIDCConfigChangedEvent) {
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
