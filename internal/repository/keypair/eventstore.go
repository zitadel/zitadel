package keypair

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, AddedEventType, AddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, AddedCertificateEventType, AddedCertificateEventMapper)
}
