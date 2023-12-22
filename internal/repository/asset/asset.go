package asset

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	AddedEventType   = "asset.added"
	RemovedEventType = "asset.removed"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	StoreKey string `json:"storeKey"`
}

func (e *AddedEvent) Payload() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func AddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "ASSET-1WEAx", "unable to unmarshal asset")
	}

	return e, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	StoreKey string `json:"storeKey"`
}

func (e *RemovedEvent) Payload() interface{} {
	return e
}

func (e *RemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func RemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "ASSET-1m9PP", "unable to unmarshal asset")
	}

	return e, nil
}
