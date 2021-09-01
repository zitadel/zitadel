package action

import "github.com/caos/zitadel/internal/eventstore"

const (
	AggregateType    = "action"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}
