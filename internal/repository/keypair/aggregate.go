package keypair

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

const (
	AggregateType    = "key_pair"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}
