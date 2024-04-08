package idpintent

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, StartedEventType, StartedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, SucceededEventType, SucceededEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, SAMLSucceededEventType, SAMLSucceededEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, SAMLRequestEventType, SAMLRequestEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, LDAPSucceededEventType, LDAPSucceededEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, FailedEventType, FailedEventMapper)
}
