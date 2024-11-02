package user

import (
	"github.com/zitadel/zitadel/internal/v2/avatar"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type HumanAvatarAddedEvent eventstore.Event[avatar.AddedPayload]

const HumanAvatarAddedType = humanPrefix + avatar.AvatarAddedTypeSuffix

var _ eventstore.TypeChecker = (*HumanAvatarAddedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *HumanAvatarAddedEvent) ActionType() string {
	return HumanAvatarAddedType
}

func HumanAvatarAddedEventFromStorage(event *eventstore.StorageEvent) (e *HumanAvatarAddedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-ddQaI", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[avatar.AddedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &HumanAvatarAddedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}

type HumanAvatarRemovedEvent eventstore.Event[avatar.RemovedPayload]

const HumanAvatarRemovedType = humanPrefix + avatar.AvatarRemovedTypeSuffix

var _ eventstore.TypeChecker = (*HumanAvatarRemovedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *HumanAvatarRemovedEvent) ActionType() string {
	return HumanAvatarRemovedType
}

func HumanAvatarRemovedEventFromStorage(event *eventstore.StorageEvent) (e *HumanAvatarRemovedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-j2CkY", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[avatar.RemovedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &HumanAvatarRemovedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
