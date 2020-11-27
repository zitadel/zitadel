package eventstore

import "context"

func PushAggregate(ctx context.Context, es *Eventstore, writeModel queryReducer, aggregate *Aggregate) error {
	err := es.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return err
	}

	events, err := es.PushAggregates(ctx, aggregate)
	if err != nil {
		return err
	}

	writeModel.AppendEvents(events...)
	return writeModel.Reduce()
}
