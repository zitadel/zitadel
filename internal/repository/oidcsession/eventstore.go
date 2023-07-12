package oidcsession

import "github.com/zitadel/zitadel/internal/eventstore"

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, AddedType, AddedEventMapper).
		RegisterFilterEventMapper(AggregateType, AccessTokenAddedType, AccessTokenAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, RefreshTokenAddedType, RefreshTokenAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, RefreshTokenRenewedType, RefreshTokenRenewedEventMapper)

}
