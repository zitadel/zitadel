package usergrant

import (
	"github.com/caos/zitadel/internal/eventstore"
)

const (
	AggregateType    = "usergrant"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}
