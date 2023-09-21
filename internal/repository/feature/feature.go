package feature

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	DefaultLoginInstanceEventType = eventTypePrefix + "default_login_instance" + setSuffix
)

type SetEvent[T SetEventType] struct {
	*eventstore.BaseEvent

	Type T
}

func (e *SetEvent[T]) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *SetEvent[T]) Data() interface{} {
	return e
}

func (e *SetEvent[T]) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

type SetEventType interface {
	Boolean
}

type Boolean struct {
	B bool
}

func NewSetEvent[T SetEventType](
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	eventType eventstore.EventType,
	setType T,
) *SetEvent[T] {
	return &SetEvent[T]{
		eventstore.NewBaseEventForPush(
			ctx, aggregate, eventType),
		setType,
	}
}
