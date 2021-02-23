package usergrant

import (
	"github.com/caos/zitadel/internal/eventstore"
)

const (
	AggregateType    = "key_pair"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}
