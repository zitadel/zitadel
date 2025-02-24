package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

func writeModelToObjectDetails(writeModel *eventstore.WriteModel) *domain.ObjectDetails {
	return &domain.ObjectDetails{
		Sequence:      writeModel.ProcessedSequence,
		ResourceOwner: writeModel.ResourceOwner,
		EventDate:     writeModel.ChangeDate,
		ID:            writeModel.AggregateID,
	}
}

func pushedEventsToObjectDetails(events []eventstore.Event) *domain.ObjectDetails {
	if len(events) == 0 {
		return &domain.ObjectDetails{}
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}
}
