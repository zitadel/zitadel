package sessionlogout

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	BackChannelLogoutRegisteredEventMapper = eventstore.GenericEventMapper[BackChannelLogoutRegisteredEvent]
	BackChannelLogoutSentEventMapper       = eventstore.GenericEventMapper[BackChannelLogoutSentEvent]
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, BackChannelLogoutRegisteredType, BackChannelLogoutRegisteredEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, BackChannelLogoutSentType, BackChannelLogoutSentEventMapper)
}
