package setup

import (
	"context"
	_ "embed"
	"slices"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type FillFieldsForProjectGrant struct {
	eventstore *eventstore.Eventstore
}

func (mig *FillFieldsForProjectGrant) Execute(ctx context.Context, _ eventstore.Event) error {
	filler := &fieldFiller{
		eventstore: mig.eventstore,
	}
	// TODO: query instance ids first
	// TODO: implement batch wise filling
	err := mig.eventstore.FilterToReducer(
		ctx,
		eventstore.
			NewSearchQueryBuilder(eventstore.ColumnsEvent).
			AddQuery().AggregateTypes("project", "org").
			Builder().
			OrderAsc(),
		filler,
	)
	if err != nil {
		return err
	}

	return filler.fill(ctx)
}

func (mig *FillFieldsForProjectGrant) String() string {
	return "28_add_search_table"
}

type fieldFiller struct {
	eventstore *eventstore.Eventstore

	events []eventstore.FillFieldsEvent
}

func (filler *fieldFiller) Reduce() error {
	return nil
}

func (filler *fieldFiller) fill(ctx context.Context) error {
	return filler.eventstore.FillFields(ctx, filler.events...)
}

func (filler *fieldFiller) AppendEvents(events ...eventstore.Event) {
	filler.events = slices.Grow(filler.events, len(events))
	for _, event := range events {
		filler.events = append(filler.events, event.(eventstore.FillFieldsEvent))
	}
}
