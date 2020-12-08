package idpprovider

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/idp/provider"
)

type IDPProviderAddedEvent struct {
	provider.AddedEvent
}

func NewIDPProviderAddedEvent(
	base *eventstore.BaseEvent,
	idpConfigID string,
	idpProviderType provider.Type,
) *IDPProviderAddedEvent {

	return &IDPProviderAddedEvent{
		AddedEvent: *provider.NewAddedEvent(
			base,
			idpConfigID,
			idpProviderType),
	}
}

func IDPProviderAddedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := provider.AddedEventEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPProviderAddedEvent{
		AddedEvent: *e.(*provider.AddedEvent),
	}, nil
}

type IDPProviderRemovedEvent struct {
	provider.RemovedEvent
}

func NewIDPProviderRemovedEvent(
	base *eventstore.BaseEvent,
	idpConfigID string,
) *IDPProviderRemovedEvent {

	return &IDPProviderRemovedEvent{
		RemovedEvent: *provider.NewRemovedEvent(base, idpConfigID),
	}
}

func IDPProviderRemovedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := provider.RemovedEventEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPProviderRemovedEvent{
		RemovedEvent: *e.(*provider.RemovedEvent),
	}, nil
}
