package instance

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	InstanceAddedEventType   = instanceEventTypePrefix + "added"
	InstanceChangedEventType = instanceEventTypePrefix + "changed"
	InstanceRemovedEventType = instanceEventTypePrefix + "removed"
)

type InstanceAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name string `json:"name,omitempty"`
}

func (e *InstanceAddedEvent) Data() interface{} {
	return e
}

func (e *InstanceAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewInstanceAddedEvent(ctx context.Context, aggregate *eventstore.Aggregate, name string) *InstanceAddedEvent {
	return &InstanceAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			InstanceAddedEventType,
		),
		Name: name,
	}
}

func InstanceAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	instanceAdded := &InstanceAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, instanceAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "INSTANCE-s9l3F", "unable to unmarshal instance added")
	}

	return instanceAdded, nil
}

type InstanceChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name string `json:"name,omitempty"`
}

func (e *InstanceChangedEvent) Data() interface{} {
	return e
}

func (e *InstanceChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewInstanceChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, newName string) *InstanceChangedEvent {
	return &InstanceChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			InstanceChangedEventType,
		),
		Name: newName,
	}
}

func InstanceChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	instanceChanged := &InstanceChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, instanceChanged)
	if err != nil {
		return nil, errors.ThrowInternal(err, "INSTANCE-3hfo8", "unable to unmarshal instance changed")
	}

	return instanceChanged, nil
}

type InstanceRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	name                 string
	domains              []string
}

func (e *InstanceRemovedEvent) Data() interface{} {
	return nil
}

func (e *InstanceRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	constraints := make([]*eventstore.EventUniqueConstraint, len(e.domains)+1)
	for i, domain := range e.domains {
		constraints[i] = NewRemoveInstanceDomainUniqueConstraint(domain)
	}
	constraints[len(e.domains)] = eventstore.NewRemoveInstanceUniqueConstraints()
	return constraints
}

func NewInstanceRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate, name string, domains []string) *InstanceRemovedEvent {
	return &InstanceRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			InstanceRemovedEventType,
		),
		name:    name,
		domains: domains,
	}
}

func InstanceRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	return &InstanceRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
