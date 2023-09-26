package feature

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	DefaultLoginInstanceEventType = EventTypeFromFeature(domain.FeatureLoginDefaultOrg)
)

func EventTypeFromFeature(feature domain.Feature) eventstore.EventType {
	return eventTypePrefix + eventstore.EventType(strings.ToLower(feature.String())) + setSuffix
}

type SetEvent[T SetEventType] struct {
	*eventstore.BaseEvent

	T T
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
	Boolean bool
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
