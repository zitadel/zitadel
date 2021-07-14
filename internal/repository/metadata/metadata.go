package metadata

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	SetEventType     = "metadata.set"
	RemovedEventType = "metadata.removed"
)

type SetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Key   string `json:"key"`
	Value string `json:"value"`
}

func (e *SetEvent) Data() interface{} {
	return e
}

func (e *SetEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewSetEvent(
	base *eventstore.BaseEvent,
	key,
	value string,
) *SetEvent {
	return &SetEvent{
		BaseEvent: *base,
		Key:       key,
		Value:     value,
	}
}

func SetEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &SetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "META-3n9fs", "unable to unmarshal metadata set")
	}

	return e, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Key string `json:"key"`
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
		Key:       key,
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "META-2m99f", "unable to unmarshal metadata removed")
	}

	return e, nil
}
