package metadata

import (
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	SetEventType        = "metadata.set"
	RemovedEventType    = "metadata.removed"
	RemovedAllEventType = "metadata.removed.all"
)

type SetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Key   string `json:"key"`
	Value []byte `json:"value"`
}

func (e *SetEvent) Payload() interface{} {
	return e
}

func (e *SetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewSetEvent(
	base *eventstore.BaseEvent,
	key string,
	value []byte,
) *SetEvent {
	return &SetEvent{
		BaseEvent: *base,
		Key:       key,
		Value:     value,
	}
}

func SetEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &SetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "META-3n9fs", "unable to unmarshal metadata set")
	}

	return e, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Key string `json:"key"`
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
		Key:       key,
	}
}

func RemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "META-2m99f", "unable to unmarshal metadata removed")
	}

	return e, nil
}

type RemovedAllEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *RemovedAllEvent) Payload() interface{} {
	return nil
}

func (e *RemovedAllEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewRemovedAllEvent(
	base *eventstore.BaseEvent,
) *RemovedAllEvent {

	return &RemovedAllEvent{
		BaseEvent: *base,
	}
}

func RemovedAllEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &RemovedAllEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
