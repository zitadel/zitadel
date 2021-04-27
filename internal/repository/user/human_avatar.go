package user

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	avatarEventPrefix      = humanEventPrefix + "avatar."
	HumanAvatarChangedType = avatarEventPrefix + "changed"
	HumanAvatarRemovedType = avatarEventPrefix + "removed"
)

type HumanAvatarChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AssetID string `json:"assetID,omitempty"`
	Avatar  []byte `json:"-"`
}

func (e *HumanAvatarChangedEvent) Data() interface{} {
	return e
}

func (e *HumanAvatarChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *HumanAvatarChangedEvent) Assets() []*eventstore.Asset {
	return []*eventstore.Asset{
		{
			ID:     e.AssetID,
			Asset:  e.Avatar,
			Action: eventstore.AssetAdd,
		},
	}
}

func NewHumanAvatarChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	assetID string,
	avatar []byte,
) *HumanAvatarChangedEvent {
	return &HumanAvatarChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanAvatarChangedType,
		),
		AssetID: assetID,
		Avatar:  avatar,
	}
}

func HumanAvatarAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	humanAvatar := &HumanAvatarChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, humanAvatar)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-3M9fs", "unable to unmarshal human avatar changed")
	}

	return humanAvatar, nil
}

type HumanAvatarRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AssetID string `json:"assetID,omitempty"`
}

func (e *HumanAvatarRemovedEvent) Data() interface{} {
	return e
}

func (e *HumanAvatarRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *HumanAvatarRemovedEvent) Assets() []*eventstore.Asset {
	return []*eventstore.Asset{
		{
			ID:     e.AssetID,
			Action: eventstore.AssetRemove,
		},
	}
}

func NewHumanAvatarRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	assetID string,
) *HumanAvatarRemovedEvent {
	return &HumanAvatarRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanAvatarRemovedType,
		),
		AssetID: assetID,
	}
}

func HumanAvatarRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	humanAvatar := &HumanAvatarRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, humanAvatar)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-2M9ds", "unable to unmarshal human avatar removed")
	}

	return humanAvatar, nil
}
