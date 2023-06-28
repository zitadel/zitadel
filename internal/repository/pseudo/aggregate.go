package pseudo

import "github.com/zitadel/zitadel/internal/eventstore"

const (
	AggregateType    = "pseudo"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}

func NewAggregate() *Aggregate {
	return &Aggregate{
		Aggregate: eventstore.Aggregate{
			Type:    AggregateType,
			Version: AggregateVersion,
		},
	}
}
