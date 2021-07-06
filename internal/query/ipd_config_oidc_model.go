package query

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/idpconfig"
)

type IDPOIDCConfigReadModel struct {
	eventstore.ReadModel

	IDPConfigID           string
	ClientID              string
	ClientSecret          *crypto.CryptoValue
	Issuer                string
	AuthorizationEndpoint string
	TokenEndpoint         string
	Scopes                []string
	IDPDisplayNameMapping domain.OIDCMappingField
	UserNameMapping       domain.OIDCMappingField
}

func (rm *IDPOIDCConfigReadModel) Reduce() error {
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

func (rm *IDPOIDCConfigReadModel) reduceConfigAddedEvent(e *idpconfig.OIDCConfigAddedEvent) {
	rm.IDPConfigID = e.IDPConfigID
	rm.ClientID = e.ClientID
	rm.ClientSecret = e.ClientSecret
	rm.Issuer = e.Issuer
	rm.AuthorizationEndpoint = e.AuthorizationEndpoint
	rm.TokenEndpoint = e.TokenEndpoint
	rm.Scopes = e.Scopes
	rm.IDPDisplayNameMapping = e.IDPDisplayNameMapping
	rm.UserNameMapping = e.UserNameMapping
}

func (rm *IDPOIDCConfigReadModel) reduceConfigChangedEvent(e *idpconfig.OIDCConfigChangedEvent) {
	if e.ClientID != nil {
		rm.ClientID = *e.ClientID
	}
	if e.Issuer != nil {
		rm.Issuer = *e.Issuer
	}
	if e.AuthorizationEndpoint != nil {
		rm.AuthorizationEndpoint = *e.AuthorizationEndpoint
	}
	if e.TokenEndpoint != nil {
		rm.TokenEndpoint = *e.TokenEndpoint
	}
	if len(e.Scopes) > 0 {
		rm.Scopes = e.Scopes
	}
	if e.IDPDisplayNameMapping != nil && e.IDPDisplayNameMapping.Valid() {
		rm.IDPDisplayNameMapping = *e.IDPDisplayNameMapping
	}
	if e.UserNameMapping != nil && e.UserNameMapping.Valid() {
		rm.UserNameMapping = *e.UserNameMapping
	}
}
