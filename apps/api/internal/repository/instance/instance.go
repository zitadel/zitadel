package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
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

func (e *InstanceAddedEvent) Payload() interface{} {
	return e
}

func (e *InstanceAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func InstanceAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	instanceAdded := &InstanceAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(instanceAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "INSTANCE-s9l3F", "unable to unmarshal instance added")
	}

	return instanceAdded, nil
}

type InstanceChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name string `json:"name,omitempty"`
}

func (e *InstanceChangedEvent) Payload() interface{} {
	return e
}

func (e *InstanceChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func InstanceChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	instanceChanged := &InstanceChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(instanceChanged)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "INSTANCE-3hfo8", "unable to unmarshal instance changed")
	}

	return instanceChanged, nil
}

type InstanceRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	name                 string
	domains              []string
}

func (e *InstanceRemovedEvent) Payload() interface{} {
	return nil
}

func (e *InstanceRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	constraints := make([]*eventstore.UniqueConstraint, len(e.domains)+1)
	for i, domain := range e.domains {
		constraints[i] = NewRemoveInstanceDomainUniqueConstraint(domain)
	}
	constraints[len(e.domains)] = eventstore.NewRemoveInstanceUniqueConstraints()
	return constraints
}

func (e *InstanceRemovedEvent) Fields() []*eventstore.FieldOperation {
	return []*eventstore.FieldOperation{
		eventstore.RemoveSearchFields(map[eventstore.FieldType]any{
			eventstore.FieldTypeInstanceID: e.Aggregate().ID,
		}),
	}
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

func InstanceRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &InstanceRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
