package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/idpconfig"
)

type JWTConfigReadModel struct {
	eventstore.ReadModel

	IDPConfigID  string
	JWTEndpoint  string
	Issuer       string
	KeysEndpoint string
}

func (rm *JWTConfigReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *idpconfig.JWTConfigAddedEvent:
			rm.reduceConfigAddedEvent(e)
		case *idpconfig.JWTConfigChangedEvent:
			rm.reduceConfigChangedEvent(e)
		}
	}

	return rm.ReadModel.Reduce()
}

func (rm *JWTConfigReadModel) reduceConfigAddedEvent(e *idpconfig.JWTConfigAddedEvent) {
	rm.IDPConfigID = e.IDPConfigID
	rm.JWTEndpoint = e.JWTEndpoint
	rm.Issuer = e.Issuer
	rm.KeysEndpoint = e.KeysEndpoint
}

func (rm *JWTConfigReadModel) reduceConfigChangedEvent(e *idpconfig.JWTConfigChangedEvent) {
	if e.JWTEndpoint != nil {
		rm.JWTEndpoint = *e.JWTEndpoint
	}
	if e.Issuer != nil {
		rm.Issuer = *e.Issuer
	}
	if e.KeysEndpoint != nil {
		rm.KeysEndpoint = *e.KeysEndpoint
	}
}
