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
	*eventstore.BaseEvent `json:"-"`

	Value T
}

func (e *SetEvent[T]) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *SetEvent[T]) Payload() interface{} {
	return e
}

func (e *SetEvent[T]) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SetEventType interface {
	Boolean
	FeatureType() domain.FeatureType
}

type EventType[T SetEventType] struct {
	eventstore.EventType
}

type Boolean struct {
	Boolean bool
}

func (b Boolean) FeatureType() domain.FeatureType {
	return domain.FeatureTypeBoolean
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
