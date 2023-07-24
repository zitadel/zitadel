package oidcsession

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, AddedType, eventstore.GenericEventMapper[AddedEvent]).
		RegisterFilterEventMapper(AggregateType, AccessTokenAddedType, eventstore.GenericEventMapper[AccessTokenAddedEvent]).
		RegisterFilterEventMapper(AggregateType, AccessTokenRevokedType, eventstore.GenericEventMapper[AccessTokenRevokedEvent]).
		RegisterFilterEventMapper(AggregateType, RefreshTokenAddedType, eventstore.GenericEventMapper[RefreshTokenAddedEvent]).
		RegisterFilterEventMapper(AggregateType, RefreshTokenRenewedType, eventstore.GenericEventMapper[RefreshTokenRenewedEvent]).
		RegisterFilterEventMapper(AggregateType, RefreshTokenRevokedType, eventstore.GenericEventMapper[RefreshTokenRevokedEvent])

}
