package user

import (
	"github.com/zitadel/zitadel/internal/v2/avatar"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type HumanAvatarAddedEvent humanAvatarAddedEvent
type humanAvatarAddedEvent = eventstore.Event[avatar.AddedPayload]

func HumanAvatarAddedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*HumanAvatarAddedEvent, error) {
	event, err := eventstore.EventFromStorage[humanAvatarAddedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*HumanAvatarAddedEvent)(event), nil
}

type HumanAvatarRemovedEvent humanAvatarRemovedEvent
type humanAvatarRemovedEvent = eventstore.Event[avatar.RemovedPayload]

func HumanAvatarRemovedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*HumanAvatarRemovedEvent, error) {
	event, err := eventstore.EventFromStorage[humanAvatarRemovedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*HumanAvatarRemovedEvent)(event), nil
}
