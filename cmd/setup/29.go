package setup

import (
	"context"
	_ "embed"
	"math"
	"slices"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type FillFieldsForProjectGrant struct {
	eventstore *eventstore.Eventstore
	config     *FillFields
}

func (mig *FillFieldsForProjectGrant) Execute(ctx context.Context, _ eventstore.Event) error {
	// TODO: query instance ids first
	filler := &fieldFiller{
		eventstore: mig.eventstore,
		batchSize:  mig.config.BatchSize,
	}

	return filler.fillInstance(
		ctx,
		eventstore.
			NewSearchQueryBuilder(eventstore.ColumnsEvent).
			AddQuery().AggregateTypes("project", "org").
			Builder().
			Limit(uint64(mig.config.BatchSize)).
			OrderAsc(),
	)
}

func (mig *FillFieldsForProjectGrant) String() string {
	return "29_fill_fields_for_project_grant"
}

type fieldFiller struct {
	eventstore *eventstore.Eventstore
	batchSize  uint32

	events []eventstore.FillFieldsEvent
}

func (filler *fieldFiller) Reduce() error {
	return nil
}

func (filler *fieldFiller) fillInstance(ctx context.Context, query *eventstore.SearchQueryBuilder) error {
	var position float64
	for {
		err := filler.eventstore.FilterToReducer(
			ctx,
			query.PositionAfter(position).Limit(uint64(filler.batchSize)),
			filler,
		)
		if err != nil {
			return err
		}
		err = filler.eventstore.FillFields(ctx, filler.events...)
		if err != nil {
			return err
		}
		if len(filler.events) < int(filler.batchSize) {
			return nil
		}
		// the math is needed because eventstore only provides position after filter and we need to ensure to miss any events
		position = math.Float64frombits(math.Float64bits(filler.events[len(filler.events)-1].Position()) - 10)
		filler.events = nil
	}
}

func (filler *fieldFiller) AppendEvents(events ...eventstore.Event) {
	filler.events = slices.Grow(filler.events, len(events))
	for _, event := range events {
		filler.events = append(filler.events, event.(eventstore.FillFieldsEvent))
	}
}
