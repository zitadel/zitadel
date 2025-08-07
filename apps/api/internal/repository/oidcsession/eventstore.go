package oidcsession

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, AddedType, eventstore.GenericEventMapper[AddedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, AccessTokenAddedType, eventstore.GenericEventMapper[AccessTokenAddedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, AccessTokenRevokedType, eventstore.GenericEventMapper[AccessTokenRevokedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, RefreshTokenAddedType, eventstore.GenericEventMapper[RefreshTokenAddedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, RefreshTokenRenewedType, eventstore.GenericEventMapper[RefreshTokenRenewedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, RefreshTokenRevokedType, eventstore.GenericEventMapper[RefreshTokenRevokedEvent])

}
