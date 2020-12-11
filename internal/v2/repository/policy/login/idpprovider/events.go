package idpprovider

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/idp/provider"
)

type AddedEvent struct {
	provider.AddedEvent
}

func NewAddedEvent(
	base *eventstore.BaseEvent,
	idpConfigID string,
	idpProviderType provider.Type,
) *AddedEvent {

	return &AddedEvent{
		AddedEvent: *provider.NewAddedEvent(
			base,
			idpConfigID,
			idpProviderType),
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := provider.AddedEventEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &AddedEvent{
		AddedEvent: *e.(*provider.AddedEvent),
	}, nil
}

type RemovedEvent struct {
	provider.RemovedEvent
}

func NewRemovedEvent(
	base *eventstore.BaseEvent,
	idpConfigID string,
) *RemovedEvent {
	return &RemovedEvent{
		RemovedEvent: *provider.NewRemovedEvent(base, idpConfigID),
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := provider.RemovedEventEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &RemovedEvent{
		RemovedEvent: *e.(*provider.RemovedEvent),
	}, nil
}
