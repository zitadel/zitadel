package user

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	AggregateType    = "user"
	AggregateVersion = "v2"
)

type Aggregate struct {
	eventstore.Aggregate
}

func NewAggregate(
	id,
	resourceOwner string,
	previousSequence uint64,
) *Aggregate {

	return &Aggregate{
		Aggregate: *eventstore.NewAggregate(
			id,
			AggregateType,
			resourceOwner,
			AggregateVersion,
			previousSequence,
		),
	}
}

func AggregateFromWriteModel(wm *eventstore.WriteModel) *Aggregate {
	return &Aggregate{
		Aggregate: *eventstore.AggregateFromWriteModel(wm, AggregateType, AggregateVersion),
	}
}

func AggregateFromReadModel(rm *ReadModel) *Aggregate {
	return &Aggregate{
		Aggregate: *eventstore.NewAggregate(
			rm.AggregateID,
			AggregateType,
			rm.ResourceOwner,
			AggregateVersion,
			rm.ProcessedSequence,
		),
	}
}
