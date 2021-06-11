package asset

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	AddedEventType   = "asset.added"
	RemovedEventType = "asset.removed"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	StoreKey string `json:"storeKey"`
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewAddedEvent(
	base *eventstore.BaseEvent,
	key string,
) *AddedEvent {

	return &AddedEvent{
		BaseEvent: *base,
		StoreKey:  key,
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ASSET-1WEAx", "unable to unmarshal asset")
	}

	return e, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	StoreKey string `json:"storeKey"`
}

func (e *RemovedEvent) Data() interface{} {
	return e
}

func (e *RemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewRemovedEvent(
	base *eventstore.BaseEvent,
	key string,
) *RemovedEvent {

	return &RemovedEvent{
		BaseEvent: *base,
		StoreKey:  key,
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ASSET-1m9PP", "unable to unmarshal asset")
	}

	return e, nil
}
