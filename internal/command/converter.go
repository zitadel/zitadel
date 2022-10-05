package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

func writeModelToObjectDetails(writeModel *eventstore.WriteModel) *domain.ObjectDetails {
	return &domain.ObjectDetails{
		ResourceOwner: writeModel.ResourceOwner,
		EventDate:     writeModel.ChangeDate,
	}
}

func pushedEventsToObjectDetails(events []eventstore.Event) *domain.ObjectDetails {
	return &domain.ObjectDetails{
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}
}
