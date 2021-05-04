package user

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/asset"
)

const (
	avatarEventPrefix      = humanEventPrefix + "avatar."
	HumanAvatarAddedType   = phoneEventPrefix + "added"
	HumanAvatarRemovedType = phoneEventPrefix + "removed"
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

func HumanAvatarAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
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

func HumanAvatarRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := asset.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &HumanAvatarRemovedEvent{RemovedEvent: *e.(*asset.RemovedEvent)}, nil
}
