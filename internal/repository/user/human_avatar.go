package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/asset"
)

const (
	avatarEventPrefix      = humanEventPrefix + "avatar."
	HumanAvatarAddedType   = avatarEventPrefix + "added"
	HumanAvatarRemovedType = avatarEventPrefix + "removed"
)

type HumanAvatarAddedEvent struct {
	asset.AddedEvent
}

func NewHumanAvatarAddedEvent(ctx context.Context, aggregate *eventstore.Aggregate, storageKey string) *HumanAvatarAddedEvent {
	return &HumanAvatarAddedEvent{
		AddedEvent: *asset.NewAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				HumanAvatarAddedType),
			storageKey),
	}
}

func HumanAvatarAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := asset.AddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &HumanAvatarAddedEvent{AddedEvent: *e.(*asset.AddedEvent)}, nil
}

type HumanAvatarRemovedEvent struct {
	asset.RemovedEvent
}

func NewHumanAvatarRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate, storageKey string) *HumanAvatarRemovedEvent {
	return &HumanAvatarRemovedEvent{
		RemovedEvent: *asset.NewRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				HumanAvatarRemovedType),
			storageKey),
	}
}

func HumanAvatarRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := asset.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &HumanAvatarRemovedEvent{RemovedEvent: *e.(*asset.RemovedEvent)}, nil
}
