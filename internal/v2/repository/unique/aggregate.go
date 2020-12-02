package user

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

type Aggregate struct {
	eventstore.Aggregate
}

func NewAggregate(
	id,
	resourceOwner string,
	previousSequence uint64,
	aggregateType eventstore.AggregateType,
	aggregateVersion eventstore.Version,
) *Aggregate {

	return &Aggregate{
		Aggregate: *eventstore.NewAggregate(
			id,
			aggregateType,
			resourceOwner,
			aggregateVersion,
			previousSequence,
		),
	}
}

func AggregateFromWriteModel(wm *eventstore.WriteModel, aggregateType eventstore.AggregateType, aggregateVersion eventstore.Version) *Aggregate {
	return &Aggregate{
		Aggregate: *eventstore.AggregateFromWriteModel(wm, aggregateType, aggregateVersion),
	}
}

func AggregateFromReadModel(rm *ReadModel, aggregateType eventstore.AggregateType, aggregateVersion eventstore.Version) *Aggregate {
	return &Aggregate{
		Aggregate: *eventstore.NewAggregate(
			rm.AggregateID,
			aggregateType,
			rm.ResourceOwner,
			aggregateVersion,
			rm.ProcessedSequence,
		),
	}
}
