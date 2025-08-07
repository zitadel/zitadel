package system

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, IDGeneratedType, eventstore.GenericEventMapper[IDGeneratedEvent])
}

const IDGeneratedType = AggregateType + ".id.generated"

type IDGeneratedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID string `json:"id"`
}

func (e *IDGeneratedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = *b
}

func (e *IDGeneratedEvent) Payload() interface{} {
	return e
}

func (e *IDGeneratedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewIDGeneratedEvent(
	ctx context.Context,
	id string,
) *IDGeneratedEvent {
	return &IDGeneratedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			eventstore.NewAggregate(ctx, AggregateOwner, AggregateType, "v1"),
			IDGeneratedType),
		ID: id,
	}
}
